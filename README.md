# order

[![Build Status](https://travis-ci.org/posener/order.svg?branch=master)](https://travis-ci.org/posener/order)
[![codecov](https://codecov.io/gh/posener/order/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/order)
[![GoDoc](https://godoc.org/github.com/posener/order?status.svg)](http://godoc.org/github.com/posener/order)
[![goreadme](https://goreadme.herokuapp.com/badge/posener/order.svg)](https://goreadme.herokuapp.com)

Package order enables more readable and easier comparison tasks.

This package provides easy value comparisons tasks (see list below) using a given three-way
comparison function of the form `func(T, T) int`.
[three-way comparison](https://en.wikipedia.org/wiki/Three-way_comparison).

* [X] Condition: get a comparable object for more readable code.

* [x] Sort - sort a slice.

* [x] Search - binary search for a value in a slice.

* [x] MinMax - get minimal and maximal values of a slice.

+ [x] Select - get the K'th greatest value of a slice.

* [x] IsSorted / IsStrictSorted - check if a slice is sorted.

The order library allowes sensible type conversions. A type `U` can be used in order function
of `T` in the following cases:

* `U` is a pointer (or pointers chain) to a `T`.

* `T` is a pointer (or pointers chain) to a `U`.

* `T` and `U` are of the same kind.

* `T` and `U` are of the same number kind group (int?, uint?, float?, complex?) and `U`'s bits
number is less or equal to `T`'s bits number.

* `U` and `T` are assignable structs.

The Go standard library provides some comparison functions, like `strings.Compare`,
`bytes.Compare`, `(time.Time).After`, and so forth. Using these functions is not that readable as
using operators such as `==`, `>`, etc. This library provides some helper functions that makes
Go code more readable.

```diff
 // Compare strings:
-if strings.Compare("a", "b") < 0 { ... }
+if order.Is("a").Less("b") { ... }

 // Compare times:
-if (a.After(b) || a.Equal(b)) && a.Before(c) { ... }
+if is := order.Is(a); is.GreaterEqual(b) && is.Less(c) { ... }

 // Sort persons (by name and then by age)
 type person struct {
 	name string
 	age  int
 }
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
```

If person was to implement:

```go
func (p person) Compare(other person) int { ... }
```

Order will use the comparison function when functions are called. For example:

```go
order.Sort(persons)
```

#### Examples

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

##### Comparable

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

##### Complex

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

##### SliceOperations

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

Created by [goreadme](https://github.com/apps/goreadme)
