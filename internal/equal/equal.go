package equal

import (
	"bytes"
	"reflect"
)

// ObjectsAreEqual determines if two objects are considered equal.
//
// This function does no assertion of any kind.
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

// ObjectsExportedFieldsAreEqual determines if the exported (public) fields of two structs are considered equal.
// If the two objects are not of the same type, or if either of them are not a struct, they are not considered equal.
//
// This function does no assertion of any kind.
func ObjectsExportedFieldsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)

	if expectedType != actualType {
		return false
	}

	if expectedType.Kind() != reflect.Struct || actualType.Kind() != reflect.Struct {
		return false
	}

	expectedValue := reflect.ValueOf(expected)
	actualValue := reflect.ValueOf(actual)

	for i := 0; i < expectedType.NumField(); i++ {
		field := expectedType.Field(i)
		isExported := field.PkgPath == "" // should use field.IsExported() but it's not available in Go 1.16.5
		if isExported {
			var equal bool
			if field.Type.Kind() == reflect.Struct {
				equal = ObjectsExportedFieldsAreEqual(expectedValue.Field(i).Interface(), actualValue.Field(i).Interface())
			} else {
				equal = ObjectsAreEqualValues(expectedValue.Field(i).Interface(), actualValue.Field(i).Interface())
			}

			if !equal {
				return false
			}
		}
	}
	return true
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
