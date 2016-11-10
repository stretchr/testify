package suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ParallelTestSuite struct {
	Suite
	testInt          int
	teardownTestInt  int
	suiteInt         int
	teardownSuiteInt int
}

func (s *ParallelTestSuite) SetupSuite() {
	s.suiteInt++
}

func (s *ParallelTestSuite) TearDownSuite() {
	s.teardownSuiteInt++
}

func (s *ParallelTestSuite) SetupTest() {
	s.testInt++
}

func (s *ParallelTestSuite) TearDownTest() {
	s.teardownTestInt++
}

func (s ParallelTestSuite) TestOne(t *testing.T) {
	t.Logf("Accessing shared suite data from TestOne:%d", s.suiteInt)
	t.Logf("Accessing shared test data from TestOne:%d", s.testInt)
	s.suiteInt++
}

func (s ParallelTestSuite) TestTwo(t *testing.T) {
	t.Logf("Accessing shared suite data from TestTwo:%d", s.suiteInt)
	t.Logf("Accessing shared test data from TestTwo:%d", s.testInt)
	s.suiteInt++
}

func (s ParallelTestSuite) TestSkip(t *testing.T) {
	t.Skip()
	t.FailNow()
}

func TestParallelSuiteRunner(t *testing.T) {
	// Run with the following to detect data races:
	// go test -v -race -run TestParallelSuiteRunner
	suite := new(ParallelTestSuite)
	RunParallel(t, suite)
	assert.Equal(t, 3, suite.testInt)
	assert.Equal(t, 3, suite.teardownTestInt)
	assert.Equal(t, 1, suite.suiteInt)
	assert.Equal(t, 1, suite.teardownSuiteInt)
}
