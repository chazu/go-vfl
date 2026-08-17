/*
Package evfl provides Extended Visual Format Language (EVFL) features on top of the standard VFL parser.

EVFL extends the standard VFL syntax with powerful features including percentage-based
layouts, mathematical expressions, view stacks, and attribute references. These extensions
make it easier to create responsive and dynamic layouts.

# Basic Usage

Create an EVFL parser with desired features enabled:

	parser := evfl.NewParser(evfl.Options{
		EnablePercentages:   true,
		EnableExpressions:   true,
		EnableStacks:        true,
		EnableAttributeRefs: true,
	})

	prog, err := parser.Parse("H:|[sidebar(25%)]-[content(75%)]|")
	if err != nil {
		log.Fatal(err)
	}

# Percentage Constraints

Define view sizes as percentages of the container:

	// Simple percentage split
	H:|[sidebar(25%)]-[content(75%)]|

	// Three-column layout
	H:|[left(20%)]-[center(60%)]-[right(20%)]|

	// Percentage with minimum
	H:|[nav(20%,>=150)]-[main(80%)]|

# Mathematical Expressions

Use expressions for dynamic calculations:

	// Basic arithmetic
	[view(100+50)]              // Addition
	[view(200-50)]              // Subtraction
	[view(50*2)]                // Multiplication
	[view(200/2)]               // Division

	// Complex expressions
	H:|[margin]-[content(100%-40)]-[margin]|

	// Nested expressions
	[calculated((100+50)*2)]

# Attribute References

Reference other view attributes:

	// Reference another view's dimension
	[view1(view2.width)]
	[view3(view4.height+20)]

	// Available attributes:
	// .width, .height, .left, .right, .top, .bottom
	// .centerX, .centerY, .leading, .trailing

# View Stacks

Create stacks of views with automatic spacing:

	// Column stack (vertical)
	[column{header,content,footer}]

	// Row stack (horizontal)
	[row{left,center,right}]

	// Grid stack
	[grid{a,b,c,d,e,f,g,h,i}]

	// Named stacks
	[sidebar:column{logo,nav,profile}]

# Expression Evaluation

Evaluate expressions with context:

	evaluator := evfl.NewExpressionEvaluator()

	context := evfl.EvalContext{
		ViewWidths: map[string]float64{
			"button": 100,
		},
		Metrics: map[string]float64{
			"margin": 20,
		},
		ContainerWidth:  800,
		ContainerHeight: 600,
	}

	expr := &vfl.Expression{
		Operator: vfl.Add,
		Left:     vfl.Value{Type: vfl.ConstantValue, Constant: 100},
		Right:    vfl.Value{Type: vfl.MetricValue, MetricName: "margin"},
	}

	result, err := evaluator.Evaluate(expr, context)

# Attribute Parsing and Resolution

Parse and resolve view attribute references:

	// Parse attribute reference
	attrParser := evfl.NewAttributeParser()
	ref, err := attrParser.Parse("button.width")

	// Resolve to actual value
	resolver := evfl.NewAttributeResolver()
	context := evfl.ResolveContext{
		Views: map[string]*evfl.ViewInfo{
			"button": {Width: 100, Height: 50},
		},
	}

	value, err := resolver.Resolve(ref, context)

# Stack Processing

Process view stack notation:

	stackParser := evfl.NewStackParser(baseParser)

	// Check if input is stack notation
	if stackParser.IsStackNotation(input) {
		stack, err := stackParser.Parse(input)

		// Generate constraints for stack
		constraints := evfl.CreateStackConstraints(stack)
	}

# Composition with Size

Compose layouts with container size for percentage calculations:

	composer := evfl.NewComposer()

	size := evfl.Size{
		Width:  800,
		Height: 600,
	}

	layout, err := composer.ComposeWithSize(prog, size)

# Integration with Standard VFL

EVFL is fully compatible with standard VFL:

	// EVFL parser can parse standard VFL
	prog, _ := evflParser.Parse("H:|[button(100)]-[text]|")

	// Use with standard VFL composer
	composer := vfl.NewComposer()
	layout, _ := composer.Compose(prog)

# Error Handling

EVFL maintains the same error handling as standard VFL:

	prog, err := parser.Parse(evflString)
	if err != nil {
		if parseErr, ok := err.(*vfl.ParseError); ok {
			// Handle parse error with partial AST
			if parseErr.PartialAST != nil {
				// Use partial AST
			}
		}
	}

# Feature Detection

Check for EVFL features in parsed programs:

	if parser.EnableFeature(evfl.PercentageFeature) == nil {
		// Percentages are supported
	}

For standard VFL features and basic usage, see the vfl package.
*/
package evfl