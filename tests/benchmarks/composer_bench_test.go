package benchmarks

import (
	"testing"
	"github.com/chazu/go-vfl/vfl"
)

func BenchmarkComposeSimple(b *testing.B) {
	parser := vfl.NewParser()
	prog, _ := parser.Parse(simpleVFL)
	composer := vfl.NewComposer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := composer.Compose(prog)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComposeComplex(b *testing.B) {
	parser := vfl.NewParser()
	prog, _ := parser.Parse(complexVFL)
	composer := vfl.NewComposer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := composer.Compose(prog)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComposeMultiplePrograms(b *testing.B) {
	parser := vfl.NewParser()
	prog1, _ := parser.Parse("H:|[a]-[b]|")
	prog2, _ := parser.Parse("V:|[a(50)]-[b]|")
	composer := vfl.NewComposer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := composer.Compose(prog1, prog2)
		if err != nil {
			b.Fatal(err)
		}
		composer.Reset()
	}
}

func BenchmarkRegisterNamedView(b *testing.B) {
	parser := vfl.NewParser()
	prog, _ := parser.Parse("[header(60)]")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		composer := vfl.NewComposer()
		err := composer.RegisterNamedView("header", prog)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComposeWithNamedViews(b *testing.B) {
	parser := vfl.NewParser()
	headerProg, _ := parser.Parse("[header(60)]")
	footerProg, _ := parser.Parse("[footer(40)]")
	mainProg, _ := parser.Parse("V:|[header]-[content]-[footer]|")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		composer := vfl.NewComposer()
		composer.RegisterNamedView("header", headerProg)
		composer.RegisterNamedView("footer", footerProg)
		_, err := composer.Compose(mainProg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComposeWithValidation(b *testing.B) {
	parser := vfl.NewParser()
	prog, _ := parser.Parse(complexVFL)
	composer := vfl.NewComposer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, errors := composer.ComposeWithValidation(prog)
		if len(errors) > 0 {
			b.Fatalf("unexpected validation errors: %v", errors)
		}
	}
}

func BenchmarkCircularReferenceDetection(b *testing.B) {
	parser := vfl.NewParser()
	prog1, _ := parser.Parse("[view1(view2.width)]")
	prog2, _ := parser.Parse("[view2(view1.width)]")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		composer := vfl.NewComposer()
		composer.RegisterNamedView("view1", prog1)
		composer.RegisterNamedView("view2", prog2)
		_, err := composer.Compose(prog1)
		if err == nil {
			b.Fatal("expected circular reference error")
		}
	}
}