// +build go1.13

/*
* CODE GENERATED AUTOMATICALLY WITH github.com/stretchr/testify/_codegen
* THIS FILE MUST NOT BE EDITED BY HAND
 */

package assert

// ErrorIs asserts that a specified error is an another error wrapper as defined by go1.13 errors package.
//
//   actualObj, err := SomeFunction()
//   a.ErrorIs(err, ErrNotFound)
//   assert.Nil(t, actualObj)
func (a *Assertions) ErrorIs(theError error, theTarget error, msgAndArgs ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return ErrorIs(a.t, theError, theTarget, msgAndArgs...)
}

// ErrorIsf asserts that a specified error is an another error wrapper as defined by go1.13 errors package.
//
//   actualObj, err := SomeFunction()
//   a.ErrorIsf(err, ErrNotFound, "error message %s", "formatted")
//   assert.Nil(t, actualObj)
func (a *Assertions) ErrorIsf(theError error, theTarget error, msg string, args ...interface{}) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return ErrorIsf(a.t, theError, theTarget, msg, args...)
}
