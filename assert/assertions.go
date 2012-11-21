package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// Comparison a custom function that returns true on success and false on failure
type Comparison func() (success bool)

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

	// Last ditch effort
	if fmt.Sprintf("%#v", a) == fmt.Sprintf("%#v", b) {
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

	file := ""
	line := 0
	ok := false

	for i := 0; ; i++ {
		_, file, line, ok = runtime.Caller(i)
		if !ok {
			return ""
		}
		parts := strings.Split(file, "/")
		dir := parts[len(parts)-2]
		file = parts[len(parts)-1]
		if (dir != "assert" && dir != "mock") || file == "mock_test.go" {
			break
		}
	}

	return fmt.Sprintf("%s:%d", file, line)
}

// getWhitespaceString returns a string that is long enough to overwrite the default
// output from the go testing framework.
func getWhitespaceString() string {

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	return strings.Repeat(" ", len(fmt.Sprintf("%s:%d:      ", file, line)))

}

func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		return msgAndArgs[0].(string)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

// Implements asserts that an object is implemented by the specified interface.
//
//    assert.Implements(t, (*MyInterface)(nil), new(MyObject), "MyObject")
func Implements(t *testing.T, interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	interfaceType := reflect.TypeOf(interfaceObject).Elem()

	if !reflect.TypeOf(object).Implements(interfaceType) {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tObject must implement: %v\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), interfaceType, message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tObject must implement: %v\n\r", getWhitespaceString(), CallerInfo(), interfaceType)
		}

		return false
	}

	return true

}

// IsType asserts that the specified objects are of the same type.
func IsType(t *testing.T, expectedType interface{}, object interface{}, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if !ObjectsAreEqual(reflect.TypeOf(object), reflect.TypeOf(expectedType)) {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tObject expected to be of type %v, but was %v\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), reflect.TypeOf(expectedType), reflect.TypeOf(object), message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tObject expected to be of type %v, but was %v\n\r", getWhitespaceString(), CallerInfo(), reflect.TypeOf(expectedType), reflect.TypeOf(object))
		}

		return false
	}

	return true
}

// Equal asserts that two objects are equal.
//
//    assert.Equal(t, 123, 123, "123 and 123 should be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func Equal(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if !ObjectsAreEqual(a, b) {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tNot equal: %#v != %#v\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), a, b, message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tNot equal: %#v != %#v\n\r", getWhitespaceString(), CallerInfo(), a, b)
		}

		return false
	}

	return true

}

// NotNil asserts that the specified object is not nil.
//
//    assert.NotNil(t, err, "err should be something")
//
// Returns whether the assertion was successful (true) or not (false).
func NotNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	var success bool = true

	if object == nil {
		success = false
	} else if reflect.ValueOf(object).IsNil() {
		success = false
	}

	if !success {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tExpected not to be nil.\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tExpected not to be nil.\n\r", getWhitespaceString(), CallerInfo())
		}

	}

	return success
}

// Nil asserts that the specified object is nil.
//
//    assert.Nil(t, err, "err should be nothing")
//
// Returns whether the assertion was successful (true) or not (false).
func Nil(t *testing.T, object interface{}, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if object == nil {
		return true
	} else if reflect.ValueOf(object).IsNil() {
		return true
	}

	if len(message) > 0 {
		t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tExpected nil, but got: %#v\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), object, message)
	} else {
		t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tExpected nil, but got: %#v\n\r", getWhitespaceString(), CallerInfo(), object)
	}

	return false
}

// True asserts that the specified value is true.
//
//    assert.True(t, myBool, "myBool should be true")
//
// Returns whether the assertion was successful (true) or not (false).
func True(t *testing.T, value bool, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if value != true {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tShould be true\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tShould be true\n\r", getWhitespaceString(), CallerInfo())
		}

		return false
	}

	return true

}

// False asserts that the specified value is true.
//
//    assert.False(t, myBool, "myBool should be false")
//
// Returns whether the assertion was successful (true) or not (false).
func False(t *testing.T, value bool, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if value != false {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tShould be false\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tShould be false\n\r", getWhitespaceString(), CallerInfo())
		}

		return false
	}

	return true

}

// NotEqual asserts that the specified values are NOT equal.
//
//    assert.NotEqual(t, obj1, obj2, "two objects shouldn't be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func NotEqual(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if ObjectsAreEqual(a, b) {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tShould not be equal.\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tShould not be equal.\n\r", getWhitespaceString(), CallerInfo())
		}

		return false
	}

	return true

}

// Contains asserts that the specified string contains the specified substring.
//
//    assert.Contains(t, "Hello World", "World", "But 'Hello World' does contain 'World'")
//
// Returns whether the assertion was successful (true) or not (false).
func Contains(t *testing.T, s, contains string, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if !strings.Contains(s, contains) {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\t\"%s\" does not contain \"%s\"\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), s, contains, message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\t\"%s\" does not contain \"%s\"\n\r", getWhitespaceString(), CallerInfo(), s, contains)
		}

		return false
	}

	return true

}

// NotContains asserts that the specified string does NOT contain the specified substring.
//
//    assert.NotContains(t, "Hello World", "Earth", "But 'Hello World' does NOT contain 'Earth'")
//
// Returns whether the assertion was successful (true) or not (false).
func NotContains(t *testing.T, s, contains string, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if strings.Contains(s, contains) {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\t\"%s\" should not contain \"%s\"\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), s, contains, message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\t\"%s\" should not contain \"%s\"\n\r", getWhitespaceString(), CallerInfo(), s, contains)
		}

		return false
	}

	return true

}

// Uses a Comparison to assert a complex condition.
func Condition(t *testing.T, comp Comparison, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	result := comp()
	if !result {
		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tCondition failed!\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tCondition failed!\n\r", getWhitespaceString(), CallerInfo())
		}
	}
	return result
}

// PanicTestFunc defines a func that should be passed to the assert.Panics and assert.NotPanics
// methods, and represents a simple func that takes no arguments, and returns nothing.
type PanicTestFunc func()

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f PanicTestFunc) (bool, interface{}) {

	var didPanic bool = false
	var message interface{}
	func() {

		defer func() {
			if message = recover(); message != nil {
				didPanic = true
			}
		}()

		// call the target function
		f()

	}()

	return didPanic, message

}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   assert.Panics(t, func(){
//     GoCrazy()
//   }, "Calling GoCrazy() should panic")
//
// Returns whether the assertion was successful (true) or not (false).
func Panics(t *testing.T, f PanicTestFunc, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if funcDidPanic, panicValue := didPanic(f); !funcDidPanic {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tfunc %#v should panic\n\r\tPanic value:\t%v\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), f, panicValue, message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tfunc %#v should panic\n\r\tPanic value:\t%v\n\r", getWhitespaceString(), CallerInfo(), f, panicValue)
		}
		return false
	}

	return true
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   assert.NotPanics(t, func(){
//     RemainCalm()
//   }, "Calling RemainCalm() should NOT panic")
//
// Returns whether the assertion was successful (true) or not (false).
func NotPanics(t *testing.T, f PanicTestFunc, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if funcDidPanic, panicValue := didPanic(f); funcDidPanic {

		if len(message) > 0 {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tfunc %#v should not panic\n\r\tPanic value:\t%v\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), f, panicValue, message)
		} else {
			t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\tfunc %#v should not panic\n\r\tPanic value:\t%v\n\r", getWhitespaceString(), CallerInfo(), f, panicValue)
		}

		return false
	}

	return true
}
