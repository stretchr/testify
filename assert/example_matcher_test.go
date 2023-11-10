package assert

import (
	"fmt"
	"strings"
	"testing"
)

// stringStartsWith is a type that implements Matcher to test the actual
// string has the expected prefix.
type stringStartsWith struct {
	expected string
}

func (s *stringStartsWith) Match(actual interface{}) bool {
	str, isStr := actual.(string)
	return isStr && strings.HasPrefix(str, s.expected)
}

func (s *stringStartsWith) Describe() string {
	return fmt.Sprintf("a string starting with %q", s.expected)
}

// StringStarting returns a Matcher that asserts that the actual value
// is a string that has the expected prefix.
//
// Wrapping the matcher type in a factory function is just a little syntactic sugar for
// its use in assert.Matches
func StringStarting(e string) Matcher {
	return &stringStartsWith{expected: e}
}

func ExampleMatcher() {
	t := &testing.T{} // provided by test

	Matches(t, StringStarting("hello"), "hello world")
}
