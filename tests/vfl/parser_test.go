package vfl_test

import (
	"testing"

	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, prog *vfl.Program)
	}{
		{
			name:  "basic horizontal layout",
			input: "H:|[button(100)]-[textField]-|",
			check: func(t *testing.T, prog *vfl.Program) {
				assert.Equal(t, vfl.Horizontal, prog.Orientation)
				assert.Len(t, prog.Statements, 1)
				assert.Len(t, prog.Statements[0].Views, 2)
				assert.Equal(t, "button", prog.Statements[0].Views[0].Name)
				assert.Equal(t, "textField", prog.Statements[0].Views[1].Name)
			},
		},
		{
			name:  "vertical layout with spacing",
			input: "V:|[top]-10-[middle]-20-[bottom]|",
			check: func(t *testing.T, prog *vfl.Program) {
				assert.Equal(t, vfl.Vertical, prog.Orientation)
				assert.Len(t, prog.Statements, 1)
				assert.Len(t, prog.Statements[0].Views, 3)
				assert.Len(t, prog.Statements[0].Connections, 2)
				assert.Equal(t, 10.0, prog.Statements[0].Connections[0].Spacing)
				assert.Equal(t, 20.0, prog.Statements[0].Connections[1].Spacing)
			},
		},
		{
			name:  "constraints with relations",
			input: "[view(>=100)]",
			check: func(t *testing.T, prog *vfl.Program) {
				view := prog.Statements[0].Views[0]
				assert.Len(t, view.Predicates, 1)
				assert.Equal(t, vfl.GreaterEqual, view.Predicates[0].Relation)
				assert.Equal(t, 100.0, view.Predicates[0].Value.Constant)
			},
		},
		{
			name:  "priority constraints",
			input: "[button(100@750)]",
			check: func(t *testing.T, prog *vfl.Program) {
				view := prog.Statements[0].Views[0]
				assert.Equal(t, 750, view.Priority)
			},
		},
		{
			name:    "invalid syntax",
			input:   "H:|[invalid(]",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := vfl.NewParser()
			prog, err := parser.Parse(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				// Check for partial AST in error
				if parseErr, ok := err.(*vfl.ParseError); ok {
					assert.NotNil(t, parseErr.PartialAST)
				}
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