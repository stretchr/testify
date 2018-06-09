package mock

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
)

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
}

// Mock is the workhorse used to track activity on another object.
// For an example of its usage, refer to the "Example Usage" section at the top
// of this document.
type Mock struct {
	// Represents the calls that are expected of
	// an object.
	ExpectedCalls []*Call

	// Holds the calls that were made to this mocked object.
	Calls []Call

	// test is An optional variable that holds the test struct, to be used when an
	// invalid mock call was made.
	test TestingT

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

// Test sets the test struct variable of the mock object
func (m *Mock) Test(t TestingT) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.test = t
}

// fail fails the current test with the given formatted format and args.
// In case that a test was defined, it uses the test APIs for failing a test,
// otherwise it uses panic.
func (m *Mock) fail(format string, args ...interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.test == nil {
		panic(fmt.Sprintf(format, args...))
	}
	m.test.Errorf(format, args...)
	m.test.FailNow()
}

// On starts a description of an expectation of the specified method
// being called.
//
//     Mock.On("MyMethod", arg1, arg2)
func (m *Mock) On(methodName string, arguments ...interface{}) *Call {
	for _, arg := range arguments {
		if v := reflect.ValueOf(arg); v.Kind() == reflect.Func {
			panic(fmt.Sprintf("cannot use Func in expectations. Use mock.AnythingOfType(\"%T\")", arg))
		}
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	c := newCall(m, methodName, assert.CallerInfo(), arguments...)
	m.ExpectedCalls = append(m.ExpectedCalls, c)
	return c
}

// /*
// 	Recording and responding to activity
// */

func (m *Mock) findExpectedCall(method string, arguments ...interface{}) (int, *Call) {
	for i, call := range m.ExpectedCalls {
		if call.Method == method && call.Repeatability > -1 {

			_, diffCount := call.Arguments.Diff(arguments)
			if diffCount == 0 {
				return i, call
			}

		}
	}
	return -1, nil
}

func (m *Mock) findClosestCall(method string, arguments ...interface{}) (*Call, string) {
	var diffCount int
	var closestCall *Call
	var err string

	for _, call := range m.expectedCalls() {
		if call.Method == method {

			errInfo, tempDiffCount := call.Arguments.Diff(arguments)
			if tempDiffCount < diffCount || diffCount == 0 {
				diffCount = tempDiffCount
				closestCall = call
				err = errInfo
			}

		}
	}

	return closestCall, err
}

func callString(method string, arguments Arguments, includeArgumentValues bool) string {

	var argValsString string
	if includeArgumentValues {
		var argVals []string
		for argIndex, arg := range arguments {
			argVals = append(argVals, fmt.Sprintf("%d: %#v", argIndex, arg))
		}
		argValsString = fmt.Sprintf("\n\t\t%s", strings.Join(argVals, "\n\t\t"))
	}

	return fmt.Sprintf("%s(%s)%s", method, arguments.String(), argValsString)
}

// Called tells the mock object that a method has been called, and gets an array
// of arguments to return.  Panics if the call is unexpected (i.e. not preceded by
// appropriate .On .Return() calls)
// If Call.WaitFor is set, blocks until the channel is closed or receives a message.
func (m *Mock) Called(arguments ...interface{}) Arguments {
	// get the calling function's name
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Couldn't get the caller information")
	}
	functionPath := runtime.FuncForPC(pc).Name()
	//Next four lines are required to use GCCGO function naming conventions.
	//For Ex:  github_com_docker_libkv_store_mock.WatchTree.pN39_github_com_docker_libkv_store_mock.Mock
	//uses interface information unlike golang github.com/docker/libkv/store/mock.(*Mock).WatchTree
	//With GCCGO we need to remove interface information starting from pN<dd>.
	re := regexp.MustCompile("\\.pN\\d+_")
	if re.MatchString(functionPath) {
		functionPath = re.Split(functionPath, -1)[0]
	}
	parts := strings.Split(functionPath, ".")
	functionName := parts[len(parts)-1]
	return m.MethodCalled(functionName, arguments...)
}

