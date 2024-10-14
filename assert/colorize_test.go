package assert

import (
	// "errors"
	"fmt"
	"testing"
)

// func TestColorize(t *testing.T) {
// 	wrapT := WrapColorize(t)
// 	assert := New(wrapT)

// 	Equal(wrapT, 1, 2)
// 	Equal(wrapT, "A", "b")
// 	Equal(t, 1, 2)
// 	assert.Equal(3, 4)

// 	Same(wrapT, new(Colorize), new(Colorize))
// 	Same(t, new(Colorize), new(Colorize))
// 	assert.Same(new(Colorize), new(Colorize))

// 	EqualValues(wrapT, uint32(123), int32(124))
// 	EqualValues(t, uint32(123), int32(124))
// 	assert.EqualValues(uint32(123), int32(124))

// 	EqualError(wrapT, errors.New("foo"), "bar")
// 	EqualError(t, errors.New("foo"), "bar")
// 	assert.EqualError(errors.New("foo"), "bar")
// }

func TestColorize_Message(t *testing.T) {
	failureMessage := fmt.Sprintf("Not equal: \n"+
		"expected: %s\n"+
		"actual  : %s%s", "expected", "actual", "diff")
	fmt.Println(WrapColorize(t).Message(failureMessage))
}
