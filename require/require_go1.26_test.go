//go:build go1.26

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

	// success: returns the matched value, does not call FailNow
	target := ErrorAsType[*requireCustomError](t, fmt.Errorf("wrap: %w", &requireCustomError{}))
	if target == nil {
		t.Error("expected non-nil target on success")
	}

	// failure: calls FailNow
	mockT := new(MockT)
	ErrorAsType[*requireCustomError](mockT, io.EOF)
	if !mockT.Failed {
		t.Error("expected FailNow to be called")
	}

	// failure on nil: calls FailNow
	mockT = new(MockT)
	ErrorAsType[*requireCustomError](mockT, nil)
	if !mockT.Failed {
		t.Error("expected FailNow to be called on nil error")
	}
}

func TestNotErrorAsType(t *testing.T) {
	t.Parallel()

	// success: does not call FailNow
	NotErrorAsType[*requireCustomError](t, io.EOF)
	NotErrorAsType[*requireCustomError](t, nil)

	// failure: calls FailNow
	mockT := new(MockT)
	NotErrorAsType[*requireCustomError](mockT, fmt.Errorf("wrap: %w", &requireCustomError{}))
	if !mockT.Failed {
		t.Error("expected FailNow to be called")
	}
}
