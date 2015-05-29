package mock

import (
	"fmt"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

/*
	Call
*/

// Call represents a method call and is used for setting expectations,
// as well as recording activity.
type Call struct {

	// The name of the method that was or will be called.
	Method string

	// Holds the arguments of the method.
	Arguments Arguments

	// Holds the arguments that should be returned when
	// this method is called.
	ReturnArguments Arguments

	// The number of times to return the return arguments when setting
	// expectations. 0 means to always return the value.
	Repeatability int

	// Holds a channel that will be used to block the Return until it either
	// recieves a message or is closed. nil means it returns immediately.
	WaitFor <-chan time.Time

	// Holds a handler used to manipulate arguments content that are passed by
	// reference. It's useful when mocking methods such as unmarshalers or
	// decoders.
	Run func(Arguments)
}

// Mock is the workhorse used to track activity on another object.
// For an example of its usage, refer to the "Example Usage" section at the top of this document.
type Mock struct {

	// The method name that is currently
	// being referred to by the On method.
	onMethodName string

	// An array of the arguments that are
	// currently being referred to by the On method.
	onMethodArguments Arguments

	// Represents the calls that are expected of
	// an object.
	ExpectedCalls []Call

	// Holds the calls that were made to this mocked object.
	Calls []Call

	// TestData holds any data that might be useful for testing.  Testify ignores
	// this data completely allowing you to do whatever you like with it.
	testData objx.Map

	mutex sync.Mutex
}

// TestData holds any data that might be useful for testing.  Testify ignores
// this data completely allowing you to do whatever you like with it.
func (m *Mock) TestData() objx.Map {

	if m.testData == nil {
		m.testData = make(objx.Map)
	}

	return m.testData
}

/*
	Setting expectations
*/

// On starts a description of an expectation of the specified method
// being called.
//
//     Mock.On("MyMethod", arg1, arg2)
func (m *Mock) On(methodName string, arguments ...interface{}) *Mock {
	m.onMethodName = methodName
	m.onMethodArguments = arguments

	for _, arg := range arguments {
		if v := reflect.ValueOf(arg); v.Kind() == reflect.Func {
			panic(fmt.Sprintf("cannot use Func in expectations. Use mock.AnythingOfType(\"%T\")", arg))
		}
	}

	return m
}

// Return finishes a description of an expectation of the method (and arguments)
// specified in the most recent On method call.
//
//     Mock.On("MyMethod", arg1, arg2).Return(returnArg1, returnArg2)
func (m *Mock) Return(returnArguments ...interface{}) *Mock {
	m.ExpectedCalls = append(m.ExpectedCalls, Call{m.onMethodName, m.onMethodArguments, returnArguments, 0, nil, nil})
	return m
}

// Once indicates that that the mock should only return the value once.
//
//    Mock.On("MyMethod", arg1, arg2).Return(returnArg1, returnArg2).Once()
func (m *Mock) Once() {
	m.ExpectedCalls[len(m.ExpectedCalls)-1].Repeatability = 1
}

// Twice indicates that that the mock should only return the value twice.
//
//    Mock.On("MyMethod", arg1, arg2).Return(returnArg1, returnArg2).Twice()
func (m *Mock) Twice() {
	m.ExpectedCalls[len(m.ExpectedCalls)-1].Repeatability = 2
}

// Times indicates that that the mock should only return the indicated number
// of times.
//
//    Mock.On("MyMethod", arg1, arg2).Return(returnArg1, returnArg2).Times(5)
func (m *Mock) Times(i int) {
	m.ExpectedCalls[len(m.ExpectedCalls)-1].Repeatability = i
}

// WaitUntil sets the channel that will block the mock's return until its closed
// or a message is received.
//
//    Mock.On("MyMethod", arg1, arg2).WaitUntil(time.After(time.Second))
func (m *Mock) WaitUntil(w <-chan time.Time) *Mock {
	m.ExpectedCalls[len(m.ExpectedCalls)-1].WaitFor = w
	return m
}

// After sets how long to block until the call returns
//
//    Mock.On("MyMethod", arg1, arg2).After(time.Second)
func (m *Mock) After(d time.Duration) *Mock {
	return m.WaitUntil(time.After(d))
}

// Run sets a handler to be called before returning. It can be used when
// mocking a method such as unmarshalers that takes a pointer to a struct and
// sets properties in such struct
//
//    Mock.On("Unmarshal", AnythingOfType("*map[string]interface{}").Return().Run(function(args Arguments) {
//    	arg := args.Get(0).(*map[string]interface{})
//    	arg["foo"] = "bar"
//    })
func (m *Mock) Run(fn func(Arguments)) *Mock {
	m.ExpectedCalls[len(m.ExpectedCalls)-1].Run = fn
	return m
}

/*
	Recording and responding to activity
*/

func (m *Mock) findExpectedCall(method string, arguments ...interface{}) (int, *Call) {
	for i, call := range m.ExpectedCalls {
		if call.Method == method && call.Repeatability > -1 {

			_, diffCount := call.Arguments.Diff(arguments)
			if diffCount == 0 {
				return i, &call
			}

		}
	}
	return -1, nil
}

func (m *Mock) findClosestCall(method string, arguments ...interface{}) (bool, *Call) {

	diffCount := 0
	var closestCall *Call = nil

	for _, call := range m.ExpectedCalls {
		if call.Method == method {

			_, tempDiffCount := call.Arguments.Diff(arguments)
			if tempDiffCount < diffCount || diffCount == 0 {
				diffCount = tempDiffCount
				closestCall = &call
			}

		}
	}

	if closestCall == nil {
		return false, nil
	}

	return true, closestCall
}

func callString(method string, arguments Arguments, includeArgumentValues bool) string {

	var argValsString string = ""
	if includeArgumentValues {
		var argVals []string
		for argIndex, arg := range arguments {
			argVals = append(argVals, fmt.Sprintf("%d: %v", argIndex, arg))
		}
		argValsString = fmt.Sprintf("\n\t\t%s", strings.Join(argVals, "\n\t\t"))
	}

	return fmt.Sprintf("%s(%s)%s", method, arguments.String(), argValsString)
}

// Called tells the mock object that a method has been called, and gets an array
// of arguments to return.  Panics if the call is unexpected (i.e. not preceeded by
// appropriate .On .Return() calls)
// If Call.WaitFor is set, blocks until the channel is closed or receives a message.
func (m *Mock) Called(arguments ...interface{}) Arguments {
	defer m.mutex.Unlock()
	m.mutex.Lock()

	// get the calling function's name
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Couldn't get the caller information")
	}
	functionPath := runtime.FuncForPC(pc).Name()
	parts := strings.Split(functionPath, ".")
	functionName := parts[len(parts)-1]

	found, call := m.findExpectedCall(functionName, arguments...)

	switch {
	case found < 0:
		// we have to fail here - because we don't know what to do
		// as the return arguments.  This is because:
		//
		//   a) this is a totally unexpected call to this method,
		//   b) the arguments are not what was expected, or
		//   c) the developer has forgotten to add an accompanying On...Return pair.

		closestFound, closestCall := m.findClosestCall(functionName, arguments...)

		if closestFound {
			panic(fmt.Sprintf("\n\nmock: Unexpected Method Call\n-----------------------------\n\n%s\n\nThe closest call I have is: \n\n%s\n", callString(functionName, arguments, true), callString(functionName, closestCall.Arguments, true)))
		} else {
			panic(fmt.Sprintf("\nassert: mock: I don't know what to return because the method call was unexpected.\n\tEither do Mock.On(\"%s\").Return(...) first, or remove the %s() call.\n\tThis method was unexpected:\n\t\t%s\n\tat: %s", functionName, functionName, callString(functionName, arguments, true), assert.CallerInfo()))
		}
	case call.Repeatability == 1:
		call.Repeatability = -1
		m.ExpectedCalls[found] = *call
	case call.Repeatability > 1:
		call.Repeatability -= 1
		m.ExpectedCalls[found] = *call
	}

	// add the call
	m.Calls = append(m.Calls, Call{functionName, arguments, make([]interface{}, 0), 0, nil, nil})

	// block if specified
	if call.WaitFor != nil {
		<-call.WaitFor
	}

	if call.Run != nil {
		call.Run(arguments)
	}

	return call.ReturnArguments

}

