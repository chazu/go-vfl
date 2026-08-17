// Package contracts defines Extended VFL specific API contracts
package contracts

import "io"

// EVFLParser extends standard VFL parsing with EVFL features
type EVFLParser interface {
	Parser // Embeds standard parser interface

	// ParseStack parses view stack notation
	ParseStack(input string) (*ViewStack, error)

	// ParseRange parses view range notation (view1..5)
	ParseRange(input string) (*ViewRange, error)

	// ParsePercentage parses percentage-based constraints
	ParsePercentage(input string) (*PercentageConstraint, error)

	// EnableFeature enables specific EVFL features
	EnableFeature(feature EVFLFeature) error
}

// ViewStack represents grouped view layouts in EVFL
type ViewStack struct {
	// Name is the stack identifier
	Name string

	// Orientation of the stack
	Orientation Orientation

	// Views in the stack
	Views []View

	// Layout defines internal arrangement
	Layout StackLayout

	// Location in source
	Location Location
}

// ViewRange represents batch view operations
type ViewRange struct {
	// Prefix is the base view name
	Prefix string

	// Start index
	Start int

	// End index (inclusive)
	End int

	// Template for generated views
	Template View

	// Location in source
	Location Location
}

// PercentageConstraint represents percentage-based sizing
type PercentageConstraint struct {
	// View being constrained
	View string

	// Percentage value (0-100)
	Percentage float64

	// RelativeTo specifies the reference view/dimension
	RelativeTo string

	// Attribute being constrained
	Attribute Attribute

	// Location in source
	Location Location
}

// StackLayout defines how views are arranged in a stack
type StackLayout int

const (
	LinearStack StackLayout = iota
	GridStack
	FlowStack
)

// EVFLFeature represents optional EVFL features
type EVFLFeature int

const (
	FeaturePercentages EVFLFeature = iota
	FeatureExpressions
	FeatureStacks
	FeatureRanges
	FeatureZOrdering
	FeatureEqualSpacers
	FeatureDisconnections
	FeatureAttributeRefs
	FeatureMultipleViews
	FeatureCombinedOrientation
)

// AttributeReference represents view.attribute notation in EVFL
type AttributeReference struct {
	// ViewName is the referenced view
	ViewName string

	// Attribute being referenced
	Attribute Attribute

	// Multiplier for the reference (default: 1.0)
	Multiplier float64

	// Offset to add to the reference
	Offset float64

	// Location in source
	Location Location
}

// EqualSpacer represents distributed spacing in EVFL
type EqualSpacer struct {
	// Views to space equally
	Views []string

	// MinSpacing if specified
	MinSpacing *float64

	// MaxSpacing if specified
	MaxSpacing *float64

	// Location in source
	Location Location
}

// ZOrderConstraint represents Z-axis layering in EVFL
type ZOrderConstraint struct {
	// Views in Z-order (front to back)
	Views []string

	// Relative indicates relative vs absolute ordering
	Relative bool

	// Location in source
	Location Location
}

// MultipleViewConstraint applies same constraint to multiple views
type MultipleViewConstraint struct {
	// Views affected by this constraint
	Views []string

	// Constraint to apply to all views
	Constraint Predicate

	// Location in source
	Location Location
}

// Disconnection represents alignment without connection in EVFL
type Disconnection struct {
	// FirstView aligned view
	FirstView string

	// SecondView target view
	SecondView string

	// Alignment type
	Alignment AlignmentType

	// Location in source
	Location Location
}

// AlignmentType for disconnections
type AlignmentType int

const (
	AlignTop AlignmentType = iota
	AlignBottom
	AlignLeft
	AlignRight
	AlignCenterX
	AlignCenterY
	AlignBaseline
)

// EVFLValidator extends validation for EVFL features
type EVFLValidator interface {
	Validator // Embeds standard validator

	// ValidatePercentages checks percentage constraint validity
	ValidatePercentages(constraints []PercentageConstraint) []ValidationError

	// ValidateExpressions checks mathematical expressions
	ValidateExpressions(expressions []*Expression) []ValidationError

	// ValidateStacks checks view stack validity
	ValidateStacks(stacks []ViewStack) []ValidationError

	// ValidateZOrder checks Z-ordering conflicts
	ValidateZOrder(constraints []ZOrderConstraint) []ValidationError
}

// EVFLComposer extends composition for EVFL features
type EVFLComposer interface {
	Composer // Embeds standard composer

	// ComposeStacks processes view stacks
	ComposeStacks(stacks ...ViewStack) (*LayoutResult, error)

	// ExpandRanges expands view ranges into individual views
	ExpandRanges(ranges ...ViewRange) ([]View, error)

	// ResolvePercentages converts percentages to constraints
	ResolvePercentages(constraints []PercentageConstraint, containerSize Size) []Constraint

	// ProcessDisconnections handles alignment without connections
	ProcessDisconnections(disconnections []Disconnection) []Constraint
}

// Size represents dimensions for percentage calculations
type Size struct {
	Width  float64
	Height float64
}

// EVFLWriter extends output formats for EVFL
type EVFLWriter interface {
	Writer // Embeds standard writer

	// WriteEVFL outputs in Extended VFL format
	WriteEVFL(w io.Writer, program *Program) error

	// WriteStacks outputs view stacks
	WriteStacks(w io.Writer, stacks []ViewStack) error

	// WritePercentages outputs percentage constraints
	WritePercentages(w io.Writer, constraints []PercentageConstraint) error
}