package evfl

import (
	"fmt"
	"github.com/chazu/go-vfl/vfl"
)

// ExpressionEvaluator evaluates mathematical expressions in EVFL
type ExpressionEvaluator interface {
	Evaluate(expr *vfl.Expression, context EvalContext) (float64, error)
	EvaluateValue(value vfl.Value, context EvalContext) (float64, error)
}

// EvalContext provides context for expression evaluation
type EvalContext struct {
	ViewWidths  map[string]float64
	ViewHeights map[string]float64
	Metrics     map[string]float64
	ContainerWidth  float64
	ContainerHeight float64
}

// expressionEvaluator implements ExpressionEvaluator
type expressionEvaluator struct{}

// NewExpressionEvaluator creates a new expression evaluator
func NewExpressionEvaluator() ExpressionEvaluator {
	return &expressionEvaluator{}
}

// Evaluate evaluates a mathematical expression
func (e *expressionEvaluator) Evaluate(expr *vfl.Expression, context EvalContext) (float64, error) {
	if expr == nil {
		return 0, fmt.Errorf("nil expression")
	}

	// Evaluate left operand
	left, err := e.EvaluateValue(expr.Left, context)
	if err != nil {
		return 0, fmt.Errorf("evaluating left operand: %w", err)
	}

	// Evaluate right operand
	right, err := e.EvaluateValue(expr.Right, context)
	if err != nil {
		return 0, fmt.Errorf("evaluating right operand: %w", err)
	}

	// Apply operator
	switch expr.Operator {
	case vfl.Add:
		return left + right, nil
	case vfl.Subtract:
		return left - right, nil
	case vfl.Multiply:
		return left * right, nil
	case vfl.Divide:
		if right == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return left / right, nil
	default:
		return 0, fmt.Errorf("unknown operator: %v", expr.Operator)
	}
}

// EvaluateValue evaluates a value which may be a constant, percentage, view reference, etc.
func (e *expressionEvaluator) EvaluateValue(value vfl.Value, context EvalContext) (float64, error) {
	switch value.Type {
	case vfl.ConstantValue:
		return value.Constant, nil

	case vfl.PercentageValue:
		// Use container width by default, could be made more sophisticated
		return value.Percentage * context.ContainerWidth / 100, nil

	case vfl.ViewRefValue:
		// Look up view dimension
		if width, ok := context.ViewWidths[value.ViewRef]; ok {
			return width, nil
		}
		return 0, fmt.Errorf("view '%s' not found in context", value.ViewRef)

	case vfl.MetricValue:
		// Look up metric
		if metric, ok := context.Metrics[value.MetricName]; ok {
			return metric, nil
		}
		return 0, fmt.Errorf("metric '%s' not found", value.MetricName)

	case vfl.ExpressionValue:
		// Recursively evaluate nested expression
		if value.Expression != nil {
			return e.Evaluate(value.Expression, context)
		}
		return 0, fmt.Errorf("expression value has nil expression")

	default:
		return 0, fmt.Errorf("unsupported value type: %v", value.Type)
	}
}

// SimplifyExpression attempts to simplify an expression to a constant if possible
func SimplifyExpression(expr *vfl.Expression) *vfl.Value {
	// If both operands are constants, compute the result
	if expr.Left.Type == vfl.ConstantValue && expr.Right.Type == vfl.ConstantValue {
		var result float64
		switch expr.Operator {
		case vfl.Add:
			result = expr.Left.Constant + expr.Right.Constant
		case vfl.Subtract:
			result = expr.Left.Constant - expr.Right.Constant
		case vfl.Multiply:
			result = expr.Left.Constant * expr.Right.Constant
		case vfl.Divide:
			if expr.Right.Constant != 0 {
				result = expr.Left.Constant / expr.Right.Constant
			} else {
				return nil // Can't simplify division by zero
			}
		default:
			return nil
		}
		return &vfl.Value{
			Type:     vfl.ConstantValue,
			Constant: result,
		}
	}
	return nil
}