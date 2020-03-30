package order

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var intFn = By(func(a, b int) int { return a - b })

type cmp1 struct{ v int }

func (c cmp1) Compare(other cmp1) int { return c.v - other.v }

type cmp2 struct{ v int }

func (c *cmp2) Compare(other *cmp2) int { return c.v - other.v }

func TestConvertTypes(t *testing.T) {
	t.Parallel()

	type myStr string
	var m1, m2 myStr
	var s1, s2 string

	tests := []struct{ a, b interface{} }{
		{s1, s2},
		{&s1, s2},
		{&s1, &s2},
		{m1, m2},
		{&m1, m2},
		{&m1, &m2},
		{s1, m2},
		{s1, &m2},
		{&s1, m2},
		{&s1, &m2},
		{cmp1{1}, cmp1{1}},
		{&cmp2{1}, &cmp2{1}},
		{&cmp1{1}, cmp1{1}},
	}

	test := func(t *testing.T, a, b interface{}) {
		t.Run(name2(a, b), func(t *testing.T) {
			assert.True(t, Is(a).Equal(b))
			assert.True(t, Is(a).GreaterEqual(b))
			assert.True(t, Is(a).LessEqual(b))
			assert.False(t, Is(a).Greater(a))
			assert.False(t, Is(a).Less(a))
		})
	}

	for _, tt := range tests {
		test(t, tt.a, tt.b)
		test(t, tt.b, tt.a)
	}
}

func TestReversed(t *testing.T) {
	t.Parallel()

	c := intFn.Reversed()

	assert.False(t, c.Is(1).Greater(0))
	assert.False(t, c.Is(1).Greater(1))
	assert.True(t, c.Is(1).Greater(2))
}

func TestSort(t *testing.T) {
	t.Parallel()

	got := []int{2, 3, 1}
	Sort(got)
	assert.Equal(t, []int{1, 2, 3}, got)
}

func TestSortStable(t *testing.T) {
	t.Parallel()

	intp := func(i int) *int { return &i }

	got := []*int{intp(2), intp(2), intp(1)}
	want := []*int{got[2], got[0], got[1]}
	SortStable(got)
	// Check the actual pointers and not the pointer values.
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("Element %d differs", i)
		}
	}
}

func TestSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		slice []int
		value int
		want  int
	}{
		{
			name:  "empty slice",
			slice: []int{},
			value: 1,
			want:  -1,
		},
		{
			name:  "odd size slice, middle value",
			slice: []int{1, 2, 3},
			value: 2,
			want:  1,
		},
		{
			name:  "odd size slice, first value",
			slice: []int{1, 2, 3},
			value: 1,
			want:  0,
		},
		{
			name:  "odd size slice, last value",
			slice: []int{1, 2, 3},
			value: 3,
			want:  2,
		},
		{
			name:  "odd size slice, not found",
			slice: []int{1, 2, 3},
			value: 4,
			want:  -1,
		},
		{
			name:  "even size slice, middle value",
			slice: []int{1, 2, 3, 4},
			value: 2,
			want:  1,
		},
		{
			name:  "even size slice, first value",
			slice: []int{1, 2, 3, 4},
			value: 1,
			want:  0,
		},
		{
			name:  "even size slice, last value",
			slice: []int{1, 2, 3, 4},
			value: 4,
			want:  3,
		},
		{
			name:  "even size slice, not found",
			slice: []int{1, 2, 3, 4},
			value: 5,
			want:  -1,
		},
		{
			name:  "not found within the slice",
			slice: []int{1, 2, 3, 5},
			value: 4,
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Search(tt.slice, tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsSorted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		slice            []int
		wantSorted       bool
		wantStrictSorted bool
	}{
		{
			name:             "empty slice",
			slice:            []int{},
			wantSorted:       true,
			wantStrictSorted: true,
		},
		{
			name:             "one element",
			slice:            []int{1},
			wantSorted:       true,
			wantStrictSorted: true,
		},
		{
			name:             "increasing",
			slice:            []int{1, 5, 5},
			wantSorted:       true,
			wantStrictSorted: false,
		},
		{
			name:             "strictly increasing",
			slice:            []int{1, 5, 10},
			wantSorted:       true,
			wantStrictSorted: true,
		},
		{
			name:             "constant",
			slice:            []int{1, 1, 1},
			wantSorted:       true,
			wantStrictSorted: false,
		},
		{
			name:             "decreasing",
			slice:            []int{10, 5, 5},
			wantSorted:       false,
			wantStrictSorted: false,
		},
		{
			name:             "strictly decreasing",
			slice:            []int{10, 5, 1},
			wantSorted:       false,
			wantStrictSorted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantSorted, IsSorted(tt.slice))
			assert.Equal(t, tt.wantStrictSorted, IsStrictSorted(tt.slice))
		})
	}
}

func TestMinMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		slice    []int
		wantMinI int
		wantMaxI int
	}{
		{
			name:     "empty slice",
			slice:    []int{},
			wantMinI: -1,
			wantMaxI: -1,
		},
		{
			name:     "single value",
			slice:    []int{1},
			wantMinI: 0,
			wantMaxI: 0,
		},
		{
			name:     "get the first minmum/maximum",
			slice:    []int{1, 1, 2, 2},
			wantMinI: 0,
			wantMaxI: 2,
		},
		{
			name:     "slice not sorted",
			slice:    []int{3, 1, 2},
			wantMinI: 1,
			wantMaxI: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMinI, gotMaxI := MinMax(tt.slice)
			assert.Equal(t, tt.wantMinI, gotMinI)
			assert.Equal(t, tt.wantMaxI, gotMaxI)
		})
	}
}

func TestBy_invalidFn(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		fns  []interface{}
	}{
		{
			name: "no functions",
			fns:  []interface{}{},
		},
		{
			name: "not a function",
			fns:  []interface{}{true},
		},
		{
			name: "1 arg",
			fns:  []interface{}{func(a int) int { return 0 }},
		},
		{
			name: "3 arg",
			fns:  []interface{}{func(a, b, c int) int { return 0 }},
		},
		{
			name: "arg types differ",
			fns:  []interface{}{func(a int, b string) int { return 0 }},
		},
		{
			name: "arg invalid types",
			fns:  []interface{}{func(a, b func()) int { return 0 }},
		},
		{
			name: "1st arg invalid types",
			fns:  []interface{}{func(a func(), b int) int { return 0 }},
		},
		{
			name: "2nd arg invalid types",
			fns:  []interface{}{func(a int, b func()) int { return 0 }},
		},
		{
			name: "no return values",
			fns:  []interface{}{func(a, b int) {}},
		},
		{
			name: "2 return values",
			fns:  []interface{}{func(a, b int) (int, int) { return 0, 0 }},
		},
		{
			name: "invalid return value type",
			fns:  []interface{}{func(a, b int) bool { return false }},
		},
		{
			name: "functions type mismatch",
			fns: []interface{}{
				func(a, b int) int { return 0 },
				func(a, b bool) int { return 0 },
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() { By(tt.fns...) })
		})
	}
}

func TestInvalidArgs(t *testing.T) {
	t.Parallel()

	fns := []func(v interface{}){
		func(v interface{}) { intFn.Sort(v) },
		func(v interface{}) { intFn.SortStable(v) },
		func(v interface{}) { intFn.Search(v, 1) },
		func(v interface{}) { intFn.IsSorted(v) },
		func(v interface{}) { intFn.IsStrictSorted(v) },
		func(v interface{}) { intFn.MinMax(v) },
		func(v interface{}) { intFn.Select(v, 0) },
	}

	for _, fn := range fns {
		t.Run("not a slice", func(t *testing.T) { assert.Panics(t, func() { fn(1) }) })
		t.Run("slice of wrong type", func(t *testing.T) { assert.Panics(t, func() { fn([]bool{true}) }) })
		t.Run("slice of invalid type", func(t *testing.T) { assert.Panics(t, func() { fn([]func(){func() {}}) }) })
	}

	// Other wrong arguments

	// Search invalid value type.
	assert.Panics(t, func() { intFn.Search([]int{}, true) })

	// Select K out of bounds.
	assert.Panics(t, func() { Select([]int{1}, -1) })
	assert.Panics(t, func() { Select([]int{1}, 1) })
	assert.Panics(t, func() { Select([]int{}, 0) })
}

func name2(a, b interface{}) string { return fmt.Sprintf("%v(%T)/%v(%T)", a, a, b, b) }
