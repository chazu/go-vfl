package evfl_test

import (
	"testing"

	"github.com/chazu/go-vfl/evfl"
	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEVFLParser_AttributeReferences(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		check   func(t *testing.T, prog *evfl.Program)
		wantErr bool
	}{
		{
			name:  "width attribute reference",
			input: "[view2(view1.width)]",
			check: func(t *testing.T, prog *evfl.Program) {
				view := prog.Statements[0].Views[0]
				pred := view.Predicates[0]

				assert.Equal(t, evfl.Equal, pred.Relation)
				assert.Equal(t, "view1", pred.Value.ViewRef)
				assert.Equal(t, "width", pred.Attribute)
			},
		},
		{
			name:  "height attribute reference",
			input: "[sidebar(==content.height)]",
			check: func(t *testing.T, prog *evfl.Program) {
				view := prog.Statements[0].Views[0]
				assert.Equal(t, "sidebar", view.Name)

				pred := view.Predicates[0]
				assert.Equal(t, "content", pred.Value.ViewRef)
				assert.Equal(t, "height", pred.Attribute)
			},
		},
		{
			name:  "centerX attribute",
			input: "[button(centerView.centerX)]",
			check: func(t *testing.T, prog *evfl.Program) {
				pred := prog.Statements[0].Views[0].Predicates[0]
				assert.Equal(t, "centerView", pred.Value.ViewRef)
				assert.Equal(t, "centerX", pred.Attribute)
			},
		},
		{
			name:  "centerY attribute",
			input: "[button(centerView.centerY)]",
			check: func(t *testing.T, prog *evfl.Program) {
				pred := prog.Statements[0].Views[0].Predicates[0]
				assert.Equal(t, "centerView", pred.Value.ViewRef)
				assert.Equal(t, "centerY", pred.Attribute)
			},
		},
		{
			name:  "leading attribute",
			input: "[view(other.leading)]",
			check: func(t *testing.T, prog *evfl.Program) {
				pred := prog.Statements[0].Views[0].Predicates[0]
				assert.Equal(t, "other", pred.Value.ViewRef)
				assert.Equal(t, "leading", pred.Attribute)
			},
		},
		{
			name:  "trailing attribute",
			input: "[view(other.trailing)]",
			check: func(t *testing.T, prog *evfl.Program) {
				pred := prog.Statements[0].Views[0].Predicates[0]
				assert.Equal(t, "other", pred.Value.ViewRef)
				assert.Equal(t, "trailing", pred.Attribute)
			},
		},
		{
			name:  "top attribute",
			input: "[view(header.top)]",
			check: func(t *testing.T, prog *evfl.Program) {
				pred := prog.Statements[0].Views[0].Predicates[0]
				assert.Equal(t, "header", pred.Value.ViewRef)
				assert.Equal(t, "top", pred.Attribute)
			},
		},
		{
			name:  "bottom attribute",
			input: "[view(footer.bottom)]",
			check: func(t *testing.T, prog *evfl.Program) {
				pred := prog.Statements[0].Views[0].Predicates[0]
				assert.Equal(t, "footer", pred.Value.ViewRef)
				assert.Equal(t, "bottom", pred.Attribute)
			},
		},
		{
			name:  "baseline attribute",
			input: "[label(text.baseline)]",
			check: func(t *testing.T, prog *evfl.Program) {
				pred := prog.Statements[0].Views[0].Predicates[0]
				assert.Equal(t, "text", pred.Value.ViewRef)
				assert.Equal(t, "baseline", pred.Attribute)
			},
		},
		{
			name:  "superview attribute reference",
			input: "[view(superview.width)]",
			check: func(t *testing.T, prog *evfl.Program) {
				pred := prog.Statements[0].Views[0].Predicates[0]
				assert.Equal(t, "superview", pred.Value.ViewRef)
				assert.Equal(t, "width", pred.Attribute)
			},
		},
		{
			name:  "attribute with multiplier and offset",
			input: "[view(==other.width*2+10)]",
			check: func(t *testing.T, prog *evfl.Program) {
				pred := prog.Statements[0].Views[0].Predicates[0]
				assert.Equal(t, vfl.ExpressionValue, pred.Value.Type)

				// Check the expression structure
				expr := pred.Value.Expression
				assert.Equal(t, evfl.Add, expr.Operator)

				// Left: other.width * 2
				leftExpr := expr.Left.Expression
				assert.Equal(t, evfl.Multiply, leftExpr.Operator)
				assert.Equal(t, "other", leftExpr.Left.ViewRef)
				assert.Equal(t, 2.0, leftExpr.Right.Constant)

				// Right: 10
				assert.Equal(t, 10.0, expr.Right.Constant)
			},
		},
		{
			name:    "invalid attribute name",
			input:   "[view(other.invalidAttr)]",
			wantErr: true,
		},
		{
			name:  "chained attribute references",
			input: "[view1(==view2.width)][view2(==view3.width)][view3(100)]",
			check: func(t *testing.T, prog *evfl.Program) {
				views := prog.Statements[0].Views

				// view1 references view2.width
				assert.Equal(t, "view2", views[0].Predicates[0].Value.ViewRef)
				assert.Equal(t, "width", views[0].Predicates[0].Attribute)

				// view2 references view3.width
				assert.Equal(t, "view3", views[1].Predicates[0].Value.ViewRef)
				assert.Equal(t, "width", views[1].Predicates[0].Attribute)

				// view3 has constant
				assert.Equal(t, 100.0, views[2].Predicates[0].Value.Constant)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := evfl.NewParser(evfl.Options{
				EnableAttributeRefs: true,
				EnableExpressions:   true, // For tests with expressions
			})

			prog, err := parser.Parse(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, prog)

			if tt.check != nil {
				tt.check(t, prog)
			}
		})
	}
}

func TestEVFLParser_AttributeReferenceValidation(t *testing.T) {
	parser := evfl.NewParser(evfl.Options{
		EnableAttributeRefs: true,
	})

	tests := []struct {
		name   string
		input  string
		errMsg string
	}{
		{
			name:   "circular attribute reference",
			input:  "[view1(view2.width)][view2(view1.width)]",
			errMsg: "circular",
		},
		{
			name:   "self reference",
			input:  "[view(view.width)]",
			errMsg: "self-reference",
		},
		{
			name:   "undefined view reference",
			input:  "[view(undefined.width)]",
			errMsg: "undefined view",
		},
		{
			name:   "incompatible attributes",
			input:  "V:[view(other.width)]", // Width in vertical orientation
			errMsg: "incompatible",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parser.Parse(tt.input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}