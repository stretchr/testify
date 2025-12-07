package suite

import (
	"runtime/debug"
	"testing"
)

func recoverAndFailOnPanic(t *testing.T) {
	t.Helper()
	r := recover()
	failOnPanic(t, r)
}

func failOnPanic(t *testing.T, r interface{}) {
	t.Helper()
	if r != nil {
		t.Errorf("test panicked: %v\n%s", r, debug.Stack())
		t.FailNow()
	}
}
