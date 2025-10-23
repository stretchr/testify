package assert

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestEventuallyFailsFast(t *testing.T) {

	type testCase struct {
		name      string
		run       func(t TestingT, tc testCase, completed *bool)
		fn        func(t TestingT) // optional
		exit      func()
		ret       bool
		expErrors []string
	}

	runFnAndExit := func(t TestingT, tc testCase, completed *bool) {
		if tc.fn != nil {
			tc.fn(t)
		}
		if tc.exit != nil {
			tc.exit()
		}
		*completed = true
	}

	evtl := func(t TestingT, tc testCase, completed *bool) {
		Eventually(t, func() bool {
			runFnAndExit(t, tc, completed)
			return tc.ret
		}, time.Hour, time.Millisecond)
	}

	withT := func(t TestingT, tc testCase, completed *bool) {
		EventuallyWithT(t, func(collect *CollectT) {
			runFnAndExit(collect, tc, completed)
		}, time.Hour, time.Millisecond)
	}

	never := func(t TestingT, tc testCase, completed *bool) {
		Never(t, func() bool {
			runFnAndExit(t, tc, completed)
			return tc.ret
		}, time.Hour, time.Millisecond)
	}

	doFail := func(t TestingT) { t.Errorf("failed") }
	exitErr := "Condition exited unexpectedly" // fail fast err on runtime.Goexit
	satisErr := "Condition satisfied"          // fail fast err on condition satisfied
	failedErr := "failed"                      // additional error from explicit failure

	cases := []testCase{
		// Fast path Eventually tests
		{
			name: "Satisfy", run: evtl,
			fn: nil, exit: nil, ret: true, // succeed fast
			expErrors: nil, // no errors expected
		},
		{
			name: "Fail", run: evtl,
			fn: doFail, exit: nil, ret: true, // fail and succeed fast
			expErrors: []string{failedErr}, // expect fail
		},
		{
			// Simulate [testing.T.FailNow], which calls
			// [testing.T.Fail] followed by [runtime.Goexit].
			name: "FailNow", run: evtl,
			fn: doFail, exit: runtime.Goexit, ret: false, // no succeed fast, but fail
			expErrors: []string{exitErr, failedErr}, // expect both errors
		},
		{
			name: "Goexit", run: evtl,
			fn: nil, exit: runtime.Goexit, ret: false, // no succeed fast, just exit
			expErrors: []string{exitErr}, // expect exit error
		},

		// Fast path EventuallyWithT tests
		{
			name: "SatisfyWithT", run: withT,
			fn: nil, exit: nil, ret: true, // succeed fast
			expErrors: nil, // no errors expected
		},
		{
			name: "GoExitWithT", run: withT,
			fn: nil, exit: runtime.Goexit, ret: false, // no succeed fast, just exit
			expErrors: []string{exitErr}, // expect exit error
		},
		// EventuallyWithT only fails fast when no errors are collected.
		// The Fail and FailNow cases are thus equivalent and will not fail fast and are not tested here.

		// Fast path Never tests
		{
			name: "SatisfyNever", run: never,
			fn: nil, exit: nil, ret: true, // fail fast by satisfying
			expErrors: []string{satisErr}, // expect satisfy error only
		},
		{
			name: "FailNowNever", run: never,
			fn: doFail, exit: runtime.Goexit, ret: false, // no satisfy, but fail + exit
			expErrors: []string{exitErr, failedErr}, // expect both errors
		},
		{
			name: "GoexitNever", run: never,
			fn: nil, exit: runtime.Goexit, ret: false, // no satisfy, just exit
			expErrors: []string{exitErr}, // expect exit error
		},
		{
			name: "FailNever", run: never,
			fn: doFail, exit: nil, ret: true, // fail then satisfy fast
			expErrors: []string{failedErr, satisErr}, // expect fail error
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			collT := &CollectT{}
			wait := make(chan struct{})
			completed := false
			var panicValue interface{}

			go func() {
				defer func() {
					panicValue = recover()
					close(wait)
				}()
				tc.run(collT, tc, &completed)
			}()

			select {
			case <-wait:
			case <-time.After(time.Second):
				FailNow(t, "test did not complete within timeout")
			}

			expFail := len(tc.expErrors) > 0

			Nil(t, panicValue, "Eventually should not panic")
			Equal(t, expFail, collT.failed(), "test state does not match expected failed state")
			Len(t, collT.errors, len(tc.expErrors), "number of collected errors does not match expectation")

		Found:
			for _, expMsg := range tc.expErrors {
				for _, err := range collT.errors {
					if strings.Contains(err.Error(), expMsg) {
						continue Found
					}
				}
				t.Errorf("expected error message %q not found in collected errors", expMsg)
			}

		})
	}
}

