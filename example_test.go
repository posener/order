package order

import (
	"fmt"
	"strings"
	"time"
)

// A simple example that shows how to use the order library with different basic types.
func Example() {
	// The order function can be used to check values equality:
	fmt.Println("now > one-second-ago ?",
		Is(time.Now()).Greater(time.Now().Add(-time.Second)))
	fmt.Println("foo == bar ?",
		Is("foo").Equal("bar"))

	// Checking if a value is within a range:
	if is := Is(3); is.GreaterEqual(3) && is.Less(4) {
		fmt.Println("3 is in [3,4)")
	}

	// Output:
	// now > one-second-ago ? true
	// foo == bar ? false
	// 3 is in [3,4)
}

func Example_sliceOperations() {
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

	// Output:
	// Sorted: [1 2 3]
	// Index of 2: 1
	// Min: 1, max: 3
	// Median: 2
}

// An example of ordering struct with multiple fields with different priorities.
func Example_complex() {
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

	// Output:
	// Reversed: [{Foo 10} {Bar 11} {Bar 10}]
	// Index of {Foo 10}: 0
}

// Define a custom type that implements `func (t T) Compare(other T) int`
type orange int

func (o orange) Compare(other orange) int { return int(o - other) }

// A type may implement a `func (t T) Compare(other T) int` function. In this case it could be just
// used with the order package functions.
func Example_comparable() {
	oranges := []orange{5, 2, 24}
	Sort(oranges)
	fmt.Println(oranges)

	// Output: [2 5 24]
}
