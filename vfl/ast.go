package vfl

// ASTNode is the base interface for all AST elements
type ASTNode interface {
	GetLocation() Location
	Accept(visitor Visitor) error
}

// Location tracks position in source string
type Location struct {
	Line   int // Line number (1-based)
	Column int // Column number (1-based)
	Offset int // Byte offset from start
}

// Program represents a parsed VFL program
type Program struct {
	Input       string
	Orientation Orientation
	Statements  []Statement
	Location    Location
}

// GetLocation returns the location of the program
func (p *Program) GetLocation() Location {
	return p.Location
}

// Accept accepts a visitor
func (p *Program) Accept(v Visitor) error {
	return v.VisitProgram(p)
}

// Statement represents a single constraint statement within a program
type Statement struct {
	Views          []View
	Connections    []Connection
	SuperviewStart bool
	SuperviewEnd   bool
	Location       Location
}

// GetLocation returns the location of the statement
func (s *Statement) GetLocation() Location {
	return s.Location
}

// Accept accepts a visitor
func (s *Statement) Accept(v Visitor) error {
	return v.VisitStatement(s)
}

// View represents a UI view with constraints
type View struct {
	Name        string
	Predicates  []Predicate
	Priority    int
	IsNamedView bool
	Location    Location
	// EVFL extensions
	IsStack          bool
	StackViews       []View
	StackConnections []Connection
	StackOrientation Orientation
	StackLayout      StackLayout
}

// GetLocation returns the location of the view
func (v *View) GetLocation() Location {
	return v.Location
}

// Accept accepts a visitor
func (v *View) Accept(visitor Visitor) error {
	return visitor.VisitView(v)
}

// Predicate represents a constraint expression for a view
type Predicate struct {
	Relation  Relation
	Value     Value
	Attribute string
	Priority  float64  // Priority of this constraint (0-1000, default 1000)
	Location  Location
}

// GetLocation returns the location of the predicate
func (p *Predicate) GetLocation() Location {
	return p.Location
}

// Accept accepts a visitor
func (p *Predicate) Accept(v Visitor) error {
	return v.VisitPredicate(p)
}

// Connection represents spacing between views
type Connection struct {
	Spacing         float64
	IsDefault       bool
	IsEqualSpace    bool
	IsDisconnection bool
	MetricName      string // Name of metric used for spacing
	Location        Location
}

// GetLocation returns the location of the connection
func (c *Connection) GetLocation() Location {
	return c.Location
}

// Accept accepts a visitor
func (c *Connection) Accept(v Visitor) error {
	return v.VisitConnection(c)
}

// Value represents different types of constraint values
type Value struct {
	Type       ValueType
	Constant   float64
	ViewRef    string
	MetricName string
	Expression *Expression
	Percentage float64
}

// Expression represents mathematical expressions in EVFL
type Expression struct {
	Operator Operator
	Left     Value
	Right    Value
	Location Location
}

// GetLocation returns the location of the expression
func (e *Expression) GetLocation() Location {
	return e.Location
}

// Accept accepts a visitor
func (e *Expression) Accept(v Visitor) error {
	return v.VisitExpression(e)
}