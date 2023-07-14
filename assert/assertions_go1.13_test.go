// +build go1.13

package assert

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

type unwrappableError struct {
	base error
}

func (e unwrappableError) Error() string { return fmt.Sprintf("wrapper: %v", e.base) }

func (e unwrappableError) Unwrap() error { return e.base }

func TestErrorIs(t *testing.T) {
	mockT := new(testing.T)

	var err error
	False(t, ErrorIs(mockT, err, os.ErrExist), "nil is not any error wrapper")

	False(t, ErrorIs(mockT, errors.New("not any error wrapper"), os.ErrExist), "New error is not any error wrapper")

	False(t, ErrorIs(mockT, fmt.Errorf("wrapper: %v", os.ErrExist), os.ErrExist), "getting err.Error is not wrapping")

	// go 1.13 new verb
	True(t, ErrorIs(mockT, fmt.Errorf("wrapper: %w", os.ErrExist), os.ErrExist), "annotating with %%w must be wrapper")

	False(t, ErrorIs(mockT, unwrappableError{base: errors.New("something")}, os.ErrNotExist), "not matching unwrap")

	False(t, ErrorIs(mockT, unwrappableError{base: os.ErrClosed}, os.ErrNotExist), "not matching unwrap")

	False(t, ErrorIs(mockT, unwrappableError{base: errors.New(os.ErrNotExist.Error())}, os.ErrNotExist), "not matching unwrap")

	True(t, ErrorIs(mockT, unwrappableError{base: os.ErrNotExist}, os.ErrNotExist), "matching unwrap")
}
