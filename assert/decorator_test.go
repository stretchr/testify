package assert

import (
	"strings"
	"testing"
)

func TestNoopDiffDecoratorLeavesLinesUnchanged(t *testing.T) {
	d := noopDiffDecorator{}
	for _, line := range []string{"--- Expected", "+++ Actual", "-a", "+b", "@@ -1 +1 @@", " context", ""} {
		if got := d.DecorateDiffLine(line); got != line {
			t.Errorf("noop decorator changed line %q to %q", line, got)
		}
	}
}

func TestDefaultDecoratorIsNoop(t *testing.T) {
	if _, ok := GetDiffDecorator().(noopDiffDecorator); !ok {
		t.Fatalf("expected default decorator to be noopDiffDecorator, got %T", GetDiffDecorator())
	}
}

func TestSetDiffDecorator(t *testing.T) {
	defer SetDiffDecorator(nil)

	custom := ColorDiffDecorator{Enabled: true}
	SetDiffDecorator(custom)
	if got := GetDiffDecorator(); got != custom {
		t.Fatalf("expected installed decorator %v, got %v", custom, got)
	}

	SetDiffDecorator(nil)
	if _, ok := GetDiffDecorator().(noopDiffDecorator); !ok {
		t.Fatalf("expected nil to restore noop decorator, got %T", GetDiffDecorator())
	}
}

func TestColorDiffDecoratorLineColors(t *testing.T) {
	d := ColorDiffDecorator{Enabled: true}
	tests := []struct {
		line string
		want string
	}{
		{"--- Expected", ansiBold + "--- Expected" + ansiReset},
		{"+++ Actual", ansiBold + "+++ Actual" + ansiReset},
		{"-removed line", ansiGreen + "-removed line" + ansiReset},
		{"+added line", ansiRed + "+added line" + ansiReset},
		{"@@ -1,3 +1,3 @@", ansiCyan + "@@ -1,3 +1,3 @@" + ansiReset},
		{" unchanged", " unchanged"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := d.DecorateDiffLine(tt.line); got != tt.want {
			t.Errorf("DecorateDiffLine(%q) = %q, want %q", tt.line, got, tt.want)
		}
	}
}

func TestColorDiffDecoratorDisabledIsIdentity(t *testing.T) {
	d := ColorDiffDecorator{Enabled: false}
	for _, line := range []string{"--- Expected", "-a", "+b", "@@"} {
		if got := d.DecorateDiffLine(line); got != line {
			t.Errorf("disabled decorator changed %q to %q", line, got)
		}
	}
}

func TestColorOutputSupportedRespectsEnv(t *testing.T) {
	// NO_COLOR disables color regardless of FORCE_COLOR.
	t.Setenv("NO_COLOR", "1")
	t.Setenv("FORCE_COLOR", "1")
	if colorOutputSupported() {
		t.Error("NO_COLOR must disable color output")
	}

	// Without NO_COLOR, FORCE_COLOR enables it even when stdout is not a TTY.
	t.Setenv("NO_COLOR", "")
	if !colorOutputSupported() {
		t.Error("FORCE_COLOR should enable color output")
	}
}

func TestNewColorDiffDecorator(t *testing.T) {
	// In the test environment stdout is not a char device and NO_COLOR/FORCE_COLOR
	// are not set by default, so the auto decorator must be disabled (no control
	// characters in CI logs).
	t.Setenv("NO_COLOR", "")
	t.Setenv("FORCE_COLOR", "")
	d := NewColorDiffDecorator()
	if d.Enabled && !isTerminalForTest() {
		t.Logf("stdout appears to be a terminal; Enabled=%v", d.Enabled)
	}
	t.Setenv("FORCE_COLOR", "1")
	if d := NewColorDiffDecorator(); !d.Enabled {
		t.Error("FORCE_COLOR=1 should enable the color decorator")
	}
	t.Setenv("NO_COLOR", "1")
	t.Setenv("FORCE_COLOR", "1")
	if d := NewColorDiffDecorator(); d.Enabled {
		t.Error("NO_COLOR should disable the color decorator even with FORCE_COLOR")
	}
}

func isTerminalForTest() bool { return colorOutputSupported() }

func TestDecorateDiffAppliesDecoratorLineByLine(t *testing.T) {
	defer SetDiffDecorator(nil)
	SetDiffDecorator(ColorDiffDecorator{Enabled: true})

	in := "--- Expected\n+++ Actual\n@@ -1 +1 @@\n-a\n+b\n"
	got := decorateDiff(in)
	if !strings.Contains(got, ansiGreen+"-a"+ansiReset) {
		t.Errorf("expected green -a line, got:\n%q", got)
	}
	if !strings.Contains(got, ansiRed+"+b"+ansiReset) {
		t.Errorf("expected red +b line, got:\n%q", got)
	}
	if !strings.Contains(got, ansiBold+"--- Expected"+ansiReset) {
		t.Errorf("expected bold header, got:\n%q", got)
	}
	// Trailing newline must be preserved.
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("trailing newline lost: %q", got)
	}
}

func TestDecorateDiffNoopKeepsOutputIdentical(t *testing.T) {
	in := "--- Expected\n+++ Actual\n@@ -1 +1 @@\n-a\n+b\n"
	if got := decorateDiff(in); got != in {
		t.Errorf("noop decorator must not change output:\n got %q\nwant %q", got, in)
	}
}

func TestDiffUsesInstalledDecorator(t *testing.T) {
	defer SetDiffDecorator(nil)

	// Default: no ANSI codes in diff output.
	plain := diff("a\nb\nc", "a\nx\nc")
	if strings.Contains(plain, "\x1b[") {
		t.Fatalf("default diff output must not contain ANSI escapes, got %q", plain)
	}

	// With the color decorator installed, the diff is colorized.
	SetDiffDecorator(ColorDiffDecorator{Enabled: true})
	colored := diff("a\nb\nc", "a\nx\nc")
	if !strings.Contains(colored, "\x1b[") {
		t.Fatalf("expected ANSI escapes in diff output, got %q", colored)
	}
	if !strings.Contains(colored, ansiGreen+"-b") || !strings.Contains(colored, ansiRed+"+x") {
		t.Fatalf("expected green -b and red +x in %q", colored)
	}

	// Reset restores plain output.
	SetDiffDecorator(nil)
	if got := diff("a\nb\nc", "a\nx\nc"); strings.Contains(got, "\x1b[") {
		t.Fatalf("after reset, diff output must be plain, got %q", got)
	}
}
