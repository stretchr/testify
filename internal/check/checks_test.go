package check

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

type S3 struct {
	Exported1 *Nested
	Exported2 *Nested
}

type S4 struct {
	Exported1 []*Nested
}

type S5 struct {
	Exported Nested
}

type S6 struct {
	Exported   string
	unexported string
}

func TestObjectsExportedFieldsAreEqual(t *testing.T) {

	intValue := 1

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

		{S3{&Nested{1, 2}, &Nested{3, 4}}, S3{&Nested{1, 2}, &Nested{3, 4}}, true},
		{S3{nil, &Nested{3, 4}}, S3{nil, &Nested{3, 4}}, true},
		{S3{&Nested{1, 2}, &Nested{3, 4}}, S3{&Nested{1, 2}, &Nested{3, "b"}}, true},
		{S3{&Nested{1, 2}, &Nested{3, 4}}, S3{&Nested{1, "a"}, &Nested{3, "b"}}, true},
		{S3{&Nested{1, 2}, &Nested{3, 4}}, S3{&Nested{"a", 2}, &Nested{3, 4}}, false},
		{S3{&Nested{1, 2}, &Nested{3, 4}}, S3{}, false},
		{S3{}, S3{}, true},

		{S4{[]*Nested{{1, 2}}}, S4{[]*Nested{{1, 2}}}, true},
		{S4{[]*Nested{{1, 2}}}, S4{[]*Nested{{1, 3}}}, true},
		{S4{[]*Nested{{1, 2}, {3, 4}}}, S4{[]*Nested{{1, "a"}, {3, "b"}}}, true},
		{S4{[]*Nested{{1, 2}, {3, 4}}}, S4{[]*Nested{{1, "a"}, {2, "b"}}}, false},

		{Nested{&intValue, 2}, Nested{&intValue, 2}, true},
		{Nested{&Nested{1, 2}, 3}, Nested{&Nested{1, "b"}, 3}, true},
		{Nested{&Nested{1, 2}, 3}, Nested{nil, 3}, false},

		{
			Nested{map[interface{}]*Nested{nil: nil}, 2},
			Nested{map[interface{}]*Nested{nil: nil}, 2},
			true,
		},
		{
			Nested{map[interface{}]*Nested{"a": nil}, 2},
			Nested{map[interface{}]*Nested{"a": nil}, 2},
			true,
		},
		{
			Nested{map[interface{}]*Nested{"a": nil}, 2},
			Nested{map[interface{}]*Nested{"a": {1, 2}}, 2},
			false,
		},
		{
			Nested{map[interface{}]Nested{"a": {1, 2}, "b": {3, 4}}, 2},
			Nested{map[interface{}]Nested{"a": {1, 5}, "b": {3, 7}}, 2},
			true,
		},
		{
			Nested{map[interface{}]Nested{"a": {1, 2}, "b": {3, 4}}, 2},
			Nested{map[interface{}]Nested{"a": {2, 2}, "b": {3, 4}}, 2},
			false,
		},
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

func TestCopyExportedFields(t *testing.T) {
	intValue := 1

	cases := []struct {
		input    interface{}
		expected interface{}
	}{
		{
			input:    Nested{"a", "b"},
			expected: Nested{"a", nil},
		},
		{
			input:    Nested{&intValue, 2},
			expected: Nested{&intValue, nil},
		},
		{
			input:    Nested{nil, 3},
			expected: Nested{nil, nil},
		},
		{
			input:    S{1, Nested{2, 3}, 4, Nested{5, 6}},
			expected: S{1, Nested{2, nil}, nil, Nested{}},
		},
		{
			input:    S3{},
			expected: S3{},
		},
		{
			input:    S3{&Nested{1, 2}, &Nested{3, 4}},
			expected: S3{&Nested{1, nil}, &Nested{3, nil}},
		},
		{
			input:    S3{Exported1: &Nested{"a", "b"}},
			expected: S3{Exported1: &Nested{"a", nil}},
		},
		{
			input: S4{[]*Nested{
				nil,
				{1, 2},
			}},
			expected: S4{[]*Nested{
				nil,
				{1, nil},
			}},
		},
		{
			input: S4{[]*Nested{
				{1, 2}},
			},
			expected: S4{[]*Nested{
				{1, nil}},
			},
		},
		{
			input: S4{[]*Nested{
				{1, 2},
				{3, 4},
			}},
			expected: S4{[]*Nested{
				{1, nil},
				{3, nil},
			}},
		},
		{
			input:    S5{Exported: Nested{"a", "b"}},
			expected: S5{Exported: Nested{"a", nil}},
		},
		{
			input:    S6{"a", "b"},
			expected: S6{"a", ""},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			output := CopyExportedFields(c.input)
			if !ObjectsAreEqualValues(c.expected, output) {
				t.Errorf("%#v, %#v should be equal", c.expected, output)
			}
		})
	}
}
