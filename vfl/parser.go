package vfl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"
)

// Parser is the main interface for parsing VFL strings
type Parser interface {
	Parse(input string) (*Program, error)
	ParseWithOptions(input string, options ParserOptions) (*Program, error)
	Validate(input string) []ValidationError
}

// ParserOptions configures parser behavior
type ParserOptions struct {
	Lookahead          int
	StrictMode         bool
	EnableEVFL         bool
	Metrics            map[string]float64
	GeneratePartialAST bool
	ValidateReferences bool
}

// parser implements the Parser interface
type parser struct {
	options ParserOptions
}

// NewParser creates a new VFL parser
func NewParser() Parser {
	return NewParserWithOptions(ParserOptions{
		Lookahead: 5,
	})
}

// NewParserWithOptions creates a new VFL parser with options
func NewParserWithOptions(options ParserOptions) Parser {
	if options.Lookahead == 0 {
		options.Lookahead = 5
	}

	return &parser{
		options: options,
	}
}

// Parse parses a VFL string and returns an AST or error with partial AST
func (p *parser) Parse(input string) (*Program, error) {
	return p.ParseWithOptions(input, p.options)
}

// ParseWithOptions parses with custom options
func (p *parser) ParseWithOptions(input string, options ParserOptions) (*Program, error) {
	// Store original options temporarily
	oldOptions := p.options
	p.options = options
	defer func() { p.options = oldOptions }()

	// Try to parse with grammar
	grammar, err := parseGrammar(input)

	if err != nil {
		// Create partial AST if not in strict mode
		if !options.StrictMode || options.GeneratePartialAST {
			partial := p.extractPartialAST(input, err)
			parseErr := &ParseError{
				Message:    err.Error(),
				Location:   p.getErrorLocation(err),
				Type:       p.classifyError(err, input),
				Input:      input,
				Partial:    partial,
				PartialAST: partial,
			}
			return partial, parseErr
		}

		return nil, &ParseError{
			Message:  err.Error(),
			Location: p.getErrorLocation(err),
			Type:     p.classifyError(err, input),
			Input:    input,
		}
	}

	// Convert grammar to public AST
	program := convertGrammarToAST(grammar, input)

	// Apply metrics if provided
	if options.Metrics != nil {
		p.applyMetrics(program, options.Metrics)
	}

	// Validate the program
	if options.ValidateReferences {
		if errors := p.validateProgram(program); len(errors) > 0 && options.StrictMode {
			return nil, &ParseError{
				Message:  errors[0].Message,
				Location: errors[0].Location,
				Input:    input,
			}
		}
	}

	return program, nil
}

// Validate checks if a VFL string is syntactically valid
func (p *parser) Validate(input string) []ValidationError {
	prog, err := p.Parse(input)

	var errors []ValidationError

	if err != nil {
		if parseErr, ok := err.(*ParseError); ok {
			errors = append(errors, ValidationError{
				Message:  parseErr.Message,
				Location: parseErr.Location,
				Severity: Error,
			})
		}
	}

	if prog != nil {
		errors = append(errors, p.validateProgram(prog)...)
	}

	return errors
}

// Helper methods

func (p *parser) extractPartialAST(input string, err error) *Program {
	// Try to parse as much as possible
	_ = strings.Split(input, "\n") // lines not used currently
	partial := &Program{
		Input:       input,
		Orientation: Horizontal, // Default
		Statements:  []Statement{},
	}

	// Simple heuristic: try to extract orientation
	if strings.HasPrefix(input, "V:") {
		partial.Orientation = Vertical
	}

	// Try to extract views that were successfully parsed
	// This is a simplified version - a real implementation would be more sophisticated
	if strings.Contains(input, "[") {
		stmt := Statement{}
		viewPattern := `\[([a-zA-Z_][a-zA-Z0-9_]*)`
		matches := findAllMatches(input, viewPattern)
		for _, match := range matches {
			stmt.Views = append(stmt.Views, View{
				Name:     match,
				Location: Location{Line: 1, Column: 1},
			})
		}
		if len(stmt.Views) > 0 {
			partial.Statements = append(partial.Statements, stmt)
		}
	}

	return partial
}

