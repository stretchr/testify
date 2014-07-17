package assert

import (
	"errors"
	"testing"
	"time"
)

func TestImplementsWrapper(t *testing.T) {
	assert := New(new(testing.T))

	if !assert.Implements((*AssertionTesterInterface)(nil), new(AssertionTesterConformingObject)) {
		t.Error("Implements method should return true: AssertionTesterConformingObject implements AssertionTesterInterface")
	}
	if assert.Implements((*AssertionTesterInterface)(nil), new(AssertionTesterNonConformingObject)) {
		t.Error("Implements method should return false: AssertionTesterNonConformingObject does not implements AssertionTesterInterface")
	}
}

func TestIsTypeWrapper(t *testing.T) {
	assert := New(new(testing.T))

	if !assert.IsType(new(AssertionTesterConformingObject), new(AssertionTesterConformingObject)) {
		t.Error("IsType should return true: AssertionTesterConformingObject is the same type as AssertionTesterConformingObject")
	}
	if assert.IsType(new(AssertionTesterConformingObject), new(AssertionTesterNonConformingObject)) {
		t.Error("IsType should return false: AssertionTesterConformingObject is not the same type as AssertionTesterNonConformingObject")
	}

}

func TestEqualWrapper(t *testing.T) {
	assert := New(new(testing.T))

	if !assert.Equal("Hello World", "Hello World") {
		t.Error("Equal should return true")
	}
	if !assert.Equal(123, 123) {
		t.Error("Equal should return true")
	}
	if !assert.Equal(123.5, 123.5) {
		t.Error("Equal should return true")
	}
	if !assert.Equal([]byte("Hello World"), []byte("Hello World")) {
		t.Error("Equal should return true")
	}
	if !assert.Equal(nil, nil) {
		t.Error("Equal should return true")
	}
}

func TestNotNilWrapper(t *testing.T) {
	assert := New(new(testing.T))

	if !assert.NotNil(new(AssertionTesterConformingObject)) {
		t.Error("NotNil should return true: object is not nil")
	}
	if assert.NotNil(nil) {
		t.Error("NotNil should return false: object is nil")
	}

}

func TestNilWrapper(t *testing.T) {
	assert := New(new(testing.T))

	if !assert.Nil(nil) {
		t.Error("Nil should return true: object is nil")
	}
	if assert.Nil(new(AssertionTesterConformingObject)) {
		t.Error("Nil should return false: object is not nil")
	}

}

func TestTrueWrapper(t *testing.T) {
	assert := New(new(testing.T))

	if !assert.True(true) {
		t.Error("True should return true")
	}
	if assert.True(false) {
		t.Error("True should return false")
	}

}

func TestFalseWrapper(t *testing.T) {
	assert := New(new(testing.T))

	if !assert.False(false) {
		t.Error("False should return true")
	}
	if assert.False(true) {
		t.Error("False should return false")
	}

}

func TestExactlyWrapper(t *testing.T) {
	assert := New(new(testing.T))

	a := float32(1)
	b := float64(1)
	c := float32(1)
	d := float32(2)

	if assert.Exactly(a, b) {
		t.Error("Exactly should return false")
	}
	if assert.Exactly(a, d) {
		t.Error("Exactly should return false")
	}
	if !assert.Exactly(a, c) {
		t.Error("Exactly should return true")
	}

	if assert.Exactly(nil, a) {
		t.Error("Exactly should return false")
	}
	if assert.Exactly(a, nil) {
		t.Error("Exactly should return false")
	}

}

func TestNotEqualWrapper(t *testing.T) {

	assert := New(new(testing.T))

	if !assert.NotEqual("Hello World", "Hello World!") {
		t.Error("NotEqual should return true")
	}
	if !assert.NotEqual(123, 1234) {
		t.Error("NotEqual should return true")
	}
	if !assert.NotEqual(123.5, 123.55) {
		t.Error("NotEqual should return true")
	}
	if !assert.NotEqual([]byte("Hello World"), []byte("Hello World!")) {
		t.Error("NotEqual should return true")
	}
	if !assert.NotEqual(nil, new(AssertionTesterConformingObject)) {
		t.Error("NotEqual should return true")
	}
}

