package evfl

import (
	"fmt"
	"strings"
	"github.com/chazu/go-vfl/vfl"
)

// StackParser parses view stack notation
type StackParser interface {
	Parse(input string) (*vfl.View, error)
	IsStackNotation(input string) bool
}

// stackParser implements StackParser
type stackParser struct {
	baseParser vfl.Parser
}

// NewStackParser creates a new stack parser
func NewStackParser(baseParser vfl.Parser) StackParser {
	return &stackParser{
		baseParser: baseParser,
	}
}

// IsStackNotation checks if input uses stack notation
func (p *stackParser) IsStackNotation(input string) bool {
	// Stack notation uses {} for grouping and :: for column/row stacks
	return strings.Contains(input, "{") || strings.Contains(input, "::")
}

// Parse parses view stack notation
func (p *stackParser) Parse(input string) (*vfl.View, error) {
	// Remove outer brackets if present
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]") {
		input = input[1 : len(input)-1]
	}

	// Check for column stack (default)
	if strings.Contains(input, "column{") || strings.Contains(input, "col{") {
		return p.parseColumnStack(input)
	}

	// Check for row stack
	if strings.Contains(input, "row{") {
		return p.parseRowStack(input)
	}

	// Check for grid stack
	if strings.Contains(input, "grid{") {
		return p.parseGridStack(input)
	}

	// Default to column stack if just using {}
	if strings.Contains(input, "{") {
		return p.parseColumnStack(input)
	}

	return nil, fmt.Errorf("not a stack notation: %s", input)
}

// parseColumnStack parses a column stack (vertical)
func (p *stackParser) parseColumnStack(input string) (*vfl.View, error) {
	stack := &vfl.View{
		Name:             p.extractStackName(input),
		IsStack:          true,
		StackOrientation: vfl.Vertical,
		StackLayout:      vfl.LinearStack,
		StackViews:       []vfl.View{},
		StackConnections: []vfl.Connection{},
	}

	// Extract view names from the stack
	viewNames := p.extractViewNames(input)
	for i, name := range viewNames {
		stack.StackViews = append(stack.StackViews, vfl.View{
			Name: name,
		})

		// Add default connection between views
		if i < len(viewNames)-1 {
			stack.StackConnections = append(stack.StackConnections, vfl.Connection{
				IsDefault: true,
				Spacing:   8,
			})
		}
	}

	return stack, nil
}

// parseRowStack parses a row stack (horizontal)
func (p *stackParser) parseRowStack(input string) (*vfl.View, error) {
	stack := &vfl.View{
		Name:             p.extractStackName(input),
		IsStack:          true,
		StackOrientation: vfl.Horizontal,
		StackLayout:      vfl.LinearStack,
		StackViews:       []vfl.View{},
		StackConnections: []vfl.Connection{},
	}

	// Extract view names from the stack
	viewNames := p.extractViewNames(input)
	for i, name := range viewNames {
		stack.StackViews = append(stack.StackViews, vfl.View{
			Name: name,
		})

		// Add default connection between views
		if i < len(viewNames)-1 {
			stack.StackConnections = append(stack.StackConnections, vfl.Connection{
				IsDefault: true,
				Spacing:   8,
			})
		}
	}

	return stack, nil
}

// parseGridStack parses a grid stack
func (p *stackParser) parseGridStack(input string) (*vfl.View, error) {
	stack := &vfl.View{
		Name:        p.extractStackName(input),
		IsStack:     true,
		StackLayout: vfl.GridStack,
		StackViews:  []vfl.View{},
	}

	// Extract view names from the stack
	viewNames := p.extractViewNames(input)
	for _, name := range viewNames {
		stack.StackViews = append(stack.StackViews, vfl.View{
			Name: name,
		})
	}

	return stack, nil
}

// extractStackName extracts the stack name from notation
func (p *stackParser) extractStackName(input string) string {
	// Look for name before {
	idx := strings.Index(input, "{")
	if idx > 0 {
		name := strings.TrimSpace(input[:idx])
		// Remove stack type prefix if present
		name = strings.TrimPrefix(name, "column")
		name = strings.TrimPrefix(name, "col")
		name = strings.TrimPrefix(name, "row")
		name = strings.TrimPrefix(name, "grid")
		name = strings.TrimSpace(name)
		if name != "" {
			return name
		}
	}
	return "stack"
}

