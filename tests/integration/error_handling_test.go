package integration_test

import (
	"testing"

	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorRecoveryWithPartialAST(t *testing.T) {
	parser := vfl.NewParser()

	tests := []struct {
		name            string
		input           string
		checkPartialAST func(t *testing.T, parseErr *vfl.ParseError)
	}{
		{
			name:  "syntax error mid-parse",
			input: "H:|[valid]-[another]-[invalid(]",
			checkPartialAST: func(t *testing.T, parseErr *vfl.ParseError) {
				require.NotNil(t, parseErr.PartialAST)

				// Should have successfully parsed up to the error
				assert.Equal(t, vfl.Horizontal, parseErr.PartialAST.Orientation)
				assert.Len(t, parseErr.PartialAST.Statements, 1)

				stmt := parseErr.PartialAST.Statements[0]
				assert.GreaterOrEqual(t, len(stmt.Views), 2)
				assert.Equal(t, "valid", stmt.Views[0].Name)
				assert.Equal(t, "another", stmt.Views[1].Name)

				// Error location should point to the problematic part
				assert.Greater(t, parseErr.Location.Column, 20)
			},
		},
		{
			name:  "missing closing bracket",
			input: "V:|[header(50)][content",
			checkPartialAST: func(t *testing.T, parseErr *vfl.ParseError) {
				require.NotNil(t, parseErr.PartialAST)

				// Should have parsed the header
				assert.Equal(t, vfl.Vertical, parseErr.PartialAST.Orientation)
				assert.Len(t, parseErr.PartialAST.Statements, 1)
				assert.GreaterOrEqual(t, len(parseErr.PartialAST.Statements[0].Views), 1)

				header := parseErr.PartialAST.Statements[0].Views[0]
				assert.Equal(t, "header", header.Name)
				assert.Equal(t, 50.0, header.Predicates[0].Value.Constant)
			},
		},
		{
			name:  "invalid operator in middle",
			input: "H:|[left]-[center]<>[right]|",
			checkPartialAST: func(t *testing.T, parseErr *vfl.ParseError) {
				require.NotNil(t, parseErr.PartialAST)

				// Should parse up to invalid operator
				stmt := parseErr.PartialAST.Statements[0]
				assert.GreaterOrEqual(t, len(stmt.Views), 2)
				assert.Equal(t, "left", stmt.Views[0].Name)
				assert.Equal(t, "center", stmt.Views[1].Name)
			},
		},
		{
			name:  "multiline with error on second line",
			input: "H:|[button]-[label]|\nV:|[invalid(]",
			checkPartialAST: func(t *testing.T, parseErr *vfl.ParseError) {
				require.NotNil(t, parseErr.PartialAST)

				// First line should be completely parsed
				assert.Len(t, parseErr.PartialAST.Statements, 1)
				firstStmt := parseErr.PartialAST.Statements[0]
				assert.Equal(t, vfl.Horizontal, parseErr.PartialAST.Orientation)
				assert.Len(t, firstStmt.Views, 2)

				// Error should be on line 2
				assert.Equal(t, 2, parseErr.Location.Line)
			},
		},
		{
			name:  "semantic error with valid syntax",
			input: "[view1(==undefined)]",
			checkPartialAST: func(t *testing.T, parseErr *vfl.ParseError) {
				// Even semantic errors should provide partial AST
				require.NotNil(t, parseErr.PartialAST)
				assert.Equal(t, vfl.SemanticError, parseErr.Type)

				// Structure should be parsed
				assert.Len(t, parseErr.PartialAST.Statements, 1)
				view := parseErr.PartialAST.Statements[0].Views[0]
				assert.Equal(t, "view1", view.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.input)
			require.Error(t, err)

			parseErr, ok := err.(*vfl.ParseError)
			require.True(t, ok, "Expected *vfl.ParseError, got %T", err)

			// All errors should have descriptive messages
			assert.NotEmpty(t, parseErr.Message)
			assert.NotEmpty(t, parseErr.Error())

			if tt.checkPartialAST != nil {
				tt.checkPartialAST(t, parseErr)
			}
		})
	}
}

func TestErrorHandlingStrategies(t *testing.T) {
	t.Run("strict mode - no partial AST", func(t *testing.T) {
		parser := vfl.NewParser()
		options := vfl.ParserOptions{
			StrictMode: true,
		}

		_, err := parser.ParseWithOptions("H:|[invalid(]", options)
		require.Error(t, err)

		if parseErr, ok := err.(*vfl.ParseError); ok {
			assert.Nil(t, parseErr.PartialAST, "Strict mode should not return partial AST")
		}
	})

	t.Run("permissive mode - maximum recovery", func(t *testing.T) {
		parser := vfl.NewParser()
		options := vfl.ParserOptions{
			StrictMode: false,
		}

		input := "H:|[view1]-[view2(invalid)]-[view3]|"
		_, err := parser.ParseWithOptions(input, options)
		require.Error(t, err)

		parseErr, ok := err.(*vfl.ParseError)
		require.True(t, ok)
		require.NotNil(t, parseErr.PartialAST)

		// Should recover as much as possible
		stmt := parseErr.PartialAST.Statements[0]
		assert.GreaterOrEqual(t, len(stmt.Views), 1)
		assert.Equal(t, "view1", stmt.Views[0].Name)
	})
}