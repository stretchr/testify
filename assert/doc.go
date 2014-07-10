// A set of comprehensive testing tools for use with the normal Go testing system.
//
// Example Usage
//
// The following is a complete example using assert in a standard test function:
//    import (
//      "testing"
//      "github.com/stretchr/testify/assert"
//    )
//
//    func TestSomething(t *testing.T) {
//
//      var a string = "Hello"
//      var b string = "Hello"
//
//      assert.Equal(t, a, b, "The two words should be the same.")
//
//    }
//
// if you assert many times, use the below:
//
//    import (
//      "testing"
//      "github.com/stretchr/testify/assert"
//    )
//
//    func TestSomething(t *testing.T) {
//      assert := assert.New(t)
//
//      var a string = "Hello"
//      var b string = "Hello"
//
//      assert.Equal(a, b, "The two words should be the same.")
//    }
//
// Assertions
//
// Assertions allow you to easily write test code, and are global funcs in the `assert` package.
// All assertion functions take, as the first argument, the `*testing.T` object provided by the
// testing framework. This allows the assertion funcs to write the failings and other details to
// the correct place.
//
// Every assertion function also takes an optional string message as the final argument,
// allowing custom error messages to be appended to the message the assertion method outputs.
//
// Here is an overview of the assert functions:
//    assert.Equal(expected, actual [, message [, format-args])
//
//    assert.NotEqual(notExpected, actual [, message [, format-args]])
//
//    assert.True(actualBool [, message [, format-args]])
//
//    assert.False(actualBool [, message [, format-args]])
//
//    assert.Nil(actualObject [, message [, format-args]])
//
//    assert.NotNil(actualObject [, message [, format-args]])
//
//    assert.Empty(actualObject [, message [, format-args]])
//
//    assert.NotEmpty(actualObject [, message [, format-args]])
//
//    assert.Error(errorObject [, message [, format-args]])
//
//    assert.NoError(errorObject [, message [, format-args]])
//
//    assert.EqualError(theError, errString [, message [, format-args]])
//
//    assert.Implements((*MyInterface)(nil), new(MyObject) [,message [, format-args]])
//
//    assert.IsType(expectedObject, actualObject [, message [, format-args]])
//
//    assert.Contains(string, substring [, message [, format-args]])
//
//    assert.NotContains(string, substring [, message [, format-args]])
//
//    assert.Panics(func(){
//
//	    // call code that should panic
//
//    } [, message [, format-args]])
//
//    assert.NotPanics(func(){
//
//	    // call code that should not panic
//
//    } [, message [, format-args]])
//
//    assert.WithinDuration(timeA, timeB, deltaTime, [, message [, format-args]])
//
//    assert.InDelta(numA, numB, delta, [, message [, format-args]])
//
//    assert.InEpsilon(numA, numB, epsilon, [, message [, format-args]])
package assert