// extractViewNames extracts view names from stack notation
func (p *stackParser) extractViewNames(input string) []string {
	// Find content between { and }
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start < 0 || end < 0 || start >= end {
		return []string{}
	}

	content := input[start+1 : end]

	// Split by common separators
	var names []string
	if strings.Contains(content, ",") {
		// Comma-separated
		parts := strings.Split(content, ",")
		for _, part := range parts {
			name := strings.TrimSpace(part)
			if name != "" {
				names = append(names, name)
			}
		}
	} else if strings.Contains(content, "::") {
		// Double-colon separated (alternative notation)
		parts := strings.Split(content, "::")
		for _, part := range parts {
			name := strings.TrimSpace(part)
			if name != "" {
				names = append(names, name)
			}
		}
	} else {
		// Space-separated or single item
		parts := strings.Fields(content)
		names = parts
	}

	return names
}

// Stack layout constants use values from vfl package

// CreateStackConstraints generates constraints for a stack view
func CreateStackConstraints(stack *vfl.View) []vfl.Constraint {
	var constraints []vfl.Constraint

	if !stack.IsStack || len(stack.StackViews) == 0 {
		return constraints
	}

	// Generate constraints based on stack orientation
	if stack.StackOrientation == vfl.Vertical {
		// Vertical stack: views are arranged top to bottom
		for i, view := range stack.StackViews {
			if i == 0 {
				// First view pins to top
				constraints = append(constraints, vfl.Constraint{
					FirstItem:      vfl.ConstraintItem{Name: view.Name},
					FirstAttribute: vfl.Top,
					Relation:       vfl.Equal,
					SecondItem:     &vfl.ConstraintItem{Name: stack.Name},
					SecondAttribute: vfl.Top,
					Constant:       0,
					Multiplier:     1.0,
					Priority:       1000,
				})
			} else {
				// Subsequent views pin to previous view's bottom
				prevView := stack.StackViews[i-1]
				spacing := float64(8) // default
				if i-1 < len(stack.StackConnections) {
					spacing = stack.StackConnections[i-1].Spacing
				}

				constraints = append(constraints, vfl.Constraint{
					FirstItem:       vfl.ConstraintItem{Name: view.Name},
					FirstAttribute:  vfl.Top,
					Relation:        vfl.Equal,
					SecondItem:      &vfl.ConstraintItem{Name: prevView.Name},
					SecondAttribute: vfl.Bottom,
					Constant:        spacing,
					Multiplier:      1.0,
					Priority:        1000,
				})
			}

			// All views have same width as stack
			constraints = append(constraints, vfl.Constraint{
				FirstItem:       vfl.ConstraintItem{Name: view.Name},
				FirstAttribute:  vfl.Width,
				Relation:        vfl.Equal,
				SecondItem:      &vfl.ConstraintItem{Name: stack.Name},
				SecondAttribute: vfl.Width,
				Constant:        0,
				Multiplier:      1.0,
				Priority:        1000,
			})
		}
	} else if stack.StackOrientation == vfl.Horizontal {
		// Horizontal stack: views are arranged left to right
		for i, view := range stack.StackViews {
			if i == 0 {
				// First view pins to left
				constraints = append(constraints, vfl.Constraint{
					FirstItem:       vfl.ConstraintItem{Name: view.Name},
					FirstAttribute:  vfl.Left,
					Relation:        vfl.Equal,
					SecondItem:      &vfl.ConstraintItem{Name: stack.Name},
					SecondAttribute: vfl.Left,
					Constant:        0,
					Multiplier:      1.0,
					Priority:        1000,
				})
			} else {
				// Subsequent views pin to previous view's right
				prevView := stack.StackViews[i-1]
				spacing := float64(8) // default
				if i-1 < len(stack.StackConnections) {
					spacing = stack.StackConnections[i-1].Spacing
				}

				constraints = append(constraints, vfl.Constraint{
					FirstItem:       vfl.ConstraintItem{Name: view.Name},
					FirstAttribute:  vfl.Left,
					Relation:        vfl.Equal,
					SecondItem:      &vfl.ConstraintItem{Name: prevView.Name},
					SecondAttribute: vfl.Right,
					Constant:        spacing,
					Multiplier:      1.0,
					Priority:        1000,
				})
			}

			// All views have same height as stack
			constraints = append(constraints, vfl.Constraint{
				FirstItem:       vfl.ConstraintItem{Name: view.Name},
				FirstAttribute:  vfl.Height,
				Relation:        vfl.Equal,
				SecondItem:      &vfl.ConstraintItem{Name: stack.Name},
				SecondAttribute: vfl.Height,
				Constant:        0,
				Multiplier:      1.0,
				Priority:        1000,
			})
		}
	}

	return constraints
}