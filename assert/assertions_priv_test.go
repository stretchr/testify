package assert

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func Test_indentMessageLines(t *testing.T) {
	tt := []struct {
		name            string
		longestLabelLen int

		// the input is constructed based on the below parameters
		bytesPerLine int
		lineCount    int
	}{
		{
			name:            "single line - over the bufio default limit",
			longestLabelLen: 11,
			bytesPerLine:    bufio.MaxScanTokenSize + 10,
			lineCount:       1,
		},
		{
			name:            "multi line - over the bufio default limit",
			longestLabelLen: 11,
			bytesPerLine:    bufio.MaxScanTokenSize + 10,
			lineCount:       3,
		},
		{
			name:            "single line - just under the bufio default limit",
			longestLabelLen: 11,
			bytesPerLine:    bufio.MaxScanTokenSize - 10,
			lineCount:       1,
		},
		{
			name:            "single line - just under the bufio default limit",
			longestLabelLen: 11,
			bytesPerLine:    bufio.MaxScanTokenSize - 10,
			lineCount:       1,
		},
		{
			name:            "single line - equal to the bufio default limit",
			longestLabelLen: 11,
			bytesPerLine:    bufio.MaxScanTokenSize,
			lineCount:       1,
		},
		{
			name:            "multi line - equal to the bufio default limit",
			longestLabelLen: 11,
			bytesPerLine:    bufio.MaxScanTokenSize,
			lineCount:       3,
		},
		{
			name:            "longest label length is zero",
			longestLabelLen: 0,
			bytesPerLine:    10,
			lineCount:       1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var input bytes.Buffer
			for i := 0; i < tc.lineCount; i++ {
				input.WriteString(strings.Repeat("#", tc.bytesPerLine))
				input.WriteRune('\n')
			}

			output := indentMessageLines(
				strings.TrimSpace(input.String()), tc.longestLabelLen,
			)
			outputLines := strings.Split(output, "\n")
			for i, line := range outputLines {
				if i > 0 {
					// count the leading white spaces. It should be equal to the longest
					// label length + 3. The +3 is to account for the 2 '\t' and 1 extra
					// space. Read the comment in the function for more context
					Equal(t, tc.longestLabelLen+3, strings.Index(line, "#"))
				}

				Len(t, strings.TrimSpace(line), tc.bytesPerLine)
			}
		})
	}
}
