package vfl

import "fmt"

// Validator checks composed layouts for issues
type Validator interface {
	ValidateLayout(layout *LayoutResult) []ValidationError
	ValidateHierarchy(tree *ViewTree) []ValidationError
	CheckConflicts(constraints []Constraint) []ConflictingConstraint
}

// validator implements the Validator interface
type validator struct{}

// NewValidator creates a new validator
func NewValidator() Validator {
	return &validator{}
}

// ValidateLayout checks a layout for constraint conflicts
func (v *validator) ValidateLayout(layout *LayoutResult) []ValidationError {
	var errors []ValidationError

	// Check for conflicts
	conflicts := v.CheckConflicts(layout.Constraints)
	for _, conflict := range conflicts {
		errors = append(errors, ValidationError{
			Message:  fmt.Sprintf("Conflicting constraints: %s", conflict.Reason),
			Severity: Error,
		})
	}

	// Validate hierarchy
	if layout.ViewHierarchy != nil {
		errors = append(errors, v.ValidateHierarchy(layout.ViewHierarchy)...)
	}

	return errors
}

// ValidateHierarchy checks view hierarchy validity
func (v *validator) ValidateHierarchy(tree *ViewTree) []ValidationError {
	var errors []ValidationError

	// Check for orphaned views
	for name, node := range tree.Views {
		if node.Parent == nil && tree.Root != nil && node != tree.Root {
			// Check if it's connected to any other view
			hasConnection := false
			for _, constraint := range node.Constraints {
				if constraint.SecondItem != nil {
					hasConnection = true
					break
				}
			}

			if !hasConnection {
				errors = append(errors, ValidationError{
					Message:  fmt.Sprintf("View '%s' is not connected to the hierarchy", name),
					Severity: Warning,
				})
			}
		}
	}

	return errors
}

// CheckConflicts detects conflicting constraints
func (v *validator) CheckConflicts(constraints []Constraint) []ConflictingConstraint {
	var conflicts []ConflictingConstraint

	// Group constraints by view and attribute
	type key struct {
		view string
		attr Attribute
	}
	constraintMap := make(map[key][]Constraint)

	for _, c := range constraints {
		k := key{view: c.FirstItem.Name, attr: c.FirstAttribute}
		constraintMap[k] = append(constraintMap[k], c)
	}

	// Check for conflicts within each group
	for k, group := range constraintMap {
		if len(group) > 1 {
			for i := 0; i < len(group)-1; i++ {
				for j := i + 1; j < len(group); j++ {
					if v.constraintsConflict(group[i], group[j]) {
						conflicts = append(conflicts, ConflictingConstraint{
							First:  group[i],
							Second: group[j],
							Reason: fmt.Sprintf("Incompatible constraints for %s.%s", k.view, k.attr),
						})
					}
				}
			}
		}
	}

	return conflicts
}

func (v *validator) constraintsConflict(c1, c2 Constraint) bool {
	// Both are equality constraints with different constants
	if c1.Relation == Equal && c2.Relation == Equal {
		if c1.SecondItem == nil && c2.SecondItem == nil {
			return c1.Constant != c2.Constant
		}
	}

	// Equal constraint conflicts with inequality
	if c1.Relation == Equal && c2.Relation == LessEqual {
		if c1.SecondItem == nil && c2.SecondItem == nil {
			return c1.Constant > c2.Constant
		}
	}

	if c1.Relation == Equal && c2.Relation == GreaterEqual {
		if c1.SecondItem == nil && c2.SecondItem == nil {
			return c1.Constant < c2.Constant
		}
	}

	// Conflicting inequalities
	if c1.Relation == GreaterEqual && c2.Relation == LessEqual {
		if c1.SecondItem == nil && c2.SecondItem == nil {
			return c1.Constant > c2.Constant
		}
	}

	return false
}