package main

import (
	"fmt"
	"log"

	"github.com/chazu/go-vfl/internal/parser"
	"github.com/davecgh/go-spew/spew"
)

func main() {

	// cases := []string{
	// 	"[Test1][Test2]",
	// 	"[Test1 >=40]",
	// 	"[Test1 >=40@10]",
	// 	"[Test1 >=40][Test2 >=Foo]",
	// 	"[Test1 >=40][Test2 >=Foo@10]",
	// 	"[Test1(>=40,<=80)]",
	// 	"[Test1 (>=40)]",
	// 	"V:[TestView]-[TestTwo]",
	// 	"V:[TestView]-50-[TestTwo]",
	// 	"V:|[TestView]-50-[TestTwo]|",
	// 	"|[Test]|",
	// 	"|[Test][TestTwo]|",
	// 	"|[Test]-[TestTwo]|",
	// 	"|[Test]-(50)-[TestTwo]|",
	// 	"|-[Test]-|",
	// 	"|-50-[Test]-|",
	// 	"|-50-[Test]-50-|",
	// 	"|-(>=50@10)-[Test]-(<=50@10)-|",
	// }

	prg := "[Test][Test2]"
	p := parser.New(parser.WithLookahead(250))

	res, err := p.ParseProgram(prg)

	if err != nil {
		log.Fatal(err)
	}

	f, err := res.Reify()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("OK\n")
	spew.Dump(f.Views)
}
