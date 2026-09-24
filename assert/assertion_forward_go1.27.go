//go:build go1.27

// Code generated with github.com/stretchr/testify/_codegen; DO NOT EDIT.

package assert

// ErrorAsType asserts that at least one of the errors in err's tree matches
// type E, using errors.AsType. On success it returns the matched error value.
// ErrorAsType avoids the need for a pre-declared target variable.
//
//	a.ErrorAsType[*json.SyntaxError](err)
func (a *Assertions) ErrorAsType[E error](err error, msgAndArgs ...any) (E, bool) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return ErrorAsType[E](a.t, err, msgAndArgs...)
}

// ErrorAsTypef asserts that at least one of the errors in err's tree matches
// type E, using errors.AsType. On success it returns the matched error value.
// ErrorAsTypef avoids the need for a pre-declared target variable.
//
//	a.ErrorAsTypef[*json.SyntaxError](err, "error message %s", "formatted")
func (a *Assertions) ErrorAsTypef[E error](err error, msg string, args ...any) (E, bool) {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return ErrorAsTypef[E](a.t, err, msg, args...)
}

// NotErrorAsType asserts that no error in err's tree matches type E.
func (a *Assertions) NotErrorAsType[E error](err error, msgAndArgs ...any) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotErrorAsType[E](a.t, err, msgAndArgs...)
}

// NotErrorAsTypef asserts that no error in err's tree matches type E.
func (a *Assertions) NotErrorAsTypef[E error](err error, msg string, args ...any) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return NotErrorAsTypef[E](a.t, err, msg, args...)
}
