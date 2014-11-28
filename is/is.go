package is

import (
	"fmt"
	"reflect"
	"testing"
)

// TB is an interface that defines an assertion-style
// testing syntax to condense and make more efficient
// the writing of tests. The output of TB is identical
// to the output of the testing framework, allowing any
// scripts, parsers, etc to parse it appropriately.
type TB interface {
	// Equal determines if the two arguments are equal.
	// If not, a failure is printed.
	Equal(actual, expected interface{}) Result
	// NotEqual determines if the two arguments are not equal.
	// If they are, a failure is printed.
	NotEqual(actual, expected interface{}) Result
}

// tb implements TB
type tb struct {
	tb testing.TB
}

// ensure tb implements TB
var _ TB = (*tb)(nil)

// New creates a new object that satisfies the TB interface.
func New(testObj testing.TB) TB {
	return &tb{tb: testObj}
}

// equal tests for equality of two objects.
func equal(actual, expected interface{}) bool {
	return reflect.DeepEqual(actual, expected)
}

// printFailure assembles a string from the arguments
// and decorates it, then prints it.
func printFailure(args ...interface{}) {
	fmt.Println(decorate(fmt.Sprintln(args...)))
}

// Equal determines if the two arguments are equal.
// If not, a failure is printed.
func (t *tb) Equal(actual, expected interface{}) Result {
	r := &result{tb: t.tb}
	if equal(actual, expected) {
		r.success = true
	} else {
		r.success = false
		printFailure(actual, "should be equal to", expected)
	}
	return r
}

// NotEqual determines if the two arguments are not equal.
// If they are, a failure is printed.
func (t *tb) NotEqual(actual, expected interface{}) Result {
	r := &result{tb: t.tb}
	if !equal(actual, expected) {
		r.success = true
	} else {
		r.success = false
		printFailure(actual, "should not be equal to", expected)
	}
	return r
}
