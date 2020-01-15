// +build go1.13

package assert

// ErrorIsf asserts that a function returned an error chain
// and that it is contains the provided error.
// Mimicking errors.Is, providing `nil` as expected error
// is the same as calling NoError
//
//   actualObj, err := SomeFunction()
//   assert.ErrorIsf(t, err,  expectedError, "error message %s", "formatted")
func ErrorIsf(t TestingT, expected error, actual error, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return ErrorIs(t, expected, actual, append([]interface{}{msg}, args...)...)
}
