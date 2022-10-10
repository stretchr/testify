package equal

import (
	"fmt"
	"testing"
	"time"
)

func TestObjectsAreEqual(t *testing.T) {
	cases := []struct {
		expected interface{}
		actual   interface{}
		result   bool
	}{
		// cases that are expected to be equal
		{"Hello World", "Hello World", true},
		{123, 123, true},
		{123.5, 123.5, true},
		{[]byte("Hello World"), []byte("Hello World"), true},
		{nil, nil, true},

		// cases that are expected not to be equal
		{map[int]int{5: 10}, map[int]int{10: 20}, false},
		{'x', "x", false},
		{"x", 'x', false},
		{0, 0.1, false},
		{0.1, 0, false},
		{time.Now, time.Now, false},
		{func() {}, func() {}, false},
		{uint32(10), int32(10), false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("ObjectsAreEqual(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := ObjectsAreEqual(c.expected, c.actual)

			if res != c.result {
				t.Errorf("ObjectsAreEqual(%#v, %#v) should return %#v", c.expected, c.actual, c.result)
			}

		})
	}

	// Cases where type differ but values are equal
	if !ObjectsAreEqualValues(uint32(10), int32(10)) {
		t.Error("ObjectsAreEqualValues should return true")
	}
	if ObjectsAreEqualValues(0, nil) {
		t.Fail()
	}
	if ObjectsAreEqualValues(nil, 0) {
		t.Fail()
	}

}
