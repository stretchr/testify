// +build go1.13

package assert

import (
	"errors"
	"fmt"
)

// ErrorIs asserts that a specified error is an another error wrapper as defined by go1.13 errors package.
//
//   actualObj, err := SomeFunction()
//   assert.ErrorIs(t, err, ErrNotFound)
//   assert.Nil(t, actualObj)
func ErrorIs(t TestingT, theError, theTarget error, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if errors.Is(theError, theTarget) {
		return true
	}

	return Fail(t, fmt.Sprintf("Error is not %v, but %v", theTarget, theError), msgAndArgs...)
}
