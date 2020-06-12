package assert

import (
	"fmt"
	"reflect"
)

type CompareType int

const (
	compareLess CompareType = iota - 1
	compareEqual
	compareGreater
)

func compare(obj1, obj2 interface{}, kind reflect.Kind) (CompareType, bool) {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		{
			intobj1 := reflect.ValueOf(obj1).Int()
			intobj2 := reflect.ValueOf(obj2).Int()
			if intobj1 > intobj2 {
				return compareGreater, true
			}
			if intobj1 == intobj2 {
				return compareEqual, true
			}
			if intobj1 < intobj2 {
				return compareLess, true
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		{
			uintobj1 := reflect.ValueOf(obj1).Uint()
			uintobj2 := reflect.ValueOf(obj2).Uint()
			if uintobj1 > uintobj2 {
				return compareGreater, true
			}
			if uintobj1 == uintobj2 {
				return compareEqual, true
			}
			if uintobj1 < uintobj2 {
				return compareLess, true
			}
		}
	case reflect.Float32, reflect.Float64:
		{
			floatobj1 := reflect.ValueOf(obj1).Float()
			float32 := reflect.ValueOf(obj2).Float()
			if floatobj1 > float32 {
				return compareGreater, true
			}
			if floatobj1 == float32 {
				return compareEqual, true
			}
			if floatobj1 < float32 {
				return compareLess, true
			}
		}
	case reflect.String:
		{
			stringobj1 := reflect.ValueOf(obj1).String()
			stringobj2 := reflect.ValueOf(obj2).String()
			if stringobj1 > stringobj2 {
				return compareGreater, true
			}
			if stringobj1 == stringobj2 {
				return compareEqual, true
			}
			if stringobj1 < stringobj2 {
				return compareLess, true
			}
		}
	}

	return compareEqual, false
}

// Greater asserts that the first element is greater than the second
//
//    assert.Greater(t, 2, 1)
//    assert.Greater(t, float64(2), float64(1))
//    assert.Greater(t, "b", "a")
func Greater(t TestingT, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	return compareTwoValues(t, e1, e2, []CompareType{compareGreater}, "\"%v\" is not greater than \"%v\"", msgAndArgs)
}

// GreaterOrEqual asserts that the first element is greater than or equal to the second
//
//    assert.GreaterOrEqual(t, 2, 1)
//    assert.GreaterOrEqual(t, 2, 2)
//    assert.GreaterOrEqual(t, "b", "a")
//    assert.GreaterOrEqual(t, "b", "b")
func GreaterOrEqual(t TestingT, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	return compareTwoValues(t, e1, e2, []CompareType{compareGreater, compareEqual}, "\"%v\" is not greater than or equal to \"%v\"", msgAndArgs)
}

// Less asserts that the first element is less than the second
//
//    assert.Less(t, 1, 2)
//    assert.Less(t, float64(1), float64(2))
//    assert.Less(t, "a", "b")
func Less(t TestingT, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	return compareTwoValues(t, e1, e2, []CompareType{compareLess}, "\"%v\" is not less than \"%v\"", msgAndArgs)
}

// LessOrEqual asserts that the first element is less than or equal to the second
//
//    assert.LessOrEqual(t, 1, 2)
//    assert.LessOrEqual(t, 2, 2)
//    assert.LessOrEqual(t, "a", "b")
//    assert.LessOrEqual(t, "b", "b")
func LessOrEqual(t TestingT, e1 interface{}, e2 interface{}, msgAndArgs ...interface{}) bool {
	return compareTwoValues(t, e1, e2, []CompareType{compareLess, compareEqual}, "\"%v\" is not less than or equal to \"%v\"", msgAndArgs)
}

func compareTwoValues(t TestingT, e1 interface{}, e2 interface{}, allowedComparesResults []CompareType, failMessage string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind {
		return Fail(t, "Elements should be the same type", msgAndArgs...)
	}

	compareResult, isComparable := compare(e1, e2, e1Kind)
	if !isComparable {
		return Fail(t, fmt.Sprintf("Can not compare type \"%s\"", reflect.TypeOf(e1)), msgAndArgs...)
	}

	if !containsValue(allowedComparesResults, compareResult) {
		return Fail(t, fmt.Sprintf(failMessage, e1, e2), msgAndArgs...)
	}

	return true
}

func containsValue(values []CompareType, value CompareType) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}
