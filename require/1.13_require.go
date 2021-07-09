// +build go1.13

package require

import "github.com/stretchr/testify/assert"

// ErrorIs asserts that a function returned an error chain
// and that it is contains the provided error.
// Mimicking errors.Is, providing `nil` as expected error
// is the same as calling NoError
//
//   actualObj, err := SomeFunction()
//   assert.ErrorIs(t, err,  expectedError)
func ErrorIs(t TestingT, expected error, actual error, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.ErrorIs(t, expected, actual, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// ErrorIsf asserts that a function returned an error chain
// and that it is contains the provided error.
// Mimicking errors.Is, providing `nil` as expected error
// is the same as calling NoError
//
//   actualObj, err := SomeFunction()
//   assert.ErrorIsf(t, err,  expectedError, "error message %s", "formatted")
func ErrorIsf(t TestingT, expected error, actual error, msg string, args ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.ErrorIsf(t, expected, actual, msg, args...) {
		return
	}
	t.FailNow()
}
