package vfl

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Grammar structures for participle parsing
type Grammar struct {
	Orientation *string     `( @("H"|"V") ":" )?`
	Statements  []*GramStmt `@@+`
}

type GramStmt struct {
	SuperStart  bool       `@"|"?`
	StartConn   *GramConn  `@@?`  // Optional connection after superview start
	Chain       *GramItem  `@@`
	EndConn     *GramConn  `@@?`  // Optional connection before superview end
	SuperEnd    bool       `@"|"?`
}

// GramConn represents a connection (spacing) without a following item
type GramConn struct {
	Dash       string    `@"-"`
	Spacing    *float64  `( @Float | @Int )?`
	MetricName *string   `( @Ident )?`
	DashEnd    string    `@"-"?`
}

type GramItem struct {
	View *GramView  `@@`
	Next *GramChain `@@?`
}

type GramChain struct {
	Dash       string    `@"-"?`  // Make dash optional for flush views
	Spacing    *float64  `( @Float | @Int )?`
	MetricName *string   `( @Ident )?`  // Support metric names
	DashEnd    string    `@"-"?`  // Trailing dash after metric/spacing
	Item       *GramItem `@@?`
}

type GramView struct {
	Open       string      `"["`
	Name       string      `@Ident`
	Predicates []*GramPred `( "(" @@ ( "," @@ )* ")" )?`
	Close      string      `"]"`
}

type GramPred struct {
	Relation  *string  `( @( "==" | ">=" | "<=" ) )?`
	Value     float64  `( @Float | @Int )`
	Priority  *float64 `( "@" ( @Float | @Int ) )?`
}

var (
	vflLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Whitespace", Pattern: `\s+`},
		{Name: "Float", Pattern: `\d+\.\d+`},
		{Name: "Int", Pattern: `\d+`},
		{Name: "Ident", Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
		{Name: "Punct", Pattern: `[:|;\[\](),.@%]`},
		{Name: "Op", Pattern: `==|>=|<=`},
		{Name: "Dash", Pattern: `-`},
	})

	grammarParser = participle.MustBuild[Grammar](
		participle.Lexer(vflLexer),
		participle.Elide("Whitespace"),
		participle.UseLookahead(10),
	)
)

// parseGrammar parses VFL input using participle
func parseGrammar(input string) (*Grammar, error) {
	return grammarParser.ParseString("", input)
}

// convertGrammarToAST converts parsed grammar to AST
func convertGrammarToAST(g *Grammar, input string) *Program {
	prog := &Program{
		Input:      input,
		Statements: []Statement{},
	}

	// Set orientation
	if g.Orientation != nil {
		if *g.Orientation == "H" {
			prog.Orientation = Horizontal
		} else {
			prog.Orientation = Vertical
		}
	} else {
		prog.Orientation = Horizontal // default
	}

	// Convert statements
	for _, stmt := range g.Statements {
		s := convertStatement(stmt)
		prog.Statements = append(prog.Statements, s)
	}

	return prog
}

func convertStatement(gs *GramStmt) Statement {
	stmt := Statement{
		SuperviewStart: gs.SuperStart,
		SuperviewEnd:   gs.SuperEnd,
		Views:          []View{},
		Connections:    []Connection{},
	}

	// Add connection from superview if present
	if gs.StartConn != nil {
		conn := convertConnection(gs.StartConn)
		stmt.Connections = append(stmt.Connections, conn)
	}

	// Process the chain of views and connections
	item := gs.Chain
	for item != nil {
		if item.View != nil {
			stmt.Views = append(stmt.Views, convertView(item.View))
		}

		if item.Next != nil {
			// Add connection
			conn := Connection{}

			// Check for metric name first
			if item.Next.MetricName != nil {
				// Look up metric value
				if metricValue, ok := GetMetric(*item.Next.MetricName); ok {
					conn.Spacing = metricValue
				} else {
					// Unknown metric, use as identifier (will be resolved later)
					conn.Spacing = 8 // Default for now
					conn.MetricName = *item.Next.MetricName
				}
			} else if item.Next.Spacing != nil {
				conn.Spacing = *item.Next.Spacing
			} else if item.Next.Dash != "" {
				// Dash present but no spacing specified = default spacing
				conn.IsDefault = true
				conn.Spacing = 8 // Default spacing
			} else {
				// No dash = flush views (0 spacing)
				conn.Spacing = 0
			}
			stmt.Connections = append(stmt.Connections, conn)

			// Move to next item if it exists
			if item.Next.Item != nil {
				item = item.Next.Item
			} else {
				// Connection to superview
				break
			}
		} else {
			break
		}
	}

	// Add connection to superview if present
	if gs.EndConn != nil {
		conn := convertConnection(gs.EndConn)
		stmt.Connections = append(stmt.Connections, conn)
	}

	return stmt
}

func convertConnection(gc *GramConn) Connection {
	conn := Connection{}

	// Check for metric name first
	if gc.MetricName != nil {
		// Look up metric value
		if metricValue, ok := GetMetric(*gc.MetricName); ok {
			conn.Spacing = metricValue
		} else {
			// Unknown metric, store name for later resolution
			conn.MetricName = *gc.MetricName
			conn.Spacing = 8 // Default for now
		}
	} else if gc.Spacing != nil {
		conn.Spacing = *gc.Spacing
	} else {
		// Dash with no spacing = default
		conn.IsDefault = true
		conn.Spacing = 8
	}

	return conn
}

func convertView(gv *GramView) View {
	v := View{
		Name:       gv.Name,
		Predicates: []Predicate{},
	}

	for _, pred := range gv.Predicates {
		p := Predicate{}

		// Set relation (default to Equal if not specified)
		if pred.Relation != nil {
			switch *pred.Relation {
			case "==":
				p.Relation = Equal
			case ">=":
				p.Relation = GreaterEqual
			case "<=":
				p.Relation = LessEqual
			}
		} else {
			p.Relation = Equal
		}

		// Set value
		p.Value = Value{
			Type:     ConstantValue,
			Constant: pred.Value,
		}

		// Set priority if specified (on the predicate, not the view)
		if pred.Priority != nil {
			p.Priority = *pred.Priority
		}

		v.Predicates = append(v.Predicates, p)
	}

	return v
}

func convertPredicate(gp *GramPred) Predicate {
	p := Predicate{}

	// Set relation (default to Equal if not specified)
	if gp.Relation != nil {
		switch *gp.Relation {
		case "==":
			p.Relation = Equal
		case ">=":
			p.Relation = GreaterEqual
		case "<=":
			p.Relation = LessEqual
		}
	} else {
		p.Relation = Equal
	}

	// Set value
	p.Value = Value{
		Type:     ConstantValue,
		Constant: gp.Value,
	}

	return p
}

func convertValue(val float64) Value {
	return Value{
		Type:     ConstantValue,
		Constant: val,
	}
}

func convertExpression(ge interface{}) *Expression {
	// Placeholder for expression conversion
	return &Expression{}
}