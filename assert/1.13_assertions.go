// +build go1.13

package assert

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorIs asserts that a function returned an error chain
// and that it is contains the provided error.
// Mimicking errors.Is, providing `nil` as expected error
// is the same as calling NoError
//
//   actualObj, err := SomeFunction()
//   assert.ErrorIs(t, err,  expectedError)
func ErrorIs(t TestingT, expected, actual error, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if expected == nil {
		return NoError(t, actual, msgAndArgs...)
	}
	if !Error(t, actual, msgAndArgs...) {
		return false
	}
	if !errors.Is(actual, expected) {
		return Fail(t, missingInChainMessage(expected, actual), msgAndArgs...)
	}
	return true
}

func missingInChainMessage(expected, actual error) string {
	builder := &strings.Builder{}
	fmt.Fprint(builder, "Expected error not in chain:\n")
	fmt.Fprintf(builder, "expected: %#v\nactual  :", expected)

	for currentErr := actual; currentErr != nil; currentErr = errors.Unwrap(currentErr) {
		fmt.Fprintf(builder, "\n  chain : %#v", currentErr)
	}

	return builder.String()
}
