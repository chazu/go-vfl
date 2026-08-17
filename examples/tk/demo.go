// Package main provides a demo of VFL parsing and constraint generation
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chazu/go-vfl/vfl"
	"github.com/chazu/go-vfl/evfl"
)

// Demo application showing VFL parsing without TK dependency
func main() {
	examples := []string{
		"H:|[button(100)]-[textField]-|",
		"V:|[header(50)]-[content]-[footer(30)]|",
		"[view(>=100@750)]",
		"H:|[sidebar(200)]-[main]-[inspector(250)]|",
	}

	fmt.Println("=== VFL Parser Demo ===\n")

	// Standard VFL parsing
	parser := vfl.NewParser()
	for i, example := range examples {
		fmt.Printf("Example %d: %s\n", i+1, example)

		prog, err := parser.Parse(example)
		if err != nil {
			log.Printf("Parse error: %v\n", err)
			continue
		}

		// Compose and generate constraints
		composer := vfl.NewComposer()
		layout, err := composer.Compose(prog)
		if err != nil {
			log.Printf("Compose error: %v\n", err)
			continue
		}

		fmt.Printf("  Generated %d constraints\n", len(layout.Constraints))

		// Write constraints in debug format
		debugWriter := vfl.NewDebugWriter(os.Stdout)
		debugWriter.WriteConstraints(layout.Constraints[:min(3, len(layout.Constraints))])
		fmt.Println()
	}

	// EVFL demonstration
	fmt.Println("\n=== EVFL Features Demo ===\n")

	evflParser := evfl.NewParser(evfl.Options{
		EnablePercentages:   true,
		EnableExpressions:   true,
		EnableStacks:        true,
		EnableAttributeRefs: true,
	})

	evflExamples := []string{
		"[column{header,content,footer}]",
		"[view(100+50)]",
		"[view1(view2.width)]",
	}

	for i, example := range evflExamples {
		fmt.Printf("EVFL Example %d: %s\n", i+1, example)

		_, err := evflParser.Parse(example)
		if err != nil {
			// Try stack parsing for stack notation
			if i == 0 {
				stack, err := evflParser.ParseStack(example)
				if err != nil {
					log.Printf("Stack parse error: %v\n", err)
				} else {
					fmt.Printf("  Parsed as stack: %s\n", stack.Name)
				}
			} else {
				log.Printf("Parse error: %v\n", err)
			}
			continue
		}

		fmt.Printf("  Successfully parsed EVFL\n")
	}

	// Expression evaluation demo
	fmt.Println("\n=== Expression Evaluation Demo ===\n")

	evaluator := evfl.NewExpressionEvaluator()
	expr := &vfl.Expression{
		Operator: vfl.Add,
		Left:     vfl.Value{Type: vfl.ConstantValue, Constant: 100},
		Right:    vfl.Value{Type: vfl.ConstantValue, Constant: 50},
	}

	context := evfl.EvalContext{
		ContainerWidth:  800,
		ContainerHeight: 600,
	}

	result, err := evaluator.Evaluate(expr, context)
	if err != nil {
		log.Printf("Evaluation error: %v\n", err)
	} else {
		fmt.Printf("100 + 50 = %.0f\n", result)
	}

	// Constraint writer demo
	fmt.Println("\n=== Constraint Writer Demo ===\n")

	prog, _ := parser.Parse("H:|[button(100)]-[text]|")
	composer := vfl.NewComposer()
	layout, _ := composer.Compose(prog)

	// Write as Auto Layout format
	fmt.Println("Auto Layout format:")
	autoWriter := vfl.NewConstraintsWriter(os.Stdout, vfl.AutoLayoutFormat)
	autoWriter.Write(layout.Constraints[:min(2, len(layout.Constraints))])

	fmt.Println("\nCassowary format:")
	cassWriter := vfl.NewConstraintsWriter(os.Stdout, vfl.CassowaryFormat)
	cassWriter.Write(layout.Constraints[:min(2, len(layout.Constraints))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}