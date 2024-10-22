package parser

import "github.com/alecthomas/participle/v2"

type Parser struct {
	Lookahead int
}

const ParserDefaultLookahead int = 10

func New(options ...func(*Parser)) *Parser {
	p := &Parser{
		Lookahead: ParserDefaultLookahead,
	}

	for _, o := range options {
		o(p)
	}

	return p
}

func WithLookahead(lookahead int) func(*Parser) {
	return func(p *Parser) {
		p.Lookahead = lookahead
	}
}

func (p *Parser) ParseProgram(program string) (*Program, error) {
	psr := participle.MustBuild[Program](
		participle.Lexer(l),
		participle.UseLookahead(p.Lookahead),
	)
	res, err := psr.ParseString("", program)
	if err != nil {
		return nil, err
	}

	return res, nil
}
