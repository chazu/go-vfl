package main

import (
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
	})
)

type View struct {
	Name string `"[" @Ident "]"`
}

type Program struct {
	Views []*View `@@*`
}

func main() {

	p := participle.MustBuild[Program](
		participle.Lexer(l),
	)
	res, err := p.ParseString("", "[Test1][Test2]")
	if err != nil {
		panic(err)
	}

	spew.Dump(res)
}
