package suite

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Use globals to test setup and teardown logic because RunParallel makes a
// copy of the suite struct. These globals allow us to count setup/teardown
// calls between suite copies.
var (
	suiteInt         int64
	teardownSuiteInt int64
	testInt          int64
	teardownTestInt  int64
)

type nested struct {
	n int
}

type complexData struct {
	nested nested
}

type ParallelTestSuite struct {
	suiteInt         int64
	setupSimpleData  int
	setupComplexData complexData
}

func (s *ParallelTestSuite) SetupSuite() {
	atomic.AddInt64(&suiteInt, 1)
	s.setupSimpleData = 1
	s.setupComplexData = complexData{nested: nested{7}}
}

func (s *ParallelTestSuite) TearDownSuite() {
	atomic.AddInt64(&teardownSuiteInt, 1)
}

func (s *ParallelTestSuite) SetupTest() {
	atomic.AddInt64(&testInt, 1)
}

func (s *ParallelTestSuite) TearDownTest() {
	atomic.AddInt64(&teardownTestInt, 1)
}

func (s *ParallelTestSuite) TestOne(t *testing.T) {
	//Access shared data from TestOne and TestTwo which should trigger race conditions.
	s.suiteInt++
}

func (s *ParallelTestSuite) TestTwo(t *testing.T) {
	s.suiteInt++
}

func (s *ParallelTestSuite) TestSkip(t *testing.T) {
	t.Skip()
	t.FailNow()
}

func TestParallelSuiteRunner(t *testing.T) {
	// Run with the following to detect data races:
	// go test -v -race -run TestParallelSuiteRunner
	suite := new(ParallelTestSuite)
	RunParallel(t, suite)
	assert.Equal(t, int64(3), testInt)
	assert.Equal(t, int64(3), teardownTestInt)
	assert.Equal(t, int64(1), suiteInt)
	assert.Equal(t, int64(1), teardownSuiteInt)
	assert.Equal(t, int64(0), suite.suiteInt)
	assert.Equal(t, 1, suite.setupSimpleData)
	assert.Equal(t, 7, suite.setupComplexData.nested.n)
}

type EmbeddedParallelTestSuite struct {
	Suite
}

func (s *EmbeddedParallelTestSuite) TestOne() {
}

func TestEmbeddedParallelSuiteRunner(t *testing.T) {
	suite := new(EmbeddedParallelTestSuite)
	RunParallel(t, suite)
}
