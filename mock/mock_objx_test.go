//go:build !testify_no_objx && !testify_no_deps

package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Mock_TestData(t *testing.T) {
	t.Parallel()

	var mockedService = new(TestExampleImplementation)

	if assert.NotNil(t, mockedService.TestData()) {

		mockedService.TestData().Set("something", 123)
		assert.Equal(t, 123, mockedService.TestData().Get("something").Data())
	}
}
