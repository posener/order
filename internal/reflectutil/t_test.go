package reflectutil

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	t1 struct{ Field int }
	t2 struct{ Field int }
	u1 struct{ OtherField int }

	myString string
)

var intT, _ = New(reflect.TypeOf(1))

func TestNew_failures(t *testing.T) {
	t.Parallel()

	var err error
	_, err = New(reflect.TypeOf([8]byte{}))
	assert.Error(t, err)
	_, err = New(reflect.TypeOf([8]int{}))
	assert.Error(t, err)
	_, err = New(reflect.TypeOf([]int{}))
	assert.Error(t, err)
	_, err = New(reflect.TypeOf(map[int]int{}))
	assert.Error(t, err)
	_, err = New(reflect.TypeOf(func() {}))
	assert.Error(t, err)
}

func TestConvert_basicTypes(t *testing.T) {
	t.Parallel()

	checkSuccess := func(t *testing.T, src, dst interface{}) {
		t.Helper()

		t1, err := New(reflect.TypeOf(dst))
		require.NoError(t, err)

		got := t1.Convert(reflect.ValueOf(src))

		assert.Equal(t, reflect.TypeOf(dst), got.Type())
		assert.Equal(t, dst, got.Interface())
	}

	checkFailure := func(t *testing.T, src, dst interface{}) {
		t.Helper()
		t1, err := New(reflect.TypeOf(dst))
		require.NoError(t, err)

		assert.Panics(t, func() { t1.Convert(reflect.ValueOf(src)) })
	}

	// Check that numbers with the same kind group can be converted to a bigger bits number and
	// can't be converted to smaller bits numbers.
	for _, values := range [][]interface{}{
		{int8(1), int16(1), int32(1), int64(1)},
		{uint8(1), uint16(1), uint32(1), uint64(1)},
		{float32(1), float64(1)},
		{complex64(1), complex128(1)},
	} {
		for i := range values {
			for j := range values {
				src := values[i]
				dst := values[j]
				t.Run(testName2(src, dst), func(t *testing.T) {
					if i <= j {
						checkSuccess(t, src, dst)
					} else {
						checkFailure(t, src, dst)
					}
				})
			}
		}
	}

	// Check both-way conversions between types.
	for _, values := range [][]interface{}{
		{int(1), intPtr(1)},
		{"a", myString("a"), stringPtr("a"), myStringPtr("a")},
		{t1{42}, t2{42}},
	} {
		for _, src := range values {
			for _, dst := range values {
				t.Run(testName2(src, dst), func(t *testing.T) {
					checkSuccess(t, src, dst)
				})
			}
		}
	}
}

func TestConvert_failures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		// src should not be able to be converted to dst.
		dst, src interface{}
	}{
		{dst: "", src: 1},
		{dst: 1, src: "1"},
		{dst: "", src: []byte("")},
		{dst: []byte(""), src: ""},
		{dst: stringPtr(""), src: 1},
		{dst: 1, src: stringPtr("")},
		{dst: stringPtr(""), src: intPtr(1)},
		{dst: intPtr(1), src: stringPtr("")},
		{dst: t1{}, src: u1{}},
		{dst: u1{}, src: t1{}},
		{dst: "", src: []string{""}},
		{dst: "", src: [1]string{""}},
		{dst: "", src: map[string]string{"": ""}},
		{dst: "", src: func() {}},
	}

	for _, tt := range tests {
		t.Run(testName2(tt.src, tt.dst), func(t *testing.T) {
			t1, err := New(reflect.TypeOf(tt.dst))
			require.NoError(t, err)

			assert.Panics(t, func() { t1.Convert(reflect.ValueOf(tt.src)) })
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func myStringPtr(s string) *myString {
	ms := myString(s)
	return &ms
}

func testName(v interface{}) string         { return fmt.Sprintf("%T(%v)", v, v) }
func testName2(src, dst interface{}) string { return fmt.Sprintf("%T(%v)/%T(%v)", src, src, dst, dst) }
