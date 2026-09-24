//go:build go1.27

// Code generated with github.com/stretchr/testify/_codegen; DO NOT EDIT.

package require

import (
	assert "github.com/stretchr/testify/assert"
)

// ErrorAsType asserts that at least one of the errors in err's tree matches
// type E, using errors.AsType. On success it returns the matched error value.
// ErrorAsType avoids the need for a pre-declared target variable.
//
//	require.ErrorAsType[*json.SyntaxError](t, err)
func ErrorAsType[E error](t TestingT, err error, msgAndArgs ...any) E {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	result0, success := assert.ErrorAsType[E](t, err, msgAndArgs...)
	if !success {
		t.FailNow()
	}
	return result0
}

// ErrorAsTypef asserts that at least one of the errors in err's tree matches
// type E, using errors.AsType. On success it returns the matched error value.
// ErrorAsTypef avoids the need for a pre-declared target variable.
//
//	require.ErrorAsTypef[*json.SyntaxError](t, err, "error message %s", "formatted")
func ErrorAsTypef[E error](t TestingT, err error, msg string, args ...any) E {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	result0, success := assert.ErrorAsTypef[E](t, err, msg, args...)
	if !success {
		t.FailNow()
	}
	return result0
}

// NotErrorAsType asserts that no error in err's tree matches type E.
func NotErrorAsType[E error](t TestingT, err error, msgAndArgs ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.NotErrorAsType[E](t, err, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// NotErrorAsTypef asserts that no error in err's tree matches type E.
func NotErrorAsTypef[E error](t TestingT, err error, msg string, args ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.NotErrorAsTypef[E](t, err, msg, args...) {
		return
	}
	t.FailNow()
}
