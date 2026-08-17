package benchmarks

import (
	"testing"
	"github.com/chazu/go-vfl/vfl"
)

// Benchmark test strings
var (
	simpleVFL   = "[button(100)]"
	mediumVFL   = "H:|[button(100)]-[textField(>=200)]-[label]|"
	complexVFL  = "H:|[sidebar(200@750)]-[content(>=300@500)]-[inspector(250@1000)]|"
	multilineVFL = `H:|[header]-[body]-[footer]|
V:|[header(60)]-[body]-[footer(40)]|`
)

func BenchmarkParseSimple(b *testing.B) {
	parser := vfl.NewParser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(simpleVFL)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseMedium(b *testing.B) {
	parser := vfl.NewParser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(mediumVFL)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseComplex(b *testing.B) {
	parser := vfl.NewParser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(complexVFL)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseMultiline(b *testing.B) {
	parser := vfl.NewParser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(multilineVFL)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseWithOptions(b *testing.B) {
	parser := vfl.NewParser()
	options := vfl.ParserOptions{
		StrictMode: true,
		Metrics: map[string]float64{
			"margin": 20,
			"padding": 10,
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.ParseWithOptions(mediumVFL, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidate(b *testing.B) {
	parser := vfl.NewParser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		errors := parser.Validate(mediumVFL)
		if len(errors) > 0 {
			b.Fatalf("unexpected validation errors: %v", errors)
		}
	}
}

// Benchmark parallel parsing
func BenchmarkParseParallel(b *testing.B) {
	parser := vfl.NewParser()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := parser.Parse(mediumVFL)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}