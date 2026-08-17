package vfl

import (
	"fmt"
)

// Composer combines multiple VFL programs with named views
type Composer interface {
	RegisterNamedView(name string, program *Program) error
	Compose(programs ...*Program) (*LayoutResult, error)
	ComposeWithValidation(programs ...*Program) (*LayoutResult, []ValidationError)
	Reset()
	HasNamedView(name string) bool
	GetNamedView(name string) (*NamedView, bool)
}

// composer implements the Composer interface
type composer struct {
	namedViews map[string]*NamedView
	visited    map[string]bool
}

// NewComposer creates a new composer
func NewComposer() Composer {
	return &composer{
		namedViews: make(map[string]*NamedView),
		visited:    make(map[string]bool),
	}
}

// RegisterNamedView adds a named view to the registry
func (c *composer) RegisterNamedView(name string, program *Program) error {
	if _, exists := c.namedViews[name]; exists {
		return fmt.Errorf("named view '%s' already registered", name)
	}

	namedView := &NamedView{
		Name:         name,
		Program:      program,
		IsComposite:  c.hasNamedViewReferences(program),
		Dependencies: c.extractDependencies(program),
	}

	c.namedViews[name] = namedView
	return nil
}

// Compose combines programs resolving named view references
func (c *composer) Compose(programs ...*Program) (*LayoutResult, error) {
	// Check for circular references before composing
	if err := c.detectCircularReferences(); err != nil {
		return nil, err
	}

	c.visited = make(map[string]bool) // Reset visited tracking

	result := &LayoutResult{
		Constraints:   []Constraint{},
		ViewHierarchy: &ViewTree{
			Views: make(map[string]*ViewNode),
		},
		Metrics: make(map[string]float64),
	}

	// Process each program
	for _, prog := range programs {
		if err := c.composeProgram(prog, result, []string{}); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// ComposeWithValidation performs composition with full validation
func (c *composer) ComposeWithValidation(programs ...*Program) (*LayoutResult, []ValidationError) {
	result, err := c.Compose(programs...)

	var validationErrors []ValidationError

	if err != nil {
		validationErrors = append(validationErrors, ValidationError{
			Message:  err.Error(),
			Severity: Error,
		})
		return result, validationErrors
	}

	// Validate the result
	validator := NewValidator()
	validationErrors = validator.ValidateLayout(result)

	return result, validationErrors
}

// Reset clears the named view registry
func (c *composer) Reset() {
	c.namedViews = make(map[string]*NamedView)
	c.visited = make(map[string]bool)
}

// HasNamedView checks if a named view exists
func (c *composer) HasNamedView(name string) bool {
	_, exists := c.namedViews[name]
	return exists
}

// GetNamedView retrieves a registered named view
func (c *composer) GetNamedView(name string) (*NamedView, bool) {
	view, exists := c.namedViews[name]
	return view, exists
}

// Private helper methods

func (c *composer) composeProgram(prog *Program, result *LayoutResult, path []string) error {
	// Process each statement
	for _, stmt := range prog.Statements {
		if err := c.composeStatement(stmt, result, path); err != nil {
			return err
		}
	}
	return nil
}

func (c *composer) composeStatement(stmt Statement, result *LayoutResult, path []string) error {
	// Create superview node if needed
	if stmt.SuperviewStart || stmt.SuperviewEnd {
		if result.ViewHierarchy.Root == nil {
			result.ViewHierarchy.Root = &ViewNode{
				Name:        "superview",
				IsSuperview: true,
				Children:    []*ViewNode{},
			}
		}
	}

	// Process views
	for i, view := range stmt.Views {
		// Check if this is a named view reference
		if c.HasNamedView(view.Name) {
			namedView, exists := c.namedViews[view.Name]
			if !exists {
				return &CompositionError{
					ViewName: view.Name,
					Message:  "named view not found",
					Type:     MissingNamedView,
				}
			}

			// Check for circular reference during composition
			for _, p := range path {
				if p == view.Name {
					return &CircularReferenceError{
						Cycle:   append(path, view.Name),
						Message: fmt.Sprintf("circular reference detected: %v", append(path, view.Name)),
					}
				}
			}

			// Compose the named view's program
			newPath := append(path, view.Name)
			if err := c.composeProgram(namedView.Program, result, newPath); err != nil {
				return err
			}
		} else {
			// Regular view - add to hierarchy
			viewNode := &ViewNode{
				Name:        view.Name,
				Constraints: []Constraint{},
			}

			// Add to view tree
			result.ViewHierarchy.Views[view.Name] = viewNode

			// If connected to superview, make it a child
			if stmt.SuperviewStart && i == 0 && result.ViewHierarchy.Root != nil {
				viewNode.Parent = result.ViewHierarchy.Root
				result.ViewHierarchy.Root.Children = append(result.ViewHierarchy.Root.Children, viewNode)
			}

			// Generate constraints from predicates
			for _, pred := range view.Predicates {
				constraint := c.predicateToConstraint(view.Name, pred)
				result.Constraints = append(result.Constraints, constraint)
				viewNode.Constraints = append(viewNode.Constraints, constraint)
			}
		}

		// Process connections
		if i < len(stmt.Connections) {
			conn := stmt.Connections[i]
			constraint := c.connectionToConstraint(stmt.Views[i].Name,
				c.getNextViewName(stmt.Views, i), conn)
			result.Constraints = append(result.Constraints, constraint)
		}
	}

	return nil
}

func (c *composer) predicateToConstraint(viewName string, pred Predicate) Constraint {
	constraint := Constraint{
		FirstItem: ConstraintItem{
			Name: viewName,
		},
		FirstAttribute: Width, // Default, should be determined by context
		Relation:       pred.Relation,
		Multiplier:     1.0,
		Priority:       250, // Default priority
	}

	// Handle different value types
	switch pred.Value.Type {
	case ConstantValue:
		constraint.Constant = pred.Value.Constant
	case ViewRefValue:
		constraint.SecondItem = &ConstraintItem{
			Name: pred.Value.ViewRef,
		}
		constraint.SecondAttribute = Width // Default
	case PercentageValue:
		// Convert percentage to actual value (needs container size)
		constraint.Constant = pred.Value.Percentage
		constraint.Multiplier = 0.01 // Percentage multiplier
	}

	// Set attribute if specified
	if pred.Attribute != "" {
		constraint.FirstAttribute = c.parseAttribute(pred.Attribute)
	}

	return constraint
}

func (c *composer) connectionToConstraint(fromView, toView string, conn Connection) Constraint {
	return Constraint{
		FirstItem: ConstraintItem{
			Name: fromView,
		},
		FirstAttribute: Trailing,
		Relation:       Equal,
		SecondItem: &ConstraintItem{
			Name: toView,
		},
		SecondAttribute: Leading,
		Constant:        conn.Spacing,
		Multiplier:      1.0,
		Priority:        1000, // Required
	}
}

func (c *composer) hasNamedViewReferences(prog *Program) bool {
	for _, stmt := range prog.Statements {
		for _, view := range stmt.Views {
			if c.HasNamedView(view.Name) {
				return true
			}
		}
	}
	return false
}

func (c *composer) extractDependencies(prog *Program) []string {
	var deps []string
	for _, stmt := range prog.Statements {
		for _, view := range stmt.Views {
			if c.HasNamedView(view.Name) {
				deps = append(deps, view.Name)
			}
		}
	}
	return deps
}

func (c *composer) detectCircularReferences() error {
	for name := range c.namedViews {
		c.visited = make(map[string]bool)
		path := []string{}
		if cycle := c.detectCycle(name, path); cycle != nil {
			return &CircularReferenceError{
				Cycle:   cycle,
				Message: fmt.Sprintf("circular reference detected: %v", cycle),
			}
		}
	}
	return nil
}

func (c *composer) detectCycle(name string, path []string) []string {
	if c.visited[name] {
		// Found cycle - return the cycle path
		for i, n := range path {
			if n == name {
				return append(path[i:], name)
			}
		}
	}

	c.visited[name] = true
	path = append(path, name)

	if namedView, exists := c.namedViews[name]; exists {
		for _, dep := range namedView.Dependencies {
			if cycle := c.detectCycle(dep, path); cycle != nil {
				return cycle
			}
		}
	}

	return nil
}

func (c *composer) getNextViewName(views []View, currentIndex int) string {
	if currentIndex+1 < len(views) {
		return views[currentIndex+1].Name
	}
	return "superview"
}

func (c *composer) parseAttribute(attr string) Attribute {
	switch attr {
	case "width":
		return Width
	case "height":
		return Height
	case "left":
		return Left
	case "right":
		return Right
	case "top":
		return Top
	case "bottom":
		return Bottom
	case "centerX":
		return CenterX
	case "centerY":
		return CenterY
	case "leading":
		return Leading
	case "trailing":
		return Trailing
	case "baseline":
		return Baseline
	default:
		return NotAnAttribute
	}
}