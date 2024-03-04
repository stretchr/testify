package assert

import (
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
			name:            "single line - extremely long",
			msg:             strings.Repeat("hello ", 20000),
			longestLabelLen: 11,
			expected:        strings.Repeat("hello ", 20000),
		},
		{
			name:            "multi line - extremely long",
			msg:             strings.Repeat("hello\n", 20000),
			longestLabelLen: 3,
			expected: strings.TrimSpace(
				strings.TrimPrefix(strings.Repeat("\thello\n\t    ", 20000), "\t"),
			),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			Equal(t, tc.expected, indentMessageLines(tc.msg, tc.longestLabelLen))
		})
	}
}
