# Contributing to go-vfl

Thank you for your interest in contributing to go-vfl! This document provides guidelines and instructions for contributing to the project.

## Code of Conduct

Please be respectful and constructive in all interactions. We aim to maintain a welcoming and inclusive environment for all contributors.

## How to Contribute

### Reporting Issues

1. Check existing issues to avoid duplicates
2. Use issue templates when available
3. Include:
   - Clear description of the problem
   - Steps to reproduce
   - Expected vs actual behavior
   - VFL input that causes the issue
   - Go version and environment

### Suggesting Features

1. Open an issue with `[Feature Request]` prefix
2. Describe the use case
3. Provide examples of VFL syntax if proposing extensions
4. Explain how it fits with existing functionality

### Submitting Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Make your changes following the coding standards
4. Add tests for new functionality
5. Ensure all tests pass: `go test ./...`
6. Run linters: `golangci-lint run`
7. Commit with descriptive messages
8. Push to your fork
9. Submit a pull request with:
   - Clear description of changes
   - Link to related issues
   - Test results

## Development Setup

```bash
# Clone the repository
git clone https://github.com/chazu/go-vfl.git
cd go-vfl

# Install dependencies
go mod download

# Run tests
go test ./...

# Run benchmarks
go test -bench=. ./tests/benchmarks/...

# Run with race detector
go test -race ./...

# Run linter
golangci-lint run
```

## Project Structure

```
go-vfl/
├── vfl/            # Core VFL parser and composer
│   ├── parser.go   # Main parser implementation
│   ├── grammar.go  # Participle grammar
│   ├── ast.go      # AST definitions
│   └── composer.go # Constraint composition
├── evfl/           # Extended VFL features
│   ├── parser.go   # EVFL parser
│   ├── expression.go # Expression evaluation
│   ├── stack.go    # View stacks
│   └── attribute.go # Attribute references
├── tests/          # Test suites
│   ├── vfl/        # Core VFL tests
│   ├── evfl/       # EVFL tests
│   └── benchmarks/ # Performance benchmarks
└── examples/       # Example applications
```

## Coding Standards

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Keep functions focused and small
- Document exported types and functions
- Use meaningful variable names

### Testing

- Write tests for all new functionality
- Follow TDD approach when possible
- Aim for >80% code coverage
- Include both positive and negative test cases
- Use table-driven tests where appropriate

Example test structure:
```go
func TestParser_Parse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *vfl.Program
        wantErr bool
    }{
        {
            name:  "simple view",
            input: "[button]",
            want:  &vfl.Program{...},
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            parser := vfl.NewParser()
            got, err := parser.Parse(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
            }
            // Assert results
        })
    }
}
```

### Documentation

- Add godoc comments to all exported types and functions
- Include usage examples in documentation
- Update README.md for user-facing changes
- Keep documentation in sync with code

### Commit Messages

Follow conventional commit format:
```
type: brief description

Longer explanation if needed

Fixes #issue-number
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions/changes
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Build/tooling changes

## Areas for Contribution

### High Priority

- Improve EVFL percentage parsing
- Add more comprehensive error messages
- Optimize parser performance
- Add support for more constraint solvers

### Good First Issues

- Add more test cases
- Improve documentation
- Add more examples
- Fix typos and formatting

### Advanced

- Implement new EVFL features
- Add constraint solver backends
- Create GUI layout renderer
- Performance optimizations

## Testing Your Changes

Before submitting:

1. Run all tests: `go test ./...`
2. Check coverage: `go test -cover ./...`
3. Run benchmarks: `go test -bench=. ./tests/benchmarks/...`
4. Test with race detector: `go test -race ./...`
5. Run linter: `golangci-lint run`
6. Test examples: `go run examples/tk/demo.go`

## Getting Help

- Open an issue for questions
- Review existing documentation
- Check the examples directory
- Look at test cases for usage patterns

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be recognized in the project README.

Thank you for contributing to go-vfl!