/*
	Assertions
*/

// AssertExpectationsForObjects asserts that everything specified with On and Return
// of the specified objects was in fact called as expected.
//
// Calls may have occurred in any order.
func AssertExpectationsForObjects(t TestingT, testObjects ...interface{}) bool {
	var success bool = true
	for _, obj := range testObjects {
		mockObj := obj.(Mock)
		success = success && mockObj.AssertExpectations(t)
	}
	return success
}

// AssertExpectations asserts that everything specified with On and Return was
// in fact called as expected.  Calls may have occurred in any order.
func (m *Mock) AssertExpectations(t TestingT) bool {

	var somethingMissing bool = false
	var failedExpectations int = 0

	// iterate through each expectation
	for _, expectedCall := range m.ExpectedCalls {
		switch {
		case !m.methodWasCalled(expectedCall.Method, expectedCall.Arguments):
			somethingMissing = true
			failedExpectations++
			t.Logf("\u274C\t%s(%s)", expectedCall.Method, expectedCall.Arguments.String())
		case expectedCall.Repeatability > 0:
			somethingMissing = true
			failedExpectations++
		default:
			t.Logf("\u2705\t%s(%s)", expectedCall.Method, expectedCall.Arguments.String())
		}
	}

	if somethingMissing {
		t.Errorf("FAIL: %d out of %d expectation(s) were met.\n\tThe code you are testing needs to make %d more call(s).\n\tat: %s", len(m.ExpectedCalls)-failedExpectations, len(m.ExpectedCalls), failedExpectations, assert.CallerInfo())
	}

	return !somethingMissing
}

