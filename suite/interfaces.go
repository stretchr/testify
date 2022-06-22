package suite

import "testing"

// TestingSuite can store and return the current *testing.T context
// generated by 'go test'.
type TestingSuite interface {
	T() *testing.T
	F() *testing.F
	SetT(*testing.T)
	SetF(*testing.F)
}

// SetupAllSuite has a SetupSuite method, which will run before the
// tests in the suite are run.
type SetupAllSuite interface {
	SetupSuite()
}

// SetupTestSuite has a SetupTest method, which will run before each
// test in the suite.
type SetupTestSuite interface {
	SetupTest()
}

// TearDownAllSuite has a TearDownSuite method, which will run after
// all the tests in the suite have been run.
type TearDownAllSuite interface {
	TearDownSuite()
}

// TearDownTestSuite has a TearDownTest method, which will run after
// each test in the suite.
type TearDownTestSuite interface {
	TearDownTest()
}

// BeforeTest has a function to be executed right before the test
// starts and receives the suite and test names as input
type BeforeTest interface {
	BeforeTest(suiteName, testName string)
}

// AfterTest has a function to be executed right after the test
// finishes and receives the suite and test names as input
type AfterTest interface {
	AfterTest(suiteName, testName string)
}

// WithStats implements HandleStats, a function that will be executed
// when a test suite is finished. The stats contain information about
// the execution of that suite and its tests.
type WithStats interface {
	HandleStats(suiteName string, stats *SuiteInformation)
}
