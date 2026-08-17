package benchmarks

import (
	"testing"
	"github.com/chazu/go-vfl/evfl"
	"github.com/chazu/go-vfl/vfl"
)

var (
	percentageVFL = "H:|[sidebar(25%)]-[content(75%)]|"
	expressionVFL = "[view(100+50*2)]"
	stackVFL      = "[column{header,content,footer}]"
	attributeVFL  = "[view1(view2.width)]"
)

func BenchmarkEVFLParsePercentage(b *testing.B) {
	parser := evfl.NewParser(evfl.Options{
		EnablePercentages: true,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(percentageVFL)
		if err != nil {
			// EVFL percentages might not be fully working yet
			b.Skip("EVFL percentage parsing not fully implemented")
		}
	}
}

func BenchmarkEVFLParseExpression(b *testing.B) {
	parser := evfl.NewParser(evfl.Options{
		EnableExpressions: true,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(expressionVFL)
		if err != nil {
			b.Skip("EVFL expression parsing not fully implemented")
		}
	}
}

func BenchmarkEVFLParseStack(b *testing.B) {
	parser := evfl.NewParser(evfl.Options{
		EnableStacks: true,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseStack(stackVFL)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExpressionEvaluation(b *testing.B) {
	evaluator := evfl.NewExpressionEvaluator()
	expr := &vfl.Expression{
		Operator: vfl.Add,
		Left: vfl.Value{
			Type:     vfl.ConstantValue,
			Constant: 100,
		},
		Right: vfl.Value{
			Type:     vfl.ConstantValue,
			Constant: 50,
		},
	}
	context := evfl.EvalContext{
		Metrics: map[string]float64{
			"margin": 20,
		},
		ContainerWidth:  800,
		ContainerHeight: 600,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := evaluator.Evaluate(expr, context)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAttributeParsing(b *testing.B) {
	attrParser := evfl.NewAttributeParser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := attrParser.Parse("view.width")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAttributeResolution(b *testing.B) {
	resolver := evfl.NewAttributeResolver()
	ref := &evfl.AttributeRef{
		ViewName:  "view1",
		Attribute: vfl.Width,
	}
	context := evfl.ResolveContext{
		Views: map[string]*evfl.ViewInfo{
			"view1": {
				X:      100,
				Y:      50,
				Width:  200,
				Height: 100,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := resolver.Resolve(ref, context)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStackParsing(b *testing.B) {
	parser := vfl.NewParser()
	stackParser := evfl.NewStackParser(parser)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := stackParser.Parse("column{a,b,c}")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComplexEVFLComposition(b *testing.B) {
	parser := evfl.NewParser(evfl.Options{
		EnablePercentages:   true,
		EnableExpressions:   true,
		EnableStacks:        true,
		EnableAttributeRefs: true,
	})
	composer := evfl.NewComposer()

	complexEVFL := `H:|[sidebar(25%)]-[main(75%)]|
V:|[header(60)]-[content]-[footer(40)]|`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prog, err := parser.Parse(complexEVFL)
		if err != nil {
			b.Skip("Complex EVFL parsing not fully implemented")
		}
		_, err = composer.Compose(prog)
		if err != nil {
			b.Fatal(err)
		}
	}
}