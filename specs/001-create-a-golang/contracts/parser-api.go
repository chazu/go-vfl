// Package contracts defines the API contracts for the VFL parser library
package contracts

import (
	"io"
)

// Parser is the main interface for parsing VFL and EVFL strings
type Parser interface {
	// Parse parses a VFL string and returns an AST or error with partial AST
	Parse(input string) (*Program, error)

	// ParseWithOptions parses with custom options
	ParseWithOptions(input string, options ParserOptions) (*Program, error)

	// Validate checks if a VFL string is syntactically valid
	Validate(input string) []ValidationError
}

// ParserOptions configures parser behavior
type ParserOptions struct {
	// Lookahead sets the parser lookahead (default: 5)
	Lookahead int

	// StrictMode disables error recovery for partial AST
	StrictMode bool

	// EnableEVFL enables Extended VFL features
	EnableEVFL bool

	// Metrics provides named spacing values
	Metrics map[string]float64
}

// Program represents a parsed VFL program
type Program struct {
	// Input is the original VFL string
	Input string

	// Orientation is the layout direction
	Orientation Orientation

	// Statements are the constraint statements
	Statements []Statement

	// Location tracks position in source
	Location Location
}

// Statement represents a single constraint line
type Statement struct {
	// Views in this statement
	Views []View

	// Connections between views
	Connections []Connection

	// SuperviewStart indicates starting at superview edge
	SuperviewStart bool

	// SuperviewEnd indicates ending at superview edge
	SuperviewEnd bool

	// Location in source
	Location Location
}

// View represents a UI view with constraints
type View struct {
	// Name is the view identifier
	Name string

	// Predicates are the size/position constraints
	Predicates []Predicate

	// Priority is the constraint priority (0-1000)
	Priority int

	// IsNamedView indicates this references an external view
	IsNamedView bool

	// Location in source
	Location Location
}

// Predicate represents a constraint expression
type Predicate struct {
	// Relation is the constraint relation
	Relation Relation

	// Value is the constraint value
	Value Value

	// Attribute for EVFL attribute references
	Attribute string

	// Location in source
	Location Location
}

// Connection represents spacing between views
type Connection struct {
	// Spacing amount (-1 for default)
	Spacing float64

	// IsDefault uses default spacing
	IsDefault bool

	// IsEqualSpace for EVFL equal spacers (~)
	IsEqualSpace bool

	// IsDisconnection for EVFL disconnections (→)
	IsDisconnection bool

	// Location in source
	Location Location
}

// Value represents different types of constraint values
type Value struct {
	// Type indicates which field is set
	Type ValueType

	// Constant for numeric values
	Constant float64

	// ViewRef for view references
	ViewRef string

	// MetricName for named metrics
	MetricName string

	// Expression for EVFL math expressions
	Expression *Expression

	// Percentage for EVFL percentages (0-100)
	Percentage float64
}

// Expression represents mathematical expressions in EVFL
type Expression struct {
	// Operator is the math operation
	Operator Operator

	// Left operand
	Left Value

	// Right operand
	Right Value

	// Location in source
	Location Location
}

// ParseError provides detailed error information
type ParseError struct {
	// Error implements the error interface
	error

	// Message describes the error
	Message string

	// Location where error occurred
	Location Location

	// Type categorizes the error
	Type ErrorType

	// PartialAST contains any successfully parsed portion
	PartialAST *Program
}

// ValidationError represents a validation issue
type ValidationError struct {
	// Message describes the issue
	Message string

	// Location of the issue
	Location Location

	// Severity of the issue
	Severity Severity
}

// Location tracks position in source string
type Location struct {
	// Line number (1-based)
	Line int

	// Column number (1-based)
	Column int

	// Offset in bytes from start
	Offset int
}

// Enums

// Orientation represents layout direction
type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
	Z
	Combined // HV: for EVFL
)

// Relation represents constraint relations
type Relation int

const (
	Equal Relation = iota
	GreaterEqual
	LessEqual
)

// ValueType indicates which Value field is set
type ValueType int

const (
	ConstantValue ValueType = iota
	ViewRefValue
	MetricValue
	ExpressionValue
	PercentageValue
)

// Operator represents math operations in EVFL
type Operator int

const (
	Add Operator = iota
	Subtract
	Multiply
	Divide
)

// ErrorType categorizes parse errors
type ErrorType int

const (
	SyntaxError ErrorType = iota
	SemanticError
	ReferenceError
)

// Severity levels for validation
type Severity int

const (
	Error Severity = iota
	Warning
	Info
)

// Visitor interface for AST traversal
type Visitor interface {
	VisitProgram(*Program) error
	VisitStatement(*Statement) error
	VisitView(*View) error
	VisitPredicate(*Predicate) error
	VisitConnection(*Connection) error
	VisitExpression(*Expression) error
}

// Node interface for all AST nodes
type Node interface {
	Accept(Visitor) error
	GetLocation() Location
}

// Writer outputs parsed AST in various formats
type Writer interface {
	// WriteJSON outputs AST as JSON
	WriteJSON(io.Writer, *Program) error

	// WriteDebug outputs human-readable AST
	WriteDebug(io.Writer, *Program) error

	// WriteConstraints outputs as constraint equations
	WriteConstraints(io.Writer, *Program) error
}