// AssertNumberOfCalls asserts that the method was called expectedCalls times.
func (m *Mock) AssertNumberOfCalls(t TestingT, methodName string, expectedCalls int) bool {
	var actualCalls int = 0
	for _, call := range m.Calls {
		if call.Method == methodName {
			actualCalls++
		}
	}
	return assert.Equal(t, actualCalls, expectedCalls, fmt.Sprintf("Expected number of calls (%d) does not match the actual number of calls (%d).", expectedCalls, actualCalls))
}

// AssertCalled asserts that the method was called.
func (m *Mock) AssertCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	if !assert.True(t, m.methodWasCalled(methodName, arguments), fmt.Sprintf("The \"%s\" method should have been called with %d argument(s), but was not.", methodName, len(arguments))) {
		t.Logf("%v", m.ExpectedCalls)
		return false
	}
	return true
}

// AssertNotCalled asserts that the method was not called.
func (m *Mock) AssertNotCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	if !assert.False(t, m.methodWasCalled(methodName, arguments), fmt.Sprintf("The \"%s\" method was called with %d argument(s), but should NOT have been.", methodName, len(arguments))) {
		t.Logf("%v", m.ExpectedCalls)
		return false
	}
	return true
}

func (m *Mock) methodWasCalled(methodName string, expected []interface{}) bool {
	for _, call := range m.Calls {
		if call.Method == methodName {

			_, differences := Arguments(expected).Diff(call.Arguments)

			if differences == 0 {
				// found the expected call
				return true
			}

		}
	}
	// we didn't find the expected call
	return false
}

/*
	Arguments
*/

// Arguments holds an array of method arguments or return values.
type Arguments []interface{}

const (
	// The "any" argument.  Used in Diff and Assert when
	// the argument being tested shouldn't be taken into consideration.
	Anything string = "mock.Anything"
)

