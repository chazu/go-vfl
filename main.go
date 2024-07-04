package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/davecgh/go-spew/spew"
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
	})
)

type Program struct {
	Orientation *Orientation `@@?`
	Views       []*View      `@@ @@?`
}

type View struct {
	Name      string        `"[" @Ident`
	Predicate PredicateList `(Space @@)? "]"`
}

type Relation struct {
	Eq  *bool `@"=="`
	Gte *bool `| @">="`
	Lte *bool `| @"<="`
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
	Value *int `At @Number`
}

type PredicateList struct {
	Predicates []*Predicate `OpenParen @@ ("," @@)* CloseParen`
	Predicate  *Predicate   `| @@`
}

type Orientation struct {
	Direction *string `(@"H"? @"V"?)! Colon`
}

type Connection struct {
	Predicates *PredicateList `"-" (@@ "-")?`
}

func main() {

	p := participle.MustBuild[Program](
		participle.Lexer(l),
	)
	cases := []string{
		"[Test]",
		"[Test 40]",
		"[Test1][Test2]",
		"[Test1 >=40]",
		"[Test1 >=40@10]",
		"[Test1 >=40][Test2 >=Foo]",
		"[Test1 >=40][Test2 >=Foo@10]",
		"[Test1 (>=40,<=80)]",
		"[Test1 (>=40)]",
		"H:[TestView]",
		"V:[TestView]",
	}
	for _, c := range cases {
		fmt.Printf("%s...", c)
		res, err := p.ParseString("", c)
		if err != nil {
			panic(err)
		}
		fmt.Printf("OK\n")
		spew.Dump(res)
	}

}
