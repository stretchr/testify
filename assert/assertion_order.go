package assert

import (
	"fmt"
	"reflect"
	"sort"
)

// isOrdered checks that collection contains orderable elements.
func isOrdered(t TestingT, object interface{}, allowedComparesResults []CompareType, failMessage string, msgAndArgs ...interface{}) bool {
	if objSortInterface, ok := object.(sort.Interface); ok {
		return isOrderedSortInterface(t, objSortInterface, allowedComparesResults, failMessage, msgAndArgs...)
	}

	objKind := reflect.TypeOf(object).Kind()
	if objKind != reflect.Slice && objKind != reflect.Array {
		return Fail(t, fmt.Sprintf("object %T is not a collection", object), msgAndArgs...)
	}

	objValue := reflect.ValueOf(object)
	objLen := objValue.Len()

	if objLen <= 1 {
		return true
	}

	value := objValue.Index(0)
	valueInterface := value.Interface()
	firstValueKind := value.Kind()

	for i := 1; i < objLen; i++ {
		prevValue := value
		prevValueInterface := valueInterface

		value = objValue.Index(i)
		valueInterface = value.Interface()

		compareResult, isComparable := compare(prevValueInterface, valueInterface, firstValueKind)

		if !isComparable {
			return Fail(t, fmt.Sprintf("Can not compare type \"%s\" and \"%s\"", reflect.TypeOf(value), reflect.TypeOf(prevValue)), msgAndArgs...)
		}

		if !containsValue(allowedComparesResults, compareResult) {
			return Fail(t, fmt.Sprintf(`"%v"`+failMessage+`"%v"`, prevValue, value), msgAndArgs...)
		}
	}

	return true
}

// isOrderedSortInterface checks that sort.Interface collection contains orderable elements.
func isOrderedSortInterface(t TestingT, object sort.Interface, allowedComparesResults []CompareType, failMessage string, msgAndArgs ...interface{}) bool {
	allowLess, allowGreater, allowEqual := false, false, false
	for _, comparison := range allowedComparesResults {
		switch comparison {
		case compareLess:
			allowLess = true
		case compareGreater:
			allowGreater = true
		case compareEqual:
			allowEqual = true
		}
	}

	n := object.Len()
	for i := n - 1; i > 0; i-- {
		if allowLess && object.Less(i-1, i) {
			continue
		}
		if allowGreater && object.Less(i, i-1) {
			continue
		}
		// if allowLess or allowGreater are true then we can assume that the less and greater tests
		// respectively have already failed and avoid making repeated calls to Less()
		if allowEqual && (allowLess || !object.Less(i-1, i)) && (allowGreater || !object.Less(i, i-1)) {
			continue
		}
		return Fail(t, fmt.Sprintf("element %d"+failMessage+"element %d", i-1, i), msgAndArgs...)
	}

	return true
}

// IsIncreasing asserts that the collection is increasing
//
//    assert.IsIncreasing(t, []int{1, 2, 3})
//    assert.IsIncreasing(t, []float{1, 2})
//    assert.IsIncreasing(t, []string{"a", "b"})
func IsIncreasing(t TestingT, object interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return isOrdered(t, object, []CompareType{compareLess}, " is not less than ", msgAndArgs)
}

// IsNonIncreasing asserts that the collection is not increasing
//
//    assert.IsNonIncreasing(t, []int{2, 1, 1})
//    assert.IsNonIncreasing(t, []float{2, 1})
//    assert.IsNonIncreasing(t, []string{"b", "a"})
func IsNonIncreasing(t TestingT, object interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return isOrdered(t, object, []CompareType{compareEqual, compareGreater}, " is not greater than or equal to ", msgAndArgs)
}

// IsDecreasing asserts that the collection is decreasing
//
//    assert.IsDecreasing(t, []int{2, 1, 0})
//    assert.IsDecreasing(t, []float{2, 1})
//    assert.IsDecreasing(t, []string{"b", "a"})
func IsDecreasing(t TestingT, object interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return isOrdered(t, object, []CompareType{compareGreater}, " is not greater than ", msgAndArgs)
}

// IsNonDecreasing asserts that the collection is not decreasing
//
//    assert.IsNonDecreasing(t, []int{1, 1, 2})
//    assert.IsNonDecreasing(t, []float{1, 2})
//    assert.IsNonDecreasing(t, []string{"a", "b"})
func IsNonDecreasing(t TestingT, object interface{}, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}
	return isOrdered(t, object, []CompareType{compareLess, compareEqual}, " is not less than or equal to ", msgAndArgs)
}
