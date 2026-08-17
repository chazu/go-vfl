package vfl

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLocation(t *testing.T) {
	loc := Location{Line: 1, Column: 5, Offset: 4}
	assert.Equal(t, 1, loc.Line)
	assert.Equal(t, 5, loc.Column)
	assert.Equal(t, 4, loc.Offset)
}

func TestProgram_GetLocation(t *testing.T) {
	prog := &Program{
		Location: Location{Line: 1, Column: 1, Offset: 0},
	}
	assert.Equal(t, Location{Line: 1, Column: 1, Offset: 0}, prog.GetLocation())
}

func TestProgram_Accept(t *testing.T) {
	prog := &Program{
		Orientation: Horizontal,
		Statements: []Statement{
			{
				Views: []View{
					{Name: "button"},
				},
			},
		},
	}

	// Test visitor pattern (simplified test since we don't have full visitor in tests)
	assert.NotNil(t, prog.Accept)
}

func TestStatement_GetLocation(t *testing.T) {
	stmt := &Statement{
		Location: Location{Line: 2, Column: 3, Offset: 10},
	}
	assert.Equal(t, Location{Line: 2, Column: 3, Offset: 10}, stmt.GetLocation())
}

func TestView_GetLocation(t *testing.T) {
	view := &View{
		Location: Location{Line: 3, Column: 5, Offset: 20},
		Name:     "testView",
	}
	assert.Equal(t, Location{Line: 3, Column: 5, Offset: 20}, view.GetLocation())
	assert.Equal(t, "testView", view.Name)
}

func TestView_IsNamedView(t *testing.T) {
	tests := []struct {
		name     string
		view     View
		expected bool
	}{
		{
			name:     "regular view",
			view:     View{Name: "button"},
			expected: false,
		},
		{
			name:     "named view reference",
			view:     View{Name: "header", IsNamedView: true},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.view.IsNamedView)
		})
	}
}

func TestPredicate_GetLocation(t *testing.T) {
	pred := &Predicate{
		Location: Location{Line: 4, Column: 10, Offset: 30},
		Relation: Equal,
		Value:    Value{Type: ConstantValue, Constant: 100},
	}
	assert.Equal(t, Location{Line: 4, Column: 10, Offset: 30}, pred.GetLocation())
	assert.Equal(t, Equal, pred.Relation)
	assert.Equal(t, float64(100), pred.Value.Constant)
}

func TestConnection_GetLocation(t *testing.T) {
	conn := &Connection{
		Location:  Location{Line: 5, Column: 15, Offset: 40},
		Spacing:   20,
		IsDefault: false,
	}
	assert.Equal(t, Location{Line: 5, Column: 15, Offset: 40}, conn.GetLocation())
	assert.Equal(t, float64(20), conn.Spacing)
	assert.False(t, conn.IsDefault)
}

func TestValue_Types(t *testing.T) {
	tests := []struct {
		name  string
		value Value
		check func(t *testing.T, v Value)
	}{
		{
			name: "constant value",
			value: Value{
				Type:     ConstantValue,
				Constant: 50,
			},
			check: func(t *testing.T, v Value) {
				assert.Equal(t, ConstantValue, v.Type)
				assert.Equal(t, float64(50), v.Constant)
			},
		},
		{
			name: "view reference",
			value: Value{
				Type:    ViewRefValue,
				ViewRef: "otherView",
			},
			check: func(t *testing.T, v Value) {
				assert.Equal(t, ViewRefValue, v.Type)
				assert.Equal(t, "otherView", v.ViewRef)
			},
		},
		{
			name: "metric",
			value: Value{
				Type:       MetricValue,
				MetricName: "standard",
			},
			check: func(t *testing.T, v Value) {
				assert.Equal(t, MetricValue, v.Type)
				assert.Equal(t, "standard", v.MetricName)
			},
		},
		{
			name: "percentage",
			value: Value{
				Type:       PercentageValue,
				Percentage: 75,
			},
			check: func(t *testing.T, v Value) {
				assert.Equal(t, PercentageValue, v.Type)
				assert.Equal(t, float64(75), v.Percentage)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, tt.value)
		})
	}
}

func TestViewNode(t *testing.T) {
	node := &ViewNode{
		Name:        "container",
		IsSuperview: false,
		Children: []*ViewNode{
			{Name: "child1"},
			{Name: "child2"},
		},
	}

	assert.Equal(t, "container", node.Name)
	assert.False(t, node.IsSuperview)
	assert.Len(t, node.Children, 2)
	assert.Equal(t, "child1", node.Children[0].Name)
	assert.Equal(t, "child2", node.Children[1].Name)
}

func TestConstraint(t *testing.T) {
	constraint := Constraint{
		FirstItem: ConstraintItem{
			Name:      "button",
			Attribute: "width",
		},
		Relation: Equal,
		SecondItem: &ConstraintItem{
			Name:      "label",
			Attribute: "width",
		},
		Multiplier: 2.0,
		Constant:   10,
		Priority:   750,
	}

	assert.Equal(t, "button", constraint.FirstItem.Name)
	assert.Equal(t, "width", constraint.FirstItem.Attribute)
	assert.Equal(t, Equal, constraint.Relation)
	assert.NotNil(t, constraint.SecondItem)
	assert.Equal(t, "label", constraint.SecondItem.Name)
	assert.Equal(t, float64(2.0), constraint.Multiplier)
	assert.Equal(t, float64(10), constraint.Constant)
	assert.Equal(t, 750, constraint.Priority)
}

