// +build go1.13

package assert

// ErrorIs asserts that a function returned an error chain
// and that it is contains the provided error.
// Mimicking errors.Is, providing `nil` as expected error
// is the same as calling NoError
//
//   actualObj, err := SomeFunction()
//   a.ErrorIs(err,  expectedError)
func (a *Assertions) ErrorIs(expected error, actual error, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return ErrorIs(a.t, expected, actual, msgAndArgs...)
}

// ErrorIsf asserts that a function returned an error chain
// and that it is contains the provided error.
// Mimicking errors.Is, providing `nil` as expected error
// is the same as calling NoError
//
//   actualObj, err := SomeFunction()
//   a.ErrorIsf(err,  expectedError, "error message %s", "formatted")
func (a *Assertions) ErrorIsf(expected error, actual error, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return ErrorIsf(a.t, expected, actual, msg, args...)
}
