// +build !go1.17

package assert

import "reflect"

// Older versions of Go does not have the reflect.Value.CanConvert
// method.
func canConvert(value reflect.Value, to reflect.Type) bool {
	return false
}
