package require

import (
	"runtime"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"
)

func TestEventuallyGoexit(t *testing.T) {
	t.Parallel()

	condition := func() bool {
		runtime.Goexit() // require.Fail(t) will also call Goexit internally
		panic("unreachable")
	}

	t.Run("WithoutMessage", func(t *testing.T) {
		outerT := new(MockT) // does not call runtime.Goexit immediately
		Eventually(outerT, condition, 100*time.Millisecond, 20*time.Millisecond)
		True(t, outerT.Failed(), "Check must fail")
		Len(t, outerT.Errors(), 1, "There must be one error recorded")
		err1 := outerT.Errors()[0]
		Contains(t, err1.Error(), "Condition exited unexpectedly", "Error message must mention unexpected exit")
	})

	t.Run("WithMessage", func(t *testing.T) {
		outerT := new(MockT) // does not call runtime.Goexit immediately
		Eventually(outerT, condition, 100*time.Millisecond, 20*time.Millisecond, "error: %s", "details")
		True(t, outerT.Failed(), "Check must fail")
		Len(t, outerT.Errors(), 1, "There must be one error recorded")
		err1 := outerT.Errors()[0]
		Contains(t, err1.Error(), "Condition exited unexpectedly", "Error message must mention unexpected exit")
		Contains(t, err1.Error(), "error: details", "Error message must contain formatted message")
	})
}

func TestEventuallyWithTGoexit(t *testing.T) {
	t.Parallel()

	condition := func(collect *assert.CollectT) {
		runtime.Goexit() // require.Fail(t) will also call Goexit internally
		panic("unreachable")
	}

	t.Run("WithoutMessage", func(t *testing.T) {
		mockT := new(MockT) // does not call runtime.Goexit immediately
		EventuallyWithT(mockT, condition, 100*time.Millisecond, 20*time.Millisecond)
		True(t, mockT.Failed(), "Check must fail")
		Len(t, mockT.Errors(), 1, "There must be one error recorded")
		Contains(t, mockT.Errors()[0].Error(), "Condition exited unexpectedly", "Error message must mention unexpected exit")
	})

	t.Run("WithMessage", func(t *testing.T) {
		mockT := new(MockT) // does not call runtime.Goexit immediately
		EventuallyWithT(mockT, condition, 100*time.Millisecond, 20*time.Millisecond, "error: %s", "details")
		True(t, mockT.Failed(), "Check must fail")
		Len(t, mockT.Errors(), 1, "There must be one error recorded")

		err1 := mockT.Errors()[0]
		Contains(t, err1.Error(), "Condition exited unexpectedly", "Error message must mention unexpected exit")
		Contains(t, err1.Error(), "error: details", "Error message must contain formatted message")
	})
}

func TestEventuallyWithTFail(t *testing.T) {
	t.Parallel()

	outerT := new(MockT)
	condition := func(collect *assert.CollectT) {
		// tick assertion failure
		assert.Fail(collect, "tick error")

		// stop the entire test immediately (outer assertion)
		outerT.FailNow() // MockT does not call Goexit internally
		runtime.Goexit() // so we need to call it here to simulate the behavior
		panic("unreachable")
	}

	EventuallyWithT(outerT, condition, 100*time.Millisecond, 20*time.Millisecond)
	True(t, outerT.Failed(), "Check must fail")
	Len(t, outerT.Errors(), 2, "There must be two errors recorded")
	err1, err2 := outerT.Errors()[0], outerT.Errors()[1]
	Contains(t, err1.Error(), "tick error", "First error must be tick error")
	Contains(t, err2.Error(), "Condition exited unexpectedly", "Second error must mention unexpected exit")
}

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

func TestEventuallyFailNow(t *testing.T) {
	t.Parallel()

	outerT := new(MockT)
	condition := func() bool {
		// tick assertion failure
		assert.Fail(outerT, "tick error")

		// stop the entire test immediately (outer assertion)
		outerT.FailNow() // MockT does not call Goexit internally
		runtime.Goexit() // so we need to call it here to simulate the behavior
		panic("unreachable")
	}

	Eventually(outerT, condition, 100*time.Millisecond, 20*time.Millisecond)
	True(t, outerT.Failed(), "Check must fail")
	True(t, outerT.calledFailNow(), "FailNow must have been called")
	Len(t, outerT.Errors(), 2, "There must be two errors recorded")
	err1, err2 := outerT.Errors()[0], outerT.Errors()[1]
	Contains(t, err1.Error(), "tick error", "First error must be tick error")
	Contains(t, err2.Error(), "Condition exited unexpectedly", "Second error must mention unexpected exit")
}
