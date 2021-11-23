package assert

import (
	"os"
	"regexp"
	"strconv"
	"sync"

	"golang.org/x/term"
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
	// Auto: Color terminal output but not piped output
	autoResponse := term.IsTerminal(int(os.Stdout.Fd()))

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

	// FIXME: Print a warning here

	return autoResponse
}

func colorize(plainText string) string {
	setShouldColorize.Do(func() {
		shouldColorize = figureOutShouldColorize()
	})

	if !shouldColorize {
		return plainText
	}

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
