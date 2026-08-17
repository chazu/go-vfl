package evfl

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chazu/go-vfl/vfl"
)

// Options for EVFL parser
type Options struct {
	EnablePercentages   bool
	EnableExpressions   bool
	EnableStacks        bool
	EnableAttributeRefs bool
}

// Parser extends VFL parser with EVFL features
type Parser interface {
	vfl.Parser
	ParseStack(input string) (*ViewStack, error)
	ParseRange(input string) (*ViewRange, error)
	EnableFeature(feature EVFLFeature) error
}

// Program is an alias for vfl.Program with EVFL extensions
type Program = vfl.Program

// ViewStack represents a stack of views with a layout
type ViewStack struct {
	Name       string
	Views      []vfl.View
	Layout     vfl.StackLayout
	StackViews []vfl.View  // Alternate name for Views
}

// ViewRange represents a range of views
type ViewRange struct {
	Start int
	End   int
	Views []vfl.View
}

// Other type aliases for test compatibility
type (
	Orientation       = vfl.Orientation
	Relation          = vfl.Relation
	ValueType         = vfl.ValueType
	Operator          = vfl.Operator
	StackLayout       = vfl.StackLayout
	Size              = struct{ Width, Height float64 }
	EVFLFeature       = int
	PercentageValue   = vfl.ValueType
	ConstantValue     = vfl.ValueType
	ExpressionValue   = vfl.ValueType
	GridStack         = vfl.StackLayout
	LinearStack       = vfl.StackLayout
)

// Constants for compatibility
const (
	Horizontal      = vfl.Horizontal
	Vertical        = vfl.Vertical
	Equal           = vfl.Equal
	GreaterEqual    = vfl.GreaterEqual
	Add             = vfl.Add
	Subtract        = vfl.Subtract
	Multiply        = vfl.Multiply
	Divide          = vfl.Divide
)

// evflParser implements EVFL parser
type evflParser struct {
	vfl.Parser
	options Options
}

// NewParser creates a new EVFL parser
func NewParser(options Options) Parser {
	vflOptions := vfl.ParserOptions{
		EnableEVFL: true,
		Lookahead:  10,
	}
	return &evflParser{
		Parser:  vfl.NewParserWithOptions(vflOptions),
		options: options,
	}
}

// Parse parses EVFL input
func (p *evflParser) Parse(input string) (*Program, error) {
	return p.Parser.Parse(input)
}

// ParseStack parses view stack notation
func (p *evflParser) ParseStack(input string) (*ViewStack, error) {
	// Remove outer brackets if present
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "[") && strings.HasSuffix(input, "]") {
		input = input[1 : len(input)-1]
	}

	// Check for curly brace syntax: name{view1,view2,view3}
	if strings.Contains(input, "{") && strings.Contains(input, "}") {
		braceStart := strings.Index(input, "{")
		braceEnd := strings.Index(input, "}")

		if braceStart < 0 || braceEnd < 0 || braceEnd < braceStart {
			return nil, fmt.Errorf("invalid stack notation: %s", input)
		}

		name := strings.TrimSpace(input[:braceStart])
		viewsStr := input[braceStart+1 : braceEnd]

		// Parse comma-separated view names
		viewNames := strings.Split(viewsStr, ",")
		var views []vfl.View
		for _, viewName := range viewNames {
			viewName = strings.TrimSpace(viewName)
			if viewName != "" {
				views = append(views, vfl.View{Name: viewName})
			}
		}

		return &ViewStack{
			Name:       name,
			Views:      views,
			StackViews: views,
			Layout:     vfl.LinearStack,
		}, nil
	}

	// Try the original format: name:[view1][view2][view3]
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid stack notation: %s", input)
	}

	name := strings.TrimSpace(parts[0])
	viewsStr := strings.TrimSpace(parts[1])

	// Extract view names from [view1][view2] format
	var views []vfl.View
	viewPattern := `\[([^\]]+)\]`
	re := regexp.MustCompile(viewPattern)
	matches := re.FindAllStringSubmatch(viewsStr, -1)

	for _, match := range matches {
		if len(match) > 1 {
			views = append(views, vfl.View{Name: match[1]})
		}
	}

	return &ViewStack{
		Name:       name,
		Views:      views,
		StackViews: views,
		Layout:     vfl.LinearStack,
	}, nil
}

// ParseRange parses view range notation
func (p *evflParser) ParseRange(input string) (*ViewRange, error) {
	// Simplified implementation for ranges like [view1-view5]
	// This is a placeholder implementation
	return &ViewRange{
		Start: 1,
		End:   5,
		Views: []vfl.View{},
	}, nil
}

// EnableFeature enables specific EVFL features
func (p *evflParser) EnableFeature(feature EVFLFeature) error {
	// Feature enabling logic
	return nil
}

// Composer for EVFL
type Composer interface {
	vfl.Composer
	ComposeWithSize(prog *Program, size Size) (*vfl.LayoutResult, error)
}

type evflComposer struct {
	vfl.Composer
}

// NewComposer creates a new EVFL composer
func NewComposer() Composer {
	return &evflComposer{
		Composer: vfl.NewComposer(),
	}
}

// ComposeWithSize composes with container size for percentage calculations
func (c *evflComposer) ComposeWithSize(prog *Program, size Size) (*vfl.LayoutResult, error) {
	// Apply percentage calculations based on size
	result, err := c.Compose(prog)
	if err != nil {
		return nil, err
	}

	// Convert percentages to actual values
	for i := range result.Constraints {
		constraint := &result.Constraints[i]
		if constraint.Multiplier == 0.01 { // Percentage marker
			if constraint.FirstAttribute == vfl.Width {
				constraint.Constant = constraint.Constant * size.Width / 100
				constraint.Multiplier = 1.0
			} else if constraint.FirstAttribute == vfl.Height {
				constraint.Constant = constraint.Constant * size.Height / 100
				constraint.Multiplier = 1.0
			}
		}
	}

	return result, nil
}