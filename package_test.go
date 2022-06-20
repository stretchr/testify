package testify

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallerInfoWithSubtests(t *testing.T) {
	assert.Equal(t, "package_test.go:11", strings.Join(assert.CallerInfo(t), " "))

	t.Run("Subtest", func(t *testing.T) {
		assert.Equal(t, "package_test.go:14 package_test.go:13", strings.Join(assert.CallerInfo(t), " "))
	})
}

func TestImports(t *testing.T) {
	if assert.Equal(t, 1, 1) != true {
		t.Error("Something is wrong.")
	}
}
