//go:build go1.27

package require

import (
	"fmt"
	"io"
	"testing"
)

type requireCustomError struct{}

func (*requireCustomError) Error() string { return "fail" }

func TestErrorAsType(t *testing.T) {
	t.Parallel()

	err := fmt.Errorf("wrap: %w", &requireCustomError{})
	if target := ErrorAsType[*requireCustomError](t, err); target == nil {
		t.Error("ErrorAsType did not return the matched error")
	}

	mockT := new(MockT)
	ErrorAsType[*requireCustomError](mockT, io.EOF)
	if !mockT.Failed {
		t.Error("expected FailNow to be called")
	}

	mockT = new(MockT)
	ErrorAsType[*requireCustomError](mockT, nil)
	if !mockT.Failed {
		t.Error("expected FailNow to be called for a nil error")
	}
}

func TestNotErrorAsType(t *testing.T) {
	t.Parallel()

	NotErrorAsType[*requireCustomError](t, io.EOF)
	NotErrorAsType[*requireCustomError](t, nil)

	mockT := new(MockT)
	NotErrorAsType[*requireCustomError](mockT, fmt.Errorf("wrap: %w", &requireCustomError{}))
	if !mockT.Failed {
		t.Error("expected FailNow to be called")
	}
}

func TestErrorAsTypeGeneratedAPIs(t *testing.T) {
	t.Parallel()

	err := fmt.Errorf("wrap: %w", &requireCustomError{})
	requirements := New(t)

	if target := ErrorAsTypef[*requireCustomError](t, err, "message"); target == nil {
		t.Error("ErrorAsTypef did not return the matched error")
	}
	if target := requirements.ErrorAsType[*requireCustomError](err); target == nil {
		t.Error("Assertions.ErrorAsType did not return the matched error")
	}
	if target := requirements.ErrorAsTypef[*requireCustomError](err, "message"); target == nil {
		t.Error("Assertions.ErrorAsTypef did not return the matched error")
	}
	NotErrorAsTypef[*requireCustomError](t, io.EOF, "message")
	requirements.NotErrorAsType[*requireCustomError](io.EOF)
	requirements.NotErrorAsTypef[*requireCustomError](io.EOF, "message")

	mockT := new(MockT)
	New(mockT).ErrorAsTypef[*requireCustomError](io.EOF, "message")
	if !mockT.Failed {
		t.Error("Assertions.ErrorAsTypef did not call FailNow")
	}

	mockT = new(MockT)
	New(mockT).NotErrorAsTypef[*requireCustomError](err, "message")
	if !mockT.Failed {
		t.Error("Assertions.NotErrorAsTypef did not call FailNow")
	}
}
