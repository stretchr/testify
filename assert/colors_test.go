package assert

import "testing"

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
