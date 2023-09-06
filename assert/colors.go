package assert

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const redMark = "\033[0;31m"
const greenMark = "\033[0;32m"
const endMark = "\033[0m"

func redColored(i interface{}) interface{} {
	if isTerminal {
		return redMark + fmt.Sprintf("%s", i) + endMark
	}
	return i
}

func greenColored(i interface{}) interface{} {
	if isTerminal {
		return greenMark + fmt.Sprintf("%s", i) + endMark
	}
	return i
}

var isTerminal = term.IsTerminal(int(os.Stdout.Fd()))