// AnythingOfTypeArgument is a string that contains the type of an argument
// for use when type checking.  Used in Diff and Assert.
type AnythingOfTypeArgument string

// AnythingOfType returns an AnythingOfTypeArgument object containing the
// name of the type to check for.  Used in Diff and Assert.
//
// For example:
//	Assert(t, AnythingOfType("string"), AnythingOfType("int"))
func AnythingOfType(t string) AnythingOfTypeArgument {
	return AnythingOfTypeArgument(t)
}

// Get Returns the argument at the specified index.
func (args Arguments) Get(index int) interface{} {
	if index+1 > len(args) {
		panic(fmt.Sprintf("assert: arguments: Cannot call Get(%d) because there are %d argument(s).", index, len(args)))
	}
	return args[index]
}

// Is gets whether the objects match the arguments specified.
func (args Arguments) Is(objects ...interface{}) bool {
	for i, obj := range args {
		if obj != objects[i] {
			return false
		}
	}
	return true
}

// Diff gets a string describing the differences between the arguments
// and the specified objects.
//
// Returns the diff string and number of differences found.
func (args Arguments) Diff(objects []interface{}) (string, int) {

	var output string = "\n"
	var differences int

	var maxArgCount int = len(args)
	if len(objects) > maxArgCount {
		maxArgCount = len(objects)
	}

	for i := 0; i < maxArgCount; i++ {
		var actual, expected interface{}

		if len(objects) <= i {
			actual = "(Missing)"
		} else {
			actual = objects[i]
		}

		if len(args) <= i {
			expected = "(Missing)"
		} else {
			expected = args[i]
		}

		if reflect.TypeOf(expected) == reflect.TypeOf((*AnythingOfTypeArgument)(nil)).Elem() {

			// type checking
			if reflect.TypeOf(actual).Name() != string(expected.(AnythingOfTypeArgument)) && reflect.TypeOf(actual).String() != string(expected.(AnythingOfTypeArgument)) {
				// not match
				differences++
				output = fmt.Sprintf("%s\t%d: \u274C  type %s != type %s - %s\n", output, i, expected, reflect.TypeOf(actual).Name(), actual)
			}

		} else {

			// normal checking

			if assert.ObjectsAreEqual(expected, Anything) || assert.ObjectsAreEqual(actual, Anything) || assert.ObjectsAreEqual(actual, expected) {
				// match
				output = fmt.Sprintf("%s\t%d: \u2705  %s == %s\n", output, i, actual, expected)
			} else {
				// not match
				differences++
				output = fmt.Sprintf("%s\t%d: \u274C  %s != %s\n", output, i, actual, expected)
			}
		}

	}

	if differences == 0 {
		return "No differences.", differences
	}

	return output, differences

}

// Assert compares the arguments with the specified objects and fails if
// they do not exactly match.
func (args Arguments) Assert(t TestingT, objects ...interface{}) bool {

	// get the differences
	diff, diffCount := args.Diff(objects)

	if diffCount == 0 {
		return true
	}

	// there are differences... report them...
	t.Logf(diff)
	t.Errorf("%sArguments do not match.", assert.CallerInfo())

	return false

}

