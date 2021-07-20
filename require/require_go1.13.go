// +build go1.13

/*
* CODE GENERATED AUTOMATICALLY WITH github.com/stretchr/testify/_codegen
* THIS FILE MUST NOT BE EDITED BY HAND
 */

package require

import (
	assert "github.com/stretchr/testify/assert"
)

var _ assert.TestingT // in case no function required assert package

// ErrorIs asserts that a specified error is an another error wrapper as defined by go1.13 errors package.
//
//   actualObj, err := SomeFunction()
//   assert.ErrorIs(t, err, ErrNotFound)
//   assert.Nil(t, actualObj)
func ErrorIs(t TestingT, theError error, theTarget error, msgAndArgs ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.ErrorIs(t, theError, theTarget, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// ErrorIsf asserts that a specified error is an another error wrapper as defined by go1.13 errors package.
//
//   actualObj, err := SomeFunction()
//   assert.ErrorIsf(t, err, ErrNotFound, "error message %s", "formatted")
//   assert.Nil(t, actualObj)
func ErrorIsf(t TestingT, theError error, theTarget error, msg string, args ...interface{}) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.ErrorIsf(t, theError, theTarget, msg, args...) {
		return
	}
	t.FailNow()
}
