package assert

import (
	"reflect"
	"sync"
)

const equalMethod = "Equal"

var equallerCacheMu sync.RWMutex
var equallerCache map[reflect.Type]bool

func init() {
	equallerCache = make(map[reflect.Type]bool, 0)
}

func isEqualler(t reflect.Type) bool {
	isEqualler, cached := isEquallerCached(t)
	if !cached {
		isEqualler = determineIsEqualler(t)
		setIsEquallerCached(t, isEqualler)
	}

	return isEqualler
}

func determineIsEqualler(t reflect.Type) bool {
	equalMethod, hasEqualMethod := t.MethodByName(equalMethod)
	if hasEqualMethod {
		// should have only 1 return value which should be a bool
		// and should have exactly 2 arguments (pointer method so first is self)
		//  of which the 2nd argument should also be of its own type
		if equalMethod.Type.NumOut() != 1 || equalMethod.Type.Out(0).Kind() != reflect.Bool {
			return false
		} else if equalMethod.Type.NumIn() != 2 || !t.ConvertibleTo(equalMethod.Type.In(1)) {
			return false
		} else {
			return true
		}
	}

	return false
}

func isEquallerCached(t reflect.Type) (bool, bool) {
	equallerCacheMu.RLock()
	defer equallerCacheMu.RUnlock()

	isEqualler, cached := equallerCache[t]

	return isEqualler, cached
}

func setIsEquallerCached(t reflect.Type, isEqualler bool) {
	equallerCacheMu.Lock()
	defer equallerCacheMu.Unlock()

	equallerCache[t] = isEqualler
}
