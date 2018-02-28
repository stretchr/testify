package assert

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
)

// compare compares two objects
func cmpEqual(expected, actual interface{}) bool {
	return cmp.Equal(expected, actual, compareOptions(expected, actual)...)
}

// diff returns a diff of both values as long as both are of the same type and
// are a struct, map, slice or array. Otherwise it returns an empty string.
func cmpDiff(expected, actual interface{}) string {
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

	diff := cmp.Diff(expected, actual, compareOptions(expected, actual)...)
	if diff != "" {
		diff = "\n\nDiff:\n--- Expected\n+++ Actual\n" + diff
	}
	return diff
}

// compareOptions are cmp.Options used for cmp.Equal and cmp.Diff to compare
// two general objects for testing purposes
func compareOptions(expected, actual interface{}) cmp.Options {
	return cmp.Options{
		deepAllowUnexported(expected, actual),
		compareIdenticalPointers,
	}
}

// deepAllowUnexported returns option for cmp.Equal or cmp.Diff in which
// all unexported fields in the two compared types (recursively) are
// allowed.
// Code from https://github.com/google/go-cmp/issues/40 with modification
// to work with cyclic struct
func deepAllowUnexported(vs ...interface{}) cmp.Option {
	var (
		// allUnexported is a set of types to be added to the unexported list
		allUnexported = make(map[reflect.Type]bool)
		// visited are list of pointer which are visited during the recursive collection
		// of the referenced types.
		// It is used to detect cycles and prevent infinite recursion.
		visited = make(map[uintptr]bool)
	)

	// Collect all types from all given objects
	for _, v := range vs {
		structTypes(reflect.ValueOf(v), allUnexported, visited)
	}

	// Collect the referenced types
	var types []interface{}
	for t := range allUnexported {
		types = append(types, reflect.New(t).Elem().Interface())
	}

	// Return cmp option which allows all unexported fields in all the collected types
	return cmp.AllowUnexported(types...)
}

// structTypes is a recursive search for all referenced types from a given object.
// It searches recursively in all the given object fields and references, and put the
// collected type in the `m` set.
// It uses the `visited` set to detect cycles and prevent infinite recursion
func structTypes(v reflect.Value, m map[reflect.Type]bool, visited map[uintptr]bool) {
	if !v.IsValid() {
		return
	}

	// dive in according to the kind of the given object
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		// prevent infinite recursion
		if visited[v.Elem().UnsafeAddr()] {
			return
		}
		// remember jumping to a pointed address
		visited[v.Elem().UnsafeAddr()] = true
		structTypes(v.Elem(), m, visited)
	case reflect.Interface:
		if v.IsNil() {
			return
		}
		// search into the object that implement the interface
		structTypes(v.Elem(), m, visited)
	case reflect.Slice, reflect.Array:
		// recursively search in all the slice/array objects
		for i := 0; i < v.Len(); i++ {
			structTypes(v.Index(i), m, visited)
		}
	case reflect.Map:
		// recursively search in all the map values
		for _, k := range v.MapKeys() {
			structTypes(v.MapIndex(k), m, visited)
		}
	case reflect.Struct:
		// add the type to the collected types.
		m[v.Type()] = true
		// recursively search in all the struct fields
		for i := 0; i < v.NumField(); i++ {
			structTypes(v.Field(i), m, visited)
		}
	}
}

// compareIdenticalPointers is a cmp option that returns true if the two compared
// objects are pointers and are pointing on the same thing.
var compareIdenticalPointers = cmp.FilterPath(func(p cmp.Path) bool {
	// Filter for pointer kinds only.
	t := p.Last().Type()
	return t != nil && t.Kind() == reflect.Ptr
}, cmp.FilterValues(func(x, y interface{}) bool {
	// Filter for pointer values that are identical.
	vx := reflect.ValueOf(x)
	vy := reflect.ValueOf(y)
	return vx.IsValid() && vy.IsValid() && vx.Pointer() == vy.Pointer()
}, cmp.Comparer(func(_, _ interface{}) bool {
	// Consider them equal no matter what.
	return true
})))

func typeAndKind(v interface{}) (reflect.Type, reflect.Kind) {
	t := reflect.TypeOf(v)
	k := t.Kind()

	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	return t, k
}
