package require

import (
	"fmt"
	"strings"
	"testing"
)

type epsilonSliceFailure struct {
	message string
	failed  bool
}

func (t *epsilonSliceFailure) Errorf(format string, args ...interface{}) {
	t.message += fmt.Sprintf(format, args...)
}

func (t *epsilonSliceFailure) FailNow() { t.failed = true }

func TestInEpsilonSliceFailureContext(t *testing.T) {
	for _, tc := range []struct {
		name     string
		expected interface{}
		actual   interface{}
	}{
		{"nil", nil, []float64{1}},
		{"not a slice", 1, 2},
		{"different types", []float64{1}, []float32{1}},
		{"length", []float64{1, 2}, []float64{1}},
		{"element", []float64{1, 2}, []float64{1, 3}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, call := range []struct {
				name string
				f    func(*epsilonSliceFailure)
			}{
				{"function", func(out *epsilonSliceFailure) {
					InEpsilonSlice(out, tc.expected, tc.actual, .01, "batch %d", 7)
				}},
				{"method", func(out *epsilonSliceFailure) {
					New(out).InEpsilonSlice(tc.expected, tc.actual, .01, "batch %d", 7)
				}},
				{"formatted function", func(out *epsilonSliceFailure) {
					InEpsilonSlicef(out, tc.expected, tc.actual, .01, "batch %d", 7)
				}},
				{"formatted method", func(out *epsilonSliceFailure) {
					New(out).InEpsilonSlicef(tc.expected, tc.actual, .01, "batch %d", 7)
				}},
			} {
				t.Run(call.name, func(t *testing.T) {
					out := new(epsilonSliceFailure)
					call.f(out)
					if !out.failed {
						t.Fatal("expected FailNow to be called")
					}
					if !strings.Contains(out.message, "batch 7") {
						t.Fatalf("caller context missing:\n%s", out.message)
					}
				})
			}
		})
	}
}
