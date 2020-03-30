package order

import (
	"fmt"
	"reflect"

	"github.com/posener/order/internal/reflectutil"
)

// Fns is a list of order functions, used to check the order between two T types.
type Fns []Fn

// Fn represent an order function.
type Fn struct {
	// fns are the 3-way functions, of the form func(T, T) int.
	fn func(lhs, rhs reflect.Value) int
	// t stores the type of the function (T).
	t reflectutil.T
}

// newFn converts a given function value to the a compare function. It also checks that the
// function `f` is of the right form (func(T, T) int) and that T is of the given type t. If the
// given type t is nil, it will be set to the type of the first argument of f.
func newFn(f reflect.Value) (Fn, error) {
	if f.Kind() != reflect.Func {
		return Fn{}, fmt.Errorf("expected function")
	}
	tp := f.Type()
	if in := tp.NumIn(); in != 2 {
		return Fn{}, fmt.Errorf("expected function with 2 arguments, got: %d", in)
	}
	// If t is not set yet, set it to the first argument of the function.
	t1, err := reflectutil.New(tp.In(0))
	if err != nil {
		return Fn{}, err
	}
	t2, err := reflectutil.New(tp.In(1))
	if err != nil {
		return Fn{}, err
	}

	if t1.Type != t2.Type {
		return Fn{}, fmt.Errorf("expected same types, got: %v, %v", t1, t2)
	}
	if out := tp.NumOut(); out != 1 {
		return Fn{}, fmt.Errorf("expected function with a single return value, got: %d", out)
	}
	if out := tp.Out(0); out.Kind() != reflect.Int {
		return Fn{}, fmt.Errorf("expected function with int return value, got: %v", out)
	}
	return Fn{
		fn: func(lhs, rhs reflect.Value) int {
			return f.Call([]reflect.Value{t1.Convert(lhs), t2.Convert(rhs)})[0].Interface().(int)
		},
		t: t1,
	}, nil
}

// compare compares two values using the comparsion functions. It starts from the first comparison
// function and continues as long as the returned value is 0.
func (fns Fns) compare(lhs, rhs reflect.Value) int {
	for _, fn := range fns {
		if cmp := fn.fn(lhs, rhs); cmp != 0 {
			return cmp
		}
	}
	return 0
}

// append a function to the function list, and check that its type agrees with the list type.
func (fns Fns) append(fn Fn) (Fns, error) {
	if len(fns) != 0 {
		if !fns.check(fn.T()) {
			return nil, fmt.Errorf("all functions should have the same type, got: %v, %v", fns.T(), fn.T())
		}
	}
	return append(fns, fn), nil
}

// T returns the type of the functions list T.
func (fns Fns) T() reflect.Type {
	return fns[0].T()
}

// T returns the type of the function T.
func (fn Fn) T() reflect.Type {
	return fn.t.Type
}

func (fns Fns) check(tp reflect.Type) bool {
	return fns[0].t.Check(tp)

}

// mustValue panics if the given value is not of type T.
func (fns Fns) mustValue(v reflect.Value) reflect.Value {
	if tp := v.Type(); !fns.check(tp) {
		panic(fmt.Sprintf("bad value type: expected: %v, got: %v", fns.T(), tp))
	}
	return v
}

// mustSlice panics if a given slice value is not a slice value or does not match T.
func (fns Fns) mustSlice(slice reflect.Value) reflectutil.Slice {
	s, err := reflectutil.NewSlice(slice)
	if err != nil {
		panic(err)
	}
	if tp := s.T(); !fns.check(tp) {
		panic(fmt.Sprintf("wrong slice type: expected []%v, got: %v", fns.T(), tp))
	}
	return s
}
