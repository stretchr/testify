package mock

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/stretchr/testify/assert"
)

/*
	Arguments
*/

// Arguments holds an array of method arguments or return values.
type Arguments []interface{}

const (
	// Anything is used in Diff and Assert when the argument being tested
	// shouldn't be taken into consideration.
	Anything = "mock.Anything"
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

// argumentMatcher performs custom argument matching, returning whether or
// not the argument is matched by the expectation fixture function.
type argumentMatcher struct {
	// fn is a function which accepts one argument, and returns a bool.
	fn reflect.Value
}

func (f argumentMatcher) Matches(argument interface{}) bool {
	expectType := f.fn.Type().In(0)
	expectTypeNilSupported := false
	switch expectType.Kind() {
	case reflect.Interface, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Ptr:
		expectTypeNilSupported = true
	}

	argType := reflect.TypeOf(argument)
	var arg reflect.Value
	if argType == nil {
		arg = reflect.New(expectType).Elem()
	} else {
		arg = reflect.ValueOf(argument)
	}

	if argType == nil && !expectTypeNilSupported {
		panic(errors.New("attempting to call matcher with nil for non-nil expected type"))
	}
	if argType == nil || argType.AssignableTo(expectType) {
		result := f.fn.Call([]reflect.Value{arg})
		return result[0].Bool()
	}
	return false
}

func (f argumentMatcher) String() string {
	return fmt.Sprintf("func(%s) bool", f.fn.Type().In(0).Name())
}

// MatchedBy can be used to match a mock call based on only certain properties
// from a complex struct or some calculation. It takes a function that will be
// evaluated with the called argument and will return true when there's a match
// and false otherwise.
//
// Example:
// m.On("Do", MatchedBy(func(req *http.Request) bool { return req.Host == "example.com" }))
//
// |fn|, must be a function accepting a single argument (of the expected type)
// which returns a bool. If |fn| doesn't match the required signature,
// MatchedBy() panics.
func MatchedBy(fn interface{}) argumentMatcher {
	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("assert: arguments: %s is not a func", fn))
	}
	if fnType.NumIn() != 1 {
		panic(fmt.Sprintf("assert: arguments: %s does not take exactly one argument", fn))
	}
	if fnType.NumOut() != 1 || fnType.Out(0).Kind() != reflect.Bool {
		panic(fmt.Sprintf("assert: arguments: %s does not return a bool", fn))
	}

	return argumentMatcher{fn: reflect.ValueOf(fn)}
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
	//TODO: could return string as error and nil for No difference

	var output = "\n"
	var differences int

	var maxArgCount = len(args)
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

		if matcher, ok := expected.(argumentMatcher); ok {
			if matcher.Matches(actual) {
				output = fmt.Sprintf("%s\t%d: PASS:  %s matched by %s\n", output, i, actual, matcher)
			} else {
				differences++
				output = fmt.Sprintf("%s\t%d: PASS:  %s not matched by %s\n", output, i, actual, matcher)
			}
		} else if reflect.TypeOf(expected) == reflect.TypeOf((*AnythingOfTypeArgument)(nil)).Elem() {

			// type checking
			if reflect.TypeOf(actual).Name() != string(expected.(AnythingOfTypeArgument)) && reflect.TypeOf(actual).String() != string(expected.(AnythingOfTypeArgument)) {
				// not match
				differences++
				output = fmt.Sprintf("%s\t%d: FAIL:  type %s != type %s - %s\n", output, i, expected, reflect.TypeOf(actual).Name(), actual)
			}

		} else {

			// normal checking

			if assert.ObjectsAreEqual(expected, Anything) || assert.ObjectsAreEqual(actual, Anything) || assert.ObjectsAreEqual(actual, expected) {
				// match
				output = fmt.Sprintf("%s\t%d: PASS:  %s == %s\n", output, i, actual, expected)
			} else {
				// not match
				differences++
				output = fmt.Sprintf("%s\t%d: FAIL:  %s != %s\n", output, i, actual, expected)
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
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

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
		var index = indexOrNil[0]
		var s string
		var ok bool
		if s, ok = args.Get(index).(string); !ok {
			panic(fmt.Sprintf("assert: arguments: String(%d) failed because object wasn't correct type: %s", index, args.Get(index)))
		}
		return s
	}

	panic(fmt.Sprintf("assert: arguments: Wrong number of arguments passed to String.  Must be 0 or 1, not %d", len(indexOrNil)))

}

// Int gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Int(index int) int {
	var s int
	var ok bool
	if s, ok = args.Get(index).(int); !ok {
		panic(fmt.Sprintf("assert: arguments: Int(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
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

// Bool gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (args Arguments) Bool(index int) bool {
	var s bool
	var ok bool
	if s, ok = args.Get(index).(bool); !ok {
		panic(fmt.Sprintf("assert: arguments: Bool(%d) failed because object wasn't correct type: %v", index, args.Get(index)))
	}
	return s
}
