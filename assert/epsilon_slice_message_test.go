package assert

import (
	"math"
	"strings"
	"testing"
)

func TestInEpsilonSliceFailureContext(t *testing.T) {
	cases := []struct {
		name     string
		expected interface{}
		actual   interface{}
		epsilon  float64
	}{
		{"nil", nil, []float64{1}, .01},
		{"not a slice", 1, 2, .01},
		{"different types", []float64{1}, []float32{1}, .01},
		{"shorter actual", []float64{1, 2}, []float64{1}, .01},
		{"longer actual", []float64{1}, []float64{1, 2}, .01},
		{"element", []float64{1, 2}, []float64{1, 3}, .01},
		{"non numerical", []string{"a"}, []string{"a"}, .01},
		{"zero expected", []float64{0}, []float64{1}, .01},
		{"NaN actual", []float64{1}, []float64{math.NaN()}, .01},
		{"NaN tolerance", []float64{1}, []float64{1}, math.NaN()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, call := range []struct {
				name string
				f    func(*bufferT) bool
			}{
				{"function", func(out *bufferT) bool {
					return InEpsilonSlice(out, tc.expected, tc.actual, tc.epsilon, "batch %d", 7)
				}},
				{"method", func(out *bufferT) bool {
					return New(out).InEpsilonSlice(tc.expected, tc.actual, tc.epsilon, "batch %d", 7)
				}},
			} {
				t.Run(call.name, func(t *testing.T) {
					out := new(bufferT)
					if call.f(out) {
						t.Fatal("expected the assertion to fail")
					}
					if !strings.Contains(out.buf.String(), "batch 7") {
						t.Fatalf("caller context missing:\n%s", out.buf.String())
					}
				})
			}
		})
	}
}

type epsilonSliceMessageProbe struct{ calls int }

func (p *epsilonSliceMessageProbe) String() string {
	p.calls++
	return "batch 7"
}

func TestInEpsilonSliceFormatsMessageOnlyOnFailure(t *testing.T) {
	for _, fail := range []bool{false, true} {
		out := new(bufferT)
		message := new(epsilonSliceMessageProbe)
		actual := []float64{1, 2, 3}
		if fail {
			actual[2] = 4
		}
		if got := InEpsilonSlice(out, []float64{1, 2, 3}, actual, .01, message); got == fail {
			t.Fatalf("assertion result = %v, want %v", got, !fail)
		}
		wantCalls := 0
		if fail {
			wantCalls = 1
			if !strings.Contains(out.buf.String(), "at index 2: batch 7") {
				t.Fatalf("index or caller context missing:\n%s", out.buf.String())
			}
		}
		if message.calls != wantCalls {
			t.Errorf("message formatted %d times, want %d (failure=%v)", message.calls, wantCalls, fail)
		}
	}
}

func TestInEpsilonSliceMessageFormatting(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []interface{}
		want string
	}{
		{"literal percent", []interface{}{"100% complete"}, "at index 1: 100% complete"},
		{"indexed format", []interface{}{"%[2]s %[1]d", 7, "batch"}, "at index 1: batch 7"},
		{"no message", nil, "at index 1\n"},
		{"empty message", []interface{}{""}, "at index 1\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bufferT)
			if InEpsilonSlice(out, []float64{1, 2}, []float64{1, 3}, .01, tc.args...) {
				t.Fatal("expected the assertion to fail")
			}
			if !strings.Contains(out.buf.String(), tc.want) {
				t.Fatalf("want %q in output:\n%s", tc.want, out.buf.String())
			}
		})
	}
}
