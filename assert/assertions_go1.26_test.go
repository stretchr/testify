//go:build go1.26

package assert

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestErrorAsType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err          error
		result       bool
		resultErrMsg string
	}{
		{
			err:    fmt.Errorf("wrap: %w", &customError{}),
			result: true,
		},
		{
			err:    io.EOF,
			result: false,
			resultErrMsg: "" +
				"Should be in error chain:\n" +
				"expected: *assert.customError\n" +
				"in chain: \"EOF\" (*errors.errorString)\n",
		},
		{
			err:    nil,
			result: false,
			resultErrMsg: "" +
				"An error is expected but got nil.\n" +
				"expected: *assert.customError\n",
		},
		{
			err:    fmt.Errorf("abc: %w", errors.New("def")),
			result: false,
			resultErrMsg: "" +
				"Should be in error chain:\n" +
				"expected: *assert.customError\n" +
				"in chain: \"abc: def\" (*fmt.wrapError)\n" +
				"\t\"def\" (*errors.errorString)\n",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("ErrorAsType[*customError](%#v)", tt.err), func(t *testing.T) {
			mockT := new(captureTestingT)
			target, ok := ErrorAsType[*customError](mockT, tt.err)
			if tt.result {
				if !ok {
					t.Error("expected ok=true but got false")
				}
				if target == nil {
					t.Error("expected non-nil target on success")
				}
			} else {
				mockT.checkResultAndErrMsg(t, false, ok, tt.resultErrMsg)
			}
		})
	}
}

func TestNotErrorAsType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err          error
		result       bool
		resultErrMsg string
	}{
		{
			err:    fmt.Errorf("wrap: %w", &customError{}),
			result: false,
			resultErrMsg: "" +
				"Target error should not be in err chain:\n" +
				"found: *assert.customError\n" +
				"in chain: \"wrap: fail\" (*fmt.wrapError)\n" +
				"\t\"fail\" (*assert.customError)\n",
		},
		{
			err:    io.EOF,
			result: true,
		},
		{
			err:    nil,
			result: true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("NotErrorAsType[*customError](%#v)", tt.err), func(t *testing.T) {
			mockT := new(captureTestingT)
			res := NotErrorAsType[*customError](mockT, tt.err)
			mockT.checkResultAndErrMsg(t, tt.result, res, tt.resultErrMsg)
		})
	}
}
