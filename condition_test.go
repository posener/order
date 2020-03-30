package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	t.Parallel()

	assert.True(t, Is(1).Equal(1))
	assert.False(t, Is(1).Equal(2))

	assert.True(t, Is(1).NotEqual(2))
	assert.False(t, Is(1).NotEqual(1))

	assert.True(t, Is(1).Greater(0))
	assert.False(t, Is(1).Greater(1))
	assert.False(t, Is(1).Greater(2))

	assert.True(t, Is(1).GreaterEqual(0))
	assert.True(t, Is(1).GreaterEqual(1))
	assert.False(t, Is(1).GreaterEqual(2))

	assert.False(t, Is(1).Less(0))
	assert.False(t, Is(1).Less(1))
	assert.True(t, Is(1).Less(2))

	assert.False(t, Is(1).LessEqual(0))
	assert.True(t, Is(1).LessEqual(1))
	assert.True(t, Is(1).LessEqual(2))
}

func TestIs_invalidArgType(t *testing.T) {
	t.Parallel()

	// Test lhs.
	assert.Panics(t, func() { By(func(a, b int) int { return 0 }).Is(true) })

	cIs := Is(1)

	// Test rhs.
	assert.Panics(t, func() { cIs.Equal(true) })
	assert.Panics(t, func() { cIs.NotEqual(true) })
	assert.Panics(t, func() { cIs.Greater(true) })
	assert.Panics(t, func() { cIs.GreaterEqual(true) })
	assert.Panics(t, func() { cIs.Less(true) })
	assert.Panics(t, func() { cIs.LessEqual(true) })
}
