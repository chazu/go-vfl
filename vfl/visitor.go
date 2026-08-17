package vfl

// Visitor interface for AST traversal
type Visitor interface {
	VisitProgram(*Program) error
	VisitStatement(*Statement) error
	VisitView(*View) error
	VisitPredicate(*Predicate) error
	VisitConnection(*Connection) error
	VisitExpression(*Expression) error
}

// BaseVisitor provides a default implementation of Visitor
type BaseVisitor struct{}

// VisitProgram visits a program node
func (v *BaseVisitor) VisitProgram(p *Program) error {
	for i := range p.Statements {
		if err := p.Statements[i].Accept(v); err != nil {
			return err
		}
	}
	return nil
}

// VisitStatement visits a statement node
func (v *BaseVisitor) VisitStatement(s *Statement) error {
	for i := range s.Views {
		if err := s.Views[i].Accept(v); err != nil {
			return err
		}
	}
	for i := range s.Connections {
		if err := s.Connections[i].Accept(v); err != nil {
			return err
		}
	}
	return nil
}

// VisitView visits a view node
func (v *BaseVisitor) VisitView(view *View) error {
	for i := range view.Predicates {
		if err := view.Predicates[i].Accept(v); err != nil {
			return err
		}
	}
	return nil
}

// VisitPredicate visits a predicate node
func (v *BaseVisitor) VisitPredicate(p *Predicate) error {
	if p.Value.Expression != nil {
		return p.Value.Expression.Accept(v)
	}
	return nil
}

// VisitConnection visits a connection node
func (v *BaseVisitor) VisitConnection(c *Connection) error {
	return nil
}

// VisitExpression visits an expression node
func (v *BaseVisitor) VisitExpression(e *Expression) error {
	if e.Left.Expression != nil {
		if err := e.Left.Expression.Accept(v); err != nil {
			return err
		}
	}
	if e.Right.Expression != nil {
		if err := e.Right.Expression.Accept(v); err != nil {
			return err
		}
	}
	return nil
}