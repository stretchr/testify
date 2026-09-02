//go:build go1.21

package require

import (
	"errors"
	"testing"
)

func TestMust(t *testing.T) {
	mockT := new(MockT)

	if v := Must(mockT, func() (int, error) { return 42, nil }); v != 42 {
		t.Errorf("Must returned %d, expected 42", v)
	}
	if mockT.Failed {
		t.Error("Must should not have failed the test")
	}

	// MockT.FailNow does not stop execution, so Must returns normally here.
	mockT = new(MockT)
	Must(mockT, func() (int, error) { return 0, errors.New("boom") })
	if !mockT.Failed {
		t.Error("Must should have failed the test")
	}
}

func TestMustf(t *testing.T) {
	mockT := new(MockT)

	if v := Mustf(mockT, func() (string, error) { return "ok", nil }, "loading %s", "config"); v != "ok" {
		t.Errorf("Mustf returned %q, expected \"ok\"", v)
	}
	if mockT.Failed {
		t.Error("Mustf should not have failed the test")
	}

	mockT = new(MockT)
	Mustf(mockT, func() (string, error) { return "", errors.New("boom") }, "loading %s", "config")
	if !mockT.Failed {
		t.Error("Mustf should have failed the test")
	}
}
