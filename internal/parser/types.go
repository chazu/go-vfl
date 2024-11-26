package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	//Custom lexer
	l = lexer.MustSimple([]lexer.SimpleRule{
		{"Ident", "[a-zA-Z_][a-zA-Z0-9_]*"},
		{"Number", "[0-9]+"},
		{"Punctuation", `[][.,]`},
		{"Relation", "==|>=|<="},
		{"Space", " +"},
		{"At", "@"},
		{"OpenParen", "\\("},
		{"CloseParen", "\\)"},
		{"Colon", ":"},
		{"Pipe", "\\|"},
		{"Dash", "-"},
	})
)

// Private structs are used with participle to parse programs - we then reify
// these into a less clumsy AST with ProgramAST. I'm sure participle provides
// ways to do this, but I'm proceeding with the naive implementation in the
// interest of actually finishing an MVP. Sue me, I have very little free time.
type program struct {
	Orientation                 orientation                 `@@?`
	LeadingSuperviewConnection  leadingSuperviewConnection  `@@?`
	Views                       []view                      `@@+`
	TrailingSuperviewConnection trailingSuperviewConnection `@@?`
}

type view struct {
	Connection connection    `@@?`
	Name       string        `"[" @Ident`
	Predicate  predicateList `(Space? @@)? "]"`
}

func (v *view) Predicates() []predicate {
	res := []predicate{}
	if len(v.Predicate.Predicates) != 0 {
		// Here be predicates
		res = v.Predicate.Predicates
	}
	if v.Predicate.Predicate != (predicate{}) {
		// It ain't zero
		res = append(res, v.Predicate.Predicate)
	}

	return res
}

type relation struct {
	Gte bool ` @">="`
	Lte bool `| @"<="`
	Eq  bool `| @"=="`
}

type predicateObject struct {
	Number   int    `@Number`
	ViewName string `| @Ident`
}

type predicate struct {
	Relation relation        `@@?`
	Object   predicateObject `@@`
	Priority priority        `@@?`
}

type priority struct {
	Value int `At @Number`
}

type predicateList struct {
	Predicates []predicate `OpenParen @@ ("," @@)* CloseParen`
	Predicate  predicate   `| @@`
}

type orientation struct {
	Direction *string `(@"H"? @"V"?)! Colon`
}

type leadingSuperviewConnection struct {
	SuperView  bool       `@Pipe`
	Connection connection `@@?`
}

type trailingSuperviewConnection struct {
	Connection connection `@@?`
	SuperView  bool       `@Pipe`
}

type connection struct {
	Predicates predicateList `Dash (@@ Dash)?`
}

// Public structs - used to hide the jankiness of the parser structs
type ProgramAST struct {
	Orientation                 orientation
	LeadingSuperviewConnection  ConnectionAST
	TrailingSuperviewConnection ConnectionAST
	Views                       []ViewAST
}

type ConnectionAST struct {
	Predicates  []predicate
	IsSuperview bool
}

type ViewAST struct {
	LeadingConnection  ConnectionAST
	TrailingConnection ConnectionAST
	Name               string
	Predicates         []predicate
}
