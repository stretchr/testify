//go:build go1.27

package require

import (
	"errors"
	"testing"
)

func TestMustWrapper(t *testing.T) {
	mockT := new(MockT)
	require := New(mockT)

	if v := require.Must(func() (int, error) { return 42, nil }); v != 42 {
		t.Errorf("Must returned %d, expected 42", v)
	}
	if mockT.Failed {
		t.Error("Must should not have failed the test")
	}

	// MockT.FailNow does not stop execution, so Must returns normally here.
	mockT = new(MockT)
	require = New(mockT)
	require.Must(func() (int, error) { return 0, errors.New("boom") })
	if !mockT.Failed {
		t.Error("Must should have failed the test")
	}
}

func TestMustfWrapper(t *testing.T) {
	mockT := new(MockT)
	require := New(mockT)

	if v := require.Mustf(func() (string, error) { return "ok", nil }, "loading %s", "config"); v != "ok" {
		t.Errorf("Mustf returned %q, expected \"ok\"", v)
	}
	if mockT.Failed {
		t.Error("Mustf should not have failed the test")
	}

	mockT = new(MockT)
	require = New(mockT)
	require.Mustf(func() (string, error) { return "", errors.New("boom") }, "loading %s", "config")
	if !mockT.Failed {
		t.Error("Mustf should have failed the test")
	}
}
