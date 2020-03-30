package order

import (
	"reflect"
)

// Condition allows comparing a given lhs value.
type Condition struct {
	Fns
	lhs reflect.Value
}

// Is returns a comparable object.
func (fns Fns) Is(lhs interface{}) Condition {
	return Condition{Fns: fns, lhs: fns.mustValue(reflect.ValueOf(lhs))}
}

// Equal tests if the compared lhs object is equal to the given rhs object.
func (c Condition) Equal(rhs interface{}) bool {
	return c.compare(c.lhs, c.mustValue(reflect.ValueOf(rhs))) == 0
}

// NotEqual tests if the compared lhs object is not equal to the given rhs object.
func (c Condition) NotEqual(rhs interface{}) bool {
	return c.compare(c.lhs, c.mustValue(reflect.ValueOf(rhs))) != 0
}

// Greater tests if the lhs object is greater than the given rhs object.
func (c Condition) Greater(rhs interface{}) bool {
	return c.compare(c.lhs, c.mustValue(reflect.ValueOf(rhs))) > 0
}

// GreaterEqual tests if the lhs object is greater than or equal to the given rhs object.
func (c Condition) GreaterEqual(rhs interface{}) bool {
	return c.compare(c.lhs, c.mustValue(reflect.ValueOf(rhs))) >= 0
}

// Less tests if the lhs object is less than the given rhs object.
func (c Condition) Less(rhs interface{}) bool {
	return c.compare(c.lhs, c.mustValue(reflect.ValueOf(rhs))) < 0
}

// LessEqual tests if the lhs object is less than or equal to the given rhs object.
func (c Condition) LessEqual(rhs interface{}) bool {
	return c.compare(c.lhs, c.mustValue(reflect.ValueOf(rhs))) <= 0
}
