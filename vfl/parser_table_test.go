package vfl

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		validate    func(t *testing.T, prog *Program)
		errContains string
	}{
		// Basic horizontal layouts
		{
			name:  "simple horizontal",
			input: "H:|[view]|",
			validate: func(t *testing.T, prog *Program) {
				assert.Equal(t, Horizontal, prog.Orientation)
				assert.Len(t, prog.Statements, 1)
				assert.Len(t, prog.Statements[0].Views, 1)
				assert.Equal(t, "view", prog.Statements[0].Views[0].Name)
			},
		},
		{
			name:  "horizontal with spacing",
			input: "H:|-[button]-|",
			validate: func(t *testing.T, prog *Program) {
				assert.Equal(t, Horizontal, prog.Orientation)
				assert.True(t, prog.Statements[0].SuperviewStart)
				assert.True(t, prog.Statements[0].SuperviewEnd)
			},
		},
		{
			name:  "horizontal with custom spacing",
			input: "H:|-20-[button]-20-|",
			validate: func(t *testing.T, prog *Program) {
				assert.Equal(t, Horizontal, prog.Orientation)
				assert.Len(t, prog.Statements[0].Connections, 2)
				assert.Equal(t, float64(20), prog.Statements[0].Connections[0].Spacing)
			},
		},
		// Vertical layouts
		{
			name:  "simple vertical",
			input: "V:|[header][body][footer]|",
			validate: func(t *testing.T, prog *Program) {
				assert.Equal(t, Vertical, prog.Orientation)
				assert.Len(t, prog.Statements[0].Views, 3)
			},
		},
		{
			name:  "vertical with spacing",
			input: "V:|-[top]-[bottom]-|",
			validate: func(t *testing.T, prog *Program) {
				assert.Equal(t, Vertical, prog.Orientation)
				assert.Len(t, prog.Statements[0].Views, 2)
				assert.Len(t, prog.Statements[0].Connections, 1)
			},
		},
		// Constraints
		{
			name:  "width constraint",
			input: "H:[button(100)]",
			validate: func(t *testing.T, prog *Program) {
				view := prog.Statements[0].Views[0]
				assert.Len(t, view.Predicates, 1)
				assert.Equal(t, Equal, view.Predicates[0].Relation)
				assert.Equal(t, float64(100), view.Predicates[0].Value.ConstantValue)
			},
		},
		{
			name:  "greater than constraint",
			input: "H:[button(>=50)]",
			validate: func(t *testing.T, prog *Program) {
				view := prog.Statements[0].Views[0]
				assert.Equal(t, GreaterThanOrEqual, view.Predicates[0].Relation)
				assert.Equal(t, float64(50), view.Predicates[0].Value.ConstantValue)
			},
		},
		{
			name:  "less than constraint",
			input: "H:[button(<=200)]",
			validate: func(t *testing.T, prog *Program) {
				view := prog.Statements[0].Views[0]
				assert.Equal(t, LessThanOrEqual, view.Predicates[0].Relation)
				assert.Equal(t, float64(200), view.Predicates[0].Value.ConstantValue)
			},
		},
		{
			name:  "priority constraint",
			input: "H:[button(100@750)]",
			validate: func(t *testing.T, prog *Program) {
				view := prog.Statements[0].Views[0]
				assert.Equal(t, 750, view.Priority)
			},
		},
		{
			name:  "view reference constraint",
			input: "H:[button(==label)]",
			validate: func(t *testing.T, prog *Program) {
				view := prog.Statements[0].Views[0]
				assert.Equal(t, ViewReference, view.Predicates[0].Value.Type)
				assert.Equal(t, "label", view.Predicates[0].Value.ViewRef)
			},
		},
		// Multiple views
		{
			name:  "multiple views horizontal",
			input: "H:[button]-[label]-[textField]",
			validate: func(t *testing.T, prog *Program) {
				assert.Len(t, prog.Statements[0].Views, 3)
				assert.Equal(t, "button", prog.Statements[0].Views[0].Name)
				assert.Equal(t, "label", prog.Statements[0].Views[1].Name)
				assert.Equal(t, "textField", prog.Statements[0].Views[2].Name)
			},
		},
		{
			name:  "flush views",
			input: "H:[button][label]",
			validate: func(t *testing.T, prog *Program) {
				assert.Len(t, prog.Statements[0].Views, 2)
				assert.Len(t, prog.Statements[0].Connections, 1)
				assert.Equal(t, float64(0), prog.Statements[0].Connections[0].Spacing)
			},
		},
		// Default orientation
		{
			name:  "default to horizontal",
			input: "[button]-[label]",
			validate: func(t *testing.T, prog *Program) {
				assert.Equal(t, Horizontal, prog.Orientation)
			},
		},
		// Metrics
		{
			name:  "standard metric",
			input: "H:|-standard-[button]-standard-|",
			validate: func(t *testing.T, prog *Program) {
				assert.Len(t, prog.Statements[0].Connections, 2)
				// Metrics would be resolved during composition
			},
		},
		// Complex layouts
		{
			name:  "complex layout",
			input: "H:|[sidebar(200)]-[content(>=100)]-[ads(100@250)]|",
			validate: func(t *testing.T, prog *Program) {
				assert.Len(t, prog.Statements[0].Views, 3)

				sidebar := prog.Statements[0].Views[0]
				assert.Equal(t, "sidebar", sidebar.Name)
				assert.Equal(t, float64(200), sidebar.Predicates[0].Value.ConstantValue)

				content := prog.Statements[0].Views[1]
				assert.Equal(t, "content", content.Name)
				assert.Equal(t, GreaterThanOrEqual, content.Predicates[0].Relation)

				ads := prog.Statements[0].Views[2]
				assert.Equal(t, "ads", ads.Name)
				assert.Equal(t, 250, ads.Priority)
			},
		},
		// Error cases
		{
			name:        "invalid syntax",
			input:       "H:|[button",
			wantErr:     true,
			errContains: "syntax error",
		},
		{
			name:        "invalid orientation",
			input:       "X:|[button]|",
			wantErr:     true,
			errContains: "invalid orientation",
		},
		{
			name:        "empty view name",
			input:       "H:|[]|",
			wantErr:     true,
			errContains: "empty view name",
		},
	}

	parser := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := parser.Parse(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, prog)
				if tt.validate != nil {
					tt.validate(t, prog)
				}
			}
		})
	}
}

