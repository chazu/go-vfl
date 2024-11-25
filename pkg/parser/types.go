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

type Program struct {
	Orientation                 *Orientation                 `@@?`
	LeadingSuperViewConnection  *LeadingSuperViewConnection  `@@?`
	Views                       []*View                      `@@+`
	TrailingSuperViewConnection *TrailingSuperViewConnection `@@?`
}

type View struct {
	Connection *Connection   `@@?`
	Name       string        `"[" @Ident`
	Predicate  PredicateList `(Space? @@)? "]"`
}

func (v *View) Predicates() []*Predicate {
	if v.Predicate.Predicate != nil {
		return []*Predicate{v.Predicate.Predicate}
	} else if v.Predicate.Predicates != nil {
		return v.Predicate.Predicates
	} else {
		return []*Predicate{}
	}
}

type Relation struct {
	Gte bool ` @">="`
	Lte bool `| @"<="`
	Eq  bool `| @"=="`
}

type PredicateObject struct {
	Number   int    `@Number`
	ViewName string `| @Ident`
}

type Predicate struct {
	Relation *Relation        `@@?`
	Object   *PredicateObject `@@`
	Priority *Priority        `@@?`
}

type Priority struct {
	Value int `At @Number`
}

type PredicateList struct {
	Predicates []*Predicate `OpenParen @@ ("," @@)* CloseParen`
	Predicate  *Predicate   `| @@`
}

type Orientation struct {
	Direction *string `(@"H"? @"V"?)! Colon`
}

type LeadingSuperViewConnection struct {
	SuperView  bool        `@Pipe`
	Connection *Connection `@@?`
}

type TrailingSuperViewConnection struct {
	Connection *Connection `@@?`
	SuperView  bool        `@Pipe`
}

type Connection struct {
	Predicates *PredicateList `Dash (@@ Dash)?`
}
