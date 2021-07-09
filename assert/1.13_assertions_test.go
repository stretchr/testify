// +build go1.13

package assert

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorIs(t *testing.T) {
	mockT := new(testing.T)

	var (
		// Similarly to EqualError, start with nil errors
		errA error
		errB error
	)
	True(t, ErrorIs(mockT, errA, errB),
		"ErrorIs should return false for expecting and receiving `nil`")

	// now set an error
	errA = errors.New("some error")
	True(t, ErrorIs(mockT, errA, errA),
		"ErrorIs should return true for same error instance")

	False(t, ErrorIs(mockT, errA, errB),
		"ErrorIs should return false for expecting error and receiving `nil`")
	False(t, ErrorIs(mockT, errB, errA),
		"ErrorIs should return false for expecting `nil` and receiving error")

	errB = errors.New("some other error")
	False(t, ErrorIs(mockT, errA, errB),
		"ErrorIs should return false for different and unrelated errors")

	// wrapping errA keeps errors.Is(err, errA) == true
	errB = fmt.Errorf("wrapping: %w", errA)
	True(t, ErrorIs(mockT, errA, errB),
		"ErrorIs should return true for wrapped error")
}