func TestParser_EdgeCases(t *testing.T) {
	parser := NewParser()

	t.Run("empty input", func(t *testing.T) {
		_, err := parser.Parse("")
		assert.Error(t, err)
	})

	t.Run("whitespace only", func(t *testing.T) {
		_, err := parser.Parse("   \n\t  ")
		assert.Error(t, err)
	})

	t.Run("multiple statements", func(t *testing.T) {
		// Most VFL parsers don't support multiple statements in one string
		_, err := parser.Parse("H:|[view1]|;V:|[view2]|")
		// This might be an error or might parse only the first statement
		// depending on implementation
		_ = err // Implementation specific
	})

	t.Run("unicode view names", func(t *testing.T) {
		prog, err := parser.Parse("H:|[按钮]-[标签]|")
		if err == nil {
			assert.NotNil(t, prog)
			// Some parsers might support unicode identifiers
		}
	})

	t.Run("very long view name", func(t *testing.T) {
		longName := "view"
		for i := 0; i < 100; i++ {
			longName += "Very"
		}
		input := "H:|[" + longName + "]|"
		prog, err := parser.Parse(input)
		if err == nil {
			assert.Equal(t, longName, prog.Statements[0].Views[0].Name)
		}
	})

	t.Run("nested brackets", func(t *testing.T) {
		// This should fail as VFL doesn't support nested brackets
		_, err := parser.Parse("H:|[[nested]]|")
		assert.Error(t, err)
	})

	t.Run("negative spacing", func(t *testing.T) {
		prog, err := parser.Parse("H:|[view1]-(-10)-[view2]|")
		if err == nil {
			// Some parsers support negative spacing for overlapping
			conn := prog.Statements[0].Connections[0]
			assert.Equal(t, float64(-10), conn.Spacing)
		}
	})
}

func TestParser_PerformanceRegression(t *testing.T) {
	parser := NewParser()

	// Test that parser doesn't have exponential behavior on certain inputs
	t.Run("deeply nested predicates", func(t *testing.T) {
		input := "H:|[view(==other.width*2+10-5*3+8)]|"
		prog, err := parser.Parse(input)
		// Should complete quickly regardless of expression complexity
		_ = prog
		_ = err
	})

	t.Run("many views", func(t *testing.T) {
		input := "H:|"
		for i := 0; i < 100; i++ {
			if i > 0 {
				input += "-"
			}
			input += "[view" + string(rune('A'+i%26)) + "]"
		}
		input += "|"

		prog, err := parser.Parse(input)
		if err == nil {
			assert.Len(t, prog.Statements[0].Views, 100)
		}
	})
}