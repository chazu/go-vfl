package evfl

import (
	"fmt"
	"strings"
	"github.com/chazu/go-vfl/vfl"
)

// AttributeParser parses attribute references like "view.width"
type AttributeParser interface {
	Parse(input string) (*AttributeRef, error)
	IsAttributeRef(input string) bool
}

// AttributeRef represents a reference to a view attribute
type AttributeRef struct {
	ViewName  string
	Attribute vfl.Attribute
}

// attributeParser implements AttributeParser
type attributeParser struct{}

// NewAttributeParser creates a new attribute parser
func NewAttributeParser() AttributeParser {
	return &attributeParser{}
}

// Parse parses an attribute reference string
func (p *attributeParser) Parse(input string) (*AttributeRef, error) {
	parts := strings.Split(input, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid attribute reference format: %s", input)
	}

	viewName := parts[0]
	attrName := parts[1]

	if viewName == "" {
		return nil, fmt.Errorf("empty view name in attribute reference")
	}

	attr := p.parseAttribute(attrName)
	if attr == vfl.NotAnAttribute {
		return nil, fmt.Errorf("unknown attribute: %s", attrName)
	}

	return &AttributeRef{
		ViewName:  viewName,
		Attribute: attr,
	}, nil
}

// IsAttributeRef checks if a string is an attribute reference
func (p *attributeParser) IsAttributeRef(input string) bool {
	// Simple check: contains a dot and doesn't start with "metrics."
	return strings.Contains(input, ".") && !strings.HasPrefix(input, "metrics.")
}

// parseAttribute converts attribute name to Attribute enum
func (p *attributeParser) parseAttribute(name string) vfl.Attribute {
	switch strings.ToLower(name) {
	case "width":
		return vfl.Width
	case "height":
		return vfl.Height
	case "left":
		return vfl.Left
	case "right":
		return vfl.Right
	case "top":
		return vfl.Top
	case "bottom":
		return vfl.Bottom
	case "centerx":
		return vfl.CenterX
	case "centery":
		return vfl.CenterY
	case "leading":
		return vfl.Leading
	case "trailing":
		return vfl.Trailing
	case "baseline":
		return vfl.Baseline
	default:
		return vfl.NotAnAttribute
	}
}

// AttributeResolver resolves attribute references to values
type AttributeResolver interface {
	Resolve(ref *AttributeRef, context ResolveContext) (float64, error)
}

// ResolveContext provides context for attribute resolution
type ResolveContext struct {
	Views      map[string]*ViewInfo
	Superview  *ViewInfo
	Metrics    map[string]float64
}

// ViewInfo contains view dimension information
type ViewInfo struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// attributeResolver implements AttributeResolver
type attributeResolver struct{}

// NewAttributeResolver creates a new attribute resolver
func NewAttributeResolver() AttributeResolver {
	return &attributeResolver{}
}

// Resolve resolves an attribute reference to a value
func (r *attributeResolver) Resolve(ref *AttributeRef, context ResolveContext) (float64, error) {
	// Special case for superview
	var view *ViewInfo
	if ref.ViewName == "superview" || ref.ViewName == "|" {
		view = context.Superview
		if view == nil {
			return 0, fmt.Errorf("superview not available in context")
		}
	} else {
		var ok bool
		view, ok = context.Views[ref.ViewName]
		if !ok {
			return 0, fmt.Errorf("view '%s' not found", ref.ViewName)
		}
	}

	// Get attribute value
	switch ref.Attribute {
	case vfl.Width:
		return view.Width, nil
	case vfl.Height:
		return view.Height, nil
	case vfl.Left:
		return view.X, nil
	case vfl.Right:
		return view.X + view.Width, nil
	case vfl.Top:
		return view.Y, nil
	case vfl.Bottom:
		return view.Y + view.Height, nil
	case vfl.CenterX:
		return view.X + view.Width/2, nil
	case vfl.CenterY:
		return view.Y + view.Height/2, nil
	case vfl.Leading:
		// Assuming LTR layout
		return view.X, nil
	case vfl.Trailing:
		// Assuming LTR layout
		return view.X + view.Width, nil
	case vfl.Baseline:
		// Simplified: use bottom for now
		return view.Y + view.Height, nil
	default:
		return 0, fmt.Errorf("unsupported attribute: %v", ref.Attribute)
	}
}