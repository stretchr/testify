package assert

import "testing"

// Verifies that invalid patterns no longer cause a panic when using Regexp/NotRegexp.
// Instead, the assertion should fail and return false.
func TestRegexp_InvalidPattern_NoPanic(t *testing.T) {
    NotPanics(t, func() {
        mockT := new(testing.T)
        False(t, Regexp(mockT, "\\C", "whatever"))
    })
}

func TestNotRegexp_InvalidPattern_NoPanic(t *testing.T) {
    NotPanics(t, func() {
        mockT := new(testing.T)
        False(t, NotRegexp(mockT, "\\C", "whatever"))
    })
}
