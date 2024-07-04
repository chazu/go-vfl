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
	})
)

type View struct {
	Name      string        `"[" @Ident`
	Predicate PredicateList `(Space @@)? "]"`
}

type Program struct {
	Views []*View `@@*`
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
	Priority *int             `("@" @Number)?`
}

type PredicateList struct {
	Predicates []*Predicate `"(" @@ ("(" "," @@ ")")?  ")"`
	Predicate  *Predicate   `| @@`
}

func main() {

	p := participle.MustBuild[Predicate](
		participle.Lexer(l),
	)
	cases := []string{
		"40",
		"Foo",
		">=Foo",
		">=40",
		"40@100",
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

	p2 := participle.MustBuild[Program](
		participle.Lexer(l),
	)
	cases = []string{
		"[Test]",
		"[Test 40]",
		"[Test1][Test2]",
		"[Test1 >=40]",
		"[Test1 >=40@10]",
		"[Test1 >=40][Test2 >=Foo]",
		"[Test1 >=40][Test2 >=Foo@10]",
	}
	for _, c := range cases {
		fmt.Printf("%s...", c)
		res, err := p2.ParseString("", c)
		if err != nil {
			panic(err)
		}
		fmt.Printf("OK\n")
		spew.Dump(res)
	}

}
