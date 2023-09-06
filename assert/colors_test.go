package assert

import (
	"fmt"
	"testing"
)

func BenchmarkColored(b *testing.B) {
	b.Run("benchMarkingString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			greenColored("helloWorld")
		}
	})
	b.Run("benchMarkingStruct", func(b *testing.B) {
		s := struct {
			a int
			b string
		}{3, "helloWorld"}
		for i := 0; i < b.N; i++ {
			greenColored(s)
		}
	})
}

// Not benchmarking `Equal` but the string formatting
// because it will make the benchmark fail.
func BenchmarkEqual(b *testing.B) {

	s := struct {
		a int
		b string
	}{3, "helloWorld"}

	b.Run("Base", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("Not equal: \n"+
				"expected: %v\n"+
				"actual  : %v%s", s, s, "")
		}
	})

	b.Run("Colored (test with and without terminal)", func(b *testing.B) {
		// This benchmark give very different results whether the output is a terminal or not.

		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("Not equal: \n"+
				"expected: %s\n"+
				"actual  : %s%s", greenColored(s), redColored(s), "")
		}
	})

}
