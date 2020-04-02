// Package order enables easier ordering and comparison tasks.
//
// This package provides functionality to easily define and apply order on values. It works out of
// the box for most primitive types and their pointer versions, and enable order of any object using
// (three-way comparison) https://en.wikipedia.org/wiki/Three-way_comparison with a given
// `func(T, T) int` function, or by implementing the generic interface: `func (T) Compare(T) int`.
//
// Supported Tasks:
//
// * [x] `Sort` / `SortStable` - sort a slice.
//
// * [x] `Search` - binary search for a value in a slice.
//
// * [x] `MinMax` - get indices of minimal and maximal values of a slice.
//
// * [X] `Is` - get a comparable object for more readable code.
//
// + [x] `Select` - get the K'th greatest value of a slice.
//
// * [x] `IsSorted` / `IsStrictSorted` - check if a slice is sorted.
//
// Types and Values
//
// Order between values can be more forgiving than strict comparison. This library allows sensible
// type conversions. A type `U` can be used in order function of type `T` in the following cases:
//
// * `U` is a pointer (or pointers chain) to a `T`.
//
// * `T` is a pointer (or pointers chain) to a `U`.
//
// * `T` and `U` are of the same kind.
//
// * `T` and `U` are of the same number kind group (int?, uint?, float?, complex?) and `U`'s bits
// number is less or equal to `T`'s bits number.
//
// * `U` and `T` are assignable structs.
//
// Usage
//
// Using this library might be less type safe - because of the usage of interfaces API, and less
// efficient - because of the use of reflection. On the other hand, this library reduce chances for
// errors by providing a well tested code and more readable code. See below how some order tasks
// can be translated to be used by this library.
//
// 	 type person struct {
// 	 	name string
// 	 	age  int
// 	 }
//
// 	 var persons []person
//
// 	 // Sort persons (by name and then by age)
// 	-lessPersons := func(i, j int) bool {
// 	-	nameCmp := strings.Compare(persons[i].name, "joe")
// 	-	if nameCmp == 0 {
// 	-		return persons[i].age < persons[i].age
// 	-	}
// 	-	return nameCmp < 0
// 	-}
// 	-sort.Slice(persons, lessPersons)
// 	+orderPersons := order.By(
// 	+	func(a, b person) int { return strings.Compare(a.name, b.name) },
// 	+	func(a, b person) int { return a.age - b.age },
// 	+)
// 	+orderPersons.Sort(persons)
//
// 	 // Search persons for "joe" at age 42:
// 	-searchPersons := func(int i) bool {
// 	-	nameCmp := strings.Compare(persons[i].name, "joe")
// 	-	if nameCmp == 0 {
// 	-		return persons[i].age >= 42
// 	-	}
// 	-	return nameCmp > 0 {
// 	-}
// 	-i := sort.Search(persons, searchPersons)
//	-// Standard library search does not guarantee equality, we should check:
// 	-if i >= len(persons) || persons[i].name != "joe" || persons[i].age != 42 {
// 	-	i := -1
// 	-}
// 	+i := orderPersons.Search(persons, person{name: "joe", age: 42})
//
// 	 // Another way is that person will implement a `Compare(T) int` method, and the order object
// 	 // will know how to handle it:
// 	+func (p person) Compare(other person) int { ... }
// 	+order.Search(persons, person{name: "joe", age: 42})
//
// 	 // Conditions can also be defined on comparable types:
// 	 var t, start, end time.Time
// 	-if (t.After(start) || t.Equal(start)) && t.Before(end) { ... }
// 	+if isT := order.Is(t); isT.GreaterEqual(start) && isT.Less(end) { ... }
package order

import (
	"fmt"
	"reflect"
	"sort"
)