// MethodCalled tells the mock object that the given method has been called, and gets
// an array of arguments to return. Panics if the call is unexpected (i.e. not preceded
// by appropriate .On .Return() calls)
// If Call.WaitFor is set, blocks until the channel is closed or receives a message.
func (m *Mock) MethodCalled(methodName string, arguments ...interface{}) Arguments {
	m.mutex.Lock()
	//TODO: could combine expected and closes in single loop
	found, call := m.findExpectedCall(methodName, arguments...)

	if found < 0 {
		// we have to fail here - because we don't know what to do
		// as the return arguments.  This is because:
		//
		//   a) this is a totally unexpected call to this method,
		//   b) the arguments are not what was expected, or
		//   c) the developer has forgotten to add an accompanying On...Return pair.

		closestCall, mismatch := m.findClosestCall(methodName, arguments...)
		m.mutex.Unlock()

		if closestCall != nil {
			m.fail("\n\nmock: Unexpected Method Call\n-----------------------------\n\n%s\n\nThe closest call I have is: \n\n%s\n\n%s\nDiff: %s",
				callString(methodName, arguments, true),
				callString(methodName, closestCall.Arguments, true),
				diffArguments(closestCall.Arguments, arguments),
				strings.TrimSpace(mismatch),
			)
		} else {
			m.fail("\nassert: mock: I don't know what to return because the method call was unexpected.\n\tEither do Mock.On(\"%s\").Return(...) first, or remove the %s() call.\n\tThis method was unexpected:\n\t\t%s\n\tat: %s", methodName, methodName, callString(methodName, arguments, true), assert.CallerInfo())
		}
	}

	if call.Repeatability == 1 {
		call.Repeatability = -1
	} else if call.Repeatability > 1 {
		call.Repeatability--
	}
	call.totalCalls++

	// add the call
	m.Calls = append(m.Calls, *newCall(m, methodName, assert.CallerInfo(), arguments...))
	m.mutex.Unlock()

	// block if specified
	if call.WaitFor != nil {
		<-call.WaitFor
	} else {
		time.Sleep(call.waitTime)
	}

	m.mutex.Lock()
	runFn := call.RunFn
	m.mutex.Unlock()

	if runFn != nil {
		runFn(arguments)
	}

	m.mutex.Lock()
	returnArgs := call.ReturnArguments
	m.mutex.Unlock()

	return returnArgs
}

/*
	Assertions
*/

type assertExpectationser interface {
	AssertExpectations(TestingT) bool
}

// AssertExpectationsForObjects asserts that everything specified with On and Return
// of the specified objects was in fact called as expected.
//
// Calls may have occurred in any order.
func AssertExpectationsForObjects(t TestingT, testObjects ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	for _, obj := range testObjects {
		if m, ok := obj.(Mock); ok {
			t.Logf("Deprecated mock.AssertExpectationsForObjects(myMock.Mock) use mock.AssertExpectationsForObjects(myMock)")
			obj = &m
		}
		m := obj.(assertExpectationser)
		if !m.AssertExpectations(t) {
			t.Logf("Expectations didn't match for Mock: %+v", reflect.TypeOf(m))
			return false
		}
	}
	return true
}

// AssertExpectations asserts that everything specified with On and Return was
// in fact called as expected.  Calls may have occurred in any order.
func (m *Mock) AssertExpectations(t TestingT) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var somethingMissing bool
	var failedExpectations int

	// iterate through each expectation
	expectedCalls := m.expectedCalls()
	for _, expectedCall := range expectedCalls {
		if !expectedCall.optional && !m.methodWasCalled(expectedCall.Method, expectedCall.Arguments) && expectedCall.totalCalls == 0 {
			somethingMissing = true
			failedExpectations++
			t.Logf("FAIL:\t%s(%s)\n\t\tat: %s", expectedCall.Method, expectedCall.Arguments.String(), expectedCall.callerInfo)
		} else {
			if expectedCall.Repeatability > 0 {
				somethingMissing = true
				failedExpectations++
				t.Logf("FAIL:\t%s(%s)\n\t\tat: %s", expectedCall.Method, expectedCall.Arguments.String(), expectedCall.callerInfo)
			} else {
				t.Logf("PASS:\t%s(%s)", expectedCall.Method, expectedCall.Arguments.String())
			}
		}
	}

	if somethingMissing {
		t.Errorf("FAIL: %d out of %d expectation(s) were met.\n\tThe code you are testing needs to make %d more call(s).\n\tat: %s", len(expectedCalls)-failedExpectations, len(expectedCalls), failedExpectations, assert.CallerInfo())
	}

	return !somethingMissing
}

