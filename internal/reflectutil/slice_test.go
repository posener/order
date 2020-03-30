package reflectutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value  interface{}
		assert func(t *testing.T, got Slice)
	}{
		{
			value: []int{42},
			assert: func(t *testing.T, got Slice) {
				assert.True(t, 42 == got.Index(0).Interface().(int))
			},
		},
		// Empty slice.
		{
			value: []int{},
			assert: func(t *testing.T, got Slice) {
				assert.Equal(t, 0, got.Len())
			},
		},
		// Nil slice.
		{
			value: []int(nil),
			assert: func(t *testing.T, got Slice) {
				assert.Equal(t, 0, got.Len())
			},
		},
		// Different convertable kind.
		{
			value: []int8{42},
			assert: func(t *testing.T, got Slice) {
				assert.True(t, 42 == got.Index(0).Interface().(int8))
			},
		},
		// Pointer to slice of pointer values.
		{
			value: &([]*int{intPtr(42)}),
			assert: func(t *testing.T, got Slice) {
				assert.True(t, 42 == *got.Index(0).Interface().(*int))
			},
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.value), func(t *testing.T) {
			got, err := NewSlice(reflect.ValueOf(tt.value))
			require.NoError(t, err)
			tt.assert(t, got)
		})
	}
}

func TestSlice_failures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value interface{}
	}{
		// Not a slice.
		{value: 1},
	}

	for _, tt := range tests {
		t.Run(testName(tt.value), func(t *testing.T) {
			_, err := NewSlice(reflect.ValueOf(tt.value))
			assert.Error(t, err)
		})
	}
}

func TestSlice_swap(t *testing.T) {
	t.Parallel()
	t.Run("swap", func(t *testing.T) {
		a := []int{1, 2}
		s, err := NewSlice(reflect.ValueOf(a))
		require.NoError(t, err)
		s.Swap(0, 1)
		assert.Equal(t, []int{2, 1}, a)
	})

	t.Run("slice and swap", func(t *testing.T) {
		a := []int{1, 2, 3}
		s, err := NewSlice(reflect.ValueOf(a))
		require.NoError(t, err)
		s.Slice(1, 3).Swap(0, 1)
		assert.Equal(t, []int{1, 3, 2}, a)
	})

	t.Run("slice3 and swap", func(t *testing.T) {
		a := []int{1, 2, 3}
		s, err := NewSlice(reflect.ValueOf(a))
		require.NoError(t, err)
		s.Slice3(1, 3, 3).Swap(0, 1)
		assert.Equal(t, []int{1, 3, 2}, a)
	})
}
