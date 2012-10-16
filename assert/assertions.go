package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

/*
	Helper functions
*/

// ObjectsAreEqual determines if two objects are considered equal.
//
// This function does no assertion of any kind.
func ObjectsAreEqual(a, b interface{}) bool {

	if reflect.DeepEqual(a, b) {
		return true
	}

	if reflect.ValueOf(a) == reflect.ValueOf(b) {
		return true
	}

	return false

}

/* CallerInfo is necessary because the assert functions use the testing object
internally, causing it to print the file:line of the assert method, rather than where
the problem actually occured in calling code.*/

// CallerInfo returns a string containing the file and line number of the assert call
// that failed. 
func CallerInfo() string {
	_, file, line, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	thisDir := parts[len(parts)-2]

	for i := 1; ; i++ {
		_, file, line, ok = runtime.Caller(i)
		if !ok {
			return ""
		}
		parts = strings.Split(file, "/")
		dir := parts[len(parts)-2]
		file = parts[len(parts)-1]
		if thisDir != dir || file == "assertions_test.go" {
			break
		}
	}

	return fmt.Sprintf("[ %s:%d ] - ", file, line)
}

// Implements asserts that an object is implemented by the specified interface.
//
// Example
//    assert.Implements(t, (*core.Codec)(nil), new(JsonCodec), "JsonCodec")
func Implements(t *testing.T, interfaceObject interface{}, object interface{}, message ...string) bool {
	interfaceType := reflect.TypeOf(interfaceObject).Elem()
	return True(t, reflect.TypeOf(object).Implements(interfaceType), fmt.Sprintf("%sObject must implement %s. %s", CallerInfo(), interfaceType, message))
}

// IsType asserts that the specified object is of the specified type.
func IsType(t *testing.T, expectedType interface{}, object interface{}, message ...string) bool {
	return Equal(t, reflect.TypeOf(object), reflect.TypeOf(expectedType), fmt.Sprintf("Object expected to be of type %s, but was %s. %s", reflect.TypeOf(expectedType), reflect.TypeOf(object), message))
}

// Equal asserts that two objects are equal.
func Equal(t *testing.T, a, b interface{}, message ...string) bool {

	if !ObjectsAreEqual(a, b) {
		t.Errorf("%s%s Not equal. %#v != %#v.", CallerInfo(), message, a, b)
		return false
	}
	return true

}

// NotNil asserts that the specified object is not nil.
func NotNil(t *testing.T, object interface{}, message ...string) bool {

	var success bool = true

	if object == nil {
		success = false
	} else if reflect.ValueOf(object).IsNil() {
		success = false
	}

	if !success {
		t.Errorf("%sExpected not to be nil. %s", CallerInfo(), message)
	}

	return success
}

// Nil asserts that the specified object is nil.
func Nil(t *testing.T, object interface{}, message ...string) bool {

	if object == nil {
		return true
	} else if reflect.ValueOf(object).IsNil() {
		return true
	}

	t.Errorf("%sExpected to be nil but was %#v. %s", CallerInfo(), object, message)

	return false
}

// True asserts that the specified value is true.
func True(t *testing.T, value bool, message ...string) bool {
	return Equal(t, true, value, message...)
}

// False asserts that the specified value is true.
func False(t *testing.T, value bool, message ...string) bool {
	return Equal(t, false, value, message...)
}

// NotEqual asserts that the specified values are NOT equal.
func NotEqual(t *testing.T, a, b interface{}, message ...string) bool {

	if ObjectsAreEqual(a, b) {
		t.Errorf("%s%s Should not be equal. %#v == %#v.", CallerInfo(), message, a, b)
		return false
	}
	return true

}

// Contains asserts that the specified string contains the specified substring.
func Contains(t *testing.T, s, contains string, message ...string) bool {

	if !strings.Contains(s, contains) {
		t.Errorf("%s %s '%s' does not contain '%s'", CallerInfo(), message, s, contains)
		return false
	}

	return true

}

// NotContains asserts that the specified string does NOT contain the specified substring.
func NotContains(t *testing.T, s, contains string, message ...string) bool {

	if strings.Contains(s, contains) {
		t.Errorf("%s%s '%s' should not contain '%s'", CallerInfo(), message, s, contains)
		return false
	}

	return true

}

// PanicTestFunc defines a type that is used for type checking and convenience
type PanicTestFunc func()

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f PanicTestFunc) bool {

	var didPanic bool = false
	func() {

		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()

		// call the target function
		f()

	}()

	return didPanic

}

// Panics asserts that the function passed to it panics
func Panics(t *testing.T, f PanicTestFunc, message ...string) bool {
	return True(t, didPanic(f), fmt.Sprintf("Func should panic but didn't. %s", message))
}

// NotPanics asserts that the function passed to it does not panic
func NotPanics(t *testing.T, f PanicTestFunc, message ...string) bool {
	return False(t, didPanic(f), fmt.Sprintf("Func should not panic. %s", message))
}
