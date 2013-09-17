package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
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

// Fail reports a failure through
func Fail(t *testing.T, failureMessage string, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)

	if len(message) > 0 {
		t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\t%s\n\r\tMessages:\t%s\n\r", getWhitespaceString(), CallerInfo(), failureMessage, message)
	} else {
		t.Errorf("\r%s\r\tLocation:\t%s\n\r\tError:\t\t%s\n\r", getWhitespaceString(), CallerInfo(), failureMessage)
	}

	return false
}

// Implements asserts that an object is implemented by the specified interface.
//
//    assert.Implements(t, (*MyInterface)(nil), new(MyObject), "MyObject")
func Implements(t *testing.T, interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {

	interfaceType := reflect.TypeOf(interfaceObject).Elem()

	if !reflect.TypeOf(object).Implements(interfaceType) {
		return Fail(t, fmt.Sprintf("Object must implement %v", interfaceType), msgAndArgs...)
	}

	return true

}

// IsType asserts that the specified objects are of the same type.
func IsType(t *testing.T, expectedType interface{}, object interface{}, msgAndArgs ...interface{}) bool {

	if !ObjectsAreEqual(reflect.TypeOf(object), reflect.TypeOf(expectedType)) {
		return Fail(t, fmt.Sprintf("Object expected to be of type %v, but was %v", reflect.TypeOf(expectedType), reflect.TypeOf(object)), msgAndArgs...)
	}

	return true
}

// Equal asserts that two objects are equal.
//
//    assert.Equal(t, 123, 123, "123 and 123 should be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func Equal(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) bool {

	if !ObjectsAreEqual(a, b) {
		return Fail(t, fmt.Sprintf("Not equal: %#v != %#v", a, b), msgAndArgs...)
	}

	return true

}

// Exactly asserts that two objects are equal is value and type.
//
//    assert.Exactly(t, int32(123), int64(123), "123 and 123 should NOT be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func Exactly(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) bool {

	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)

	if aType != bType {
		return Fail(t, "Types expected to match exactly", "%v != %v", aType, bType)
	}

	return Equal(t, a, b, msgAndArgs...)

}

// NotNil asserts that the specified object is not nil.
//
//    assert.NotNil(t, err, "err should be something")
//
// Returns whether the assertion was successful (true) or not (false).
func NotNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) bool {

	var success bool = true

	if object == nil {
		success = false
	} else {
		value := reflect.ValueOf(object)
		kind := value.Kind()
		if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
			success = false
		}
	}

	if !success {
		Fail(t, "Expected not to be nil.", msgAndArgs...)
	}

	return success
}

// Nil asserts that the specified object is nil.
//
//    assert.Nil(t, err, "err should be nothing")
//
// Returns whether the assertion was successful (true) or not (false).
func Nil(t *testing.T, object interface{}, msgAndArgs ...interface{}) bool {

	if object == nil {
		return true
	} else {
		value := reflect.ValueOf(object)
		kind := value.Kind()
		if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
			return true
		}
	}

	return Fail(t, fmt.Sprintf("Expected nil, but got: %#v", object), msgAndArgs...)
}

// isEmpty gets whether the specified object is considered empty or not.
func isEmpty(object interface{}) bool {

	if object == nil {
		return true
	} else if object == "" {
		return true
	} else if object == 0 {
		return true
	} else if object == false {
		return true
	}

	objValue := reflect.ValueOf(object)
	switch objValue.Kind() {
	case reflect.Slice:
		{
			return (objValue.Len() == 0)
		}
	}

	return false

}

// Empty asserts that the specified object is empty.  I.e. nil, "", false, 0 or a
// slice with len == 0.
//
// assert.Empty(t, obj)
//
// Returns whether the assertion was successful (true) or not (false).
func Empty(t *testing.T, object interface{}, msgAndArgs ...interface{}) bool {

	pass := isEmpty(object)
	if !pass {
		Fail(t, fmt.Sprintf("Should be empty, but was %v", object), msgAndArgs...)
	}

	return pass

}

// Empty asserts that the specified object is NOT empty.  I.e. not nil, "", false, 0 or a
// slice with len == 0.
//
// if assert.NotEmpty(t, obj) {
//   assert.Equal(t, "two", obj[1])
// }
//
// Returns whether the assertion was successful (true) or not (false).
func NotEmpty(t *testing.T, object interface{}, msgAndArgs ...interface{}) bool {

	pass := !isEmpty(object)
	if !pass {
		Fail(t, fmt.Sprintf("Should NOT be empty, but was %v", object), msgAndArgs...)
	}

	return pass

}

