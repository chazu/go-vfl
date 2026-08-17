package vfl

import (
	"fmt"
	"io"
	"strings"
)

// ConstraintsWriter writes constraints in various formats
type ConstraintsWriter struct {
	writer io.Writer
	format ConstraintFormat
}

// ConstraintFormat specifies the output format for constraints
type ConstraintFormat int

const (
	// AutoLayoutFormat outputs iOS Auto Layout format
	AutoLayoutFormat ConstraintFormat = iota
	// CassowaryFormat outputs Cassowary constraint format
	CassowaryFormat
	// KiwiFormat outputs Kiwi constraint solver format
	KiwiFormat
)

// NewConstraintsWriter creates a new constraints writer
func NewConstraintsWriter(w io.Writer, format ConstraintFormat) *ConstraintsWriter {
	return &ConstraintsWriter{
		writer: w,
		format: format,
	}
}

// Write writes constraints in the specified format
func (w *ConstraintsWriter) Write(constraints []Constraint) error {
	switch w.format {
	case AutoLayoutFormat:
		return w.writeAutoLayout(constraints)
	case CassowaryFormat:
		return w.writeCassowary(constraints)
	case KiwiFormat:
		return w.writeKiwi(constraints)
	default:
		return fmt.Errorf("unsupported constraint format: %v", w.format)
	}
}

// writeAutoLayout writes constraints in iOS Auto Layout format
func (w *ConstraintsWriter) writeAutoLayout(constraints []Constraint) error {
	for _, c := range constraints {
		line := w.formatAutoLayoutConstraint(c)
		if _, err := fmt.Fprintln(w.writer, line); err != nil {
			return err
		}
	}
	return nil
}

// formatAutoLayoutConstraint formats a single constraint for Auto Layout
func (w *ConstraintsWriter) formatAutoLayoutConstraint(c Constraint) string {
	var sb strings.Builder

	// NSLayoutConstraint format
	sb.WriteString("NSLayoutConstraint(item: ")
	sb.WriteString(c.FirstItem.Name)
	sb.WriteString(", attribute: .")
	sb.WriteString(w.formatAutoLayoutAttribute(c.FirstAttribute))
	sb.WriteString(", relatedBy: .")
	sb.WriteString(w.formatAutoLayoutRelation(c.Relation))
	sb.WriteString(", toItem: ")

	if c.SecondItem != nil {
		sb.WriteString(c.SecondItem.Name)
		sb.WriteString(", attribute: .")
		sb.WriteString(w.formatAutoLayoutAttribute(c.SecondAttribute))
		sb.WriteString(", multiplier: ")
		sb.WriteString(fmt.Sprintf("%.2f", c.Multiplier))
		sb.WriteString(", constant: ")
		sb.WriteString(fmt.Sprintf("%.0f", c.Constant))
	} else {
		sb.WriteString("nil, attribute: .notAnAttribute")
		sb.WriteString(", multiplier: 1.0")
		sb.WriteString(", constant: ")
		sb.WriteString(fmt.Sprintf("%.0f", c.Constant))
	}

	sb.WriteString(")")

	if c.Priority != 1000 {
		sb.WriteString(fmt.Sprintf(".priority(%d)", c.Priority))
	}

	return sb.String()
}

// formatAutoLayoutAttribute converts attribute to Auto Layout format
func (w *ConstraintsWriter) formatAutoLayoutAttribute(attr Attribute) string {
	switch attr {
	case Width:
		return "width"
	case Height:
		return "height"
	case Left:
		return "left"
	case Right:
		return "right"
	case Top:
		return "top"
	case Bottom:
		return "bottom"
	case Leading:
		return "leading"
	case Trailing:
		return "trailing"
	case CenterX:
		return "centerX"
	case CenterY:
		return "centerY"
	case Baseline:
		return "baseline"
	default:
		return "notAnAttribute"
	}
}

// formatAutoLayoutRelation converts relation to Auto Layout format
func (w *ConstraintsWriter) formatAutoLayoutRelation(rel Relation) string {
	switch rel {
	case Equal:
		return "equal"
	case GreaterEqual:
		return "greaterThanOrEqual"
	case LessEqual:
		return "lessThanOrEqual"
	default:
		return "equal"
	}
}

// writeCassowary writes constraints in Cassowary format
func (w *ConstraintsWriter) writeCassowary(constraints []Constraint) error {
	for _, c := range constraints {
		line := w.formatCassowaryConstraint(c)
		if _, err := fmt.Fprintln(w.writer, line); err != nil {
			return err
		}
	}
	return nil
}

// formatCassowaryConstraint formats a constraint for Cassowary
func (w *ConstraintsWriter) formatCassowaryConstraint(c Constraint) string {
	var sb strings.Builder

	// Cassowary linear constraint format: var == expression
	sb.WriteString(c.FirstItem.Name)
	sb.WriteString(".")
	sb.WriteString(strings.ToLower(c.FirstAttribute.String()))

	switch c.Relation {
	case Equal:
		sb.WriteString(" == ")
	case GreaterEqual:
		sb.WriteString(" >= ")
	case LessEqual:
		sb.WriteString(" <= ")
	}

	if c.SecondItem != nil {
		if c.Multiplier != 1.0 {
			sb.WriteString(fmt.Sprintf("%.2f * ", c.Multiplier))
		}
		sb.WriteString(c.SecondItem.Name)
		sb.WriteString(".")
		sb.WriteString(strings.ToLower(c.SecondAttribute.String()))
		if c.Constant != 0 {
			if c.Constant > 0 {
				sb.WriteString(fmt.Sprintf(" + %.0f", c.Constant))
			} else {
				sb.WriteString(fmt.Sprintf(" - %.0f", -c.Constant))
			}
		}
	} else {
		sb.WriteString(fmt.Sprintf("%.0f", c.Constant))
	}

	// Add strength for Cassowary
	if c.Priority < 1000 {
		strength := "weak"
		if c.Priority >= 750 {
			strength = "strong"
		} else if c.Priority >= 500 {
			strength = "medium"
		}
		sb.WriteString(fmt.Sprintf(" :%s", strength))
	}

	return sb.String()
}

// writeKiwi writes constraints in Kiwi format
func (w *ConstraintsWriter) writeKiwi(constraints []Constraint) error {
	for _, c := range constraints {
		line := w.formatKiwiConstraint(c)
		if _, err := fmt.Fprintln(w.writer, line); err != nil {
			return err
		}
	}
	return nil
}

// formatKiwiConstraint formats a constraint for Kiwi solver
func (w *ConstraintsWriter) formatKiwiConstraint(c Constraint) string {
	// Kiwi format is similar to Cassowary
	return w.formatCassowaryConstraint(c)
}

// WriteProgram writes constraints from a parsed program
func (w *ConstraintsWriter) WriteProgram(prog *Program) error {
	composer := NewComposer()
	result, err := composer.Compose(prog)
	if err != nil {
		return err
	}
	return w.Write(result.Constraints)
}