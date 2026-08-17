package vfl_test

import (
	"testing"

	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLayoutResult_Generation(t *testing.T) {
	composer := vfl.NewComposer()
	parser := vfl.NewParser()

	// Create a program with various constraints
	prog, err := parser.Parse("H:|[button(100)]-[label(>=50)]-[icon(20@750)]|")
	require.NoError(t, err)

	result, err := composer.Compose(prog)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Check constraints generation
	assert.NotEmpty(t, result.Constraints)

	// Verify constraint structure
	for _, constraint := range result.Constraints {
		// Each constraint should have valid items
		assert.NotEmpty(t, constraint.FirstItem.Name)
		assert.NotEqual(t, vfl.NotAnAttribute, constraint.FirstAttribute)

		// Check relations
		assert.Contains(t, []vfl.Relation{vfl.Equal, vfl.GreaterEqual, vfl.LessEqual}, constraint.Relation)

		// Priority should be within valid range
		assert.GreaterOrEqual(t, constraint.Priority, 0)
		assert.LessOrEqual(t, constraint.Priority, 1000)
	}
}

func TestLayoutResult_ViewHierarchy(t *testing.T) {
	composer := vfl.NewComposer()
	parser := vfl.NewParser()

	tests := []struct {
		name        string
		input       string
		checkTree   func(t *testing.T, tree *vfl.ViewTree)
	}{
		{
			name:  "flat hierarchy",
			input: "H:|[view1]-[view2]-[view3]|",
			checkTree: func(t *testing.T, tree *vfl.ViewTree) {
				assert.Len(t, tree.Views, 3)
				assert.Contains(t, tree.Views, "view1")
				assert.Contains(t, tree.Views, "view2")
				assert.Contains(t, tree.Views, "view3")

				// All views should be at same level (no parent-child)
				for _, node := range tree.Views {
					assert.Empty(t, node.Children)
				}
			},
		},
		{
			name:  "with superview constraints",
			input: "H:|[content]|",
			checkTree: func(t *testing.T, tree *vfl.ViewTree) {
				// Should have superview as root
				assert.NotNil(t, tree.Root)
				assert.True(t, tree.Root.IsSuperview)

				// Content should be child of superview
				contentNode := tree.Views["content"]
				assert.NotNil(t, contentNode)
				assert.Equal(t, tree.Root, contentNode.Parent)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := parser.Parse(tt.input)
			require.NoError(t, err)

			result, err := composer.Compose(prog)
			require.NoError(t, err)
			require.NotNil(t, result.ViewHierarchy)

			if tt.checkTree != nil {
				tt.checkTree(t, result.ViewHierarchy)
			}
		})
	}
}

func TestLayoutResult_Metrics(t *testing.T) {
	composer := vfl.NewComposer()
	parser := vfl.NewParser()

	options := vfl.ParserOptions{
		Metrics: map[string]float64{
			"standardSpace": 8,
			"buttonWidth":   100,
			"iconSize":      24,
		},
	}

	prog, err := parser.ParseWithOptions(
		"H:|[button(buttonWidth)]-standardSpace-[icon(iconSize)]|",
		options,
	)
	require.NoError(t, err)

	result, err := composer.Compose(prog)
	require.NoError(t, err)

	// Metrics should be available in result
	assert.Equal(t, options.Metrics, result.Metrics)

	// Constraints should use resolved metric values
	hasButtonWidthConstraint := false
	hasIconSizeConstraint := false

	for _, c := range result.Constraints {
		if c.FirstItem.Name == "button" && c.Constant == 100 {
			hasButtonWidthConstraint = true
		}
		if c.FirstItem.Name == "icon" && c.Constant == 24 {
			hasIconSizeConstraint = true
		}
	}

	assert.True(t, hasButtonWidthConstraint, "Should have button width constraint")
	assert.True(t, hasIconSizeConstraint, "Should have icon size constraint")
}