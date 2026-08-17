package vfl_test

import (
	"testing"

	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseError_PartialAST(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		checkPartialAST  func(t *testing.T, err *vfl.ParseError)
	}{
		{
			name:  "partial AST for incomplete input",
			input: "H:|[valid]-[another]-[invalid(",
			checkPartialAST: func(t *testing.T, err *vfl.ParseError) {
				require.NotNil(t, err.PartialAST)
				assert.Equal(t, vfl.Horizontal, err.PartialAST.Orientation)
				// Should have parsed the first two views
				assert.GreaterOrEqual(t, len(err.PartialAST.Statements[0].Views), 2)
				assert.Equal(t, "valid", err.PartialAST.Statements[0].Views[0].Name)
				assert.Equal(t, "another", err.PartialAST.Statements[0].Views[1].Name)
			},
		},
		{
			name:  "error location tracking",
			input: "H:|[view1]-\n[view2(invalid)]",
			checkPartialAST: func(t *testing.T, err *vfl.ParseError) {
				require.NotNil(t, err.Location)
				// Error should be on line 2
				assert.Equal(t, 2, err.Location.Line)
				assert.Greater(t, err.Location.Column, 0)
				assert.Greater(t, err.Location.Offset, 0)
			},
		},
		{
			name:  "error type classification",
			input: "[view(==noMetric)]",
			checkPartialAST: func(t *testing.T, err *vfl.ParseError) {
				// Should be a reference error for undefined metric
				assert.Equal(t, vfl.ReferenceError, err.Type)
				assert.Contains(t, err.Message, "metric")
			},
		},
		{
			name:  "semantic error",
			input: "[view1(==view2)]", // view2 not defined
			checkPartialAST: func(t *testing.T, err *vfl.ParseError) {
				assert.Equal(t, vfl.SemanticError, err.Type)
				assert.Contains(t, err.Message, "view2")
			},
		},
		{
			name:  "syntax error",
			input: "H:|[view1]--==--[view2]|",
			checkPartialAST: func(t *testing.T, err *vfl.ParseError) {
				assert.Equal(t, vfl.SyntaxError, err.Type)
				assert.NotEmpty(t, err.Message)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := vfl.NewParser()
			_, err := parser.Parse(tt.input)

			require.Error(t, err)
			parseErr, ok := err.(*vfl.ParseError)
			require.True(t, ok, "Expected *vfl.ParseError, got %T", err)

			if tt.checkPartialAST != nil {
				tt.checkPartialAST(t, parseErr)
			}
		})
	}
}

func TestParseError_ErrorInterface(t *testing.T) {
	parseErr := &vfl.ParseError{
		Message: "Unexpected token",
		Location: vfl.Location{
			Line:   2,
			Column: 5,
			Offset: 20,
		},
		Type: vfl.SyntaxError,
	}

	// Should implement error interface
	var err error = parseErr
	assert.NotEmpty(t, err.Error())
	assert.Contains(t, err.Error(), "Unexpected token")
	assert.Contains(t, err.Error(), "line 2")
}