package suite

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// This suite is intended to store values to make sure that only
// testing-suite-related methods are run.
type SuiteTester struct {
	Suite
	SetupSuiteRunCount int
	TearDownSuiteRunCount int
	SetupTestRunCount int
	TearDownTestRunCount int
	TestOneRunCount int
	TestTwoRunCount int
	NonTestMethodRunCount int
}

func (suite *SuiteTester) SetupSuite() {
	suite.SetupSuiteRunCount++
}

func (suite *SuiteTester) TearDownSuite() {
	suite.TearDownSuiteRunCount++
}

func (suite *SuiteTester) SetupTest() {
	suite.SetupTestRunCount++
}

func (suite *SuiteTester) TearDownTest() {
	suite.TearDownTestRunCount++
}

func (suite *SuiteTester) TestOne() {
	suite.TestOneRunCount++
}

func (suite *SuiteTester) TestTwo() {
	suite.TestTwoRunCount++
}

func (suite *SuiteTester) NonTestMethod() {
	suite.NonTestMethodRunCount++
}

func TestSuiteLogic(t *testing.T) {
	suiteTester := new(SuiteTester)
	Run(t, suiteTester)

	// The suite was only run once, so the SetupSuite and TearDownSuite
	// methods should have each been run only once.
	assert.Equal(t, suiteTester.SetupSuiteRunCount, 1)
	assert.Equal(t, suiteTester.TearDownSuiteRunCount, 1)

	// There are two test methods (TestOne and TestTwo), so the
	// SetupTest and TearDownTest methods (which should be run once for
	// each test) should have been run twice.
	assert.Equal(t, suiteTester.SetupTestRunCount, 2)
	assert.Equal(t, suiteTester.TearDownTestRunCount, 2)

	// Each test should have been run once.
	assert.Equal(t, suiteTester.TestOneRunCount, 1)
	assert.Equal(t, suiteTester.TestTwoRunCount, 1)

	// Methods that don't match the test method identifier shouldn't
	// have been run at all.
	assert.Equal(t, suiteTester.NonTestMethodRunCount, 0)
}
