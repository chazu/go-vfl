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

type program struct {
	Orientation                 *orientation                 `@@?`
	LeadingSuperViewConnection  *leadingSuperViewConnection  `@@?`
	Views                       []*view                      `@@+`
	TrailingSuperViewConnection *trailingSuperViewConnection `@@?`
}

type view struct {
	Connection *connection   `@@?`
	Name       string        `"[" @Ident`
	Predicate  predicateList `(Space? @@)? "]"`
}

func (v *view) Predicates() []*predicate {
	if v.Predicate.Predicate != nil {
		return []*predicate{v.Predicate.Predicate}
	} else if v.Predicate.Predicates != nil {
		return v.Predicate.Predicates
	} else {
		return []*predicate{}
	}
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
	Relation *relation        `@@?`
	Object   *predicateObject `@@`
	Priority *priority        `@@?`
}

type priority struct {
	Value int `At @Number`
}

type predicateList struct {
	Predicates []*predicate `OpenParen @@ ("," @@)* CloseParen`
	Predicate  *predicate   `| @@`
}

type orientation struct {
	Direction *string `(@"H"? @"V"?)! Colon`
}

type leadingSuperViewConnection struct {
	SuperView  bool        `@Pipe`
	Connection *connection `@@?`
}

type trailingSuperViewConnection struct {
	Connection *connection `@@?`
	SuperView  bool        `@Pipe`
}

type connection struct {
	Predicates *predicateList `Dash (@@ Dash)?`
}
