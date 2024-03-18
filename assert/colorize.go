package assert

import (
	"regexp"
)

// Colorize struct to colorize failure message
type Colorize struct {
	t TestingT
}

// Errorf implements for TestingT interface
func (c Colorize) Errorf(format string, args ...interface{}) {
	c.t.Errorf(format, args...)
}

// Message implements colourize actual and expected at the very least
func (Colorize) Message(failureMessage string) string {
	re := regexp.MustCompile(`(.*)\n(expected.*)\n(actual.*)`)

	return string(
		re.ReplaceAll(
			[]byte(failureMessage),
			[]byte("\033[33m${1}\033[0m\n\033[32m${2}\033[0m\n\033[31m${3}\033[0m"),
		),
	)
}

// WrapColorize wrappers around TestingT to Colorize
func WrapColorize(t TestingT) Colorize {
	return Colorize{
		t: t,
	}
}

// Retrieve is a helper function to retrieve field of TestingT and Colorize or return self and nil
func Retrieve(t TestingT) (TestingT, *Colorize) {
	if c, ok := t.(Colorize); ok {
		return c.t, &c
	}

	return t, nil
}
