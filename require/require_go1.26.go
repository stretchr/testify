//go:build go1.26

package require

import (
	assert "github.com/stretchr/testify/assert"
)

// ErrorAsType asserts that at least one of the errors in err's tree matches
// type E, using errors.AsType. On success it returns the matched error value.
// This is a Go 1.26+ generic alternative to ErrorAs that avoids the need for
// a pre-declared target variable.
//
// If the assertion fails, FailNow is called.
func ErrorAsType[E error](t TestingT, err error, msgAndArgs ...any) E {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	target, ok := assert.ErrorAsType[E](t, err, msgAndArgs...)
	if !ok {
		t.FailNow()
	}
	return target
}

// ErrorAsTypef asserts that at least one of the errors in err's tree matches
// type E, using errors.AsType. On success it returns the matched error value.
// This is a Go 1.26+ generic alternative to ErrorAs that avoids the need for
// a pre-declared target variable.
//
// If the assertion fails, FailNow is called.
func ErrorAsTypef[E error](t TestingT, err error, msg string, args ...any) E {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	target, ok := assert.ErrorAsTypef[E](t, err, msg, args...)
	if !ok {
		t.FailNow()
	}
	return target
}

// NotErrorAsType asserts that no error in err's tree matches type E.
// This is a Go 1.26+ generic alternative to NotErrorAs.
//
// If the assertion fails, FailNow is called.
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
// This is a Go 1.26+ generic alternative to NotErrorAs.
//
// If the assertion fails, FailNow is called.
func NotErrorAsTypef[E error](t TestingT, err error, msg string, args ...any) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	if assert.NotErrorAsTypef[E](t, err, msg, args...) {
		return
	}
	t.FailNow()
}