func TestContainsWrapper(t *testing.T) {

	assert := New(new(testing.T))

	if !assert.Contains("Hello World", "Hello") {
		t.Error("Contains should return true: \"Hello World\" contains \"Hello\"")
	}
	if assert.Contains("Hello World", "Salut") {
		t.Error("Contains should return false: \"Hello World\" does not contain \"Salut\"")
	}

}

func TestNotContainsWrapper(t *testing.T) {

	assert := New(new(testing.T))

	if !assert.NotContains("Hello World", "Hello!") {
		t.Error("NotContains should return true: \"Hello World\" does not contain \"Hello!\"")
	}
	if assert.NotContains("Hello World", "Hello") {
		t.Error("NotContains should return false: \"Hello World\" contains \"Hello\"")
	}

}

func TestDidPanicWrapper(t *testing.T) {

	if funcDidPanic, _ := didPanic(func() {
		panic("Panic!")
	}); !funcDidPanic {
		t.Error("didPanic should return true")
	}

	if funcDidPanic, _ := didPanic(func() {
	}); funcDidPanic {
		t.Error("didPanic should return false")
	}

}

func TestPanicsWrapper(t *testing.T) {

	assert := New(new(testing.T))

	if !assert.Panics(func() {
		panic("Panic!")
	}) {
		t.Error("Panics should return true")
	}

	if assert.Panics(func() {
	}) {
		t.Error("Panics should return false")
	}

}

func TestNotPanicsWrapper(t *testing.T) {

	assert := New(new(testing.T))

	if !assert.NotPanics(func() {
	}) {
		t.Error("NotPanics should return true")
	}

	if assert.NotPanics(func() {
		panic("Panic!")
	}) {
		t.Error("NotPanics should return false")
	}

}

func TestEqualWrapper_Funcs(t *testing.T) {

	assert := New(t)

	type f func() int
	var f1 f = func() int { return 1 }
	var f2 f = func() int { return 2 }

	var f1_copy f = f1

	assert.Equal(f1_copy, f1, "Funcs are the same and should be considered equal")
	assert.NotEqual(f1, f2, "f1 and f2 are different")

}

func TestNoErrorWrapper(t *testing.T) {
	assert := New(t)
	mockAssert := New(new(testing.T))

	// start with a nil error
	var err error = nil

	assert.True(mockAssert.NoError(err), "NoError should return True for nil arg")

	// now set an error
	err = errors.New("Some error")

	assert.False(mockAssert.NoError(err), "NoError with error should return False")

}

func TestErrorWrapper(t *testing.T) {
	assert := New(t)
	mockAssert := New(new(testing.T))

	// start with a nil error
	var err error = nil

	assert.False(mockAssert.Error(err), "Error should return False for nil arg")

	// now set an error
	err = errors.New("Some error")

	assert.True(mockAssert.Error(err), "Error with error should return True")

}

func TestEqualErrorWrapper(t *testing.T) {
	assert := New(t)
	mockAssert := New(new(testing.T))

	// start with a nil error
	var err error
	assert.False(mockAssert.EqualError(err, ""),
		"EqualError should return false for nil arg")

	// now set an error
	err = errors.New("some error")
	assert.False(mockAssert.EqualError(err, "Not some error"),
		"EqualError should return false for different error string")
	assert.True(mockAssert.EqualError(err, "some error"),
		"EqualError should return true")
}

