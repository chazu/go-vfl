package evfl_test

import (
	"testing"

	"github.com/chazu/go-vfl/evfl"
	"github.com/chazu/go-vfl/vfl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEVFLParser_Expressions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		check   func(t *testing.T, prog *evfl.Program)
		wantErr bool
	}{
		{
			name:  "simple arithmetic expression",
			input: "[view1(==view2.width*2-10)]",
			check: func(t *testing.T, prog *evfl.Program) {
				view := prog.Statements[0].Views[0]
				pred := view.Predicates[0]

				assert.Equal(t, evfl.Equal, pred.Relation)
				assert.Equal(t, vfl.ExpressionValue, pred.Value.Type)

				expr := pred.Value.Expression
				assert.NotNil(t, expr)
				assert.Equal(t, evfl.Subtract, expr.Operator)

				// Left side: view2.width * 2
				leftExpr := expr.Left.Expression
				assert.Equal(t, evfl.Multiply, leftExpr.Operator)
				assert.Equal(t, "view2", leftExpr.Left.ViewRef)
				assert.Equal(t, "width", pred.Attribute)
				assert.Equal(t, 2.0, leftExpr.Right.Constant)

				// Right side: 10
				assert.Equal(t, 10.0, expr.Right.Constant)
			},
		},
		{
			name:  "addition expression",
			input: "[view(==otherView+50)]",
			check: func(t *testing.T, prog *evfl.Program) {
				expr := prog.Statements[0].Views[0].Predicates[0].Value.Expression
				assert.Equal(t, evfl.Add, expr.Operator)
				assert.Equal(t, "otherView", expr.Left.ViewRef)
				assert.Equal(t, 50.0, expr.Right.Constant)
			},
		},
		{
			name:  "division expression",
			input: "[half(==container/2)]",
			check: func(t *testing.T, prog *evfl.Program) {
				expr := prog.Statements[0].Views[0].Predicates[0].Value.Expression
				assert.Equal(t, evfl.Divide, expr.Operator)
				assert.Equal(t, "container", expr.Left.ViewRef)
				assert.Equal(t, 2.0, expr.Right.Constant)
			},
		},
		{
			name:  "complex nested expression",
			input: "[view(==(superview.width-20)/2+10)]",
			check: func(t *testing.T, prog *evfl.Program) {
				// Should parse: ((superview.width - 20) / 2) + 10
				expr := prog.Statements[0].Views[0].Predicates[0].Value.Expression
				assert.Equal(t, evfl.Add, expr.Operator)

				// Left: (superview.width - 20) / 2
				leftExpr := expr.Left.Expression
				assert.Equal(t, evfl.Divide, leftExpr.Operator)

				// Left-left: superview.width - 20
				leftLeftExpr := leftExpr.Left.Expression
				assert.Equal(t, evfl.Subtract, leftLeftExpr.Operator)
				assert.Contains(t, leftLeftExpr.Left.ViewRef, "superview")
				assert.Equal(t, 20.0, leftLeftExpr.Right.Constant)

				// Left-right: 2
				assert.Equal(t, 2.0, leftExpr.Right.Constant)

				// Right: 10
				assert.Equal(t, 10.0, expr.Right.Constant)
			},
		},
		{
			name:    "division by zero should be caught",
			input:   "[view(==other/0)]",
			wantErr: true,
		},
		{
			name:  "expression with percentage",
			input: "[view(==50%*2)]",
			check: func(t *testing.T, prog *evfl.Program) {
				expr := prog.Statements[0].Views[0].Predicates[0].Value.Expression
				assert.Equal(t, evfl.Multiply, expr.Operator)
				assert.Equal(t, 50.0, expr.Left.Percentage)
				assert.Equal(t, 2.0, expr.Right.Constant)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := evfl.NewParser(evfl.Options{
				EnableExpressions: true,
				EnablePercentages: true,
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

func TestEVFLParser_ExpressionValidation(t *testing.T) {
	parser := evfl.NewParser(evfl.Options{
		EnableExpressions: true,
	})

	tests := []struct {
		name    string
		input   string
		errMsg  string
	}{
		{
			name:   "undefined view in expression",
			input:  "[view(==undefined.width*2)]",
			errMsg: "undefined view",
		},
		{
			name:   "invalid attribute",
			input:  "[view(==other.invalid)]",
			errMsg: "invalid attribute",
		},
		{
			name:   "type mismatch in expression",
			input:  "[view(==\"string\"+10)]",
			errMsg: "type mismatch",
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