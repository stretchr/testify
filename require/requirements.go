package require

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
}

type tHelper interface {
	Helper()
}

// ComparisonAssertionFunc is a common function prototype when comparing two values.  Can be useful
// for table driven tests.
type ComparisonAssertionFunc func(TestingT, interface{}, interface{}, ...interface{})

// ValueAssertionFunc is a common function prototype when validating a single value.  Can be useful
// for table driven tests.
type ValueAssertionFunc func(TestingT, interface{}, ...interface{})

// BoolAssertionFunc is a common function prototype when validating a bool value.  Can be useful
// for table driven tests.
type BoolAssertionFunc func(TestingT, bool, ...interface{})

// ErrorAssertionFunc is a common function prototype when validating an error value.  Can be useful
// for table driven tests.
type ErrorAssertionFunc func(TestingT, error, ...interface{})

// ErrorIsFunc returns an [ErrorAssertionFunc] which tests if the error wraps target.
func ErrorIsFor(expectedError error) ErrorAssertionFunc {
	return func(t TestingT, err error, msgsAndArgs ...interface{}) {
		if h, ok := t.(tHelper); ok {
			h.Helper()
		}

		ErrorIs(t, err, expectedError, msgsAndArgs...)
	}
}

//go:generate sh -c "cd ../_codegen && go build && cd - && ../_codegen/_codegen -output-package=require -template=require.go.tmpl -include-format-funcs"
