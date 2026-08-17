# VFL Parser Library Quick Start Guide

## Installation

```bash
go get github.com/chazu/go-vfl
```

## Basic Usage

### 1. Parse Standard VFL

```go
package main

import (
    "fmt"
    "log"

    "github.com/chazu/go-vfl/vfl"
)

func main() {
    // Create a parser instance
    parser := vfl.NewParser()

    // Parse a VFL string
    input := "H:|[button(100)]-[textField]-|"
    program, err := parser.Parse(input)
    if err != nil {
        // Check for partial AST in error
        if parseErr, ok := err.(*vfl.ParseError); ok && parseErr.PartialAST != nil {
            fmt.Printf("Partial parse succeeded up to: %v\n", parseErr.Location)
            // Use partial AST...
        }
        log.Fatal(err)
    }

    // Use the parsed AST
    fmt.Printf("Parsed %d statements\n", len(program.Statements))
}
```

### 2. Parse Extended VFL

```go
package main

import (
    "fmt"
    "log"

    "github.com/chazu/go-vfl/evfl"
)

func main() {
    // Create EVFL parser with options
    parser := evfl.NewParser(evfl.Options{
        EnablePercentages: true,
        EnableExpressions: true,
    })

    // Parse EVFL with percentages and expressions
    input := "H:|[sidebar(==20%)]-[content(==view1.width*2-10)]|"
    program, err := parser.Parse(input)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Parsed EVFL with %d views\n", len(program.Statements[0].Views))
}
```

### 3. Compose Multiple Programs with Named Views

```go
package main

import (
    "log"

    "github.com/chazu/go-vfl/vfl"
)

func main() {
    parser := vfl.NewParser()
    composer := vfl.NewComposer()

    // Register a named view
    headerProgram, _ := parser.Parse("H:|[logo(50)]-[title]-[menu(100)]|")
    composer.RegisterNamedView("header", headerProgram)

    // Parse main layout referencing the named view
    mainProgram, _ := parser.Parse("V:|[header]-[content]-|")

    // Compose programs resolving named views
    layout, err := composer.Compose(mainProgram)
    if err != nil {
        // Check for circular references
        if circErr, ok := err.(*vfl.CircularReferenceError); ok {
            log.Fatalf("Circular reference detected: %v", circErr.Cycle)
        }
        log.Fatal(err)
    }

    // Use composed layout
    fmt.Printf("Composed layout with %d constraints\n", len(layout.Constraints))
}
```

### 4. UI Integration with modernc.org/tk

```go
package main

import (
    "log"

    "github.com/chazu/go-vfl/examples/tk"
    "github.com/chazu/go-vfl/vfl"
    "modernc.org/tk"
)

func main() {
    // Parse VFL
    parser := vfl.NewParser()
    program, err := parser.Parse("H:|-[button(100)]-[label]-|")
    if err != nil {
        log.Fatal(err)
    }

    // Create TK application
    app := tk.NewApp()

    // Apply VFL constraints to TK widgets
    renderer := tk.NewRenderer()
    widgets := map[string]*tk.Widget{
        "button": app.NewButton("Click Me"),
        "label":  app.NewLabel("Hello VFL"),
    }

    err = renderer.Apply(program, widgets)
    if err != nil {
        log.Fatal(err)
    }

    // Run the UI
    app.Run()
}
```

## Testing Your Integration

### Test Case 1: Basic Horizontal Layout
```go
func TestHorizontalLayout(t *testing.T) {
    parser := vfl.NewParser()
    input := "H:|[view1(100)]-20-[view2(>=50)]|"

    program, err := parser.Parse(input)
    assert.NoError(t, err)
    assert.Equal(t, 2, len(program.Statements[0].Views))
    assert.Equal(t, 100.0, program.Statements[0].Views[0].Predicates[0].Value.Constant)
}
```

### Test Case 2: Circular Reference Detection
```go
func TestCircularReference(t *testing.T) {
    composer := vfl.NewComposer()
    parser := vfl.NewParser()

    // Create circular reference
    prog1, _ := parser.Parse("[view2]")
    prog2, _ := parser.Parse("[view1]")

    composer.RegisterNamedView("view1", prog1)
    composer.RegisterNamedView("view2", prog2)

    _, err := composer.Compose(prog1)
    assert.Error(t, err)
    assert.IsType(t, &vfl.CircularReferenceError{}, err)
}
```

### Test Case 3: Error Recovery with Partial AST
```go
func TestPartialAST(t *testing.T) {
    parser := vfl.NewParser()
    input := "H:|[valid]-[invalid(]" // Syntax error

    _, err := parser.Parse(input)
    assert.Error(t, err)

    parseErr, ok := err.(*vfl.ParseError)
    assert.True(t, ok)
    assert.NotNil(t, parseErr.PartialAST)
    assert.Equal(t, 1, len(parseErr.PartialAST.Statements[0].Views))
}
```

### Test Case 4: EVFL Percentages
```go
func TestPercentageConstraints(t *testing.T) {
    parser := evfl.NewParser(evfl.Options{
        EnablePercentages: true,
    })

    input := "H:|[sidebar(==25%)]-[content(==75%)]|"
    program, err := parser.Parse(input)

    assert.NoError(t, err)
    assert.Equal(t, 25.0, program.Statements[0].Views[0].Predicates[0].Value.Percentage)
    assert.Equal(t, 75.0, program.Statements[0].Views[1].Predicates[0].Value.Percentage)
}
```

## Running the Examples

### Basic Examples
```bash
# Run parser example
go run examples/basic/main.go

# Run TK UI example
go run examples/tk/main.go

# Run all 10 comprehensive examples
go run examples/showcase/main.go
```

### Example Layouts Included

1. **Simple Horizontal**: `H:|[button]-[label]|`
2. **Fixed Spacing**: `H:|−20−[view1]−20−|`
3. **Priority Constraints**: `[button(100@750)]`
4. **Vertical Stack**: `V:|[header(50)]-[content]-[footer(50)]|`
5. **Named Views**: Complex multi-program composition
6. **EVFL Percentages**: `[sidebar(==20%)]-[content(==80%)]`
7. **EVFL Expressions**: `[view2(==view1.width*2-10)]`
8. **View Stacks**: `V:|[column:[item1][item2][item3]]|`
9. **Z-Ordering**: `Z:|[background][content][overlay]|`
10. **Complex Dashboard**: Complete application layout

## Performance Testing

```go
func BenchmarkParser(b *testing.B) {
    parser := vfl.NewParser()
    input := "H:|[v1(100)]-[v2(>=50)]-[v3(==v1)]-[v4]-|"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = parser.Parse(input)
    }
}
```

## Validation and Debugging

```go
// Enable debug output
parser := vfl.NewParser(vfl.Options{
    Debug: true,
})

// Validate without parsing
errors := parser.Validate(input)
for _, err := range errors {
    fmt.Printf("Line %d, Col %d: %s\n",
        err.Location.Line, err.Location.Column, err.Message)
}

// Write AST for debugging
writer := vfl.NewWriter()
writer.WriteDebug(os.Stdout, program)
```

## Next Steps

1. Explore the [API Documentation](https://pkg.go.dev/github.com/chazu/go-vfl)
2. Review [Complete Examples](./examples/)
3. Read the [VFL Specification](https://developer.apple.com/library/archive/documentation/UserExperience/Conceptual/AutolayoutPG/VisualFormatLanguage.html)
4. Learn about [Extended VFL](https://github.com/lume/autolayout#extended-visual-format-language-evfl)
5. Contribute to the [GitHub Repository](https://github.com/chazu/go-vfl)