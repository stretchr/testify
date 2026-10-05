package mock_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

func zeroArgumentOptionA() {}
func zeroArgumentOptionB() {}

type optionsTestingT struct {
	failed bool
}

func (t *optionsTestingT) FailNow()                      { t.failed = true }
func (t *optionsTestingT) Logf(string, ...interface{})   {}
func (t *optionsTestingT) Errorf(string, ...interface{}) { t.failed = true }

func TestFunctionalOptionsZeroArgumentComparison(t *testing.T) {
	tests := []struct {
		name        string
		expected    []func()
		actual      []func()
		differences int
	}{
		{"equal", []func(){zeroArgumentOptionA, zeroArgumentOptionB}, []func(){zeroArgumentOptionA, zeroArgumentOptionB}, 0},
		{"different_first", []func(){zeroArgumentOptionA, zeroArgumentOptionB}, []func(){zeroArgumentOptionB, zeroArgumentOptionB}, 1},
		{"different_second", []func(){zeroArgumentOptionA, zeroArgumentOptionA}, []func(){zeroArgumentOptionA, zeroArgumentOptionB}, 1},
		{"different_third", []func(){zeroArgumentOptionA, zeroArgumentOptionB, zeroArgumentOptionA}, []func(){zeroArgumentOptionA, zeroArgumentOptionB, zeroArgumentOptionB}, 1},
		{"different_length", []func(){zeroArgumentOptionA}, []func(){zeroArgumentOptionA, zeroArgumentOptionB}, 1},
		{"empty", nil, nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := make([]interface{}, len(tt.expected))
			for i, value := range tt.expected {
				values[i] = value
			}
			args := mock.Arguments{mock.FunctionalOptions(values...)}
			diff, differences := args.Diff([]interface{}{tt.actual})
			if differences != tt.differences {
				t.Errorf("Diff returned %d differences, want %d: %s", differences, tt.differences, diff)
			}
			reporter := new(optionsTestingT)
			matched := args.Assert(reporter, tt.actual)
			if matched != (tt.differences == 0) || reporter.failed != (tt.differences != 0) {
				t.Errorf("Assert returned %t, failed=%t; want match=%t", matched, reporter.failed, tt.differences == 0)
			}
		})
	}
}
