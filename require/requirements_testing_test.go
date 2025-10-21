//go:build !novet
// +build !novet

package require

import (
	"testing"
	"time"
)

func TestTestingTFailNow(t *testing.T) {
	t.Parallel()

	tt := new(testing.T)
	done := make(chan struct{})
	// Run in a separate goroutine to capture the Goexit behavior.
	// This avoid test panics from the unssupported Goexit call in the main test goroutine.
	// Note that this will trigger linter warnings about goroutines in tests. (SA2002)
	go func(tt *testing.T) {
		defer close(done)
		defer func() {
			r := recover() // [runtime.Goexit] does not trigger a panic
			// If we see a panic here, the condition function misbehaved
			Nil(t, r, "Condition function must not panic: %v", r)
		}()
		tt.Errorf("test error")
		tt.FailNow()
		panic("unreachable")
	}(tt)
	<-done
	True(t, tt.Failed(), "testing.T must be marked as failed")
}

func TestEventuallyTestingTFailNow(t *testing.T) {
	tt := new(testing.T)

	count := 0
	done := make(chan struct{})

	// Run Eventually in a separate goroutine to capture the Goexit behavior.
	// This avoid test panics from the unssupported Goexit call in the main test goroutine.
	// Note that this will trigger linter warnings about goroutines in tests. (SA2002)
	go func(tt *testing.T) {
		defer close(done)
		defer func() {
			r := recover() // [runtime.Goexit] does not trigger a panic
			// If we see a panic here, the condition function misbehaved
			Nil(t, r, "Condition function must not panic: %v", r)
		}()
		condition := func() bool {
			// tick assertion failure
			count++
			tt.Error("tick error")
			tt.FailNow()
			panic("unreachable")
		}
		Eventually(tt, condition, 100*time.Millisecond, 20*time.Millisecond)
	}(tt)
	<-done

	True(t, tt.Failed(), "Check must fail")
	Equal(t, 1, count, "Condition function must have been called once")
}
