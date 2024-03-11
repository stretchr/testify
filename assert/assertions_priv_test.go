package assert

import (
	"bufio"
	"strings"
	"testing"
)

func Test_indentMessageLines(t *testing.T) {
	tt := []struct {
		name            string
		msg             string
		longestLabelLen int
		expected        string
	}{
		{
			name:            "single line",
			msg:             "Hello World\n",
			longestLabelLen: 11,
			expected:        "Hello World",
		},
		{
			name:            "multi line",
			msg:             "Hello\nWorld\n",
			longestLabelLen: 11,
			expected:        "Hello\n\t            \tWorld",
		},
		{
			name:            "single line - over the bufio default limit",
			msg:             strings.Repeat("hello ", bufio.MaxScanTokenSize+10),
			longestLabelLen: 11,
			expected:        strings.Repeat("hello ", bufio.MaxScanTokenSize+10),
		},
		{
			name:            "multi line - over the bufio default limit",
			msg:             strings.Repeat("hello\n", bufio.MaxScanTokenSize+10),
			longestLabelLen: 3,
			expected: strings.TrimSpace(
				strings.TrimPrefix(strings.Repeat("\thello\n\t    ", bufio.MaxScanTokenSize+10), "\t"),
			),
		},
		{
			name:            "single line - just under the bufio default limit",
			msg:             strings.Repeat("hello ", bufio.MaxScanTokenSize-10),
			longestLabelLen: 11,
			expected:        strings.Repeat("hello ", bufio.MaxScanTokenSize-10),
		},
		{
			name:            "multi line - just under the bufio default limit",
			msg:             strings.Repeat("hello\n", bufio.MaxScanTokenSize-10),
			longestLabelLen: 3,
			expected: strings.TrimSpace(
				strings.TrimPrefix(strings.Repeat("\thello\n\t    ", bufio.MaxScanTokenSize-10), "\t"),
			),
		},
		{
			name:            "single line - equal to the bufio default limit",
			msg:             strings.Repeat("hello ", bufio.MaxScanTokenSize),
			longestLabelLen: 11,
			expected:        strings.Repeat("hello ", bufio.MaxScanTokenSize),
		},
		{
			name:            "multi line - equal to the bufio default limit",
			msg:             strings.Repeat("hello\n", bufio.MaxScanTokenSize),
			longestLabelLen: 3,
			expected: strings.TrimSpace(
				strings.TrimPrefix(strings.Repeat("\thello\n\t    ", bufio.MaxScanTokenSize), "\t"),
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			Equal(t, tc.expected, indentMessageLines(tc.msg, tc.longestLabelLen))
		})
	}
}
