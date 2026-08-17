package integration_test

import (
	"testing"

	"github.com/chazu/go-vfl/evfl"
	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPercentageConstraints(t *testing.T) {
	parser := evfl.NewParser(evfl.Options{
		EnablePercentages: true,
		EnableExpressions: true,
	})

	tests := []struct {
		name         string
		input        string
		containerSize evfl.Size
		check        func(t *testing.T, result *vfl.LayoutResult)
	}{
		{
			name:  "basic percentage layout",
			input: "H:|[sidebar(==25%)]-[content(==75%)]|",
			containerSize: evfl.Size{
				Width:  1000,
				Height: 600,
			},
			check: func(t *testing.T, result *vfl.LayoutResult) {
				// Should have constraints for percentage-based widths
				sidebarWidth := findConstraintForView(result, "sidebar", vfl.Width)
				contentWidth := findConstraintForView(result, "content", vfl.Width)

				assert.NotNil(t, sidebarWidth)
				assert.NotNil(t, contentWidth)

				// Check calculated values
				assert.Equal(t, 250.0, sidebarWidth.Constant) // 25% of 1000
				assert.Equal(t, 750.0, contentWidth.Constant) // 75% of 1000
			},
		},
		{
			name:  "percentage with minimum constraint",
			input: "H:|[sidebar(>=20%,<=300)]-[content]|",
			containerSize: evfl.Size{
				Width:  1000,
				Height: 600,
			},
			check: func(t *testing.T, result *vfl.LayoutResult) {
				sidebarConstraints := findAllConstraintsForView(result, "sidebar", vfl.Width)
				assert.GreaterOrEqual(t, len(sidebarConstraints), 2)

				// Should have both minimum percentage and maximum fixed
				hasMinPercentage := false
				hasMaxFixed := false

				for _, c := range sidebarConstraints {
					if c.Relation == vfl.GreaterEqual && c.Constant == 200.0 {
						hasMinPercentage = true
					}
					if c.Relation == vfl.LessEqual && c.Constant == 300.0 {
						hasMaxFixed = true
					}
				}

				assert.True(t, hasMinPercentage, "Should have >=20% (200px) constraint")
				assert.True(t, hasMaxFixed, "Should have <=300px constraint")
			},
		},
		{
			name:  "nested percentage calculations",
			input: "H:|[container(==50%):[left(==50%)][right(==50%)]]|",
			containerSize: evfl.Size{
				Width:  1000,
				Height: 600,
			},
			check: func(t *testing.T, result *vfl.LayoutResult) {
				// Container should be 50% of parent (500px)
				containerWidth := findConstraintForView(result, "container", vfl.Width)
				assert.Equal(t, 500.0, containerWidth.Constant)

				// Nested views should be 50% of container (250px each)
				leftWidth := findConstraintForView(result, "left", vfl.Width)
				rightWidth := findConstraintForView(result, "right", vfl.Width)
				assert.Equal(t, 250.0, leftWidth.Constant)
				assert.Equal(t, 250.0, rightWidth.Constant)
			},
		},
		{
			name:  "percentage with expression",
			input: "[view(==50%+20)]",
			containerSize: evfl.Size{
				Width:  1000,
				Height: 600,
			},
			check: func(t *testing.T, result *vfl.LayoutResult) {
				// 50% of width + 20 = 500 + 20 = 520
				viewWidth := findConstraintForView(result, "view", vfl.Width)
				assert.Equal(t, 520.0, viewWidth.Constant)
			},
		},
		{
			name:  "responsive layout with percentages",
			input: "V:|[header(==10%)]-[content(==80%)]-[footer(==10%)]|",
			containerSize: evfl.Size{
				Width:  800,
				Height: 1000,
			},
			check: func(t *testing.T, result *vfl.LayoutResult) {
				headerHeight := findConstraintForView(result, "header", vfl.Height)
				contentHeight := findConstraintForView(result, "content", vfl.Height)
				footerHeight := findConstraintForView(result, "footer", vfl.Height)

				assert.Equal(t, 100.0, headerHeight.Constant)  // 10% of 1000
				assert.Equal(t, 800.0, contentHeight.Constant) // 80% of 1000
				assert.Equal(t, 100.0, footerHeight.Constant)  // 10% of 1000
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := parser.Parse(tt.input)
			require.NoError(t, err)

			composer := evfl.NewComposer()
			result, err := composer.ComposeWithSize(prog, tt.containerSize)
			require.NoError(t, err)
			require.NotNil(t, result)

			if tt.check != nil {
				tt.check(t, result)
			}
		})
	}
}

func TestPercentageIntegrationWithComposer(t *testing.T) {
	parser := evfl.NewParser(evfl.Options{
		EnablePercentages: true,
		EnableStacks:      true,
	})
	composer := evfl.NewComposer()

	// Register a named view with percentages
	headerProg, err := parser.Parse("H:|[logo(==10%)]-[title(==60%)]-[menu(==30%)]|")
	require.NoError(t, err)
	err = composer.RegisterNamedView("header", headerProg)
	require.NoError(t, err)

	// Main layout using the named view
	mainProg, err := parser.Parse("V:|[header(==15%)]-[content(==85%)]|")
	require.NoError(t, err)

	// Compose with container size
	containerSize := evfl.Size{Width: 1200, Height: 800}
	result, err := composer.ComposeWithSize(mainProg, containerSize)
	require.NoError(t, err)

	// Check header height (15% of 800 = 120)
	headerConstraint := findConstraintForView(result, "header", vfl.Height)
	assert.Equal(t, 120.0, headerConstraint.Constant)

	// Check nested logo width (10% of 1200 = 120)
	logoConstraint := findConstraintForView(result, "logo", vfl.Width)
	assert.Equal(t, 120.0, logoConstraint.Constant)

	// Check nested title width (60% of 1200 = 720)
	titleConstraint := findConstraintForView(result, "title", vfl.Width)
	assert.Equal(t, 720.0, titleConstraint.Constant)
}

// Helper functions for finding constraints
func findConstraintForView(result *vfl.LayoutResult, viewName string, attr vfl.Attribute) *vfl.Constraint {
	for _, c := range result.Constraints {
		if c.FirstItem.Name == viewName && c.FirstAttribute == attr {
			return &c
		}
	}
	return nil
}

func findAllConstraintsForView(result *vfl.LayoutResult, viewName string, attr vfl.Attribute) []vfl.Constraint {
	var constraints []vfl.Constraint
	for _, c := range result.Constraints {
		if c.FirstItem.Name == viewName && c.FirstAttribute == attr {
			constraints = append(constraints, c)
		}
	}
	return constraints
}