package main

import (
	"go/doc"
	"testing"
)

func TestCommentFormat(t *testing.T) {
	tests := []struct {
		name string
		doc  string
		want string
	}{
		{
			name: "nested calls",
			doc:  "Equal compares values.\n\n\tassert.Equal(t, outer(inner(1)), want)",
			want: "// Equalf compares values.\n// \n// \tassert.Equalf(t, outer(inner(1)), want, \"error message %s\", \"formatted\")",
		},
		{
			name: "parenthesis in string",
			doc:  "Equal compares values.\n\n\tassert.Equal(t, \")\", got)",
			want: "// Equalf compares values.\n// \n// \tassert.Equalf(t, \")\", got, \"error message %s\", \"formatted\")",
		},
		{
			name: "parenthesis in rune",
			doc:  "Equal compares values.\n\n\tassert.Equal(t, ')', got)",
			want: "// Equalf compares values.\n// \n// \tassert.Equalf(t, ')', got, \"error message %s\", \"formatted\")",
		},
		{
			name: "multiline call",
			doc:  "Equal compares values.\n\n\tassert.Equal(t, outer(\n\t\tinner(1)), want)",
			want: "// Equalf compares values.\n// \n// \tassert.Equalf(t, outer(\n// \t\tinner(1)), want, \"error message %s\", \"formatted\")",
		},
		{
			name: "multiple examples",
			doc:  "Equal compares values.\n\n\tassert.Equal(t, 1, got)\n\tassert.Equal(t, 2, got)",
			want: "// Equalf compares values.\n// \n// \tassert.Equalf(t, 1, got, \"error message %s\", \"formatted\")\n// \tassert.Equalf(t, 2, got, \"error message %s\", \"formatted\")",
		},
		{
			name: "parenthesis in raw string",
			doc:  "Equal compares values.\n\n\tassert.Equal(t, `)`, got)",
			want: "// Equalf compares values.\n// \n// \tassert.Equalf(t, `)`, got, \"error message %s\", \"formatted\")",
		},
		{
			name: "parenthesis in comment",
			doc:  "Equal compares values.\n\n\tassert.Equal(t, expected /* ) */, got)",
			want: "// Equalf compares values.\n// \n// \tassert.Equalf(t, expected /* ) */, got, \"error message %s\", \"formatted\")",
		},
		{
			name: "nested assertion",
			doc:  "Equal compares values.\n\n\tassert.Equal(t, true, assert.Equal(t, expected, got))",
			want: "// Equalf compares values.\n// \n// \tassert.Equalf(t, true, assert.Equalf(t, expected, got, \"error message %s\", \"formatted\"), \"error message %s\", \"formatted\")",
		},
		{
			name: "incomplete call",
			doc:  "Equal compares values.\n\n\tassert.Equal(t, outer(inner(1))",
			want: "// Equalf compares values.\n// \n// \tassert.Equalf(t, outer(inner(1))",
		},
		{
			name: "multiline trailing comma",
			doc:  "Equal compares values.\n\n\tassert.Equal(t, expected, got,\n\t)",
			want: "// Equalf compares values.\n// \n// \tassert.Equalf(t, expected, got,\n// \t\"error message %s\", \"formatted\")",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := testFunc{DocInfo: &doc.Func{Name: "Equal", Doc: tt.doc}}
			if got := f.CommentFormat(); got != tt.want {
				t.Errorf("CommentFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}
