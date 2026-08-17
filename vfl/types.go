package vfl

// Orientation represents layout direction
type Orientation int

const (
	Horizontal Orientation = iota
	Vertical
	Z
	Combined // HV: for EVFL
)

// String returns the string representation of the orientation
func (o Orientation) String() string {
	switch o {
	case Horizontal:
		return "H"
	case Vertical:
		return "V"
	case Z:
		return "Z"
	case Combined:
		return "HV"
	default:
		return "unknown"
	}
}

// Relation represents constraint relations
type Relation int

const (
	Equal Relation = iota
	GreaterEqual
	LessEqual
)

// String returns the string representation of the relation
func (r Relation) String() string {
	switch r {
	case Equal:
		return "=="
	case GreaterEqual:
		return ">="
	case LessEqual:
		return "<="
	default:
		return "unknown"
	}
}

// ValueType indicates which Value field is set
type ValueType int

const (
	ConstantValue ValueType = iota
	ViewRefValue
	MetricValue
	ExpressionValue
	PercentageValue
)

// String returns the string representation of the value type
func (v ValueType) String() string {
	switch v {
	case ConstantValue:
		return "constant"
	case ViewRefValue:
		return "view reference"
	case MetricValue:
		return "metric"
	case ExpressionValue:
		return "expression"
	case PercentageValue:
		return "percentage"
	default:
		return "unknown"
	}
}

// Operator represents math operations in EVFL
type Operator int

const (
	Add Operator = iota
	Subtract
	Multiply
	Divide
)

// String returns the string representation of the operator
func (o Operator) String() string {
	switch o {
	case Add:
		return "+"
	case Subtract:
		return "-"
	case Multiply:
		return "*"
	case Divide:
		return "/"
	default:
		return "unknown"
	}
}

// Attribute represents a view attribute
type Attribute int

const (
	Left Attribute = iota
	Right
	Top
	Bottom
	Leading
	Trailing
	Width
	Height
	CenterX
	CenterY
	Baseline
	NotAnAttribute
)

// String returns the string representation of the attribute
func (a Attribute) String() string {
	switch a {
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
	case Width:
		return "width"
	case Height:
		return "height"
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

// StackLayout defines how views are arranged in a stack
type StackLayout int

const (
	LinearStack StackLayout = iota
	GridStack
	FlowStack
)

// String returns the string representation of the stack layout
func (s StackLayout) String() string {
	switch s {
	case LinearStack:
		return "linear"
	case GridStack:
		return "grid"
	case FlowStack:
		return "flow"
	default:
		return "unknown"
	}
}

// AlignmentType for disconnections
type AlignmentType int

const (
	AlignTop AlignmentType = iota
	AlignBottom
	AlignLeft
	AlignRight
	AlignCenterX
	AlignCenterY
	AlignBaseline
)

// String returns the string representation of the alignment type
func (a AlignmentType) String() string {
	switch a {
	case AlignTop:
		return "top"
	case AlignBottom:
		return "bottom"
	case AlignLeft:
		return "left"
	case AlignRight:
		return "right"
	case AlignCenterX:
		return "centerX"
	case AlignCenterY:
		return "centerY"
	case AlignBaseline:
		return "baseline"
	default:
		return "unknown"
	}
}