package assert

import (
	"os"
	"strings"
	"sync"
)

// DiffDecorator decorates lines of the unified diff that is printed when an
// assertion comparing two values fails (see the "Diff:" section of failure
// output). Implementations can, for example, colorize the diff to make it
// easier to spot the differences.
//
// The default decorator leaves every line undecorated, preserving testify's
// historical output. Use SetDiffDecorator to install a custom decorator, e.g.
// NewColorDiffDecorator to colorize "expected" (green) vs "actual" (red).
type DiffDecorator interface {
	// DecorateDiffLine decorates a single line of a unified diff (without the
	// trailing newline) and returns the decorated line.
	DecorateDiffLine(line string) string
}

// noopDiffDecorator is the default DiffDecorator: it returns lines unchanged.
type noopDiffDecorator struct{}

func (noopDiffDecorator) DecorateDiffLine(line string) string { return line }

var (
	diffDecoratorMu sync.RWMutex
	diffDecorator   DiffDecorator = noopDiffDecorator{}
)

// SetDiffDecorator installs the DiffDecorator used for diff output of failing
// assertions. Passing nil restores the default (undecorated) behavior.
func SetDiffDecorator(d DiffDecorator) {
	diffDecoratorMu.Lock()
	defer diffDecoratorMu.Unlock()
	if d == nil {
		diffDecorator = noopDiffDecorator{}
		return
	}
	diffDecorator = d
}

// GetDiffDecorator returns the currently installed DiffDecorator.
func GetDiffDecorator() DiffDecorator {
	diffDecoratorMu.RLock()
	defer diffDecoratorMu.RUnlock()
	return diffDecorator
}

// ANSI escape codes used by ColorDiffDecorator.
const (
	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiRed    = "\x1b[31m"
	ansiGreen  = "\x1b[32m"
	ansiCyan   = "\x1b[36m"
	ansiYellow = "\x1b[33m"
)

// ColorDiffDecorator is a DiffDecorator that colorizes unified diff lines with
// ANSI escape codes: the "Expected" header and removed ("-") lines are green,
// the "Actual" header and added ("+") lines are red, and hunk headers ("@@")
// are cyan.
//
// When Enabled is false, lines are returned unchanged.
type ColorDiffDecorator struct {
	Enabled bool
}

// NewColorDiffDecorator returns a ColorDiffDecorator that is enabled only when
// color output is appropriate: the NO_COLOR environment variable
// (https://no-color.org) is unset and standard output is a terminal. In
// non-interactive environments (CI logs, redirected output) it stays disabled
// so no control characters end up in the logs.
func NewColorDiffDecorator() *ColorDiffDecorator {
	return &ColorDiffDecorator{Enabled: colorOutputSupported()}
}

// colorOutputSupported reports whether colorized output should be emitted.
func colorOutputSupported() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// DecorateDiffLine implements DiffDecorator.
func (c ColorDiffDecorator) DecorateDiffLine(line string) string {
	if !c.Enabled {
		return line
	}
	switch {
	case strings.HasPrefix(line, "--- Expected"), strings.HasPrefix(line, "+++ Actual"):
		return ansiBold + line + ansiReset
	case strings.HasPrefix(line, "-"):
		return ansiGreen + line + ansiReset
	case strings.HasPrefix(line, "+"):
		return ansiRed + line + ansiReset
	case strings.HasPrefix(line, "@@"):
		return ansiCyan + line + ansiReset
	default:
		return line
	}
}

// decorateDiff applies the installed DiffDecorator to every line of a unified
// diff string.
func decorateDiff(unifiedDiff string) string {
	decorator := GetDiffDecorator()
	if _, isNoop := decorator.(noopDiffDecorator); isNoop {
		return unifiedDiff
	}
	lines := strings.Split(unifiedDiff, "\n")
	for i, line := range lines {
		// Split keeps no trailing newline on the last element; difflib lines
		// keep theirs, so strip, decorate and re-attach.
		nl := ""
		if strings.HasSuffix(line, "\n") {
			line, nl = strings.TrimSuffix(line, "\n"), "\n"
		}
		lines[i] = decorator.DecorateDiffLine(line) + nl
	}
	return strings.Join(lines, "\n")
}
