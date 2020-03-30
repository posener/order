package reflectutil

import (
	"fmt"
	"reflect"
)

// Slice is a wrapper around reflect.Value that extends the functionality of slice values.
type Slice struct {
	// Value holds the slice object.
	reflect.Value

	// swap function that swaps two elements in the slice. This swap function is created by
	// `reflect.Swapper` and acts on the original slice.
	swap func(i, j int)
	// swapOffset holds offset from original slice to adjust the swap function, in case that `Slice`
	// or `Slice3` functions were called and moved the slice starting point.
	swapOffset int
}

func NewSlice(slice reflect.Value) (Slice, error) {
	// Check slice type.
	s, ok := getSliceValue(slice)
	if !ok {
		return Slice{}, fmt.Errorf("not a slice: %v", slice.Type())
	}
	return Slice{
		Value: s,
		swap:  reflect.Swapper(s.Interface()),
	}, nil
}

func (s Slice) T() reflect.Type {
	return s.Type().Elem()
}

// Slice does reflect.Value.Slice
func (s Slice) Slice(i, j int) Slice {
	s.Value = s.Value.Slice(i, j)
	s.swapOffset += i
	return s
}

// Slice3 does reflect.Value.Slice3
func (s Slice) Slice3(i, j, k int) Slice {
	s.Value = s.Value.Slice3(i, j, k)
	s.swapOffset += i
	return s
}

// Swap swaps elements in position i and j.
func (s Slice) Swap(i, j int) {
	s.swap(i+s.swapOffset, j+s.swapOffset)
}

// getSliceValue returns the slice reflect.Value of the given slice, or a pointer to a slice.
func getSliceValue(s reflect.Value) (reflect.Value, bool) {
	for {
		switch s.Kind() {
		case reflect.Slice:
			return s, true
		case reflect.Ptr:
			s = s.Elem()
		default:
			return reflect.Value{}, false
		}
	}
}
