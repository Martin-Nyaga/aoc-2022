package pqueue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPqueueMin(t *testing.T) {
	q := NewPqueue[int, string](MinQueue)
	assert.True(t, q.Empty())

	q.Push(2, "priority 2")
	assert.Equal(t, 1, q.Len())

	q.Push(3, "priority 3")
	assert.Equal(t, 2, q.Len())

	q.Push(3, "priority 3")
	assert.Equal(t, 3, q.Len())

	val, err := q.Pop()
	assert.Nil(t, err)
	assert.Equal(t, "priority 2", val)
	assert.Equal(t, 2, q.Len())

	q.Push(1, "priority 1")
	assert.Equal(t, 3, q.Len())

	val, err = q.Pop()
	assert.Nil(t, err)
	assert.Equal(t, "priority 1", val)
	assert.Equal(t, 2, q.Len())

	val, err = q.Pop()
	assert.Equal(t, "priority 3", val)
	assert.Nil(t, err)

	val, err = q.Pop()
	assert.Equal(t, "priority 3", val)
	assert.Nil(t, err)

	val, err = q.Pop()
	assert.NotNil(t, err)
}

func TestPqueueMax(t *testing.T) {
	q := NewPqueue[int, string](MaxQueue)
	assert.True(t, q.Empty())

	q.Push(2, "priority 2")
	assert.Equal(t, 1, q.Len())

	q.Push(3, "priority 3")
	assert.Equal(t, 2, q.Len())

	q.Push(3, "priority 3")
	assert.Equal(t, 3, q.Len())

	val, err := q.Pop()
	assert.Nil(t, err)
	assert.Equal(t, "priority 3", val)
	assert.Equal(t, 2, q.Len())

	q.Push(1, "priority 1")
	assert.Equal(t, 3, q.Len())

	val, err = q.Pop()
	assert.Nil(t, err)
	assert.Equal(t, "priority 3", val)
	assert.Equal(t, 2, q.Len())

	val, err = q.Pop()
	assert.Equal(t, "priority 2", val)
	assert.Nil(t, err)

	val, err = q.Pop()
	assert.Equal(t, "priority 1", val)
	assert.Nil(t, err)

	val, err = q.Pop()
	assert.NotNil(t, err)
}
