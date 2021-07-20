// +build go1.13

/*
* CODE GENERATED AUTOMATICALLY WITH github.com/stretchr/testify/_codegen
* THIS FILE MUST NOT BE EDITED BY HAND
 */

package assert

// ErrorIsf asserts that a specified error is an another error wrapper as defined by go1.13 errors package.
//
//   actualObj, err := SomeFunction()
//   assert.ErrorIsf(t, err, ErrNotFound, "error message %s", "formatted")
//   assert.Nil(t, actualObj)
func ErrorIsf(t TestingT, theError error, theTarget error, msg string, args ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return ErrorIs(t, theError, theTarget, append([]interface{}{msg}, args...)...)
}