// Bool gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Bool(index int) bool {
	var s bool
	var ok bool
	if s, ok = args.Get(index).(bool); !ok {
		panic(fmt.Sprintf("assert: arguments: Bool(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Int8 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Int8(index int) int8 {
	var s int8
	var ok bool
	if s, ok = args.Get(index).(int8); !ok {
		panic(fmt.Sprintf("assert: arguments: Int8(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Int16 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Int16(index int) int16 {
	var s int16
	var ok bool
	if s, ok = args.Get(index).(int16); !ok {
		panic(fmt.Sprintf("assert: arguments: Int16(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Int32 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Int32(index int) int32 {
	var s int32
	var ok bool
	if s, ok = args.Get(index).(int32); !ok {
		panic(fmt.Sprintf("assert: arguments: Int32(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Int64 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Int64(index int) int64 {
	var s int64
	var ok bool
	if s, ok = args.Get(index).(int64); !ok {
		panic(fmt.Sprintf("assert: arguments: Int64(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Uint8 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Uint8(index int) uint8 {
	var s uint8
	var ok bool
	if s, ok = args.Get(index).(uint8); !ok {
		panic(fmt.Sprintf("assert: arguments: Uint8(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Uint16 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Uint16(index int) uint16 {
	var s uint16
	var ok bool
	if s, ok = args.Get(index).(uint16); !ok {
		panic(fmt.Sprintf("assert: arguments: Uint16(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Uint32 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Uint32(index int) uint32 {
	var s uint32
	var ok bool
	if s, ok = args.Get(index).(uint32); !ok {
		panic(fmt.Sprintf("assert: arguments: Uint32(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Uint64 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Uint64(index int) uint64 {
	var s uint64
	var ok bool
	if s, ok = args.Get(index).(uint64); !ok {
		panic(fmt.Sprintf("assert: arguments: Uint64(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Float32 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Float32(index int) float32 {
	var s float32
	var ok bool
	if s, ok = args.Get(index).(float32); !ok {
		panic(fmt.Sprintf("assert: arguments: Float32(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Float64 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Float64(index int) float64 {
	var s float64
	var ok bool
	if s, ok = args.Get(index).(float64); !ok {
		panic(fmt.Sprintf("assert: arguments: Float64(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Complex64 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Complex64(index int) complex64 {
	var s complex64
	var ok bool
	if s, ok = args.Get(index).(complex64); !ok {
		panic(fmt.Sprintf("assert: arguments: Complex64(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Complex128 gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Complex128(index int) complex128 {
	var s complex128
	var ok bool
	if s, ok = args.Get(index).(complex128); !ok {
		panic(fmt.Sprintf("assert: arguments: Complex128(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Byte gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Byte(index int) byte {
	var s byte
	var ok bool
	if s, ok = args.Get(index).(byte); !ok {
		panic(fmt.Sprintf("assert: arguments: Byte(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Rune gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Rune(index int) rune {
	var s rune
	var ok bool
	if s, ok = args.Get(index).(rune); !ok {
		panic(fmt.Sprintf("assert: arguments: Rune(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Int gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Int(index int) int {
	var s int
	var ok bool
	if s, ok = args.Get(index).(int); !ok {
		panic(fmt.Sprintf("assert: arguments: Int(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Uint gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Uint(index int) uint {
	var s uint
	var ok bool
	if s, ok = args.Get(index).(uint); !ok {
		panic(fmt.Sprintf("assert: arguments: Uint(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// Uintptr gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Uintptr(index int) uintptr {
	var s uintptr
	var ok bool
	if s, ok = args.Get(index).(uintptr); !ok {
		panic(fmt.Sprintf("assert: arguments: Uintptr(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
	}
	return s
}

// String gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
//
// If no index is provided, String() returns a complete string representation
// of the arguments.
func (args Arguments) String(indexOrNil ...int) string {

	if len(indexOrNil) == 0 {
		// normal String() method - return a string representation of the args
		var argsStr []string
		for _, arg := range args {
			argsStr = append(argsStr, fmt.Sprintf("%s", reflect.TypeOf(arg)))
		}
		return strings.Join(argsStr, ",")
	} else if len(indexOrNil) == 1 {
		// Index has been specified - get the argument at that index
		var index int = indexOrNil[0]
		var s string
		var ok bool
		if s, ok = args.Get(index).(string); !ok {
			panic(fmt.Sprintf("assert: arguments: String(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
		}
		return s
	}

	panic(fmt.Sprintf("assert: arguments: Wrong number of arguments passed to String.  Must be 0 or 1, not %d", len(indexOrNil)))

}

// Error gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Error(index int) error {
	obj := args.Get(index)
	var s error
	var ok bool
	if obj == nil {
		return nil
	}
	if s, ok = obj.(error); !ok {
		panic(fmt.Sprintf("assert: arguments: Error(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}
