package suite_test

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type My struct {
	suite.Suite
}

func (m *My) TestSequntial() {
	m.Assert().True(true)
	// passes during suite
}

func (m *My) ParallelTestPass(t *testing.T) {
	runtime.Gosched()
	m.Assert().True(true)
}

func (m *My) ParallelTestFail(t *testing.T) {
	runtime.Gosched()
	assert.True(t, false)
}

func TestKataSuite(t *testing.T) {
	suite.Run(t, &My{})
}
