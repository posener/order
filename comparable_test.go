package order

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPredefinedTypes(t *testing.T) {
	t.Parallel()

	assert.True(t, Is(1).Greater(0))
	assert.True(t, Is(1).Equal(1))
	assert.True(t, Is(1).Less(2))

	assert.True(t, Is("b").Greater("a"))
	assert.True(t, Is("b").Equal("b"))
	assert.True(t, Is("b").Less("c"))

	assert.True(t, Is([]byte{1}).Greater([]byte{0}))
	assert.True(t, Is([]byte{1}).Equal([]byte{1}))
	assert.True(t, Is([]byte{1}).Less([]byte{2}))

	assert.True(t, Is(true).Greater(false))
	assert.True(t, Is(true).Equal(true))
	assert.True(t, Is(false).Less(true))

	assert.True(t, Is(time.Unix(1, 0)).Greater(time.Unix(0, 0)))
	assert.True(t, Is(time.Unix(1, 0)).Equal(time.Unix(1, 0)))
	assert.True(t, Is(time.Unix(1, 0)).Less(time.Unix(2, 0)))

	assert.True(t, Is(1*time.Nanosecond).Greater(0*time.Nanosecond))
	assert.True(t, Is(1*time.Nanosecond).Equal(1*time.Nanosecond))
	assert.True(t, Is(1*time.Nanosecond).Less(2*time.Nanosecond))
}

type notComparable struct{}

type wrong1 struct{}

func (w wrong1) Compare(other wrong1) bool { return false }

func TestComparable_invalid(t *testing.T) {
	t.Parallel()

	fns := []func(v interface{}){
		func(v interface{}) { Sort(v) },
		func(v interface{}) { SortStable(v) },
		func(v interface{}) { Search(v, 1) },
		func(v interface{}) { IsSorted(v) },
		func(v interface{}) { IsStrictSorted(v) },
		func(v interface{}) { MinMax(v) },
		func(v interface{}) { Select(v, 0) },
	}

	for _, fn := range fns {
		t.Run("not a slice", func(t *testing.T) {
			assert.Panics(t, func() {
				fn(1)
			})
		})
		t.Run("not a comparable", func(t *testing.T) {
			assert.Panics(t, func() {
				fn([]notComparable{})
			})
		})
		t.Run("wrong signature", func(t *testing.T) {
			assert.Panics(t, func() {
				fn([]wrong1{})
			})
		})
		t.Run("slice of invalid type", func(t *testing.T) {
			assert.Panics(t, func() {
				fn([]func(){func() {}})
			})
		})
	}

	// Other wrong arguments

	// Search invalid value type.
	assert.Panics(t, func() { intFn.Search([]int{}, true) })

	// Select K out of bounds.
	assert.Panics(t, func() { Select([]int{1}, -1) })
	assert.Panics(t, func() { Select([]int{1}, 1) })
	assert.Panics(t, func() { Select([]int{}, 0) })
}
