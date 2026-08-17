package vfl

import "fmt"

// ParseError provides detailed error information with partial AST support
type ParseError struct {
	Message    string
	Location   Location
	Type       ErrorType
	Input      string   // The input that failed to parse
	Partial    *Program // Partial AST if available (alias for PartialAST)
	PartialAST *Program
	cause      error
}

// Error implements the error interface
func (e *ParseError) Error() string {
	if e.Location.Line > 0 {
		return fmt.Sprintf("%s at line %d, column %d: %s",
			e.Type.String(), e.Location.Line, e.Location.Column, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Type.String(), e.Message)
}

// Unwrap returns the underlying error
func (e *ParseError) Unwrap() error {
	return e.cause
}

// NewParseError creates a new parse error
func NewParseError(message string, location Location, errorType ErrorType) *ParseError {
	return &ParseError{
		Message:  message,
		Location: location,
		Type:     errorType,
	}
}

// WithPartialAST adds a partial AST to the error
func (e *ParseError) WithPartialAST(ast *Program) *ParseError {
	e.PartialAST = ast
	e.Partial = ast
	return e
}

// WithCause adds an underlying cause to the error
func (e *ParseError) WithCause(err error) *ParseError {
	e.cause = err
	return e
}

// ValidationError represents a validation issue
type ValidationError struct {
	Message  string
	Location Location
	Severity Severity
}

// Error implements the error interface
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Severity.String(), e.Message)
}

// CircularReferenceError indicates a cycle in named views
type CircularReferenceError struct {
	Cycle   []string
	Message string
}

// Error implements the error interface
func (e *CircularReferenceError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("circular reference detected: %v", e.Cycle)
}

// CompositionError represents errors during composition
type CompositionError struct {
	ViewName string
	Message  string
	Type     CompositionErrorType
}

// Error implements the error interface
func (e *CompositionError) Error() string {
	return fmt.Sprintf("%s: %s: %s", e.Type.String(), e.ViewName, e.Message)
}

// ErrorType categorizes parse errors
type ErrorType int

const (
	SyntaxError ErrorType = iota
	SemanticError
	ReferenceError
)

// String returns the string representation of the error type
func (e ErrorType) String() string {
	switch e {
	case SyntaxError:
		return "syntax error"
	case SemanticError:
		return "semantic error"
	case ReferenceError:
		return "reference error"
	default:
		return "unknown error"
	}
}

// Severity levels for validation
type Severity int

const (
	Error Severity = iota
	Warning
	Info
)

// String returns the string representation of the severity
func (s Severity) String() string {
	switch s {
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Info:
		return "info"
	default:
		return "unknown"
	}
}

// CompositionErrorType categorizes composition errors
type CompositionErrorType int

const (
	MissingNamedView CompositionErrorType = iota
	CircularReference
	ConflictingConstraints
	InvalidHierarchy
)

// String returns the string representation of the composition error type
func (c CompositionErrorType) String() string {
	switch c {
	case MissingNamedView:
		return "missing named view"
	case CircularReference:
		return "circular reference"
	case ConflictingConstraints:
		return "conflicting constraints"
	case InvalidHierarchy:
		return "invalid hierarchy"
	default:
		return "unknown composition error"
	}
}