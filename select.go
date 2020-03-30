package order

import (
	"fmt"
	"reflect"

	"github.com/posener/order/internal/reflectutil"
)

// Select applies select-k algorithm on the given slice and k index. After invoking this method,
// the k'th greatest element according to the comparison function will be available in the k'th
// index.
// As a side effect, the slice will be partitioned according to the k'th index:
//
// 	{slice[i] <= slice[k] | i < k}
// 	{slice[i] >= slice[k] | i > k}
//
// This function will panic if k is out of the bounds of slice.
func (fns Fns) Select(slice interface{}, k int) {
	s := fns.mustSlice(reflect.ValueOf(slice))
	if k < 0 || k >= s.Len() {
		panic(fmt.Sprintf("k value %d out of bounds: [0, %d)", k, s.Len()))
	}
	for {
		fns.pivot(s)
		pivot := fns.partition(s, 0)
		switch {
		case pivot == k:
			return
		case pivot < k:
			k -= pivot + 1
			s = s.Slice(pivot+1, s.Len())
		default: // pivot > k
			s = s.Slice(0, pivot)
		}
	}
}

// pivot puts the median-of-medians in the index 0 of the slice.
func (fns Fns) pivot(s reflectutil.Slice) {
	const size = 5

	for s.Len() > 0 {
		n := s.Len()
		// For 5 or less elements return the median.
		if n <= size {
			fns.sortSmallSlice(s)
			s.Swap((n-1)/2, 0)
			return
		}

		// Move the medians of 5 elements groups to the beginning of the slice.
		medLen := 0
		for left := 0; left < n; left += size {
			// Sort the group of 5 elements.
			right := minInt(left+size, n)
			fns.sortSmallSlice(s.Slice(left, right))

			// Move the middle element to the beginning of the slice.
			s.Swap((left+right-1)/2, medLen)
			medLen++
		}

		// Update the slice to point only on the medians slice, such that in the next iterations the
		// medians of these medians will be found.
		s = s.Slice(0, medLen)
	}
}

// partition updates the slice according to a given pivot index. It returns a new pivot index such
// that all elements left to the new pivot index are smaller then s[pivot] and all elements left to
// the new pivot index are greater than or equal to the pivot value.
func (fns Fns) partition(s reflectutil.Slice, p int) int {
	n := s.Len()

	// Put the pivot at the end of the slice.
	s.Swap(p, n-1)
	pivot := s.Index(n - 1)

	// Iterate over the slice and move to cursor location all values that are smaller than the pivot
	// value.
	cursor := 0
	for i := 0; i < n-1; i++ {
		if fns.compare(s.Index(i), pivot) < 0 {
			s.Swap(cursor, i)
			cursor++
		}
	}

	// Move the pivot value back to the cursor location.
	s.Swap(cursor, n-1)

	return cursor
}

// sortSmallSlice simply and inefficiently insertion-sorts a small slice.
func (fns Fns) sortSmallSlice(s reflectutil.Slice) {
	for i := 1; i < s.Len(); i++ {
		for j := i; j > 0 && fns.compare(s.Index(j-1), s.Index(j)) > 0; j-- {
			s.Swap(j-1, j)
		}
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
