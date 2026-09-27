//go:build go1.27

package require

// Must calls f and returns its value, requiring that f returned a nil error.
// If f returns a non-nil error the test fails immediately, as with
// [Assertions.NoError].
//
//	a := require.New(t)
//	cfg := a.Must(func() (Config, error) { return LoadConfig(path) })
//
// Because the test is stopped with [testing.T.FailNow], Must must be called
// from the goroutine running the test function.
func (a *Assertions) Must[T any](f func() (T, error), msgAndArgs ...interface{}) T {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Must(a.t, f, msgAndArgs...)
}

// Mustf is like [Assertions.Must] but uses a formatted message.
func (a *Assertions) Mustf[T any](f func() (T, error), msg string, args ...interface{}) T {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return Mustf(a.t, f, msg, args...)
}
