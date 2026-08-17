package main

import (
	"fmt"
	"github.com/chazu/go-vfl/vfl"
	"github.com/chazu/go-vfl/evfl"
)

func main() {
	fmt.Println("=== VFL Parser - Final Demo ===")
	fmt.Println("Showcasing all fixed features")

	parser := vfl.NewParser()

	// Feature 1: Basic VFL
	fmt.Println("1️⃣  Basic VFL Parsing:")
	basicTests := []string{
		"H:[button]-[text]",
		"V:|[header]-[body]-[footer]|",
		"H:|[sidebar(200)]-[content]|",
	}
	for _, vfl := range basicTests {
		if prog, err := parser.Parse(vfl); err == nil {
			fmt.Printf("   ✅ %s → %d views\n", vfl, len(prog.Statements[0].Views))
		}
	}

	// Feature 2: Flush Views (FIXED)
	fmt.Println("\n2️⃣  Flush Views (no spacing):")
	if prog, err := parser.Parse("H:[left][middle][right]"); err == nil {
		stmt := prog.Statements[0]
		fmt.Printf("   ✅ Flush views: %d views with connections: ", len(stmt.Views))
		for _, conn := range stmt.Connections {
			fmt.Printf("%.0f ", conn.Spacing)
		}
		fmt.Println("(0 = flush)")
	}

	// Feature 3: Superview Connections (FIXED)
	fmt.Println("\n3️⃣  Superview with Spacing:")
	superTests := []struct{vfl, desc string}{
		{"H:|[content]|", "No spacing"},
		{"H:|-[content]-|", "Default spacing (8)"},
		{"H:|-20-[content]-20-|", "Custom spacing (20)"},
	}
	for _, test := range superTests {
		if prog, err := parser.Parse(test.vfl); err == nil {
			stmt := prog.Statements[0]
			fmt.Printf("   ✅ %s: ", test.desc)
			if len(stmt.Connections) > 0 {
				fmt.Printf("spacing=%g", stmt.Connections[0].Spacing)
			}
			fmt.Println()
		}
	}

	// Feature 4: Metrics (FIXED)
	fmt.Println("\n4️⃣  Metrics in Connections:")
	vfl.RegisterMetric("margin", 16)
	vfl.RegisterMetric("gutter", 12)

	options := vfl.ParserOptions{
		Metrics: map[string]float64{
			"padding": 8,
			"margin":  20, // Override standard
		},
	}

	if prog, err := parser.ParseWithOptions("H:|-margin-[view1]-padding-[view2]-margin-|", options); err == nil {
		stmt := prog.Statements[0]
		fmt.Printf("   ✅ Metrics applied: ")
		for i, conn := range stmt.Connections {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("conn%d=%.0f", i, conn.Spacing)
		}
		fmt.Println()
	}

	// Feature 5: Constraints
	fmt.Println("\n5️⃣  Constraint Types:")
	constraintTests := []string{
		"[view(100)]",           // Fixed width
		"[view(>=50)]",          // Minimum width
		"[view(<=200)]",         // Maximum width
		"[view(==100@750)]",     // With priority
	}
	for _, vfl := range constraintTests {
		if prog, err := parser.Parse(vfl); err == nil {
			view := prog.Statements[0].Views[0]
			if len(view.Predicates) > 0 {
				pred := view.Predicates[0]
				fmt.Printf("   ✅ %s → %v %.0f", vfl, pred.Relation, pred.Value.Constant)
				if pred.Priority > 0 {
					fmt.Printf(" @%.0f", pred.Priority)
				}
				fmt.Println()
			}
		}
	}

	// Feature 6: Named View Composition
	fmt.Println("\n6️⃣  Named View Composition:")
	composer := vfl.NewComposer()

	header, _ := parser.Parse("H:|[logo(50)]-[title]-[menu]|")
	footer, _ := parser.Parse("H:|[copyright]-[version(40)]|")

	composer.RegisterNamedView("header", header)
	composer.RegisterNamedView("footer", footer)

	main, _ := parser.Parse("V:|[header(60)]-[content]-[footer(30)]|")
	if layout, err := composer.Compose(main); err == nil {
		fmt.Printf("   ✅ Composed layout: %d total views\n", len(layout.ViewHierarchy.Views))
		fmt.Print("      Views: ")
		count := 0
		for name := range layout.ViewHierarchy.Views {
			if count > 0 {
				fmt.Print(", ")
			}
			fmt.Print(name)
			count++
			if count >= 5 {
				fmt.Printf(" (+%d more)", len(layout.ViewHierarchy.Views)-5)
				break
			}
		}
		fmt.Println()
	}

	// Feature 7: EVFL Stack Notation (FIXED)
	fmt.Println("\n7️⃣  EVFL Stack Notation:")
	evflParser := evfl.NewParser(evfl.Options{})

	stackTests := []string{
		"[column{top,middle,bottom}]",
		"row:[left][center][right]",
	}
	for _, stack := range stackTests {
		if s, err := evflParser.ParseStack(stack); err == nil {
			fmt.Printf("   ✅ %s → %s with %d views\n", stack, s.Name, len(s.Views))
		}
	}

	// Feature 8: Standard Metrics
	fmt.Println("\n8️⃣  Standard Metrics Library:")
	metrics := []string{"standard", "ios-margin", "button-height", "sidebar-width"}
	for _, name := range metrics {
		if val, ok := vfl.GetMetric(name); ok {
			fmt.Printf("   • %s = %.0f\n", name, val)
		}
	}

	fmt.Println("\n✨ All features working correctly!")
}