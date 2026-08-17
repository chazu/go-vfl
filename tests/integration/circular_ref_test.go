package integration_test

import (
	"testing"

	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCircularReferenceDetection(t *testing.T) {
	composer := vfl.NewComposer()
	parser := vfl.NewParser()

	// Create circular reference scenario
	prog1, err := parser.Parse("[view2]")
	require.NoError(t, err)

	prog2, err := parser.Parse("[view1]")
	require.NoError(t, err)

	// Register named views
	err = composer.RegisterNamedView("view1", prog1)
	require.NoError(t, err)

	err = composer.RegisterNamedView("view2", prog2)
	require.NoError(t, err)

	// Attempt composition - should detect circular reference
	_, err = composer.Compose(prog1)
	require.Error(t, err)

	// Verify error type
	circErr, ok := err.(*vfl.CircularReferenceError)
	require.True(t, ok, "Expected CircularReferenceError, got %T", err)

	// Check error details
	assert.NotEmpty(t, circErr.Cycle)
	assert.Contains(t, circErr.Cycle, "view1")
	assert.Contains(t, circErr.Cycle, "view2")
	assert.Contains(t, circErr.Message, "circular")
}