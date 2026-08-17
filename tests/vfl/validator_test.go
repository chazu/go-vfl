package vfl_test

import (
	"testing"

	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
)

func TestParser_Validate(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedErrors int
		checkErrors    func(t *testing.T, errors []vfl.ValidationError)
	}{
		{
			name:           "valid VFL",
			input:          "H:|[button]-[label]|",
			expectedErrors: 0,
		},
		{
			name:           "missing closing bracket",
			input:          "H:|[button",
			expectedErrors: 1,
			checkErrors: func(t *testing.T, errors []vfl.ValidationError) {
				assert.Equal(t, vfl.Error, errors[0].Severity)
				assert.Contains(t, errors[0].Message, "bracket")
			},
		},
		{
			name:           "invalid orientation",
			input:          "X:|[view]|",
			expectedErrors: 1,
			checkErrors: func(t *testing.T, errors []vfl.ValidationError) {
				assert.Contains(t, errors[0].Message, "orientation")
			},
		},
		{
			name:           "multiple errors",
			input:          "H:|[view1(][view2(invalid)]",
			expectedErrors: 2,
			checkErrors: func(t *testing.T, errors []vfl.ValidationError) {
				assert.Len(t, errors, 2)
				for _, err := range errors {
					assert.NotEmpty(t, err.Message)
					assert.NotZero(t, err.Location.Line)
					assert.NotZero(t, err.Location.Column)
				}
			},
		},
		{
			name:           "warning for ambiguous spacing",
			input:          "H:|[view1]--[view2]|", // Double dash might be ambiguous
			expectedErrors: 1,
			checkErrors: func(t *testing.T, errors []vfl.ValidationError) {
				assert.Equal(t, vfl.Warning, errors[0].Severity)
			},
		},
		{
			name:           "invalid priority value",
			input:          "H:|[button(100@1500)]|", // Priority > 1000
			expectedErrors: 1,
			checkErrors: func(t *testing.T, errors []vfl.ValidationError) {
				assert.Contains(t, errors[0].Message, "priority")
			},
		},
		{
			name:           "conflicting constraints",
			input:          "[view(==100,<=50)]",
			expectedErrors: 1,
			checkErrors: func(t *testing.T, errors []vfl.ValidationError) {
				assert.Contains(t, errors[0].Message, "conflict")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := vfl.NewParser()
			errors := parser.Validate(tt.input)

			assert.Len(t, errors, tt.expectedErrors)

			if tt.checkErrors != nil && len(errors) > 0 {
				tt.checkErrors(t, errors)
			}
		})
	}
}