package reflectutil

import (
	"fmt"
	"reflect"
	"strings"
)

// T represents any type T.
type T struct {
	// The underlying non-pointer type T.
	reflect.Type
	// Counts how many times the given type was pointing on an underlying non-pointer type T.
	ptrCount int
}

func (t T) String() string {
	return strings.Repeat("*", t.ptrCount) + t.Type.String()
}

// New returns a T for the given type. Not all types are supported as T, if the type is not
// supported this function will return an error.
func New(tp reflect.Type) (T, error) {
	var t T
loop:
	for {
		switch tp.Kind() {
		case reflect.Ptr:
			// If the type is a pointer, get the enderlying type and increment the pointer counter.
			tp = tp.Elem()
			t.ptrCount++
		case reflect.Slice:
			// Only allow slice of []byte.
			if tp.Elem().Kind() == reflect.Uint8 {
				break loop
			}
			return t, fmt.Errorf("slice (besides []byte) is not supported for T.")
		case reflect.Array, reflect.Map, reflect.Func:
			return t, fmt.Errorf("%v is not supported for T.", tp.Kind())
		default:
			break loop
		}
	}
	t.Type = tp
	return t, nil
}

// Convert returns the given value as T. If the conversion is not possible, it returns false as the
// second argument. It panics when the value can't be converted.
func (t T) Convert(v reflect.Value) reflect.Value {
	ok := t.convert(v.Type(), &v)
	if !ok {
		panic(fmt.Sprintf("type %v can't be converted to: %v", v.Type(), t.Type))
	}
	return v
}

// Check if another type is convertable to T.
func (t T) Check(tp reflect.Type) bool {
	return t.convert(tp, nil)
}

// converts checks if src can be converted to T and applies the conversion on v if given.
func (t T) convert(src reflect.Type, v *reflect.Value) (ok bool) {
	dst := t.Type
	// If the conversion was successful set v to be a pointer to T according to the T.ptrCount.
	defer func() {
		if !ok || v == nil {
			return
		}
		for i := 0; i < t.ptrCount; i++ {
			*v = ptrTo(*v)
		}
	}()
	for {
		switch {
		case src == dst:
			// Exactly the same types.
			ok = true
			return
		case kindConversionAllowed(src, dst):
			// The conversion between src to dst is allowed.
			if v != nil {
				*v = v.Convert(dst)
			}
			ok = true
			return
		case src.Kind() == reflect.Ptr:
			// src might be a pointer to dst, take the underlying object and look for dst.
			if v != nil {
				*v = v.Elem()
			}
			src = src.Elem()
		default:
			return
		}
	}
}

// kindConversionAllowed checks if the conversion from src to dst is allowed.
func kindConversionAllowed(src reflect.Type, dst reflect.Type) bool {
	// If the same kind return true, with an exception for struct in which src should be
	// convertable to dst.
	if src.Kind() == dst.Kind() && (dst.Kind() != reflect.Struct || src.ConvertibleTo(dst)) {
		return true
	}

	// For numerical kinds, allow converting the same numerical group where dst has number of bits
	// greater or equal to src.
	srcKindGroup := numKindOf(src.Kind())
	dstKindGroup := numKindOf(dst.Kind())
	return srcKindGroup != numNot && srcKindGroup == dstKindGroup && src.Bits() <= dst.Bits()
}

// numKind represents a group of numerical kinds.
type numKind int

const (
	numNot     numKind = iota // Not a number
	numInt                    // int* kinds.
	numUint                   // uint* kinds.
	numFloat                  // float* kinds.
	numComplex                // complex* kinds.
)

// numKindOf returns the numerical kind of a given kind.
func numKindOf(k reflect.Kind) numKind {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return numInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return numUint
	case reflect.Float32, reflect.Float64:
		return numFloat
	case reflect.Complex64, reflect.Complex128:
		return numFloat
	default:
		return numNot
	}
}

// ptrTo returns a value which is the pointer to the given value.
func ptrTo(v reflect.Value) reflect.Value {
	p := reflect.New(v.Type())
	p.Elem().Set(v)
	return p
}
