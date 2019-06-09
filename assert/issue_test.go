package assert

import (
	"runtime/debug"
	"testing"
)

func TestEqual_Issue(t *testing.T) {

	var stack []byte
	recovery := func() {
		if r := recover(); r != nil {
			stack = debug.Stack()
			panic(r)
		}
	}

	t.Run("User type, string", func(t *testing.T) {

		type MyType string

		const (
			One   MyType = "one"
			Two   MyType = "two"
			Three MyType = "three"
			Four  MyType = "four"
		)

		tested := &testing.T{}
		assert, assertT := New(t), New(tested)

		assertT.Equal(One, One)
		assert.False(tested.Failed(), "this is ok")

		assertT.Equal(string(One), string(Two))
		assert.True(tested.Failed(), "should just fail")
		if !assert.NotPanics(func() {
			defer recovery()
			assertT.Equal(Three, Four)
		}, "shouldn't panic!!") {
			t.Logf("watchout in `assertions.go:1352`!!\n%s", stack)
		}
	})

	t.Run("User type, int", func(t *testing.T) {

		type MyType int

		const (
			One   MyType = 1
			Two   MyType = 2
			Three MyType = 3
			Four  MyType = 4
		)

		tested := &testing.T{}
		assert, assertT := New(t), New(tested)

		assertT.Equal(One, One)
		assert.False(tested.Failed(), "this is ok")

		assertT.Equal(int(One), int(Two))
		assert.True(tested.Failed(), "should just fail")

		assert.NotPanics(func() { assertT.Equal(Three, Four) }, "shouldn't panic!!")
	})
}
