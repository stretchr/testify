package assert

import (
	"testing"
)

type issue1955Shape interface{ Area() int }

type issue1955Square struct{ Side int }

func (s *issue1955Square) Area() int { return s.Side * s.Side }

func TestIssue1955EmptyPointerToInterface(t *testing.T) {
	var s issue1955Shape = &issue1955Square{}
	var nilShape issue1955Shape
	var nonEmpty issue1955Shape = &issue1955Square{Side: 2}

	tests := []struct {
		name  string
		value interface{}
		empty bool
	}{
		{"unwrapped value", issue1955Square{}, true},
		{"pointer to value", &issue1955Square{}, true},
		{"interface", s, true},
		{"pointer to interface holding empty value", &s, true},
		{"pointer to empty interface", &nilShape, true},
		{"pointer to interface holding non-empty value", &nonEmpty, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := new(testing.T)
			got := Empty(mockT, tt.value)
			if got != tt.empty {
				t.Errorf("Empty(%T) = %v, want %v", tt.value, got, tt.empty)
			}
			notGot := NotEmpty(mockT, tt.value)
			if notGot == tt.empty {
				t.Errorf("NotEmpty(%T) = %v, want %v", tt.value, notGot, !tt.empty)
			}
		})
	}
}
