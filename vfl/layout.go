package vfl

// NamedView represents a reusable view definition
type NamedView struct {
	Name         string
	Program      *Program
	IsComposite  bool
	Dependencies []string
}

// LayoutResult is the final composed AST ready for rendering
type LayoutResult struct {
	Constraints   []Constraint
	ViewHierarchy *ViewTree
	Errors        []ParseError
	Metrics       map[string]float64
}

// Constraint represents a resolved layout constraint
type Constraint struct {
	FirstItem       ConstraintItem
	FirstAttribute  Attribute
	Relation        Relation
	SecondItem      *ConstraintItem
	SecondAttribute Attribute
	Constant        float64
	Multiplier      float64
	Priority        int
}

// ConstraintItem represents a view in a constraint
type ConstraintItem struct {
	Name        string
	IsSuperview bool
}

// ViewTree represents the view hierarchy
type ViewTree struct {
	Root  *ViewNode
	Views map[string]*ViewNode
}

// ViewNode represents a view in the hierarchy
type ViewNode struct {
	Name        string
	Children    []*ViewNode
	Parent      *ViewNode
	ZIndex      int
	Constraints []Constraint
	IsSuperview bool
}

// ConflictingConstraint represents conflicting constraints
type ConflictingConstraint struct {
	First  Constraint
	Second Constraint
	Reason string
}