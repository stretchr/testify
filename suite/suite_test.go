package suite

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// This suite is intended to store values to make sure that only
// testing-suite-related methods are run.
type SuiteTester struct {
	Suite
	BeforeSuiteRunCount int
	AfterSuiteRunCount int
	BeforeTestRunCount int
	AfterTestRunCount int
	TestOneRunCount int
	TestTwoRunCount int
	NonTestMethodRunCount int
}

func (suite *SuiteTester) BeforeSuite() {
	suite.BeforeSuiteRunCount++
}

func (suite *SuiteTester) AfterSuite() {
	suite.AfterSuiteRunCount++
}

func (suite *SuiteTester) BeforeTest() {
	suite.BeforeTestRunCount++
}

func (suite *SuiteTester) AfterTest() {
	suite.AfterTestRunCount++
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

	// The suite was only run once, so the BeforeSuite and AfterSuite
	// methods should have each been run only once.
	assert.Equal(t, suiteTester.BeforeSuiteRunCount, 1)
	assert.Equal(t, suiteTester.AfterSuiteRunCount, 1)

	// There are two test methods (TestOne and TestTwo), so the
	// BeforeTest and AfterTest methods (which should be run once for
	// each test) should have been run twice.
	assert.Equal(t, suiteTester.BeforeTestRunCount, 2)
	assert.Equal(t, suiteTester.AfterTestRunCount, 2)

	// Each test should have been run once.
	assert.Equal(t, suiteTester.TestOneRunCount, 1)
	assert.Equal(t, suiteTester.TestTwoRunCount, 1)

	// Methods that don't match the test method identifier shouldn't
	// have been run at all.
	assert.Equal(t, suiteTester.NonTestMethodRunCount, 0)
}
