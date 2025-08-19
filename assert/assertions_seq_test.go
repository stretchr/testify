//go:build go1.23 || goexperiment.rangefunc

package assert

import (
	"fmt"
	"testing"
)

// go.mod version is set to 1.17, which precludes the use of generics (even though this file wouldn't be taken into
// account per the build tags).

func intSeq(s ...int) func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for _, elem := range s {
			if !yield(elem) {
				break
			}
		}
	}
}

func strSeq(s ...string) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		for _, elem := range s {
			if !yield(elem) {
				break
			}
		}
	}
}

func TestElementsMatch_Seq(t *testing.T) {
	mockT := new(testing.T)

	cases := []struct {
		expected interface{}
		actual   interface{}
		result   bool
	}{
		{intSeq(), intSeq(), true},
		{intSeq(1), intSeq(1), true},
		{intSeq(1, 1), intSeq(1, 1), true},
		{intSeq(1, 2), intSeq(1, 2), true},
		{intSeq(1, 2), intSeq(2, 1), true},
		{strSeq("hello", "world"), strSeq("world", "hello"), true},
		{strSeq("hello", "hello"), strSeq("hello", "hello"), true},
		{strSeq("hello", "hello", "world"), strSeq("hello", "world", "hello"), true},
		{intSeq(), nil, true},

		// not matching
		{intSeq(1), intSeq(1, 1), false},
		{intSeq(1, 2), intSeq(2, 2), false},
		{strSeq("hello", "hello"), strSeq("hello"), false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("ElementsMatch(%#v, %#v)", seqToSlice(c.expected), seqToSlice(c.actual)), func(t *testing.T) {
			res := ElementsMatch(mockT, c.actual, c.expected)

			if res != c.result {
				t.Errorf("ElementsMatch(%#v, %#v) should return %v", seqToSlice(c.actual), seqToSlice(c.expected), c.result)
			}
		})
	}
}

func TestNotElementsMatch_Seq(t *testing.T) {
	mockT := new(testing.T)

	cases := []struct {
		expected interface{}
		actual   interface{}
		result   bool
	}{
		// not matching
		{intSeq(1), intSeq(), true},
		{intSeq(), intSeq(2), true},
		{intSeq(1), intSeq(2), true},
		{intSeq(1), intSeq(1, 1), true},
		{intSeq(1, 2), intSeq(3, 4), true},
		{intSeq(3, 4), intSeq(1, 2), true},
		{intSeq(1, 1, 2, 3), intSeq(1, 2, 3), true},
		{strSeq("hello"), strSeq("world"), true},
		{strSeq("hello", "hello"), strSeq("world", "world"), true},

		// matching
		{intSeq(), nil, false},
		{intSeq(), intSeq(), false},
		{intSeq(1), intSeq(1), false},
		{intSeq(1, 1), intSeq(1, 1), false},
		{intSeq(1, 2), intSeq(2, 1), false},
		{intSeq(1, 1, 2), intSeq(1, 2, 1), false},
		{strSeq("hello", "world"), strSeq("world", "hello"), false},
		{strSeq("hello", "hello"), strSeq("hello", "hello"), false},
		{strSeq("hello", "hello", "world"), strSeq("hello", "world", "hello"), false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("NotElementsMatch(%#v, %#v)", seqToSlice(c.expected), seqToSlice(c.actual)), func(t *testing.T) {
			res := NotElementsMatch(mockT, c.actual, c.expected)

			if res != c.result {
				t.Errorf("NotElementsMatch(%#v, %#v) should return %v", seqToSlice(c.actual), seqToSlice(c.expected), c.result)
			}
		})
	}
}

func TestContainsNotContains_Seq(t *testing.T) {

	type A struct {
		Name, Value string
	}
	complexSeq := func(s ...*A) func(yield func(*A) bool) {
		return func(yield func(*A) bool) {
			for _, elem := range s {
				if !yield(elem) {
					break
				}
			}
		}
	}

	list := []string{"Foo", "Bar"}

	complexList := []*A{
		{"b", "c"},
		{"d", "e"},
		{"g", "h"},
		{"j", "k"},
	}

	cases := []struct {
		expected interface{}
		actual   interface{}
		result   bool
	}{
		{strSeq(list...), "Bar", true},
		{strSeq(list...), "Salut", false},
		{complexSeq(complexList...), &A{"g", "h"}, true},
		{complexSeq(complexList...), &A{"g", "e"}, false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Contains(%#v, %#v)", seqToSlice(c.expected), seqToSlice(c.actual)), func(t *testing.T) {
			mockT := new(testing.T)
			res := Contains(mockT, c.expected, c.actual)

			if res != c.result {
				if res {
					t.Errorf(
						"Contains(%#v, %#v) should return true:\n\t%#v contains %#v",
						seqToSlice(c.expected), seqToSlice(c.actual), seqToSlice(c.expected), seqToSlice(c.actual))
				} else {
					t.Errorf(
						"Contains(%#v, %#v) should return false:\n\t%#v does not contain %#v",
						seqToSlice(c.expected), seqToSlice(c.actual), seqToSlice(c.expected), seqToSlice(c.actual))
				}
			}
		})
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("NotContains(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			mockT := new(testing.T)
			res := NotContains(mockT, c.expected, c.actual)

			// NotContains should be inverse of Contains. If it's not, something is wrong
			if res == Contains(mockT, c.expected, c.actual) {
				if res {
					t.Errorf("NotContains(%#v, %#v) should return true:\n\t%#v does not contains %#v", c.expected, c.actual, c.expected, c.actual)
				} else {
					t.Errorf("NotContains(%#v, %#v) should return false:\n\t%#v contains %#v", c.expected, c.actual, c.expected, c.actual)
				}
			}
		})
	}
}
