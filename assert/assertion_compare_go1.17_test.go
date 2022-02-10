//go:build go1.17
// +build go1.17

// TODO: once support for Go 1.16 is dropped, this file can be
//       merged/removed with assertion_compare_can_convert.go and
//       assertion_compare_legacy.go

package assert

import (
	"reflect"
	"testing"
	"time"
)

func TestCompare17(t *testing.T) {
	type customTime time.Time
	for _, currCase := range []struct {
		less    interface{}
		greater interface{}
		cType   string
	}{
		{less: time.Now(), greater: time.Now().Add(time.Hour), cType: "time.Time"},
		{less: customTime(time.Now()), greater: customTime(time.Now().Add(time.Hour)), cType: "time.Time"},
	} {
		resLess, isComparable := compare(currCase.less, currCase.greater, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object should be comparable for type " + currCase.cType)
		}

		if resLess != compareLess {
			t.Errorf("object less (%v) should be less than greater (%v) for type "+currCase.cType,
				currCase.less, currCase.greater)
		}

		resGreater, isComparable := compare(currCase.greater, currCase.less, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object are comparable for type " + currCase.cType)
		}

		if resGreater != compareGreater {
			t.Errorf("object greater should be greater than less for type " + currCase.cType)
		}

		resEqual, isComparable := compare(currCase.less, currCase.less, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object are comparable for type " + currCase.cType)
		}

		if resEqual != 0 {
			t.Errorf("objects should be equal for type " + currCase.cType)
		}
	}
}