// True asserts that the specified value is true.
//
//    assert.True(t, myBool, "myBool should be true")
//
// Returns whether the assertion was successful (true) or not (false).
func True(t *testing.T, value bool, msgAndArgs ...interface{}) bool {

	if value != true {
		return Fail(t, "Should be true", msgAndArgs...)
	}

	return true

}

// False asserts that the specified value is true.
//
//    assert.False(t, myBool, "myBool should be false")
//
// Returns whether the assertion was successful (true) or not (false).
func False(t *testing.T, value bool, msgAndArgs ...interface{}) bool {

	if value != false {
		return Fail(t, "Should be false", msgAndArgs...)
	}

	return true

}

// NotEqual asserts that the specified values are NOT equal.
//
//    assert.NotEqual(t, obj1, obj2, "two objects shouldn't be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func NotEqual(t *testing.T, a, b interface{}, msgAndArgs ...interface{}) bool {

	if ObjectsAreEqual(a, b) {
		return Fail(t, "Should not be equal", msgAndArgs...)
	}

	return true

}

// Contains asserts that the specified string contains the specified substring.
//
//    assert.Contains(t, "Hello World", "World", "But 'Hello World' does contain 'World'")
//
// Returns whether the assertion was successful (true) or not (false).
func Contains(t *testing.T, s, contains string, msgAndArgs ...interface{}) bool {

	if !strings.Contains(s, contains) {
		return Fail(t, fmt.Sprintf("\"%s\" does not contain \"%s\"", s, contains), msgAndArgs...)
	}

	return true

}

// NotContains asserts that the specified string does NOT contain the specified substring.
//
//    assert.NotContains(t, "Hello World", "Earth", "But 'Hello World' does NOT contain 'Earth'")
//
// Returns whether the assertion was successful (true) or not (false).
func NotContains(t *testing.T, s, contains string, msgAndArgs ...interface{}) bool {

	if strings.Contains(s, contains) {
		return Fail(t, fmt.Sprintf("\"%s\" should not contain \"%s\"", s, contains), msgAndArgs...)
	}

	return true

}

// Uses a Comparison to assert a complex condition.
func Condition(t *testing.T, comp Comparison, msgAndArgs ...interface{}) bool {
	result := comp()
	if !result {
		Fail(t, "Condition failed!", msgAndArgs...)
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

	if funcDidPanic, panicValue := didPanic(f); !funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should panic\n\r\tPanic value:\t%v", f, panicValue), msgAndArgs...)
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

	if funcDidPanic, panicValue := didPanic(f); funcDidPanic {
		return Fail(t, fmt.Sprintf("func %#v should not panic\n\r\tPanic value:\t%v", f, panicValue), msgAndArgs...)
	}

	return true
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
//   assert.WithinDuration(t, time.Now(), time.Now(), 10*time.Second, "The difference should not be more than 10s")
//
// Returns whether the assertion was successful (true) or not (false).
func WithinDuration(t *testing.T, a, b time.Time, delta time.Duration, msgAndArgs ...interface{}) bool {

	dt := a.Sub(b)
	if dt < -delta || dt > delta {
		return Fail(t, fmt.Sprintf("Max difference between %v and %v allowed is %v, but difference was %v", a, b, dt, delta), msgAndArgs...)
	}

	return true
}

/*
	Errors
*/

// NoError asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.NoError(t, err) {
//	   assert.Equal(t, actualObj, expectedObj)
//   }
//
// Returns whether the assertion was successful (true) or not (false).
func NoError(t *testing.T, theError error, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)
	return Nil(t, theError, "No error is expected but got %v %s", theError, message)

}

// Error asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.Error(t, err, "An error was expected") {
//	   assert.Equal(t, err, expectedError)
//   }
//
// Returns whether the assertion was successful (true) or not (false).
func Error(t *testing.T, theError error, msgAndArgs ...interface{}) bool {

	message := messageFromMsgAndArgs(msgAndArgs...)
	return NotNil(t, theError, "An error is expected but got nil. %s", message)

}
