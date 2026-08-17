package vfl_test

import (
	"testing"

	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_ParseWithOptions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		options vfl.ParserOptions
		check   func(t *testing.T, prog *vfl.Program)
	}{
		{
			name:  "with custom lookahead",
			input: "H:|[view1]-[view2]-[view3]-[view4]-[view5]|",
			options: vfl.ParserOptions{
				Lookahead: 10,
			},
			check: func(t *testing.T, prog *vfl.Program) {
				assert.Len(t, prog.Statements[0].Views, 5)
			},
		},
		{
			name:  "with metrics",
			input: "H:|[button(buttonWidth)]-standardSpace-[label]|",
			options: vfl.ParserOptions{
				Metrics: map[string]float64{
					"buttonWidth":    100,
					"standardSpace":  8,
				},
			},
			check: func(t *testing.T, prog *vfl.Program) {
				// Check metric resolution
				view := prog.Statements[0].Views[0]
				assert.Len(t, view.Predicates, 1)
				assert.Equal(t, "buttonWidth", view.Predicates[0].Value.MetricName)

				// Check connection metric
				conn := prog.Statements[0].Connections[0]
				assert.Equal(t, 8.0, conn.Spacing)
			},
		},
		{
			name:  "strict mode - no partial AST on error",
			input: "H:|[invalid(]",
			options: vfl.ParserOptions{
				StrictMode: true,
			},
			check: func(t *testing.T, prog *vfl.Program) {
				// In strict mode, should fail without partial AST
				assert.Nil(t, prog)
			},
		},
		{
			name:  "EVFL disabled by default",
			input: "H:|[view(==50%)]|",
			options: vfl.ParserOptions{
				EnableEVFL: false,
			},
			check: func(t *testing.T, prog *vfl.Program) {
				// Should fail to parse percentage without EVFL
				assert.Nil(t, prog)
			},
		},
		{
			name:  "EVFL enabled",
			input: "H:|[view(==50%)]|",
			options: vfl.ParserOptions{
				EnableEVFL: true,
			},
			check: func(t *testing.T, prog *vfl.Program) {
				view := prog.Statements[0].Views[0]
				assert.Equal(t, 50.0, view.Predicates[0].Value.Percentage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := vfl.NewParser()
			prog, err := parser.ParseWithOptions(tt.input, tt.options)

			if tt.name == "strict mode - no partial AST on error" || tt.name == "EVFL disabled by default" {
				assert.Error(t, err)
				if tt.options.StrictMode {
					if parseErr, ok := err.(*vfl.ParseError); ok {
						assert.Nil(t, parseErr.PartialAST)
					}
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