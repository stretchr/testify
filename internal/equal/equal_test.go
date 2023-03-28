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

func TestObjectsExportedFieldsAreEqual(t *testing.T) {
	type Nested struct {
		Exported    interface{}
		notExported interface{}
	}

	type S struct {
		Exported1    interface{}
		Exported2    Nested
		notExported1 interface{}
		notExported2 Nested
	}

	type S2 struct {
		foo interface{}
	}

	cases := []struct {
		expected interface{}
		actual   interface{}
		result   bool
	}{
		{S{1, Nested{2, 3}, 4, Nested{5, 6}}, S{1, Nested{2, 3}, 4, Nested{5, 6}}, true},
		{S{1, Nested{2, 3}, 4, Nested{5, 6}}, S{1, Nested{2, 3}, "a", Nested{5, 6}}, true},
		{S{1, Nested{2, 3}, 4, Nested{5, 6}}, S{1, Nested{2, 3}, 4, Nested{5, "a"}}, true},
		{S{1, Nested{2, 3}, 4, Nested{5, 6}}, S{1, Nested{2, 3}, 4, Nested{"a", "a"}}, true},
		{S{1, Nested{2, 3}, 4, Nested{5, 6}}, S{1, Nested{2, "a"}, 4, Nested{5, 6}}, true},
		{S{1, Nested{2, 3}, 4, Nested{5, 6}}, S{"a", Nested{2, 3}, 4, Nested{5, 6}}, false},
		{S{1, Nested{2, 3}, 4, Nested{5, 6}}, S{1, Nested{"a", 3}, 4, Nested{5, 6}}, false},
		{S{1, Nested{2, 3}, 4, Nested{5, 6}}, S2{1}, false},
		{1, S{1, Nested{2, 3}, 4, Nested{5, 6}}, false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("ObjectsExportedFieldsAreEqual(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := ObjectsExportedFieldsAreEqual(c.expected, c.actual)

			if res != c.result {
				t.Errorf("ObjectsExportedFieldsAreEqual(%#v, %#v) should return %#v", c.expected, c.actual, c.result)
			}

		})
	}
}
