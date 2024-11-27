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

func (p *Parser) ParseProgram(pgm string) (*program, error) {
	psr := participle.MustBuild[program](
		participle.Lexer(l),
		participle.UseLookahead(p.Lookahead),
	)
	res, err := psr.ParseString("", pgm)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func unifyPredicateList(pl predicateList) []predicate {
	res := []predicate{}

	if pl.Predicate != (predicate{}) {
		res = append(pl.Predicates, pl.Predicate)
	} else {
		res = pl.Predicates
	}

	return res
}

func viewToViewAST(v view) ViewAST {
	res := ViewAST{}

	c := ConnectionAST{
		Predicates:  unifyPredicateList(v.Connection.Predicates),
		IsSuperview: false,
	}

	res.Name = v.Name
	res.LeadingConnection = c
	res.Predicates = unifyPredicateList(v.Predicate)

	return res
}

func (p *program) Reify() (ProgramAST, error) {
	res := ProgramAST{}
	res.Orientation = p.Orientation

	lsvc := ConnectionAST{
		Predicates:  unifyPredicateList(p.LeadingSuperviewConnection.Connection.Predicates),
		IsSuperview: true,
	}
	res.LeadingSuperviewConnection = lsvc

	tsvc := ConnectionAST{
		Predicates:  unifyPredicateList(p.TrailingSuperviewConnection.Connection.Predicates),
		IsSuperview: true,
	}
	res.TrailingSuperviewConnection = tsvc

	for i, v := range p.Views {
		// TODO make connections bi-directional
		vast := viewToViewAST(v)
		if i != 0 {
			res.Views[i-1].TrailingConnection = vast.LeadingConnection
		}
		res.Views = append(res.Views, viewToViewAST(v))
	}

	return res, nil
}
