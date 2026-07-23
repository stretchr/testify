package assert

import (
	"testing"
)

type money struct {
	cents int
	label string // display only; equality uses cents
}

func (m money) Equal(other interface{}) bool {
	o, ok := other.(money)
	if !ok {
		return false
	}
	return m.cents == o.cents
}

func TestComparerTopLevel(t *testing.T) {
	mockT := new(testing.T)
	a := money{cents: 100, label: "$1.00"}
	b := money{cents: 100, label: "1 dollar"}
	c := money{cents: 200, label: "$2.00"}

	if !Equal(mockT, a, b) {
		t.Fatal("Comparer should treat same cents as equal despite label")
	}
	if Equal(mockT, a, c) {
		t.Fatal("different cents must not be equal")
	}
}

func TestComparerNotEqual(t *testing.T) {
	mockT := new(testing.T)
	a := money{cents: 100, label: "$1"}
	b := money{cents: 100, label: "1 USD"}
	if !NotEqual(mockT, a, money{cents: 101, label: a.label}) {
		t.Fatal("NotEqual should use Comparer")
	}
	// same cents → Equal, so NotEqual fails (returns false)
	if NotEqual(mockT, a, b) {
		t.Fatal("NotEqual should be false when Comparer says equal")
	}
}