// By enables ordering values of type T by a given list of three-way comparison functions of the
// form `func(T, T) int`. Each function compares two values (`lhs`, `rhs`) of type T, and returns a
// value `c` of type int as follows:
//
// If lhs >  rhs then c > 0.
// If lhs == rhs then c = 0.
// If lhs <  rhs then c < 0.
//
// The list of functions is used in order to define multiple orderings. When two values are
// compared, the first function is evaluated, if the comparison value is not zero, the value is
// returned. Otherwise, the following function is evaluated until a non-zero value is returned.
// If all the comparison functions returned zero, the returned value is also zero.
func By(fns ...interface{}) Fns {
	if len(fns) == 0 {
		panic("Expected at least one comparison function")
	}
	cmpFns := make(Fns, 0, len(fns))
	for i, fn := range fns {
		cmpFn, err := newFn(reflect.ValueOf(fn))
		if err != nil {
			panic(fmt.Sprintf("Invalid function %d: %s", i, err))
		}
		cmpFns, err = cmpFns.append(cmpFn)
		if err != nil {
			panic(err)
		}
	}
	return cmpFns
}

// Reversed returns a reversed comparison of the original function.
func (fns Fns) Reversed() Fns {
	newFns := make(Fns, len(fns))
	for i := range fns {
		original := fns[i] // Copy.
		newFns[i] = Fn{
			fn: func(lhs, rhs reflect.Value) int { return -original.fn(lhs, rhs) },
			t:  original.t,
		}
	}
	return newFns
}

// Sort sorts a given slice according to the comparison function.
func (fns Fns) Sort(slice interface{}) {
	sort.Slice(slice, fns.less(reflect.ValueOf(slice)))
}

// SortStable sorts a given slice according to the comparison function, while keeping the original
// order of equal elements.
func (fns Fns) SortStable(slice interface{}) {
	sort.SliceStable(slice, fns.less(reflect.ValueOf(slice)))
}

// less return a comparison function for a given slice to be used with sort.Slice and
// sort.SliceStable.
func (fns Fns) less(slice reflect.Value) func(i, j int) bool {
	s := fns.mustSlice(slice)

	return func(i, j int) bool {
		return fns.compare(s.Index(i), s.Index(j)) < 0
	}
}

// Search searches the given slice for a value. The given slice should be sorted relative to the
// comparsion function. It returns an index of an element that is equal to the given value. It
// returns -1 if no element was found that is equal to the given value.
func (fns Fns) Search(slice, value interface{}) int {
	s := fns.mustSlice(reflect.ValueOf(slice))
	v := fns.mustValue(reflect.ValueOf(value))

	start, end := 0, s.Len()-1
	if start > end {
		return -1
	}
	for {
		i := int(uint(start+end) >> 1) // Avoid overflow when computing i.
		cmp := fns.compare(s.Index(i), v)
		switch {
		case cmp == 0: // Found.
			return i
		case start == end: // Not found.
			return -1
		case cmp < 0: // slice[i] < value
			start = i + 1
		default: // slice[i] > value
			end = i - 1
		}
	}
}

// MinMax returns the indices of the minimal and maximal values in the given slice. It returns
// values (-1, -1) if the slice is empty. If there are several minimal/maximal values, this function
// will return the index of the first of them.
func (fns Fns) MinMax(slice interface{}) (min, max int) {
	s := fns.mustSlice(reflect.ValueOf(slice))

	if s.Len() == 0 {
		return -1, -1
	}
	for i := 1; i < s.Len(); i++ {
		if fns.compare(s.Index(min), s.Index(i)) > 0 {
			min = i
		}
		if fns.compare(s.Index(max), s.Index(i)) < 0 {
			max = i
		}
	}
	return
}

// IsSorted returns whether the slice is in an increasing order, according to the comparsion
// function.
//
// To check if a slice is in a decreasing order, it is possible to `fn.Reversed().IsSorted(slice)`.
func (fns Fns) IsSorted(slice interface{}) bool {
	return fns.isSorted(reflect.ValueOf(slice), false)
}

// IsStrictSorted returns whether the slice is in a strictly increasing order, according to the
// comparsion function.
//
// To check if a slice is in a strictly decreasing order, it is possible to
// `fn.Reversed().IsStrictSorted(slice)`.
func (fns Fns) IsStrictSorted(slice interface{}) bool {
	return fns.isSorted(reflect.ValueOf(slice), true)
}

// isSorted checks if the slice is sorted.
func (fns Fns) isSorted(slice reflect.Value, strict bool) bool {
	s := fns.mustSlice(slice)

	for i := s.Len() - 1; i > 0; i-- {
		cmp := fns.compare(s.Index(i-1), s.Index(i))
		if cmp > 0 || (cmp == 0 && strict) {
			return false
		}
	}
	return true
}
