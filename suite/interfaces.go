package suite

import (
	"testing"
	"time"
)

// *testing.T interface for < go1.15
type testingT interface {
	testing.TB

	Parallel()
	Run(name string, f func(t *testing.T)) bool
}

// *testing.T interface for go1.15
type testingT115 interface {
	TempDir() string
	Deadline() (deadline time.Time, ok bool)
}

// *testing.T interface for go1.17
type testingT117 interface {
	Setenv(key, value string)
}

// SetupAllSuite has a SetupSuite method, which will run before the
// tests in the suite are run.
type SetupAllSuite interface {
	SetupSuite(t *T)
}

// SetupTestSuite has a SetupTest method, which will run before each
// test in the suite.
type SetupTestSuite interface {
	SetupTest(t *T)
}

// TearDownAllSuite has a TearDownSuite method, which will run after
// all the tests in the suite have been run.
type TearDownAllSuite interface {
	TearDownSuite(t *T)
}

// TearDownTestSuite has a TearDownTest method, which will run after
// each test in the suite.
type TearDownTestSuite interface {
	TearDownTest(t *T)
}

// BeforeTest has a function to be executed right before the test
// starts and receives the suite and test names as input
type BeforeTest interface {
	BeforeTest(t *T, suiteName, testName string)
}

// AfterTest has a function to be executed right after the test
// finishes and receives the suite and test names as input
type AfterTest interface {
	AfterTest(t *T, suiteName, testName string)
}

// WithStats implements HandleStats, a function that will be executed
// when a test suite is finished. The stats contain information about
// the execution of that suite and its tests.
type WithStats interface {
	HandleStats(t *T, suiteName string, stats *SuiteInformation)
}