// AssertNumberOfCalls asserts that the method was called expectedCalls times.
func (m *Mock) AssertNumberOfCalls(t TestingT, methodName string, expectedCalls int) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var actualCalls int
	for _, call := range m.calls() {
		if call.Method == methodName {
			actualCalls++
		}
	}
	return assert.Equal(t, expectedCalls, actualCalls, fmt.Sprintf("Expected number of calls (%d) does not match the actual number of calls (%d).", expectedCalls, actualCalls))
}

// AssertCalled asserts that the method was called.
// It can produce a false result when an argument is a pointer type and the underlying value changed after calling the mocked method.
func (m *Mock) AssertCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if !m.methodWasCalled(methodName, arguments) {
		var calledWithArgs []string
		for _, call := range m.calls() {
			calledWithArgs = append(calledWithArgs, fmt.Sprintf("%v", call.Arguments))
		}
		if len(calledWithArgs) == 0 {
			return assert.Fail(t, "Should have called with given arguments",
				fmt.Sprintf("Expected %q to have been called with:\n%v\nbut no actual calls happened", methodName, arguments))
		}
		return assert.Fail(t, "Should have called with given arguments",
			fmt.Sprintf("Expected %q to have been called with:\n%v\nbut actual calls were:\n        %v", methodName, arguments, strings.Join(calledWithArgs, "\n")))
	}
	return true
}

// AssertNotCalled asserts that the method was not called.
// It can produce a false result when an argument is a pointer type and the underlying value changed after calling the mocked method.
func (m *Mock) AssertNotCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.methodWasCalled(methodName, arguments) {
		return assert.Fail(t, "Should not have called with given arguments",
			fmt.Sprintf("Expected %q to not have been called with:\n%v\nbut actually it was.", methodName, arguments))
	}
	return true
}

func (m *Mock) methodWasCalled(methodName string, expected []interface{}) bool {
	for _, call := range m.calls() {
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

func (m *Mock) expectedCalls() []*Call {
	return append([]*Call{}, m.ExpectedCalls...)
}

func (m *Mock) calls() []Call {
	return append([]Call{}, m.Calls...)
}

func typeAndKind(v interface{}) (reflect.Type, reflect.Kind) {
	t := reflect.TypeOf(v)
	k := t.Kind()

	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	return t, k
}

func diffArguments(expected Arguments, actual Arguments) string {
	if len(expected) != len(actual) {
		return fmt.Sprintf("Provided %v arguments, mocked for %v arguments", len(expected), len(actual))
	}

	for x := range expected {
		if diffString := diff(expected[x], actual[x]); diffString != "" {
			return fmt.Sprintf("Difference found in argument %v:\n\n%s", x, diffString)
		}
	}

	return ""
}

// diff returns a diff of both values as long as both are of the same type and
// are a struct, map, slice or array. Otherwise it returns an empty string.
func diff(expected interface{}, actual interface{}) string {
	if expected == nil || actual == nil {
		return ""
	}

	et, ek := typeAndKind(expected)
	at, _ := typeAndKind(actual)

	if et != at {
		return ""
	}

	if ek != reflect.Struct && ek != reflect.Map && ek != reflect.Slice && ek != reflect.Array {
		return ""
	}

	e := spewConfig.Sdump(expected)
	a := spewConfig.Sdump(actual)

	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(e),
		B:        difflib.SplitLines(a),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  1,
	})

	return diff
}

var spewConfig = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
}

type tHelper interface {
	Helper()
}
