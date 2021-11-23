package assert

import "testing"

func TestColorize(t *testing.T) {
	Equal(t, "hello", doColorize("hello"))

	Equal(t, "\x1b[33mZ\x1b[0m\nexpected:\x1b[32m X\x1b[0m\nactual  :\x1b[31m Y\x1b[0m",
		doColorize("Z\nexpected: X\nactual  : Y"))
}