func TestEventuallyCompletes(t *testing.T) {
	t.Parallel()
	mockT := &mockTestingT{}
	Eventually(mockT, func() bool {
		return true
	}, time.Second, time.Millisecond)
	False(t, mockT.Failed(), "test should not fail")

	mockT = &mockTestingT{}
	EventuallyWithT(mockT, func(collect *CollectT) {
		// no assertion failures
	}, time.Second, time.Millisecond)
	False(t, mockT.Failed(), "test should not fail")

	mockT = &mockTestingT{}
	Never(mockT, func() bool {
		return false
	}, time.Second, time.Millisecond)
	False(t, mockT.Failed(), "test should not fail")
}

func TestEventuallyHandlesUnexpectedExit(t *testing.T) {
	t.Parallel()
	collT := &CollectT{}
	Eventually(collT, func() bool {
		runtime.Goexit()
		panic("unreachable")
	}, time.Second, time.Millisecond)
	True(t, collT.failed(), "test should fail")
	Len(t, collT.errors, 1, "should have one error")
	Contains(t, collT.errors[0].Error(), "Condition exited unexpectedly")

	collT = &CollectT{}
	EventuallyWithT(collT, func(collect *CollectT) {
		runtime.Goexit()
		panic("unreachable")
	}, time.Second, time.Millisecond)
	True(t, collT.failed(), "test should fail")
	Len(t, collT.errors, 1, "should have one error")
	Contains(t, collT.errors[0].Error(), "Condition exited unexpectedly")

	collT = &CollectT{}
	Never(collT, func() bool {
		runtime.Goexit()
		panic("unreachable")
	}, time.Second, time.Millisecond)
	True(t, collT.failed(), "test should fail")
	Len(t, collT.errors, 1, "should have one error")
	Contains(t, collT.errors[0].Error(), "Condition exited unexpectedly")
}

func TestPanicInEventuallyNotRecovered(t *testing.T) {
	testPanicUnrecoverable(t, func() {
		Eventually(t, func() bool {
			panic("demo panic")
		}, time.Minute, time.Millisecond)
	})
}

func TestPanicInEventuallyWithTNotRecovered(t *testing.T) {
	testPanicUnrecoverable(t, func() {
		EventuallyWithT(t, func(collect *CollectT) {
			panic("demo panic")
		}, time.Minute, time.Millisecond)
	})
}

func TestPanicInNeverNotRecovered(t *testing.T) {
	testPanicUnrecoverable(t, func() {
		Never(t, func() bool {
			panic("demo panic")
		}, time.Minute, time.Millisecond)
	})
}

// testPanicUnrecoverable ensures current goroutine panic behavior.
//
// Currently, [Eventually] runs the condition function in a separate goroutine.
// If that goroutine panics, the panic is not recovered, and the entire test process
// is terminated.
//
// In the future, this behavior may change so that panics in the condition are caught
// and handled more gracefully. For now we ensure such panics are not unrecoved.
//
// To run this test, set the environment variable TestPanic=1.
// The test is successful if it panics and fails the test process and does NOT print
// "UNREACHABLE CODE!" after the initial log messages.
func testPanicUnrecoverable(t *testing.T, failingDemoTest func()) {
	if os.Getenv("TestPanic") == "" {
		t.Skip("Skipping test, set TestPanic=1 to run")
	}
	// Use fmt.Println instead of t.Log because t.Log output may be suppressed.
	fmt.Println("⚠️ This test must fail by a panic in a goroutine.")
	fmt.Println("⚠️ If you see the text 'UNREACHABLE CODE!' after this point, this means the test exited in an unintended way")
	defer func() {
		// defer statements are not run when a goroutine panics, so this code is
		// only reachable if the panic was somehow recovered.
		fmt.Println("❌ UNREACHABLE CODE!")
		fmt.Println("❌ If you see this, the test has not failed as expected.")
	}()
	failingDemoTest()
}
