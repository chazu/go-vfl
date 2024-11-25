package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var psr = New()

func TestBasicProgram(t *testing.T) {
	program := "[Test]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 1)
	assert.Equal(t, res.Views[0].Name, "Test")
}

func TestTwoViews(t *testing.T) {
	program := "[Test1][Test2]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 2)
	assert.Equal(t, res.Views[0].Name, "Test1")
	assert.Equal(t, res.Views[1].Name, "Test2")
}

func TestRelationGreaterThan(t *testing.T) {
	program := "[Test1 >=40]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 1)
	assert.Equal(t, res.Views[0].Name, "Test1")

	assert.NotNil(t, res.Views[0].Predicate.Predicate)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Object.Number, 40)
	assert.NotNil(t, res.Views[0].Predicate.Predicate.Relation)
	assert.True(t, res.Views[0].Predicate.Predicate.Relation.Gte)
}

func TestRelationGreaterThanPriority(t *testing.T) {
	program := "[Test1 >=40@10]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 1)
	assert.Equal(t, res.Views[0].Name, "Test1")

	assert.NotNil(t, res.Views[0].Predicate.Predicate)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Object.Number, 40)
	assert.NotNil(t, res.Views[0].Predicate.Predicate.Relation)
	assert.True(t, res.Views[0].Predicate.Predicate.Relation.Gte)
	assert.NotNil(t, res.Views[0].Predicate.Predicate.Priority)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Priority.Value, 10)
}

func TestTwoRelations(t *testing.T) {
	program := "[Test1 >=40][Test2 >=Foo]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 2)
	assert.Equal(t, res.Views[0].Name, "Test1")
	assert.Equal(t, res.Views[1].Name, "Test2")

	assert.NotNil(t, res.Views[0].Predicate.Predicate)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Object.Number, 40)
	assert.NotNil(t, res.Views[0].Predicate.Predicate.Relation)
	assert.True(t, res.Views[0].Predicate.Predicate.Relation.Gte)

	assert.NotNil(t, res.Views[1].Predicate.Predicate)
	assert.Equal(t, res.Views[1].Predicate.Predicate.Object.ViewName, "Foo")
	assert.True(t, res.Views[1].Predicate.Predicate.Relation.Gte)
}

func TestPriorityOnViewReference(t *testing.T) {
	program := "[Test1 >=40][Test2 >=Foo@10]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 2)
	assert.Equal(t, res.Views[0].Name, "Test1")
	assert.Equal(t, res.Views[1].Name, "Test2")

	assert.NotNil(t, res.Views[0].Predicate.Predicate)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Object.Number, 40)
	assert.NotNil(t, res.Views[0].Predicate.Predicate.Relation)
	assert.True(t, res.Views[0].Predicate.Predicate.Relation.Gte)

	assert.NotNil(t, res.Views[1].Predicate.Predicate)
	assert.Equal(t, res.Views[1].Predicate.Predicate.Object.ViewName, "Foo")
	assert.True(t, res.Views[1].Predicate.Predicate.Relation.Gte)
	assert.NotNil(t, res.Views[1].Predicate.Predicate.Priority)
	assert.Equal(t, res.Views[1].Predicate.Predicate.Priority.Value, 10)
}

func TestTwoRelationsInParentheses(t *testing.T) {
	program := "[Test1 (>=40,<=80)]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 1)
	assert.Equal(t, res.Views[0].Name, "Test1")

	assert.Len(t, res.Views[0].Predicate.Predicates, 2)
	assert.Equal(t, res.Views[0].Predicate.Predicates[0].Object.Number, 40)
	assert.True(t, res.Views[0].Predicate.Predicates[0].Relation.Gte)

	assert.Equal(t, res.Views[0].Predicate.Predicates[1].Object.Number, 80)
	assert.True(t, res.Views[0].Predicate.Predicates[1].Relation.Lte)
}

func TestBasicPredicate(t *testing.T) {
	program := "[Test 40]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 1)
	assert.Equal(t, res.Views[0].Name, "Test")

	assert.NotNil(t, res.Views[0].Predicate.Predicate)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Object.Number, 40)

	assert.Nil(t, res.Views[0].Predicate.Predicate.Relation)
}

