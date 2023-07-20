package check

import (
	"bytes"
	"reflect"
)

// ObjectsAreEqual determines if two objects are considered equal.
func ObjectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

// CopyExportedFields iterates downward through nested data structures and creates a copy
// that only contains the exported struct fields.
func CopyExportedFields(expected interface{}) interface{} {
	if IsNil(expected) {
		return expected
	}

	expectedType := reflect.TypeOf(expected)
	expectedKind := expectedType.Kind()
	expectedValue := reflect.ValueOf(expected)

	switch expectedKind {
	case reflect.Struct:
		result := reflect.New(expectedType).Elem()
		for i := 0; i < expectedType.NumField(); i++ {
			field := expectedType.Field(i)
			isExported := field.IsExported()
			if isExported {
				fieldValue := expectedValue.Field(i)
				if IsNil(fieldValue) || IsNil(fieldValue.Interface()) {
					continue
				}
				newValue := CopyExportedFields(fieldValue.Interface())
				result.Field(i).Set(reflect.ValueOf(newValue))
			}
		}
		return result.Interface()

	case reflect.Ptr:
		result := reflect.New(expectedType.Elem())
		unexportedRemoved := CopyExportedFields(expectedValue.Elem().Interface())
		result.Elem().Set(reflect.ValueOf(unexportedRemoved))
		return result.Interface()

	case reflect.Array, reflect.Slice:
		result := reflect.MakeSlice(expectedType, expectedValue.Len(), expectedValue.Len())
		for i := 0; i < expectedValue.Len(); i++ {
			index := expectedValue.Index(i)
			if IsNil(index) {
				continue
			}
			unexportedRemoved := CopyExportedFields(index.Interface())
			result.Index(i).Set(reflect.ValueOf(unexportedRemoved))
		}
		return result.Interface()

	case reflect.Map:
		result := reflect.MakeMap(expectedType)
		for _, k := range expectedValue.MapKeys() {
			index := expectedValue.MapIndex(k)
			unexportedRemoved := CopyExportedFields(index.Interface())
			result.SetMapIndex(k, reflect.ValueOf(unexportedRemoved))
		}
		return result.Interface()

	default:
		return expected
	}
}

// ObjectsExportedFieldsAreEqual determines if the exported (public) fields of two objects are
// considered equal. This comparison of only exported fields is applied recursively to nested data
// structures.
func ObjectsExportedFieldsAreEqual(expected, actual interface{}) bool {
	expectedCleaned := CopyExportedFields(expected)
	actualCleaned := CopyExportedFields(actual)
	return ObjectsAreEqualValues(expectedCleaned, actualCleaned)
}

// ObjectsAreEqualValues gets whether two objects are equal, or if their
// values are equal.
func ObjectsAreEqualValues(expected, actual interface{}) bool {
	if ObjectsAreEqual(expected, actual) {
		return true
	}

	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		return false
	}
	expectedValue := reflect.ValueOf(expected)
	if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
		// Attempt comparison after type conversion
		return reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual)
	}

	return false
}

// containsKind checks if a specified kind in the slice of kinds.
func containsKind(kinds []reflect.Kind, kind reflect.Kind) bool {
	for i := 0; i < len(kinds); i++ {
		if kind == kinds[i] {
			return true
		}
	}

	return false
}

// IsNil checks if a specified object is nil or not, without Failing.
func IsNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	isNilableKind := containsKind(
		[]reflect.Kind{
			reflect.Chan, reflect.Func,
			reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Slice, reflect.UnsafePointer},
		kind)

	if isNilableKind && value.IsNil() {
		return true
	}

	return false
}
