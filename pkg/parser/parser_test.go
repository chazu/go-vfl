package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasicProgram(t *testing.T) {
	program := "[Test]"
	psr := New()
	res, err := psr.ParseProgram(program)
	assert.Nil(t, err)
	assert.Len(t, res.Views, 1)
	assert.Equal(t, res.Views[0].Name, "Test")
}

func TestBasicPredicate(t *testing.T) {
	program := "[Test 40]"
	psr := New()
	res, err := psr.ParseProgram(program)

	assert.Nil(t, err)
	assert.Len(t, res.Views, 1)
	assert.Equal(t, res.Views[0].Name, "Test")

	assert.NotNil(t, res.Views[0].Predicate.Predicate)
	assert.Equal(t, res.Views[0].Predicate.Predicate.Object.Number, 40)

	assert.True(t, res.Views[0].Predicate.Predicate.Relation.Eq == true)

}
