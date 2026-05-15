//go:build go1.26

package assert

import (
	"errors"
	"fmt"
	"reflect"
)

// ErrorAsType asserts that at least one of the errors in err's tree matches
// type E, using errors.AsType. On success it returns the matched error value.
// This is a Go 1.26+ generic alternative to ErrorAs that avoids the need for
// a pre-declared target variable.
//
//	assert.ErrorAsType[*json.SyntaxError](t, err)
func ErrorAsType[E error](t TestingT, err error, msgAndArgs ...any) (E, bool) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if target, ok := errors.AsType[E](err); ok {
		return target, true
	}

	expectedType := reflect.TypeFor[E]().String()
	if err == nil {
		Fail(t, fmt.Sprintf("An error is expected but got nil.\n"+
			"expected: %s", expectedType), msgAndArgs...)
		var zero E
		return zero, false
	}

	chain := buildErrorChainString(err, true)
	Fail(t, fmt.Sprintf("Should be in error chain:\n"+
		"expected: %s\n"+
		"in chain: %s", expectedType, truncatingFormat("%s", chain),
	), msgAndArgs...)
	var zero E
	return zero, false
}

// ErrorAsTypef asserts that at least one of the errors in err's tree matches
// type E, using errors.AsType. On success it returns the matched error value.
// This is a Go 1.26+ generic alternative to ErrorAs that avoids the need for
// a pre-declared target variable.
func ErrorAsTypef[E error](t TestingT, err error, msg string, args ...any) (E, bool) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return ErrorAsType[E](t, err, append([]any{msg}, args...)...)
}

// NotErrorAsType asserts that no error in err's tree matches type E.
// This is a Go 1.26+ generic alternative to NotErrorAs.
func NotErrorAsType[E error](t TestingT, err error, msgAndArgs ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	if _, ok := errors.AsType[E](err); !ok {
		return true
	}

	chain := buildErrorChainString(err, true)
	return Fail(t, fmt.Sprintf("Target error should not be in err chain:\n"+
		"found: %s\n"+
		"in chain: %s", reflect.TypeFor[E]().String(), truncatingFormat("%s", chain),
	), msgAndArgs...)
}

// NotErrorAsTypef asserts that no error in err's tree matches type E.
// This is a Go 1.26+ generic alternative to NotErrorAs.
func NotErrorAsTypef[E error](t TestingT, err error, msg string, args ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NotErrorAsType[E](t, err, append([]any{msg}, args...)...)
}
