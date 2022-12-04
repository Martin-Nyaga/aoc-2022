package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptySet(t *testing.T) {
	set := NewSet[int]()
	assert.True(t, set.Empty())
	assert.Equal(t, 0, set.Len())
}

func TestNonEmptySet(t *testing.T) {
	set := NewSet([]int{1, 2}...)
	assert.False(t, set.Empty())
	assert.Equal(t, 2, set.Len())
	assert.True(t, set.Has(1))
}

func TestAdd(t *testing.T) {
	set := NewSet[int]()
	assert.True(t, set.Empty())
	set.Add(1)
	assert.Equal(t, 1, set.Len())
	assert.True(t, set.Has(1))
	assert.False(t, set.Empty())
}

func TestIntersection(t *testing.T) {
	a := NewSet([]int{1, 2}...)
	b := NewSet([]int{2, 3}...)
	c := NewSet[int]()

	intersection1 := a.Intersection(&b)
	intersection2 := a.Intersection(&c)
	assert.Equal(t, 1, intersection1.Len())
	assert.True(t, intersection1.Has(2))
	assert.True(t, intersection2.Empty())
}

func TestRemove(t *testing.T) {
	set := NewSet([]int{2}...)
	err := set.Remove(2)
	assert.Nil(t, err)
	assert.True(t, set.Empty())
	err = set.Remove(1)
	assert.NotNil(t, err)
}

func TestPop(t *testing.T) {
	set := NewSet([]int{1}...)
	v, err := set.PopAny()
	assert.Equal(t, 1, v)
	assert.Nil(t, err)
	_, err = set.PopAny()
	assert.NotNil(t, err)
}
