package is

import "testing"

func TestIsEqual(t *testing.T) {
	mockT := &testing.T{}
	is := New(mockT)
	if is == nil {
		t.Fatal("is should not be nil")
	}
	r := is.Equal(1, 1)
	if r == nil {
		t.Fatal("Equal should return a result object")
	}
	if !r.Success() {
		t.Fatal("should be equal")
	}
	if is.Equal(1, 2).Success() {
		t.Fatal("1 is not 2")
	}
	// is.Equal(1,2).Require() cannot be tested as it
	// either fails the test or causes a panic and we cannot
	// mock testing.TB as there is a private method to prevent that.
}

func TestIsNotEqual(t *testing.T) {
	mockT := &testing.T{}
	is := New(mockT)
	if is == nil {
		t.Fatal("is should not be nil")
	}
	r := is.NotEqual(1, 2)
	if r == nil {
		t.Fatal("NotEqual should return a result object")
	}
	if !r.Success() {
		t.Fatal("should not be equal")
	}
	if is.NotEqual(1, 1).Success() {
		t.Fatal("1 is 1")
	}
	// is.Equal(1,2).Require() cannot be tested as it
	// either fails the test or causes a panic and we cannot
	// mock testing.TB as there is a private method to prevent that.
}
