package difflib

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
)

func splitChars(s string) []string {
	chars := make([]string, 0, len(s))
	// Assume ASCII inputs
	for i := 0; i != len(s); i++ {
		chars = append(chars, string(s[i]))
	}
	return chars
}

func TestSequenceMatcherRatio(t *testing.T) {
	s := NewMatcher(splitChars("abcd"), splitChars("bcde"))
	ratio := s.Ratio()
	if ratio != 0.75 {
		t.Errorf("unexpected ratio: %v", ratio)
	}

	ratio = s.QuickRatio()
	if ratio != 0.75 {
		t.Errorf("unexpected quick ratio: %v", ratio)
	}

	ratio = s.RealQuickRatio()
	if ratio != 1.0 {
		t.Errorf("unexpected real quick ratio: %v", ratio)
	}
}

func TestGetOptCodes(t *testing.T) {
	a := "qabxcd"
	b := "abycdf"
	s := NewMatcher(splitChars(a), splitChars(b))
	w := &bytes.Buffer{}
	for _, op := range s.GetOpCodes() {
		fmt.Fprintf(w, "%s a[%d:%d], (%s) b[%d:%d] (%s)\n", string(op.Tag),
			op.I1, op.I2, a[op.I1:op.I2], op.J1, op.J2, b[op.J1:op.J2])
	}
	result := string(w.Bytes())
	expected := `d a[0:1], (q) b[0:0] ()
e a[1:3], (ab) b[0:2] (ab)
r a[3:4], (x) b[2:3] (y)
e a[4:6], (cd) b[3:5] (cd)
i a[6:6], () b[5:6] (f)
`
	if expected != result {
		t.Errorf("unexpected op codes: \n%s", result)
	}
}

func TestGroupedOpCodes(t *testing.T) {
	a := []string{}
	for i := 0; i != 39; i++ {
		a = append(a, fmt.Sprintf("%02d", i))
	}
	b := []string{}
	b = append(b, a[:8]...)
	b = append(b, " i")
	b = append(b, a[8:19]...)
	b = append(b, " x")
	b = append(b, a[20:22]...)
	b = append(b, a[27:34]...)
	b = append(b, " y")
	b = append(b, a[35:]...)
	s := NewMatcher(a, b)
	w := &bytes.Buffer{}
	for _, g := range s.GetGroupedOpCodes(-1) {
		fmt.Fprintf(w, "group\n")
		for _, op := range g {
			fmt.Fprintf(w, "  %s, %d, %d, %d, %d\n", string(op.Tag),
				op.I1, op.I2, op.J1, op.J2)
		}
	}
	result := string(w.Bytes())
	expected := `group
  e, 5, 8, 5, 8
  i, 8, 8, 8, 9
  e, 8, 11, 9, 12
group
  e, 16, 19, 17, 20
  r, 19, 20, 20, 21
  e, 20, 22, 21, 23
  d, 22, 27, 23, 23
  e, 27, 30, 23, 26
group
  e, 31, 34, 27, 30
  r, 34, 35, 30, 31
  e, 35, 38, 31, 34
`
	if expected != result {
		t.Errorf("unexpected op codes: \n%s", result)
	}
}

func splitLines(s string) []string {
	lines := []string{}
	for _, line := range strings.Split(s, "\n") {
		lines = append(lines, line+"\n")
	}
	return lines
}

func TestUnifiedDiff(t *testing.T) {
	a := `one
two
three
four`
	b := `zero
one
three
four`
	diff := UnifiedDiff{
		A:        splitLines(a),
		B:        splitLines(b),
		FromFile: "Original",
		FromDate: "2005-01-26 23:30:50",
		ToFile:   "Current",
		ToDate:   "2010-04-02 10:20:52",
		Context:  3,
	}
	result, err := GetUnifiedDiffString(diff)
	if err != nil {
		t.Errorf("unified diff failed: %s", err)
	}
	expected := `--- Original\t2005-01-26 23:30:50
+++ Current\t2010-04-02 10:20:52
@@ -1,4 +1,4 @@
+zero
 one
-two
 three
 four
`
	// TABs are a pain to preserve through editors
	expected = strings.Replace(expected, "\\t", "\t", -1)
	if expected != result {
		t.Errorf("unexpected diff result:\n%s", result)
	}
}

func assertAlmostEqual(t *testing.T, a, b float64, places int) {
	if math.Abs(a-b) > math.Pow10(-places) {
		t.Errorf("%.7f != %.7f", a, b)
	}
}

func assertEqual(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%v != %v", a, b)
	}
}

func rep(s string, count int) string {
	return strings.Repeat(s, count)
}

func TestWithAsciiOneInsert(t *testing.T) {
	sm := NewMatcher(splitChars(rep("b", 100)),
		splitChars("a"+rep("b", 100)))
	assertAlmostEqual(t, sm.Ratio(), 0.995, 3)
	assertEqual(t, sm.GetOpCodes(),
		[]OpCode{{'i', 0, 0, 0, 1}, {'e', 0, 100, 1, 101}})
	assertEqual(t, len(sm.bPopular), 0)

	sm = NewMatcher(splitChars(rep("b", 100)),
		splitChars(rep("b", 50)+"a"+rep("b", 50)))
	assertAlmostEqual(t, sm.Ratio(), 0.995, 3)
	assertEqual(t, sm.GetOpCodes(),
		[]OpCode{{'e', 0, 50, 0, 50}, {'i', 50, 50, 50, 51}, {'e', 50, 100, 51, 101}})
	assertEqual(t, len(sm.bPopular), 0)
}

func TestWithAsciiOnDelete(t *testing.T) {
	sm := NewMatcher(splitChars(rep("a", 40)+"c"+rep("b", 40)),
		splitChars(rep("a", 40)+rep("b", 40)))
	assertAlmostEqual(t, sm.Ratio(), 0.994, 3)
	assertEqual(t, sm.GetOpCodes(),
		[]OpCode{{'e', 0, 40, 0, 40}, {'d', 40, 41, 40, 40}, {'e', 41, 81, 40, 80}})
}

func TestWithAsciiBJunk(t *testing.T) {
	isJunk := func(s string) bool {
		return s == " "
	}
	sm := NewMatcherWithJunk(splitChars(rep("a", 40)+rep("b", 40)),
		splitChars(rep("a", 44)+rep("b", 40)), true, isJunk)
	assertEqual(t, sm.bJunk, map[string]struct{}{})

	sm = NewMatcherWithJunk(splitChars(rep("a", 40)+rep("b", 40)),
		splitChars(rep("a", 44)+rep("b", 40)+rep(" ", 20)), false, isJunk)
	assertEqual(t, sm.bJunk, map[string]struct{}{" ": struct{}{}})

	isJunk = func(s string) bool {
		return s == " " || s == "b"
	}
	sm = NewMatcherWithJunk(splitChars(rep("a", 40)+rep("b", 40)),
		splitChars(rep("a", 44)+rep("b", 40)+rep(" ", 20)), false, isJunk)
	assertEqual(t, sm.bJunk, map[string]struct{}{" ": struct{}{}, "b": struct{}{}})
}