func (p *parser) getErrorLocation(err error) Location {
	// Extract location from participle error if possible
	if parseErr, ok := err.(*participle.ParseError); ok {
		if parseErr.Pos.Line > 0 {
			return Location{
				Line:   parseErr.Pos.Line,
				Column: parseErr.Pos.Column,
				Offset: parseErr.Pos.Offset,
			}
		}
	}
	return Location{Line: 1, Column: 1}
}

func (p *parser) classifyError(err error, input string) ErrorType {
	errMsg := err.Error()

	// Check for reference errors (undefined metrics or views)
	if strings.Contains(errMsg, "metric") || strings.Contains(errMsg, "noMetric") {
		return ReferenceError
	}

	// Check for semantic errors (invalid view references)
	if strings.Contains(errMsg, "view") && strings.Contains(input, "==") {
		return SemanticError
	}

	// Check for syntax errors
	if strings.Contains(errMsg, "expected") || strings.Contains(errMsg, "unexpected") ||
	   strings.Contains(errMsg, "--==--") || strings.Contains(errMsg, "must match") {
		return SyntaxError
	}

	// Default to syntax error
	return SyntaxError
}

func findAllMatches(input, pattern string) []string {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(input, -1)
	result := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			result = append(result, match[1])
		}
	}
	return result
}


func (p *parser) applyMetrics(program *Program, metrics map[string]float64) {
	// Merge standard metrics with provided metrics
	allMetrics := make(map[string]float64)
	for k, v := range StandardMetrics {
		allMetrics[k] = v
	}
	for k, v := range metrics {
		allMetrics[k] = v // User metrics override standard ones
	}

	for i := range program.Statements {
		for j := range program.Statements[i].Views {
			for k := range program.Statements[i].Views[j].Predicates {
				pred := &program.Statements[i].Views[j].Predicates[k]
				if pred.Value.Type == MetricValue {
					if val, ok := allMetrics[pred.Value.MetricName]; ok {
						pred.Value.Type = ConstantValue
						pred.Value.Constant = val
					}
				}
			}
		}

		for k := range program.Statements[i].Connections {
			conn := &program.Statements[i].Connections[k]
			// Apply metrics to connections
			if conn.MetricName != "" {
				if val, ok := allMetrics[conn.MetricName]; ok {
					conn.Spacing = val
				}
			} else if conn.IsDefault {
				conn.Spacing = StandardMetrics["standard"]
			}
		}
	}
}

func (p *parser) validateProgram(program *Program) []ValidationError {
	var errors []ValidationError

	// Validate priorities
	for _, stmt := range program.Statements {
		for _, view := range stmt.Views {
			if view.Priority < 0 || view.Priority > 1000 {
				errors = append(errors, ValidationError{
					Message:  fmt.Sprintf("priority %d out of range (0-1000)", view.Priority),
					Location: view.Location,
					Severity: Error,
				})
			}

			// Check for conflicting constraints
			for i, pred1 := range view.Predicates {
				for j, pred2 := range view.Predicates {
					if i != j && p.constraintsConflict(pred1, pred2) {
						errors = append(errors, ValidationError{
							Message:  "conflicting constraints",
							Location: view.Location,
							Severity: Error,
						})
					}
				}
			}
		}
	}

	return errors
}

func (p *parser) constraintsConflict(p1, p2 Predicate) bool {
	// Simple conflict detection
	if p1.Value.Type == ConstantValue && p2.Value.Type == ConstantValue {
		if p1.Relation == Equal && p2.Relation == Equal {
			return p1.Value.Constant != p2.Value.Constant
		}
		if p1.Relation == Equal && p2.Relation == LessEqual {
			return p1.Value.Constant > p2.Value.Constant
		}
		if p1.Relation == Equal && p2.Relation == GreaterEqual {
			return p1.Value.Constant < p2.Value.Constant
		}
	}
	return false
}

func isValidIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Check first character
	if !((s[0] >= 'a' && s[0] <= 'z') || (s[0] >= 'A' && s[0] <= 'Z') || s[0] == '_') {
		return false
	}
	// Check remaining characters
	for i := 1; i < len(s); i++ {
		if !((s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') ||
			(s[i] >= '0' && s[i] <= '9') || s[i] == '_') {
			return false
		}
	}
	return true
}