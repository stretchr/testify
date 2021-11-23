package assert

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
)

const (
	expected_colon = "expected:"
	actual_colon   = "actual  :"
)

var (
	shouldColorize    bool
	setShouldColorize sync.Once
)

func figureOutShouldColorize() bool {
	// FIXME: The auto response should be true on terminals and false on
	// not-terminals, but my attempts at checking os.Stdout() terminalness have
	// all failed. Patches welcome!
	autoResponse := false

	hint := os.Getenv("TESTIFY_COLOR")
	if len(hint) == 0 {
		// Environment variable empty or missing: Auto
		return autoResponse
	}

	if hint == "auto" {
		return autoResponse
	}

	value, err := strconv.ParseBool(hint)
	if err == nil {
		return value
	}

	fmt.Printf("WARNING: TESTIFY_COLOR should be either true or false, was: \"%s\"\n", hint)

	return autoResponse
}

// Colorize string or not depending on what figureOutShouldColorize() says
func colorize(plainText string) string {
	setShouldColorize.Do(func() {
		shouldColorize = figureOutShouldColorize()
	})

	if shouldColorize {
		return doColorize(plainText)
	}
	return plainText
}

// Colorize unconditionally. You probably should call colorize() instead of this
// function.
func doColorize(plainText string) string {
	re := regexp.MustCompile(`(.*)\n` + expected_colon + `(.*)\n` + actual_colon + `(.*)`)

	return string(
		re.ReplaceAll(
			[]byte(plainText),
			[]byte("\033[33m${1}\033[0m\n"+
				expected_colon+"\033[32m${2}\033[0m\n"+
				actual_colon+"\033[31m${3}\033[0m"),
		),
	)
}
