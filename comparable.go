package order

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/posener/order/internal/reflectutil"
)

// Convenient comparison function for standard types.
var (
// cmpInts    fn3way = func(a, b reflect.Value) int { return a.Interface().(int) - b.Interface().(int) }
// cmpStrings fn3way = func(a, b reflect.Value) int { return strings.Compare(a.Interface().(string), b.Interface().(string)) }
// cmpBytes   fn3way = func(a, b reflect.Value) int { return bytes.Compare(a.Interface().([]byte), b.Interface().([]byte)) }
// cmpBools   fn3way = func(a, b reflect.Value) int {
// 	aa := a.Interface().(bool)
// 	bb := b.Interface().(bool)
// 	switch {
// 	case aa == bb:
// 		return 0
// 	case aa:
// 		return 1
// 	default:
// 		return -1
// 	}
// }
// cmpTimes fn3way = func(a, b reflect.Value) int {
// 	aa := a.Interface().(time.Time)
// 	bb := b.Interface().(time.Time)
// 	switch {
// 	case aa.Equal(bb):
// 		return 0
// 	case aa.After(bb):
// 		return 1
// 	default:
// 		return -1
// 	}
// }
)

// Is returns a Condition<T> for type T the implements a `func (T) Compare(T) int`.  It panics if
// value does not implement the compare function.
func Is(value interface{}) Condition {
	return compareableFn(reflect.TypeOf(value)).Is(value)
}

// Sort a Slice<T> if T implements a `func (T) Compare(T) int`. See Fn.Sort. It panics if slice does
// not implement the compare function.
func Sort(slice interface{}) {
	compareableSlice(reflect.ValueOf(slice)).Sort(slice)
}

// SortStable a Slice<T> if T implements a `func (T) Compare(T) int`. See Fn.SortStable.  It panics
// if slice does not implement the compare function.
func SortStable(slice interface{}) {
	compareableSlice(reflect.ValueOf(slice)).SortStable(slice)
}

// Search a Slice<T> if T implements a `func (T) Compare(T) int` for a value. See Fn.Search.
func Search(slice, value interface{}) int {
	return compareableSlice(reflect.ValueOf(slice)).Search(slice, value)
}

// MinMax returns the indices of the minimal and maximal values in a Slice<T> if T implements a
// `func (T) Compare(T) int` for a value. See Fn.MinMax. It panics if slice does not implement the
// compare function.
func MinMax(slice interface{}) (min, max int) {
	return compareableSlice(reflect.ValueOf(slice)).MinMax(slice)
}

// IsSorted returns whether a Slice<T> if T implements a `func (T) Compare(T) int` is sorted. See
// Fn.IsSorted. It panics if slice does not implement the compare function.
func IsSorted(slice interface{}) bool {
	return compareableSlice(reflect.ValueOf(slice)).IsSorted(slice)
}

// IsStrictSorted returns whether a Slice<T> if T implements a `func (T) Compare(T) int` is strictly
// sorted. See Fn.IsStrictSorted. It panics if slice does not implement the compare function.
func IsStrictSorted(slice interface{}) bool {
	return compareableSlice(reflect.ValueOf(slice)).IsStrictSorted(slice)
}

// Select applies select-k algorithm on a Slice<T> if T implements a `func (T) Compare(T) int`. See
// Fn.Select. It panics if slice does not implement the compare function.
func Select(slice interface{}, k int) {
	compareableSlice(reflect.ValueOf(slice)).Select(slice, k)
}

func compareableFn(tp reflect.Type) Fns {
	f, err := fnOfComparableT(tp)
	if err != nil {
		panic(err)
	}
	return f
}

// Return a compare function for a given slice.
func compareableSlice(slice reflect.Value) Fns {
	s, err := reflectutil.NewSlice(slice)
	if err != nil {
		panic(err)
	}
	return compareableFn(s.T())
}

var predefined = []Fns{
	By(func(a, b int64) int { return int(a - b) }),
	By(func(a, b uint64) int { return int(a - b) }),
	By(strings.Compare),
	By(bytes.Compare),
	By(func(a, b bool) int {
		switch {
		case a == b:
			return 0
		case a:
			return 1
		default:
			return -1
		}
	}),
	By(func(a, b time.Time) int {
		switch {
		case a.Equal(b):
			return 0
		case a.After(b):
			return 1
		default:
			return -1
		}
	}),
}

func fnOfComparableT(tp reflect.Type) (Fns, error) {
	ss := fmt.Sprintf("%v", tp)
	_ = ss
	method, ok := tp.MethodByName("Compare")
	if ok {
		fn, err := newFn(method.Func)
		if err != nil {
			return nil, fmt.Errorf("invalid `Compare` signature: %s", err)
		}
		return Fns{fn}, nil
	}

	for _, fn := range predefined {
		if fn.check(tp) {
			return fn, nil
		}
	}

	return nil, fmt.Errorf("Type %v should have a method 'Compare'", tp)
}
