//go:build go1.19
// +build go1.19

package assert

import (
	"reflect"
	"strings"
)

// objectsAreEqualAtomic compares loaded pointers without copying atomic values
// or dereferencing unsafe.Pointer. Only top-level *atomic.Pointer[T] values are
// handled; all other values keep reflect.DeepEqual's behavior.
func objectsAreEqualAtomic(expected, actual interface{}) (equal, handled bool) {
	typ := reflect.TypeOf(expected)
	if typ == nil || typ != reflect.TypeOf(actual) || typ.Kind() != reflect.Ptr || typ.Name() != "" {
		return false, false
	}
	elem := typ.Elem()
	if elem.Kind() != reflect.Struct || elem.PkgPath() != "sync/atomic" ||
		!strings.HasPrefix(elem.Name(), "Pointer[") || !strings.HasSuffix(elem.Name(), "]") {
		return false, false
	}
	load, ok := typ.MethodByName("Load")
	if !ok || load.PkgPath != "" || load.Type.IsVariadic() || load.Type.NumIn() != 1 ||
		load.Type.In(0) != typ || load.Type.NumOut() != 1 || load.Type.Out(0).Kind() != reflect.Ptr {
		return false, false
	}
	// Preserve pointer identity, including nil wrappers and non-reflexive values.
	if expected == actual {
		return true, true
	}
	left, right := reflect.ValueOf(expected), reflect.ValueOf(actual)
	if left.IsNil() || right.IsNil() {
		return false, true
	}
	// Each Load reads the original wrapper once. DeepEqual retains the loaded
	// pointer's identity semantics and handles cycles in the pointed-to values.
	a := load.Func.Call([]reflect.Value{left})[0].Interface()
	b := load.Func.Call([]reflect.Value{right})[0].Interface()
	return reflect.DeepEqual(a, b), true
}
