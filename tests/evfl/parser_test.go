package evfl_test

import (
	"testing"

	"github.com/chazu/go-vfl/evfl"
	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEVFLParser_ParsePercentage(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		check   func(t *testing.T, prog *vfl.Program)
	}{
		{
			name:  "simple percentage",
			input: "H:|[sidebar(==20%)]-[content(==80%)]|",
			check: func(t *testing.T, prog *vfl.Program) {
				sidebar := prog.Statements[0].Views[0]
				content := prog.Statements[0].Views[1]

				assert.Equal(t, 20.0, sidebar.Predicates[0].Value.Percentage)
				assert.Equal(t, vfl.PercentageValue, sidebar.Predicates[0].Value.Type)

				assert.Equal(t, 80.0, content.Predicates[0].Value.Percentage)
			},
		},
		{
			name:  "percentage with relation",
			input: "[view(>=50%)]",
			check: func(t *testing.T, prog *vfl.Program) {
				view := prog.Statements[0].Views[0]
				assert.Equal(t, evfl.GreaterEqual, view.Predicates[0].Relation)
				assert.Equal(t, 50.0, view.Predicates[0].Value.Percentage)
			},
		},
		{
			name:  "mixed percentage and fixed",
			input: "H:|[fixed(100)]-[flex(==50%)]|",
			check: func(t *testing.T, prog *vfl.Program) {
				fixed := prog.Statements[0].Views[0]
				flex := prog.Statements[0].Views[1]

				assert.Equal(t, vfl.ConstantValue, fixed.Predicates[0].Value.Type)
				assert.Equal(t, 100.0, fixed.Predicates[0].Value.Constant)

				assert.Equal(t, vfl.PercentageValue, flex.Predicates[0].Value.Type)
				assert.Equal(t, 50.0, flex.Predicates[0].Value.Percentage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := evfl.NewParser(evfl.Options{
				EnablePercentages: true,
			})

			prog, err := parser.Parse(tt.input)
			require.NoError(t, err)
			require.NotNil(t, prog)

			if tt.check != nil {
				tt.check(t, prog)
			}
		})
	}
}

func TestEVFLParser_DisabledFeatures(t *testing.T) {
	// Without enabling percentages, should fail
	parser := evfl.NewParser(evfl.Options{
		EnablePercentages: false,
	})

	_, err := parser.Parse("H:|[view(==50%)]|")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "percentage")
}