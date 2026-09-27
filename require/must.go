// This file uses generics, which require language version go1.18 or later.
// The go.mod go directive is go1.17 so that testify keeps building on old
// toolchains; a //go:build constraint can only raise the language version on
// Go 1.21 and later, so that is the lowest version this file can be gated on.
//go:build go1.21

package require

import (
	assert "github.com/stretchr/testify/assert"
)

// Must calls f and returns its value, requiring that f returned a nil error.
// If f returns a non-nil error the test fails immediately, as with [NoError].
//
//	cfg := require.Must(t, func() (Config, error) { return LoadConfig(path) })
//
// Because the test is stopped with [testing.T.FailNow], Must must be called
// from the goroutine running the test function.
func Must[T any](t TestingT, f func() (T, error), msgAndArgs ...interface{}) T {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	v, err := f()
	if !assert.NoError(t, err, msgAndArgs...) {
		t.FailNow()
	}
	return v
}

// Mustf is like [Must] but uses a formatted message.
func Mustf[T any](t TestingT, f func() (T, error), msg string, args ...interface{}) T {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	v, err := f()
	if !assert.NoErrorf(t, err, msg, args...) {
		t.FailNow()
	}
	return v
}
