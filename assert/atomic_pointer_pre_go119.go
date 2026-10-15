//go:build !go1.19
// +build !go1.19

package assert

func objectsAreEqualAtomic(expected, actual interface{}) (equal, handled bool) {
	return false, false
}
