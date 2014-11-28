package is

import "testing"

// Result is an interface that defines methods for
// acting on the result of a test.
type Result interface {
	// Success returns the success value of the last test.
	Success() bool
	// Require causes the testing procedure to abort
	// if the test failed.
	Require()
}

// result implements Result
type result struct {
	success bool
	tb      testing.TB
}

// ensure result implements Result
var _ Result = (*result)(nil)

// Success returns the success value of the test that
// returned it.
func (r *result) Success() bool {
	return r.success
}

// Require causes the testing procedure to abort
// if the test failed.
func (r *result) Require() {
	if r.success == false {
		r.tb.FailNow()
	}
}
