package order

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/posener/order/internal/reflectutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	t.Parallel()

	tests := []struct {
		slice []int
	}{
		{slice: []int{1}},
		{slice: []int{4, 1, 3, 2}},
		{slice: []int{5, 20, 3, 10, 100}},
		{slice: []int{10, 1001, 23, 12, 43, 65, 504, 34, 123, 101, 21, 24, 11, -10, 999, 666, 1212}},
	}

	for _, tt := range tests {
		for k := range tt.slice {
			// k := 7
			t.Run(fmt.Sprintf("slice: %v/k: %v", tt.slice, k), func(t *testing.T) {
				slice := copySlice(tt.slice)

				// Apply the select algorithm.
				Select(slice, k)
				assert.ElementsMatch(t, tt.slice, slice)
				gotSelect := slice[k]

				t.Logf("Selected slice: %v", slice)

				// Calculate the k'th element by sorting the slice.
				wantSelect := calcKValue(tt.slice, k)
				assert.Equal(t, wantSelect, gotSelect)

				// Test partition side effect
				for _, v := range slice[:k] {
					assert.LessOrEqual(t, v, gotSelect)
				}
				for _, v := range slice[k:] {
					assert.GreaterOrEqual(t, v, gotSelect)
				}
			})
		}
	}
}

func TestSelect_partition(t *testing.T) {
	t.Parallel()

	a := []int{5, 4, 2, 3, 1}
	s, err := reflectutil.NewSlice(reflect.ValueOf(a))
	require.NoError(t, err)
	k := intFn.partition(s, 3)
	assert.Equal(t, 2, k)
	assert.Equal(t, []int{2, 1, 3, 4, 5}, a)
}

func TestSelect_sortSmallSlice(t *testing.T) {
	t.Parallel()

	a := []int{5, 1, -2, 10, 4}
	s, err := reflectutil.NewSlice(reflect.ValueOf(a))
	require.NoError(t, err)
	intFn.sortSmallSlice(s)
	assert.Equal(t, []int{-2, 1, 4, 5, 10}, a)
}

func TestSelect_pivot(t *testing.T) {
	t.Parallel()
	tests := []struct {
		slice   []int
		wantMed int
	}{
		{
			slice:   []int{5, 4, 2, 3, 1},
			wantMed: 3,
		},
		{
			slice:   []int{5, 4, 2, 3, 1, 10, 9, 8, 7},
			wantMed: 3,
		},
		{
			slice:   []int{5, 4, 2, 3, 1, 10, 9, 8, 7, 6, 15, 14, 13},
			wantMed: 8,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%b", tt.slice), func(t *testing.T) {
			s, err := reflectutil.NewSlice(reflect.ValueOf(tt.slice))
			require.NoError(t, err)
			intFn.pivot(s)
			t.Logf("slice: %v", tt.slice)
			assert.Equal(t, tt.wantMed, tt.slice[0])
		})
	}
}

func copySlice(s []int) []int {
	cp := make([]int, len(s))
	copy(cp, s)
	return cp
}

func calcKValue(s []int, k int) int {
	cp := copySlice(s)
	sort.Ints(cp)
	return cp[k]
}
