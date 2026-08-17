// Package contracts defines the API contracts for VFL program composition
package contracts

// Composer combines multiple VFL programs with named views
type Composer interface {
	// RegisterNamedView adds a named view to the registry
	RegisterNamedView(name string, program *Program) error

	// Compose combines programs resolving named view references
	Compose(programs ...*Program) (*LayoutResult, error)

	// ComposeWithValidation performs composition with full validation
	ComposeWithValidation(programs ...*Program) (*LayoutResult, []ValidationError)

	// Reset clears the named view registry
	Reset()

	// HasNamedView checks if a named view exists
	HasNamedView(name string) bool

	// GetNamedView retrieves a registered named view
	GetNamedView(name string) (*NamedView, bool)
}

// NamedView represents a reusable view definition
type NamedView struct {
	// Name is the global identifier
	Name string

	// Program is the view definition
	Program *Program

	// IsComposite indicates if it contains other named views
	IsComposite bool

	// Dependencies lists other named views this one requires
	Dependencies []string
}

// LayoutResult is the final composed AST ready for rendering
type LayoutResult struct {
	// Constraints are all resolved constraints
	Constraints []Constraint

	// ViewHierarchy represents view relationships
	ViewHierarchy *ViewTree

	// Errors contains any partial errors encountered
	Errors []ParseError

	// Metrics used during composition
	Metrics map[string]float64
}

// Constraint represents a resolved layout constraint
type Constraint struct {
	// FirstItem is the first view or superview
	FirstItem ConstraintItem

	// FirstAttribute is the attribute being constrained
	FirstAttribute Attribute

	// Relation is the constraint relation
	Relation Relation

	// SecondItem is the second view (optional)
	SecondItem *ConstraintItem

	// SecondAttribute is the second attribute (optional)
	SecondAttribute Attribute

	// Constant is the constant value
	Constant float64

	// Multiplier for the constraint equation
	Multiplier float64

	// Priority is the constraint priority
	Priority int
}

// ConstraintItem represents a view in a constraint
type ConstraintItem struct {
	// Name is the view identifier
	Name string

	// IsSuperview indicates this is the parent view
	IsSuperview bool
}

// ViewTree represents the view hierarchy
type ViewTree struct {
	// Root is the top-level view
	Root *ViewNode

	// Views maps names to nodes for quick lookup
	Views map[string]*ViewNode
}

// ViewNode represents a view in the hierarchy
type ViewNode struct {
	// Name is the view identifier
	Name string

	// Children are nested views
	Children []*ViewNode

	// Parent is the containing view
	Parent *ViewNode

	// ZIndex for Z-ordering (EVFL)
	ZIndex int

	// Constraints affecting this view
	Constraints []Constraint
}

// Attribute represents a view attribute
type Attribute int

const (
	Left Attribute = iota
	Right
	Top
	Bottom
	Leading
	Trailing
	Width
	Height
	CenterX
	CenterY
	Baseline
	NotAnAttribute
)

// CircularReferenceError indicates a cycle in named views
type CircularReferenceError struct {
	// Error implements the error interface
	error

	// Cycle lists the views forming the cycle
	Cycle []string

	// Message describes the error
	Message string
}

// CompositionError represents errors during composition
type CompositionError struct {
	// Error implements the error interface
	error

	// ViewName that caused the error
	ViewName string

	// Message describes the error
	Message string

	// Type of composition error
	Type CompositionErrorType
}

// CompositionErrorType categorizes composition errors
type CompositionErrorType int

const (
	MissingNamedView CompositionErrorType = iota
	CircularReference
	ConflictingConstraints
	InvalidHierarchy
)

// Validator checks composed layouts for issues
type Validator interface {
	// ValidateLayout checks a layout for constraint conflicts
	ValidateLayout(layout *LayoutResult) []ValidationError

	// ValidateHierarchy checks view hierarchy validity
	ValidateHierarchy(tree *ViewTree) []ValidationError

	// CheckConflicts detects conflicting constraints
	CheckConflicts(constraints []Constraint) []ConflictingConstraint
}

// ConflictingConstraint represents conflicting constraints
type ConflictingConstraint struct {
	// First constraint
	First Constraint

	// Second conflicting constraint
	Second Constraint

	// Reason explains the conflict
	Reason string
}

// Registry manages named views
type Registry interface {
	// Register adds a named view
	Register(name string, view *NamedView) error

	// Unregister removes a named view
	Unregister(name string) bool

	// Get retrieves a named view
	Get(name string) (*NamedView, bool)

	// List returns all registered view names
	List() []string

	// Clear removes all named views
	Clear()

	// DetectCycles checks for circular references
	DetectCycles() ([][]string, error)
}