func TestPriority(t *testing.T) {
	program := "[Test 40@10]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 1)
	assert.Equal(t, res.Views[0].Name, "Test")

	assert.NotNil(t, res.Views[0].Predicate.Predicate)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Object.Number, 40)
	assert.Nil(t, res.Views[0].Predicate.Predicate.Relation)
	assert.NotNil(t, res.Views[0].Predicate.Predicate.Priority)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Priority.Value, 10)
}

func TestPriorityOnRelation(t *testing.T) {
	program := "[Test1 >=40@10]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 1)
	assert.Equal(t, res.Views[0].Name, "Test1")

	assert.NotNil(t, res.Views[0].Predicate.Predicate)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Object.Number, 40)
	assert.NotNil(t, res.Views[0].Predicate.Predicate.Relation)
	assert.True(t, res.Views[0].Predicate.Predicate.Relation.Gte)
	assert.NotNil(t, res.Views[0].Predicate.Predicate.Priority)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Priority.Value, 10)
}

func TestOrientationHorizontal(t *testing.T) {
	program := "H:[TestView]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.NotNil(t, res.Orientation)
	assert.Equal(t, *res.Orientation.Direction, "H")
}

func TestOrientationVertical(t *testing.T) {
	program := "V:[TestView]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.NotNil(t, res.Orientation)
	assert.Equal(t, *res.Orientation.Direction, "V")
}

func TestConnectionBetweenViews(t *testing.T) {
	program := "[TestView]-[TestTwo]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.NotNil(t, res.Views[1].Connection)
	// TODO we'd like this to be not nil so
	// lets make the connection bidirectional
	assert.Nil(t, res.Views[0].Connection)
}

func TestConnectionBetweenViewsWithConstant(t *testing.T) {
	program := "V:[TestView]-50-[TestTwo]"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.NotNil(t, res.Views[1].Connection)
	assert.Equal(t, res.Views[1].Connection.Predicates.Predicate.Object.Number, 50)
	//assert.Nil(t, res.Views[1].Connection)
}

func TestLeadingAndTrailingSuperviewConnections(t *testing.T) {
	program := "|-[Test]-|"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.NotNil(t, res.LeadingSuperviewConnection)
	assert.NotNil(t, res.TrailingSuperviewConnection)
}

func TestSuperviewConnectionPredicatesSimple(t *testing.T) {
	program := "|-50-[Test]-50-|"
	res, err := psr.ParseProgram(program)
	assert.Nil(t, err)
	assert.Equal(t, res.LeadingSuperviewConnection.Connection.Predicates.Predicate.Object.Number, 50)
	assert.Equal(t, res.TrailingSuperviewConnection.Connection.Predicates.Predicate.Object.Number, 50)
}

func TestSuperviewConnectionPredicatesComplex(t *testing.T) {
	program := "|-(>=50@10)-[Test]-(<=50@10)-|"
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	lsvCon := res.LeadingSuperviewConnection.Connection
	tsvCon := res.TrailingSuperviewConnection.Connection
	assert.Len(t, lsvCon.Predicates.Predicates, 1)
	assert.Len(t, tsvCon.Predicates.Predicates, 1)

	assert.True(t, lsvCon.Predicates.Predicates[0].Relation.Gte)
	assert.False(t, lsvCon.Predicates.Predicates[0].Relation.Lte)
	assert.False(t, lsvCon.Predicates.Predicates[0].Relation.Eq)

	assert.False(t, tsvCon.Predicates.Predicates[0].Relation.Gte)
	assert.True(t, tsvCon.Predicates.Predicates[0].Relation.Lte)
	assert.False(t, tsvCon.Predicates.Predicates[0].Relation.Eq)

	assert.Equal(t, lsvCon.Predicates.Predicates[0].Priority.Value, 10)
	assert.Equal(t, tsvCon.Predicates.Predicates[0].Priority.Value, 10)

	assert.Equal(t, lsvCon.Predicates.Predicates[0].Object.Number, 50)
	assert.Equal(t, tsvCon.Predicates.Predicates[0].Object.Number, 50)
}
