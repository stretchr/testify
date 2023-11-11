package assert

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	redMark   = "\033[0;31m"
	greenMark = "\033[0;32m"
	endMark   = "\033[0m"
)

var isTerminal = term.IsTerminal(int(os.Stdout.Fd()))

func redColored(i interface{}) interface{} {
	if isTerminal {
		return fmt.Sprintf(redMark+"%s"+endMark, i)
	}
	return i
}

func greenColored(i interface{}) interface{} {
	if isTerminal {
		return fmt.Sprintf(greenMark+"%s"+endMark, i)
	}
	return i
}
