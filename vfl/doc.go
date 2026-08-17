/*
Package vfl implements a Visual Format Language (VFL) parser and composer for Go.

VFL is a domain-specific language created by Apple for defining Auto Layout constraints
in a compact, visual manner. This package provides a complete implementation of VFL
parsing, validation, and composition with support for extended features.

# Basic Usage

Parse a VFL string and generate layout constraints:

	parser := vfl.NewParser()
	prog, err := parser.Parse("H:|[button(100)]-[textField]-|")
	if err != nil {
		log.Fatal(err)
	}

	composer := vfl.NewComposer()
	layout, err := composer.Compose(prog)
	if err != nil {
		log.Fatal(err)
	}

	// Use the generated constraints
	for _, constraint := range layout.Constraints {
		fmt.Printf("%s.%s %s %f\n",
			constraint.FirstItem.Name,
			constraint.FirstAttribute,
			constraint.Relation,
			constraint.Constant)
	}

# VFL Syntax

The VFL syntax supports:

Orientations:
	H:  // Horizontal orientation (default)
	V:  // Vertical orientation

Views:
	[viewName]                   // Simple view
	[viewName(100)]             // Fixed width/height
	[viewName(>=50)]            // Minimum constraint
	[viewName(<=200)]           // Maximum constraint
	[viewName(==100)]           // Exact constraint
	[viewName(100@750)]         // Priority constraint

Connections:
	-        // Standard spacing (8 points)
	-10-     // Fixed spacing
	->=20-   // Minimum spacing

Superview:
	|        // Superview edge

# Parser Options

Configure parser behavior with options:

	options := vfl.ParserOptions{
		StrictMode: true,           // Strict validation
		Metrics: map[string]float64{
			"margin": 20,
			"padding": 10,
		},
	}
	prog, err := parser.ParseWithOptions(input, options)

# Composition

Compose multiple VFL programs with named views:

	composer := vfl.NewComposer()

	// Register named views
	headerProg, _ := parser.Parse("[header(60)]")
	composer.RegisterNamedView("header", headerProg)

	// Compose with references
	mainProg, _ := parser.Parse("V:|[header]-[content]-[footer]|")
	layout, err := composer.Compose(mainProg)

# Validation

Validate VFL syntax and constraints:

	errors := parser.Validate(vflString)
	for _, err := range errors {
		fmt.Printf("%s at line %d: %s\n",
			err.Severity,
			err.Location.Line,
			err.Message)
	}

# Writers

Export constraints in various formats:

	// JSON output
	jsonWriter := vfl.NewJSONWriter(os.Stdout)
	jsonWriter.Write(prog)

	// Debug output
	debugWriter := vfl.NewDebugWriter(os.Stdout)
	debugWriter.Write(prog)

	// Constraint formats
	constraintWriter := vfl.NewConstraintsWriter(os.Stdout, vfl.AutoLayoutFormat)
	constraintWriter.Write(layout.Constraints)

# Error Handling

Parse errors include partial AST for recovery:

	prog, err := parser.Parse(invalidVFL)
	if err != nil {
		if parseErr, ok := err.(*vfl.ParseError); ok {
			// Access partial AST
			if parseErr.PartialAST != nil {
				// Use partial AST for error recovery
			}
			fmt.Printf("Error at %d:%d: %s\n",
				parseErr.Location.Line,
				parseErr.Location.Column,
				parseErr.Message)
		}
	}

# Visitor Pattern

Traverse and transform AST with the visitor pattern:

	type MyVisitor struct {
		vfl.BaseVisitor
	}

	func (v *MyVisitor) VisitView(view *vfl.View) error {
		fmt.Printf("Found view: %s\n", view.Name)
		return nil
	}

	visitor := &MyVisitor{}
	prog.Accept(visitor)

For extended VFL features including percentages, expressions, and stacks,
see the evfl package.
*/
package vfl