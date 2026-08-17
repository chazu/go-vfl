package vfl

import (
	"errors"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseError(t *testing.T) {
	tests := []struct {
		name        string
		err         *ParseError
		expectedMsg string
		hasPartial  bool
	}{
		{
			name: "syntax error with location",
			err: &ParseError{
				Message:   "unexpected token",
				Location:  Location{Line: 2, Column: 5, Offset: 15},
				ErrorType: SyntaxError,
			},
			expectedMsg: "syntax error at 2:5: unexpected token",
			hasPartial:  false,
		},
		{
			name: "semantic error with partial AST",
			err: &ParseError{
				Message:   "undefined view reference",
				Location:  Location{Line: 3, Column: 10, Offset: 30},
				ErrorType: SemanticError,
				PartialAST: &Program{
					Orientation: Horizontal,
					Statements: []Statement{
						{Views: []View{{Name: "partial"}}},
					},
				},
			},
			expectedMsg: "semantic error at 3:10: undefined view reference",
			hasPartial:  true,
		},
		{
			name: "reference error",
			err: &ParseError{
				Message:   "circular reference detected",
				Location:  Location{Line: 1, Column: 1, Offset: 0},
				ErrorType: ReferenceError,
			},
			expectedMsg: "reference error at 1:1: circular reference detected",
			hasPartial:  false,
		},
		{
			name: "error with expected and actual tokens",
			err: &ParseError{
				Message:        "unexpected token",
				Location:       Location{Line: 4, Column: 8, Offset: 40},
				ErrorType:      SyntaxError,
				ExpectedTokens: []string{"[", "("},
				ActualToken:    "}",
			},
			expectedMsg: "syntax error at 4:8: unexpected token (expected: [, (; found: })",
			hasPartial:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedMsg, tt.err.Error())
			assert.Equal(t, tt.hasPartial, tt.err.PartialAST != nil)
		})
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Message:  "invalid constraint",
		Location: Location{Line: 5, Column: 12, Offset: 50},
		Severity: ErrorSeverity,
	}

	assert.Equal(t, "validation error at 5:12: invalid constraint", err.Error())
	assert.Equal(t, ErrorSeverity, err.Severity)
}

func TestCompositionError(t *testing.T) {
	tests := []struct {
		name        string
		err         *CompositionError
		expectedMsg string
	}{
		{
			name: "missing named view",
			err: &CompositionError{
				Message: "named view not found: header",
				Type:    MissingNamedView,
			},
			expectedMsg: "composition error (missing named view): named view not found: header",
		},
		{
			name: "circular reference with cycle",
			err: &CompositionError{
				Message: "circular reference detected",
				Type:    CircularReference,
				Cycle:   []string{"viewA", "viewB", "viewC", "viewA"},
			},
			expectedMsg: "composition error (circular reference): circular reference detected (cycle: viewA -> viewB -> viewC -> viewA)",
		},
		{
			name: "conflicting constraints",
			err: &CompositionError{
				Message: "conflicting width constraints",
				Type:    ConflictingConstraints,
			},
			expectedMsg: "composition error (conflicting constraints): conflicting width constraints",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedMsg, tt.err.Error())
		})
	}
}

func TestCircularReferenceError(t *testing.T) {
	err := &CircularReferenceError{
		Cycle: []string{"view1", "view2", "view3", "view1"},
	}

	expectedMsg := "circular reference detected: view1 -> view2 -> view3 -> view1"
	assert.Equal(t, expectedMsg, err.Error())
}

func TestMissingNamedViewError(t *testing.T) {
	err := &MissingNamedViewError{
		Name: "headerView",
	}

	assert.Equal(t, "missing named view: headerView", err.Error())
}

func TestErrorHelpers(t *testing.T) {
	t.Run("IsParseError", func(t *testing.T) {
		parseErr := &ParseError{Message: "test"}
		normalErr := errors.New("normal error")
		wrappedErr := WrapError(normalErr, Location{Line: 1, Column: 1})

		assert.True(t, IsParseError(parseErr))
		assert.False(t, IsParseError(normalErr))
		assert.True(t, IsParseError(wrappedErr))
	})

	t.Run("IsValidationError", func(t *testing.T) {
		validErr := &ValidationError{Message: "test"}
		normalErr := errors.New("normal error")

		assert.True(t, IsValidationError(validErr))
		assert.False(t, IsValidationError(normalErr))
	})

	t.Run("IsCompositionError", func(t *testing.T) {
		compErr := &CompositionError{Message: "test"}
		normalErr := errors.New("normal error")

		assert.True(t, IsCompositionError(compErr))
		assert.False(t, IsCompositionError(normalErr))
	})

	t.Run("GetErrorLocation", func(t *testing.T) {
		loc := Location{Line: 10, Column: 20, Offset: 100}
		parseErr := &ParseError{Location: loc}
		normalErr := errors.New("normal error")

		parseLoc, hasLoc := GetErrorLocation(parseErr)
		assert.True(t, hasLoc)
		assert.Equal(t, loc, parseLoc)

		_, hasLoc = GetErrorLocation(normalErr)
		assert.False(t, hasLoc)
	})
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	loc := Location{Line: 7, Column: 15, Offset: 75}

	wrappedErr := WrapError(originalErr, loc)
	parseErr, ok := wrappedErr.(*ParseError)

	assert.True(t, ok)
	assert.Equal(t, "original error", parseErr.Message)
	assert.Equal(t, loc, parseErr.Location)
	assert.Equal(t, SyntaxError, parseErr.ErrorType)
}

func TestNewParseError(t *testing.T) {
	err := NewParseError("test error", Location{Line: 1, Column: 1}, SyntaxError)

	assert.Equal(t, "test error", err.Message)
	assert.Equal(t, Location{Line: 1, Column: 1}, err.Location)
	assert.Equal(t, SyntaxError, err.ErrorType)
	assert.Nil(t, err.PartialAST)
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("validation failed", Location{Line: 2, Column: 3}, WarningSeverity)

	assert.Equal(t, "validation failed", err.Message)
	assert.Equal(t, Location{Line: 2, Column: 3}, err.Location)
	assert.Equal(t, WarningSeverity, err.Severity)
}