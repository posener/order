# order

[![codecov](https://codecov.io/gh/posener/order/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/order)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/posener/order)

Package order enables easier ordering and comparison tasks.

This package provides functionality to easily define and apply order on values. It works out of
the box for most primitive types and their pointer versions, and enable order of any object using
[three-way comparison](https://en.wikipedia.org/wiki/Three-way_comparison) with a given
`func(T, T) int` function, or by implementing the generic interface: `func (T) Compare(T) int`.

Supported Tasks:

* [x] `Sort` / `SortStable` - sort a slice.

* [x] `Search` - binary search for a value in a slice.

* [x] `MinMax` - get indices of minimal and maximal values of a slice.

* [X] `Is` - get a comparable object for more readable code.

+ [x] `Select` - get the K'th greatest value of a slice.

* [x] `IsSorted` / `IsStrictSorted` - check if a slice is sorted.

## Types and Values

Order between values can be more forgiving than strict comparison. This library allows sensible
type conversions. A type `U` can be used in order function of type `T` in the following cases:

* `U` is a pointer (or pointers chain) to a `T`.

* `T` is a pointer (or pointers chain) to a `U`.

* `T` and `U` are of the same kind.

* `T` and `U` are of the same number kind group (int?, uint?, float?, complex?) and `U`'s bits
number is less or equal to `T`'s bits number.

* `U` and `T` are assignable structs.

## Usage

Using this library might be less type safe - because of the usage of interfaces API, and less
efficient - because of the use of reflection. On the other hand, this library reduce chances for
errors by providing a well tested code and more readable code. See below how some order tasks
can be translated to be used by this library.

```diff
 type person struct {
 	name string
 	age  int
 }

 var persons []person

 // Sort persons (by name and then by age)
-lessPersons := func(i, j int) bool {
-	nameCmp := strings.Compare(persons[i].name, "joe")
-	if nameCmp == 0 {
-		return persons[i].age < persons[i].age
-	}
-	return nameCmp < 0
-}
-sort.Slice(persons, lessPersons)
+orderPersons := order.By(
+	func(a, b person) int { return strings.Compare(a.name, b.name) },
+	func(a, b person) int { return a.age - b.age },
+)
+orderPersons.Sort(persons)

 // Search persons for "joe" at age 42:
-searchPersons := func(int i) bool {
-	nameCmp := strings.Compare(persons[i].name, "joe")
-	if nameCmp == 0 {
-		return persons[i].age >= 42
-	}
-	return nameCmp > 0 {
-}
-i := sort.Search(persons, searchPersons)
-// Standard library search does not guarantee equality, we should check:
-if i >= len(persons) || persons[i].name != "joe" || persons[i].age != 42 {
-	i := -1
-}
+i := orderPersons.Search(persons, person{name: "joe", age: 42})

 // Another way is that person will implement a `Compare(T) int` method, and the order object
 // will know how to handle it:
+func (p person) Compare(other person) int { ... }
+order.Search(persons, person{name: "joe", age: 42})

 // Conditions can also be defined on comparable types:
 var t, start, end time.Time
-if (t.After(start) || t.Equal(start)) && t.Before(end) { ... }
+if isT := order.Is(t); isT.GreaterEqual(start) && isT.Less(end) { ... }
```

## Examples

A simple example that shows how to use the order library with different basic types.

```golang
// The order function can be used to check values equality:
fmt.Println("now > one-second-ago ?",
    Is(time.Now()).Greater(time.Now().Add(-time.Second)))
fmt.Println("foo == bar ?",
    Is("foo").Equal("bar"))

// Checking if a value is within a range:
if is := Is(3); is.GreaterEqual(3) && is.Less(4) {
    fmt.Println("3 is in [3,4)")
}
```

 Output:

```
now > one-second-ago ? true
foo == bar ? false
3 is in [3,4)
```

### Comparable

A type may implement a `func (t T) Compare(other T) int` function. In this case it could be just
used with the order package functions.

```golang
oranges := []orange{5, 2, 24}
Sort(oranges)
fmt.Println(oranges)
```

 Output:

```
[2 5 24]
```

### Complex

An example of ordering struct with multiple fields with different priorities.

```golang
// Define a struct with fields of different types.
type person struct {
    name string
    age  int
}
// Order persons: first by name and then by age - reversed.
orderPersons := By(
    func(a, b person) int { return strings.Compare(a.name, b.name) },
    func(a, b person) int { return a.age - b.age },
).Reversed()

// Sort a list of persons in reversed order.
list := []person{
    {"Bar", 10},
    {"Foo", 10},
    {"Bar", 11},
}
orderPersons.Sort(list)
fmt.Println("Reversed:", list)

// Search for a specific person in the sorted list.
fmt.Println("Index of {Foo 10}:", orderPersons.Search(list, person{"Foo", 10}))
```

 Output:

```
Reversed: [{Foo 10} {Bar 11} {Bar 10}]
Index of {Foo 10}: 0
```

### SliceOperations

```golang
// The order function can be used to sort lists:
list := []int{2, 1, 3}
Sort(list)
fmt.Println("Sorted:", list)

// Values can be looked up in sorted lists using a binary search:
fmt.Println("Index of 2:", Search(list, 2))

// Get the minimal and maximal values:
minI, maxI := MinMax(list)
fmt.Printf("Min: %d, max: %d\n", list[minI], list[maxI])

// Get the k'th greatest value:
Select(list, len(list)/2)
fmt.Printf("Median: %d\n", list[1])
```

 Output:

```
Sorted: [1 2 3]
Index of 2: 1
Min: 1, max: 3
Median: 2
```

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
