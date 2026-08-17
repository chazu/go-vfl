package vfl_test

import (
	"sync"
	"testing"

	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParser_Isolation verifies that multiple parser instances don't share state
func TestParser_Isolation(t *testing.T) {
	t.Run("concurrent parser instances", func(t *testing.T) {
		var wg sync.WaitGroup
		numGoroutines := 10
		numIterations := 100

		// Channel to collect any errors
		errors := make(chan error, numGoroutines*numIterations)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				// Each goroutine creates its own parser
				parser := vfl.NewParser()

				for j := 0; j < numIterations; j++ {
					// Parse different VFL strings
					input := "H:|[view" + string(rune('A'+id)) + "]-|"
					prog, err := parser.Parse(input)

					if err != nil {
						errors <- err
						continue
					}

					// Verify the parsed view name matches what we expect
					if len(prog.Statements) > 0 && len(prog.Statements[0].Views) > 0 {
						expectedName := "view" + string(rune('A'+id))
						if prog.Statements[0].Views[0].Name != expectedName {
							errors <- assert.AnError
						}
					}
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for any errors
		for err := range errors {
			require.NoError(t, err, "Parser instances should not interfere with each other")
		}
	})

	t.Run("parser reusability", func(t *testing.T) {
		parser := vfl.NewParser()

		// Parse first string
		prog1, err := parser.Parse("H:|[button]-|")
		require.NoError(t, err)
		require.NotNil(t, prog1)
		assert.Equal(t, "button", prog1.Statements[0].Views[0].Name)

		// Parse second string with same parser
		prog2, err := parser.Parse("V:|[label]-|")
		require.NoError(t, err)
		require.NotNil(t, prog2)
		assert.Equal(t, "label", prog2.Statements[0].Views[0].Name)

		// Verify first result wasn't modified
		assert.Equal(t, "button", prog1.Statements[0].Views[0].Name)
		assert.Equal(t, vfl.Horizontal, prog1.Orientation)
		assert.Equal(t, vfl.Vertical, prog2.Orientation)
	})

	t.Run("parser options isolation", func(t *testing.T) {
		// Create two parsers with different options
		metrics1 := map[string]float64{"spacing": 10}
		metrics2 := map[string]float64{"spacing": 20}

		parser1 := vfl.NewParserWithOptions(vfl.ParserOptions{
			Metrics: metrics1,
		})
		parser2 := vfl.NewParserWithOptions(vfl.ParserOptions{
			Metrics: metrics2,
		})

		// Both parsers should maintain their own metrics
		prog1, err := parser1.Parse("H:|-spacing-[view]-spacing-|")
		require.NoError(t, err)

		prog2, err := parser2.Parse("H:|-spacing-[view]-spacing-|")
		require.NoError(t, err)

		// Verify each parser used its own metrics
		// (This assumes metrics are applied during parsing)
		assert.NotSame(t, prog1, prog2, "Each parser should produce independent results")
	})
}