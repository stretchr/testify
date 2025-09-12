package main

import (
	"testing"
)

func TestRequireCommentParseIf(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single line without if",
			input:    "// Simple comment line",
			expected: "// Simple comment line",
		},
		{
			name:     "simple if require block transformation",
			input:    "//\tif require.NotEmpty(t, obj) {\n//\t  require.Equal(t, \"two\", obj[1])\n//\t}",
			expected: "//\trequire.NotEmpty(t, obj) \n//\trequire.Equal(t, \"two\", obj[1])",
		},
		{
			name:     "no if block - should remain unchanged",
			input:    "// Contains function\n//\trequire.Contains(t, \"Hello World\", \"World\")",
			expected: "// Contains function\n//\trequire.Contains(t, \"Hello World\", \"World\")",
		},
		{
			name:     "mixed content with if block",
			input:    "//\t  actualObj, err := SomeFunction()\n//\tif require.NoError(t, err) {\n//\t\t do something\n//\t}",
			expected: "//\t  actualObj, err := SomeFunction()\n//\t  require.NoError(t, err) \n//\t  do something",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := requireCommentParseIf(tt.input)
			if result != tt.expected {
				t.Errorf("requireCommentParseIf() failed:\nInput: %q\nGot: %q\nWant: %q", tt.input, result, tt.expected)
			}
		})
	}
}
