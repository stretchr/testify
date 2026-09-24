//go:build go1.27

// Code generated with github.com/stretchr/testify/_codegen; DO NOT EDIT.

package assert

// ErrorAsTypef asserts that at least one of the errors in err's tree matches
// type E, using errors.AsType. On success it returns the matched error value.
// ErrorAsTypef avoids the need for a pre-declared target variable.
//
//	assert.ErrorAsTypef[*json.SyntaxError](t, err, "error message %s", "formatted")
func ErrorAsTypef[E error](t TestingT, err error, msg string, args ...any) (E, bool) {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return ErrorAsType[E](t, err, append([]interface{}{msg}, args...)...)
}

// NotErrorAsTypef asserts that no error in err's tree matches type E.
func NotErrorAsTypef[E error](t TestingT, err error, msg string, args ...any) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return NotErrorAsType[E](t, err, append([]interface{}{msg}, args...)...)
}
