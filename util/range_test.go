package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeCovers(t *testing.T) {
	r1 := Range{0, 5}
	r2 := Range{2, 3}
	assert.True(t, r1.Covers(r2))
	assert.False(t, r2.Covers(r1))
}

func TestRangeIntersects(t *testing.T) {
	r1 := Range{0, 5}
	r2 := Range{4, 3}
	r3 := Range{9, 15}
	assert.True(t, r1.Intersects(r2))
	assert.True(t, r2.Intersects(r1))
	assert.False(t, r1.Intersects(r3))
}
