package vfl_test

import (
	"fmt"
	"log"

	"github.com/chazu/go-vfl/vfl"
)

// ExampleParser_Parse demonstrates basic VFL parsing.
func ExampleParser_Parse() {
	parser := vfl.NewParser()

	// Parse a simple horizontal layout
	program, err := parser.Parse("H:|[button(100)]-[textField]-|")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Orientation: %d\n", program.Orientation)
	fmt.Printf("Number of views: %d\n", len(program.Statements[0].Views))
	fmt.Printf("First view: %s\n", program.Statements[0].Views[0].Name)

	// Output:
	// Orientation: 0
	// Number of views: 2
	// First view: button
}

// ExampleParser_ParseWithOptions demonstrates parsing with custom metrics.
func ExampleParser_ParseWithOptions() {
	parser := vfl.NewParser()

	options := vfl.ParserOptions{
		Metrics: map[string]float64{
			"margin": 20,
			"padding": 10,
		},
		ValidateReferences: true,
	}

	// Parse with custom metrics
	program, err := parser.ParseWithOptions("H:|-margin-[view]-margin-|", options)
	if err != nil {
		log.Fatal(err)
	}

	// The metric "margin" will be replaced with 20
	conn := program.Statements[0].Connections[0]
	fmt.Printf("Connection spacing: %.0f\n", conn.Spacing)

	// Output:
	// Connection spacing: 20
}

// ExampleComposer_Compose demonstrates composing multiple VFL programs with named views.
func ExampleComposer_Compose() {
	parser := vfl.NewParser()
	composer := vfl.NewComposer()

	// Register a named view for header
	headerProgram, _ := parser.Parse("H:|[logo(50)]-[title]-[menu(100)]|")
	composer.RegisterNamedView("header", headerProgram)

	// Register a named view for footer
	footerProgram, _ := parser.Parse("H:|[copyright]-[links]|")
	composer.RegisterNamedView("footer", footerProgram)

	// Compose a layout using the named views
	mainProgram, _ := parser.Parse("V:|[header]-[content]-[footer]|")
	layout, err := composer.Compose(mainProgram)
	if err != nil {
		log.Fatal(err)
	}

	// The composed layout will include all views from header and footer
	fmt.Printf("Total views in hierarchy: %d\n", len(layout.ViewHierarchy.Views))

	// Check if specific views are present
	_, hasLogo := layout.ViewHierarchy.Views["logo"]
	_, hasTitle := layout.ViewHierarchy.Views["title"]
	fmt.Printf("Has logo: %v, Has title: %v\n", hasLogo, hasTitle)

	// Output:
	// Total views in hierarchy: 6
	// Has logo: true, Has title: true
}

// ExampleComposer_RegisterNamedView shows how to register reusable view components.
func ExampleComposer_RegisterNamedView() {
	parser := vfl.NewParser()
	composer := vfl.NewComposer()

	// Create a reusable button group
	buttonGroup, err := parser.Parse("H:[cancel(80)]-[ok(80)]")
	if err != nil {
		log.Fatal(err)
	}

	// Register it as a named view
	err = composer.RegisterNamedView("buttonGroup", buttonGroup)
	if err != nil {
		log.Fatal(err)
	}

	// Check if it was registered
	fmt.Printf("Has buttonGroup: %v\n", composer.HasNamedView("buttonGroup"))

	// Try to register again (will fail)
	err = composer.RegisterNamedView("buttonGroup", buttonGroup)
	fmt.Printf("Duplicate registration error: %v\n", err != nil)

	// Output:
	// Has buttonGroup: true
	// Duplicate registration error: true
}

// ExampleParseError demonstrates error handling with partial AST recovery.
func ExampleParseError() {
	parser := vfl.NewParser()

	// Parse invalid VFL with error recovery
	program, err := parser.Parse("H:|[view1]-[view2(invalid]")

	if err != nil {
		parseErr, ok := err.(*vfl.ParseError)
		if ok {
			fmt.Printf("Error type: %d\n", parseErr.Type)
			fmt.Printf("Has partial AST: %v\n", parseErr.PartialAST != nil)

			// Can still access partial results
			if parseErr.PartialAST != nil && len(parseErr.PartialAST.Statements) > 0 {
				fmt.Printf("Parsed views before error: %d\n",
					len(parseErr.PartialAST.Statements[0].Views))
			}
		}
	}

	_ = program // May be nil or partial

	// Output:
	// Error type: 0
	// Has partial AST: true
	// Parsed views before error: 2
}