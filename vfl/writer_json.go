package vfl

import (
	"encoding/json"
	"io"
)

// JSONWriter writes VFL AST to JSON format
type JSONWriter struct {
	writer io.Writer
	indent string
}

// NewJSONWriter creates a new JSON writer
func NewJSONWriter(w io.Writer) *JSONWriter {
	return &JSONWriter{
		writer: w,
		indent: "  ",
	}
}

// Write writes the program as JSON
func (w *JSONWriter) Write(prog *Program) error {
	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", w.indent)
	return encoder.Encode(prog)
}

// WriteCompact writes compact JSON without indentation
func (w *JSONWriter) WriteCompact(prog *Program) error {
	encoder := json.NewEncoder(w.writer)
	return encoder.Encode(prog)
}

// WriteConstraints writes constraints as JSON
func (w *JSONWriter) WriteConstraints(constraints []Constraint) error {
	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", w.indent)
	return encoder.Encode(constraints)
}

// WriteLayoutResult writes a layout result as JSON
func (w *JSONWriter) WriteLayoutResult(result *LayoutResult) error {
	encoder := json.NewEncoder(w.writer)
	encoder.SetIndent("", w.indent)
	return encoder.Encode(result)
}