func TestEmptyWrapper(t *testing.T) {
	assert := New(t)
	mockAssert := New(new(testing.T))

	assert.True(mockAssert.Empty(""), "Empty string is empty")
	assert.True(mockAssert.Empty(nil), "Nil is empty")
	assert.True(mockAssert.Empty([]string{}), "Empty string array is empty")
	assert.True(mockAssert.Empty(0), "Zero int value is empty")
	assert.True(mockAssert.Empty(false), "False value is empty")

	assert.False(mockAssert.Empty("something"), "Non Empty string is not empty")
	assert.False(mockAssert.Empty(errors.New("something")), "Non nil object is not empty")
	assert.False(mockAssert.Empty([]string{"something"}), "Non empty string array is not empty")
	assert.False(mockAssert.Empty(1), "Non-zero int value is not empty")
	assert.False(mockAssert.Empty(true), "True value is not empty")

}

func TestNotEmptyWrapper(t *testing.T) {
	assert := New(t)
	mockAssert := New(new(testing.T))

	assert.False(mockAssert.NotEmpty(""), "Empty string is empty")
	assert.False(mockAssert.NotEmpty(nil), "Nil is empty")
	assert.False(mockAssert.NotEmpty([]string{}), "Empty string array is empty")
	assert.False(mockAssert.NotEmpty(0), "Zero int value is empty")
	assert.False(mockAssert.NotEmpty(false), "False value is empty")

	assert.True(mockAssert.NotEmpty("something"), "Non Empty string is not empty")
	assert.True(mockAssert.NotEmpty(errors.New("something")), "Non nil object is not empty")
	assert.True(mockAssert.NotEmpty([]string{"something"}), "Non empty string array is not empty")
	assert.True(mockAssert.NotEmpty(1), "Non-zero int value is not empty")
	assert.True(mockAssert.NotEmpty(true), "True value is not empty")

}

func TestLenWrapper(t *testing.T) {
	assert := New(t)
	mockAssert := New(new(testing.T))

	assert.False(mockAssert.Len(nil, 0), "nil does not have length")
	assert.False(mockAssert.Len(0, 0), "int does not have length")
	assert.False(mockAssert.Len(true, 0), "true does not have length")
	assert.False(mockAssert.Len(false, 0), "false does not have length")
	assert.False(mockAssert.Len('A', 0), "Rune does not have length")
	assert.False(mockAssert.Len(struct{}{}, 0), "Struct does not have length")

	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3

	cases := []struct {
		v interface{}
		l int
	}{
		{[]int{1, 2, 3}, 3},
		{[...]int{1, 2, 3}, 3},
		{"ABC", 3},
		{map[int]int{1: 2, 2: 4, 3: 6}, 3},
		{ch, 3},

		{[]int{}, 0},
		{map[int]int{}, 0},
		{make(chan int), 0},

		{[]int(nil), 0},
		{map[int]int(nil), 0},
		{(chan int)(nil), 0},
	}

	for _, c := range cases {
		assert.True(mockAssert.Len(c.v, c.l), "%#v have %d items", c.v, c.l)
	}
}

func TestWithinDurationWrapper(t *testing.T) {
	assert := New(t)
	mockAssert := New(new(testing.T))
	a := time.Now()
	b := a.Add(10 * time.Second)

	assert.True(mockAssert.WithinDuration(a, b, 10*time.Second), "A 10s difference is within a 10s time difference")
	assert.True(mockAssert.WithinDuration(b, a, 10*time.Second), "A 10s difference is within a 10s time difference")

	assert.False(mockAssert.WithinDuration(a, b, 9*time.Second), "A 10s difference is not within a 9s time difference")
	assert.False(mockAssert.WithinDuration(b, a, 9*time.Second), "A 10s difference is not within a 9s time difference")

	assert.False(mockAssert.WithinDuration(a, b, -9*time.Second), "A 10s difference is not within a 9s time difference")
	assert.False(mockAssert.WithinDuration(b, a, -9*time.Second), "A 10s difference is not within a 9s time difference")

	assert.False(mockAssert.WithinDuration(a, b, -11*time.Second), "A 10s difference is not within a 9s time difference")
	assert.False(mockAssert.WithinDuration(b, a, -11*time.Second), "A 10s difference is not within a 9s time difference")
}
