package assert

import (
	"errors"
)

// ErrExample is an error instance useful for testing.  If the code does not care
// about error specifics, and only needs to return the error for example, this
// error should be used to make the test code more readable.
var ErrExample = errors.New("assert.ErrExample general error for testing")
