package evfl_test

import (
	"testing"

	"github.com/chazu/go-vfl/evfl"
	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEVFLParser_ViewStacks(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		check   func(t *testing.T, prog *evfl.Program)
		wantErr bool
	}{
		{
			name:  "simple vertical stack",
			input: "V:|[column:[header][content][footer]]|",
			check: func(t *testing.T, prog *evfl.Program) {
				// Should parse as a view stack
				assert.Len(t, prog.Statements, 1)

				// The column should be a stack container
				stack := prog.Statements[0].Views[0]
				assert.Equal(t, "column", stack.Name)
				assert.True(t, stack.IsStack)

				// Check nested views
				assert.Len(t, stack.StackViews, 3)
				assert.Equal(t, "header", stack.StackViews[0].Name)
				assert.Equal(t, "content", stack.StackViews[1].Name)
				assert.Equal(t, "footer", stack.StackViews[2].Name)

				// Stack should have vertical orientation
				assert.Equal(t, evfl.Vertical, stack.StackOrientation)
			},
		},
		{
			name:  "horizontal stack",
			input: "H:|[row:[left][center][right]]|",
			check: func(t *testing.T, prog *evfl.Program) {
				stack := prog.Statements[0].Views[0]
				assert.Equal(t, "row", stack.Name)
				assert.True(t, stack.IsStack)
				assert.Equal(t, evfl.Horizontal, stack.StackOrientation)
				assert.Len(t, stack.StackViews, 3)
			},
		},
		{
			name:  "stack with constraints",
			input: "V:|[column:[header(50)][content(>=100)][footer(==header)]]|",
			check: func(t *testing.T, prog *evfl.Program) {
				stack := prog.Statements[0].Views[0]

				// Header with fixed height
				header := stack.StackViews[0]
				assert.Len(t, header.Predicates, 1)
				assert.Equal(t, 50.0, header.Predicates[0].Value.Constant)

				// Content with minimum height
				content := stack.StackViews[1]
				assert.Equal(t, evfl.GreaterEqual, content.Predicates[0].Relation)
				assert.Equal(t, 100.0, content.Predicates[0].Value.Constant)

				// Footer equal to header
				footer := stack.StackViews[2]
				assert.Equal(t, "header", footer.Predicates[0].Value.ViewRef)
			},
		},
		{
			name:  "nested stacks",
			input: "V:|[main:[header][body:[sidebar][content]][footer]]|",
			check: func(t *testing.T, prog *evfl.Program) {
				mainStack := prog.Statements[0].Views[0]
				assert.Equal(t, "main", mainStack.Name)
				assert.Len(t, mainStack.StackViews, 3)

				// Body should also be a stack
				body := mainStack.StackViews[1]
				assert.Equal(t, "body", body.Name)
				assert.True(t, body.IsStack)
				assert.Len(t, body.StackViews, 2)
				assert.Equal(t, "sidebar", body.StackViews[0].Name)
				assert.Equal(t, "content", body.StackViews[1].Name)
			},
		},
		{
			name:  "stack with spacing",
			input: "V:|[column:[view1]-10-[view2]-20-[view3]]|",
			check: func(t *testing.T, prog *evfl.Program) {
				stack := prog.Statements[0].Views[0]
				assert.Len(t, stack.StackConnections, 2)
				assert.Equal(t, 10.0, stack.StackConnections[0].Spacing)
				assert.Equal(t, 20.0, stack.StackConnections[1].Spacing)
			},
		},
		{
			name:  "grid stack layout",
			input: "HV:[grid:[row1:[cell1][cell2]][row2:[cell3][cell4]]]",
			check: func(t *testing.T, prog *evfl.Program) {
				grid := prog.Statements[0].Views[0]
				assert.Equal(t, vfl.GridStack, grid.StackLayout)
				assert.Len(t, grid.StackViews, 2)
			},
		},
		{
			name:    "invalid stack syntax",
			input:   "V:|[column:[]]|",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := evfl.NewParser(evfl.Options{
				EnableStacks: true,
			})

			prog, err := parser.Parse(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, prog)

			if tt.check != nil {
				tt.check(t, prog)
			}
		})
	}
}

func TestEVFLParser_ParseStack(t *testing.T) {
	parser := evfl.NewParser(evfl.Options{
		EnableStacks: true,
	})

	// Test the ParseStack method specifically
	stack, err := parser.ParseStack("column:[header][content][footer]")
	require.NoError(t, err)
	require.NotNil(t, stack)

	assert.Equal(t, "column", stack.Name)
	assert.Len(t, stack.Views, 3)
	assert.Equal(t, vfl.LinearStack, stack.Layout)
}