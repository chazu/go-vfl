package integration_test

import (
	"testing"

	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicHorizontalLayout(t *testing.T) {
	parser := vfl.NewParser()

	input := "H:|[button(100)]-[textField]-|"
	program, err := parser.Parse(input)

	require.NoError(t, err)
	require.NotNil(t, program)

	// Verify horizontal orientation
	assert.Equal(t, vfl.Horizontal, program.Orientation)

	// Verify structure
	assert.Len(t, program.Statements, 1)
	stmt := program.Statements[0]

	// Check superview constraints
	assert.True(t, stmt.SuperviewStart)
	assert.True(t, stmt.SuperviewEnd)

	// Check views
	assert.Len(t, stmt.Views, 2)
	assert.Equal(t, "button", stmt.Views[0].Name)
	assert.Equal(t, "textField", stmt.Views[1].Name)

	// Check button constraint
	assert.Len(t, stmt.Views[0].Predicates, 1)
	assert.Equal(t, 100.0, stmt.Views[0].Predicates[0].Value.Constant)
}