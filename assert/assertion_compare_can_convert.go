// +build go1.17

package assert

import "reflect"

// Wrapper around reflect.Value.CanConvert, for compatability
// reasons.
func canConvert(value reflect.Value, to reflect.Type) bool {
	return value.CanConvert(to)
}
