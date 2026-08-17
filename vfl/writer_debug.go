package vfl

import (
	"fmt"
	"io"
	"strings"
)

// DebugWriter writes human-readable debug output for VFL AST
type DebugWriter struct {
	writer io.Writer
	indent int
}

// NewDebugWriter creates a new debug writer
func NewDebugWriter(w io.Writer) *DebugWriter {
	return &DebugWriter{
		writer: w,
		indent: 0,
	}
}

// Write writes the program in debug format
func (w *DebugWriter) Write(prog *Program) error {
	_, err := fmt.Fprintf(w.writer, "Program:\n")
	if err != nil {
		return err
	}

	w.indent++
	defer func() { w.indent-- }()

	if err := w.writeLine("Input: %q", prog.Input); err != nil {
		return err
	}
	if err := w.writeLine("Orientation: %s", prog.Orientation); err != nil {
		return err
	}
	if err := w.writeLine("Statements: %d", len(prog.Statements)); err != nil {
		return err
	}

	for i, stmt := range prog.Statements {
		if err := w.writeStatement(i, stmt); err != nil {
			return err
		}
	}

	return nil
}

// writeStatement writes a statement in debug format
func (w *DebugWriter) writeStatement(index int, stmt Statement) error {
	if err := w.writeLine("Statement %d:", index); err != nil {
		return err
	}

	w.indent++
	defer func() { w.indent-- }()

	if stmt.SuperviewStart {
		if err := w.writeLine("SuperviewStart: true"); err != nil {
			return err
		}
	}
	if stmt.SuperviewEnd {
		if err := w.writeLine("SuperviewEnd: true"); err != nil {
			return err
		}
	}

	// Write views
	if len(stmt.Views) > 0 {
		if err := w.writeLine("Views:"); err != nil {
			return err
		}
		w.indent++
		for i, view := range stmt.Views {
			if err := w.writeView(i, view); err != nil {
				return err
			}
		}
		w.indent--
	}

	// Write connections
	if len(stmt.Connections) > 0 {
		if err := w.writeLine("Connections:"); err != nil {
			return err
		}
		w.indent++
		for i, conn := range stmt.Connections {
			if err := w.writeConnection(i, conn); err != nil {
				return err
			}
		}
		w.indent--
	}

	return nil
}

// writeView writes a view in debug format
func (w *DebugWriter) writeView(index int, view View) error {
	if err := w.writeLine("View %d: %s", index, view.Name); err != nil {
		return err
	}

	w.indent++
	defer func() { w.indent-- }()

	if view.Priority != 0 {
		if err := w.writeLine("Priority: %d", view.Priority); err != nil {
			return err
		}
	}

	if view.IsNamedView {
		if err := w.writeLine("IsNamedView: true"); err != nil {
			return err
		}
	}

	if view.IsStack {
		if err := w.writeLine("IsStack: true"); err != nil {
			return err
		}
		if err := w.writeLine("StackOrientation: %s", view.StackOrientation); err != nil {
			return err
		}
		if err := w.writeLine("StackLayout: %s", view.StackLayout); err != nil {
			return err
		}
	}

	// Write predicates
	for i, pred := range view.Predicates {
		if err := w.writePredicate(i, pred); err != nil {
			return err
		}
	}

	return nil
}

// writePredicate writes a predicate in debug format
func (w *DebugWriter) writePredicate(index int, pred Predicate) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Predicate %d: ", index))

	if pred.Attribute != "" {
		sb.WriteString(fmt.Sprintf("%s.", pred.Attribute))
	}

	sb.WriteString(pred.Relation.String())

	switch pred.Value.Type {
	case ConstantValue:
		sb.WriteString(fmt.Sprintf(" %.0f", pred.Value.Constant))
	case ViewRefValue:
		sb.WriteString(fmt.Sprintf(" %s", pred.Value.ViewRef))
	case MetricValue:
		sb.WriteString(fmt.Sprintf(" metrics.%s", pred.Value.MetricName))
	case PercentageValue:
		sb.WriteString(fmt.Sprintf(" %.0f%%", pred.Value.Percentage))
	case ExpressionValue:
		sb.WriteString(" (expression)")
	}

	return w.writeLine(sb.String())
}

// writeConnection writes a connection in debug format
func (w *DebugWriter) writeConnection(index int, conn Connection) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Connection %d: ", index))

	if conn.IsDefault {
		sb.WriteString("default")
	} else if conn.IsEqualSpace {
		sb.WriteString("equal-space")
	} else if conn.IsDisconnection {
		sb.WriteString("disconnection")
	} else {
		sb.WriteString(fmt.Sprintf("%.0f", conn.Spacing))
	}

	return w.writeLine(sb.String())
}

// WriteConstraints writes constraints in debug format
func (w *DebugWriter) WriteConstraints(constraints []Constraint) error {
	_, err := fmt.Fprintf(w.writer, "Constraints: %d\n", len(constraints))
	if err != nil {
		return err
	}

	w.indent++
	defer func() { w.indent-- }()

	for i, constraint := range constraints {
		if err := w.writeConstraint(i, constraint); err != nil {
			return err
		}
	}

	return nil
}

// writeConstraint writes a single constraint in debug format
func (w *DebugWriter) writeConstraint(index int, c Constraint) error {
	var sb strings.Builder

	// First item
	sb.WriteString(c.FirstItem.Name)
	sb.WriteString(".")
	sb.WriteString(c.FirstAttribute.String())

	// Relation
	sb.WriteString(" ")
	sb.WriteString(c.Relation.String())
	sb.WriteString(" ")

	// Second item (if any)
	if c.SecondItem != nil {
		sb.WriteString(c.SecondItem.Name)
		sb.WriteString(".")
		sb.WriteString(c.SecondAttribute.String())

		if c.Multiplier != 1.0 {
			sb.WriteString(fmt.Sprintf(" * %.2f", c.Multiplier))
		}

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

	if c.Priority != 1000 {
		sb.WriteString(fmt.Sprintf(" @%d", c.Priority))
	}

	return w.writeLine("Constraint %d: %s", index, sb.String())
}

// writeLine writes an indented line
func (w *DebugWriter) writeLine(format string, args ...interface{}) error {
	indent := strings.Repeat("  ", w.indent)
	_, err := fmt.Fprintf(w.writer, "%s%s\n", indent, fmt.Sprintf(format, args...))
	return err
}