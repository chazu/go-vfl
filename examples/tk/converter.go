package main

import (
	"github.com/chazu/go-vfl/vfl"
)

// ConstraintConverter converts VFL constraints to TK layout commands
type ConstraintConverter struct {
	constraints []vfl.Constraint
	views       map[string]*ViewLayout
}

// ViewLayout represents the calculated layout for a view
type ViewLayout struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// NewConstraintConverter creates a new converter
func NewConstraintConverter(constraints []vfl.Constraint) *ConstraintConverter {
	return &ConstraintConverter{
		constraints: constraints,
		views:       make(map[string]*ViewLayout),
	}
}

// Convert converts constraints to view layouts
func (c *ConstraintConverter) Convert() map[string]*ViewLayout {
	// Initialize views with default values
	c.initializeViews()

	// Apply constraints iteratively until stable
	maxIterations := 100
	for i := 0; i < maxIterations; i++ {
		changed := false
		for _, constraint := range c.constraints {
			if c.applyConstraint(constraint) {
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	return c.views
}

// initializeViews creates initial view layouts
func (c *ConstraintConverter) initializeViews() {
	viewNames := make(map[string]bool)

	// Collect all view names from constraints
	for _, constraint := range c.constraints {
		viewNames[constraint.FirstItem.Name] = true
		if constraint.SecondItem != nil {
			viewNames[constraint.SecondItem.Name] = true
		}
	}

	// Initialize each view with default values
	for name := range viewNames {
		if name == "superview" || name == "|" {
			// Skip superview, it has fixed dimensions
			continue
		}
		c.views[name] = &ViewLayout{
			X:      0,
			Y:      0,
			Width:  100,
			Height: 50,
		}
	}

	// Add superview with window dimensions
	c.views["superview"] = &ViewLayout{
		X:      0,
		Y:      0,
		Width:  800,
		Height: 600,
	}
}

// applyConstraint applies a single constraint
func (c *ConstraintConverter) applyConstraint(constraint vfl.Constraint) bool {
	firstView := c.getOrCreateView(constraint.FirstItem.Name)

	if constraint.SecondItem == nil {
		// Absolute constraint
		return c.applyAbsoluteConstraint(firstView, constraint)
	}

	// Relative constraint
	secondView := c.getOrCreateView(constraint.SecondItem.Name)
	return c.applyRelativeConstraint(firstView, secondView, constraint)
}

// getOrCreateView gets or creates a view layout
func (c *ConstraintConverter) getOrCreateView(name string) *ViewLayout {
	if name == "|" {
		name = "superview"
	}

	if view, ok := c.views[name]; ok {
		return view
	}

	// Create new view with defaults
	view := &ViewLayout{
		X:      0,
		Y:      0,
		Width:  100,
		Height: 50,
	}
	c.views[name] = view
	return view
}

// applyAbsoluteConstraint applies a constraint with no second item
func (c *ConstraintConverter) applyAbsoluteConstraint(view *ViewLayout, constraint vfl.Constraint) bool {
	oldValue := c.getAttributeValue(view, constraint.FirstAttribute)
	newValue := constraint.Constant

	if !c.shouldUpdateValue(oldValue, newValue, constraint.Relation) {
		return false
	}

	c.setAttributeValue(view, constraint.FirstAttribute, newValue)
	return true
}

// applyRelativeConstraint applies a constraint between two views
func (c *ConstraintConverter) applyRelativeConstraint(firstView, secondView *ViewLayout, constraint vfl.Constraint) bool {
	secondValue := c.getAttributeValue(secondView, constraint.SecondAttribute)
	targetValue := secondValue*constraint.Multiplier + constraint.Constant

	oldValue := c.getAttributeValue(firstView, constraint.FirstAttribute)

	if !c.shouldUpdateValue(oldValue, targetValue, constraint.Relation) {
		return false
	}

	c.setAttributeValue(firstView, constraint.FirstAttribute, targetValue)
	return true
}

// getAttributeValue gets an attribute value from a view
func (c *ConstraintConverter) getAttributeValue(view *ViewLayout, attr vfl.Attribute) float64 {
	switch attr {
	case vfl.Width:
		return view.Width
	case vfl.Height:
		return view.Height
	case vfl.Left, vfl.Leading:
		return view.X
	case vfl.Right, vfl.Trailing:
		return view.X + view.Width
	case vfl.Top:
		return view.Y
	case vfl.Bottom:
		return view.Y + view.Height
	case vfl.CenterX:
		return view.X + view.Width/2
	case vfl.CenterY:
		return view.Y + view.Height/2
	default:
		return 0
	}
}

// setAttributeValue sets an attribute value on a view
func (c *ConstraintConverter) setAttributeValue(view *ViewLayout, attr vfl.Attribute, value float64) {
	switch attr {
	case vfl.Width:
		view.Width = value
	case vfl.Height:
		view.Height = value
	case vfl.Left, vfl.Leading:
		view.X = value
	case vfl.Right, vfl.Trailing:
		view.X = value - view.Width
	case vfl.Top:
		view.Y = value
	case vfl.Bottom:
		view.Y = value - view.Height
	case vfl.CenterX:
		view.X = value - view.Width/2
	case vfl.CenterY:
		view.Y = value - view.Height/2
	}
}

// shouldUpdateValue checks if a value should be updated based on relation
func (c *ConstraintConverter) shouldUpdateValue(currentValue, targetValue float64, relation vfl.Relation) bool {
	const epsilon = 0.001

	switch relation {
	case vfl.Equal:
		return abs(currentValue-targetValue) > epsilon
	case vfl.GreaterEqual:
		return currentValue < targetValue-epsilon
	case vfl.LessEqual:
		return currentValue > targetValue+epsilon
	default:
		return false
	}
}

// abs returns absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}