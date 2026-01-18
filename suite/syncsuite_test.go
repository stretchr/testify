//go:build go1.25

package suite

import (
	"bytes"
	"flag"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SuiteRequireTwice is intended to test the usage of suite.Require in two
// different tests
type SyncSuiteRequireTwice struct{ SyncSuite }

// TestSuiteRequireTwice checks for regressions of issue #149 where
// suite.requirements was not initialized in suite.SetT()
// A regression would result on these tests panicking rather than failing.
func TestSyncSuiteRequireTwice(t *testing.T) {
	ok := testing.RunTests(
		allTestsFilter,
		[]testing.InternalTest{{
			Name: t.Name() + "/SyncSuiteRequireTwice",
			F: func(t *testing.T) {
				suite := new(SuiteRequireTwice)
				Run(t, suite)
			},
		}},
	)
	assert.False(t, ok)
}

func (s *SyncSuiteRequireTwice) TestRequireOne() {
	r := s.Require()
	r.Equal(1, 2)
}

func (s *SyncSuiteRequireTwice) TestRequireTwo() {
	r := s.Require()
	r.Equal(1, 2)
}

type panickingSyncSuite struct {
	SyncSuite
	panicInSetupTest    bool
	panicInBeforeTest   bool
	panicInTest         bool
	panicInAfterTest    bool
	panicInTearDownTest bool
}

func (s *panickingSyncSuite) SetupTest() {
	if s.panicInSetupTest {
		panic("oops in setup test")
	}
}

func (s *panickingSyncSuite) BeforeTest(_, _ string) {
	if s.panicInBeforeTest {
		panic("oops in before test")
	}
}

func (s *panickingSyncSuite) Test() {
	if s.panicInTest {
		panic("oops in test")
	}
}

func (s *panickingSyncSuite) AfterTest(_, _ string) {
	if s.panicInAfterTest {
		panic("oops in after test")
	}
}

func (s *panickingSyncSuite) TearDownTest() {
	if s.panicInTearDownTest {
		panic("oops in tear down test")
	}
}

func TestSyncSuiteRecoverPanic(t *testing.T) {
	ok := true
	panickingTests := []testing.InternalTest{
		{
			Name: t.Name() + "/InSetupTest",
			F:    func(t *testing.T) { Run(t, &panickingSyncSuite{panicInSetupTest: true}) },
		},
		{
			Name: t.Name() + "InBeforeTest",
			F:    func(t *testing.T) { Run(t, &panickingSyncSuite{panicInBeforeTest: true}) },
		},
		{
			Name: t.Name() + "/InTest",
			F:    func(t *testing.T) { Run(t, &panickingSyncSuite{panicInTest: true}) },
		},
		{
			Name: t.Name() + "/InAfterTest",
			F:    func(t *testing.T) { Run(t, &panickingSyncSuite{panicInAfterTest: true}) },
		},
		{
			Name: t.Name() + "/InTearDownTest",
			F:    func(t *testing.T) { Run(t, &panickingSyncSuite{panicInTearDownTest: true}) },
		},
	}

	require.NotPanics(t, func() {
		ok = testing.RunTests(allTestsFilter, panickingTests)
	})

	assert.False(t, ok)
}

// This suite is intended to store values to make sure that only
// testing-suite-related methods are run.  It's also a fully
// functional example of a testing suite, using setup/teardown methods
// and a helper method that is ignored by testify.  To make this look
// more like a real world example, all tests in the suite perform some
// type of assertion.
type SyncSuiteTester struct {
	// Include our basic suite logic.
	SyncSuite

	// Keep counts of how many times each method is run.
	SetupTestRunCount     int
	TearDownTestRunCount  int
	TestOneRunCount       int
	TestTwoRunCount       int
	NonTestMethodRunCount int

	SuiteNameBefore []string
	TestNameBefore  []string

	SuiteNameAfter []string
	TestNameAfter  []string

	TimeBefore []time.Time
	TimeAfter  []time.Time
}

func (suite *SyncSuiteTester) BeforeTest(suiteName, testName string) {
	suite.SuiteNameBefore = append(suite.SuiteNameBefore, suiteName)
	suite.TestNameBefore = append(suite.TestNameBefore, testName)
	suite.TimeBefore = append(suite.TimeBefore, time.Now())
}

func (suite *SyncSuiteTester) AfterTest(suiteName, testName string) {
	suite.SuiteNameAfter = append(suite.SuiteNameAfter, suiteName)
	suite.TestNameAfter = append(suite.TestNameAfter, testName)
	suite.TimeAfter = append(suite.TimeAfter, time.Now())
}

// The SetupTest method will be run before every test in the suite.
func (suite *SyncSuiteTester) SetupTest() {
	suite.SetupTestRunCount++
}

// The TearDownTest method will be run after every test in the suite.
func (suite *SyncSuiteTester) TearDownTest() {
	suite.TearDownTestRunCount++
}

// Every method in a testing suite that begins with "Test" will be run
// as a test.  TestOne is an example of a test.  For the purposes of
// this example, we've included assertions in the tests, since most
// tests will issue assertions.
func (suite *SyncSuiteTester) TestOne() {
	beforeCount := suite.TestOneRunCount
	suite.TestOneRunCount++
	assert.Equal(suite.T(), suite.TestOneRunCount, beforeCount+1)
	suite.Equal(suite.TestOneRunCount, beforeCount+1)
}

// TestTwo is another example of a test.
func (suite *SyncSuiteTester) TestTwo() {
	beforeCount := suite.TestTwoRunCount
	suite.TestTwoRunCount++
	assert.NotEqual(suite.T(), suite.TestTwoRunCount, beforeCount)
	suite.NotEqual(suite.TestTwoRunCount, beforeCount)
}

func (suite *SyncSuiteTester) TestSkip() {
	suite.T().Skip()
}

// NonTestMethod does not begin with "Test", so it will not be run by
// testify as a test in the suite.  This is useful for creating helper
// methods for your tests.
func (suite *SyncSuiteTester) NonTestMethod() {
	suite.NonTestMethodRunCount++
}

// TestRunSuite will be run by the 'go test' command, so within it, we
// can run our suite using the Run(*testing.T, TestingSuite) function.
func TestRunSyncSuite(t *testing.T) {
	suiteTester := new(SyncSuiteTester)
	Run(t, suiteTester)

	// Normally, the test would end here.  The following are simply
	// some assertions to ensure that the Run function is working as
	// intended - they are not part of the example.

	assert.Len(t, suiteTester.SuiteNameAfter, 3)
	assert.Len(t, suiteTester.SuiteNameBefore, 3)
	assert.Len(t, suiteTester.TestNameAfter, 3)
	assert.Len(t, suiteTester.TestNameBefore, 3)

	assert.Contains(t, suiteTester.TestNameAfter, "TestOne")
	assert.Contains(t, suiteTester.TestNameAfter, "TestTwo")
	assert.Contains(t, suiteTester.TestNameAfter, "TestSkip")

	assert.Contains(t, suiteTester.TestNameBefore, "TestOne")
	assert.Contains(t, suiteTester.TestNameBefore, "TestTwo")
	assert.Contains(t, suiteTester.TestNameBefore, "TestSkip")

	for _, suiteName := range suiteTester.SuiteNameAfter {
		assert.Equal(t, "SyncSuiteTester", suiteName)
	}

	for _, suiteName := range suiteTester.SuiteNameBefore {
		assert.Equal(t, "SyncSuiteTester", suiteName)
	}

	for _, when := range suiteTester.TimeAfter {
		assert.False(t, when.IsZero())
	}

	for _, when := range suiteTester.TimeBefore {
		assert.False(t, when.IsZero())
	}

	// There are three test methods (TestOne, TestTwo and TestSkip), so
	// the SetupTest and TearDownTest methods (which should be run once for
	// each test) should have been run three times.
	assert.Equal(t, 3, suiteTester.SetupTestRunCount)
	assert.Equal(t, 3, suiteTester.TearDownTestRunCount)

	// Each test should have been run once.
	assert.Equal(t, 1, suiteTester.TestOneRunCount)
	assert.Equal(t, 1, suiteTester.TestTwoRunCount)

	// Methods that don't match the test method identifier shouldn't
	// have been run at all.
	assert.Equal(t, 0, suiteTester.NonTestMethodRunCount)
}

// This suite has no Test... methods. It's setup and teardown must be skipped.
type SyncSuiteSetupSkipTester struct {
	SyncSuite

	setUp    bool
	toreDown bool
}

func (s *SyncSuiteSetupSkipTester) SetupSuite() {
	s.setUp = true
}

func (s *SyncSuiteSetupSkipTester) NonTestMethod() {
}

func (s *SyncSuiteSetupSkipTester) TearDownSuite() {
	s.toreDown = true
}

func TestSkippingSyncSuiteSetup(t *testing.T) {
	suiteTester := new(SyncSuiteSetupSkipTester)
	Run(t, suiteTester)
	assert.False(t, suiteTester.setUp)
	assert.False(t, suiteTester.toreDown)
}

func TestSyncSuiteGetters(t *testing.T) {
	suite := new(SyncSuiteTester)
	suite.SetT(t)
	assert.NotNil(t, suite.Assert())
	assert.Equal(t, suite.Assertions, suite.Assert())
	assert.NotNil(t, suite.Require())
	assert.Equal(t, suite.require, suite.Require())
}

type SyncSuiteLoggingTester struct {
	SyncSuite
}

func (s *SyncSuiteLoggingTester) TestLoggingPass() {
	s.T().Log("TESTLOGPASS")
}

func (s *SyncSuiteLoggingTester) TestLoggingFail() {
	s.T().Log("TESTLOGFAIL")
	assert.NotNil(s.T(), nil) // expected to fail
}

func TestSyncSuiteLogging(t *testing.T) {
	suiteLoggingTester := new(SyncSuiteLoggingTester)
	capture := StdoutCapture{}
	internalTest := testing.InternalTest{
		Name: t.Name() + "/SyncSuiteLoggingTester",
		F: func(subT *testing.T) {
			Run(subT, suiteLoggingTester)
		},
	}
	capture.StartCapture()
	testing.RunTests(allTestsFilter, []testing.InternalTest{internalTest})
	output, err := capture.StopCapture()
	require.NoError(t, err, "Got an error trying to capture stdout and stderr!")
	require.NotEmpty(t, output, "output content must not be empty")

	// Failed tests' output is always printed
	assert.Contains(t, output, "TESTLOGFAIL")

	if testing.Verbose() {
		// In verbose mode, output from successful tests is also printed
		assert.Contains(t, output, "TESTLOGPASS")
	} else {
		assert.NotContains(t, output, "TESTLOGPASS")
	}
}

type syncSuiteWithStats struct {
	SyncSuite
	wasCalled bool
	stats     *SuiteInformation
}

func (s *syncSuiteWithStats) HandleStats(suiteName string, stats *SuiteInformation) {
	s.wasCalled = true
	s.stats = stats
}

func (s *syncSuiteWithStats) TestSomething() {
	s.Equal(1, 1)
}

func (s *syncSuiteWithStats) TestPanic() {
	panic("oops")
}

func TestSyncSuiteWithStats(t *testing.T) {
	syncSuiteWithStats := new(syncSuiteWithStats)

	suiteSuccess := testing.RunTests(allTestsFilter, []testing.InternalTest{
		{
			Name: t.Name() + "/syncSuiteWithStats",
			F: func(t *testing.T) {
				Run(t, syncSuiteWithStats)
			},
		},
	})
	require.False(t, suiteSuccess, "syncSuiteWithStats should report test failure because of panic in TestPanic")

	assert.True(t, syncSuiteWithStats.wasCalled)
	assert.NotZero(t, syncSuiteWithStats.stats.Start)
	assert.NotZero(t, syncSuiteWithStats.stats.End)
	assert.False(t, syncSuiteWithStats.stats.Passed())

	testStats := syncSuiteWithStats.stats.TestStats

	assert.NotZero(t, testStats["TestSomething"].Start)
	assert.NotZero(t, testStats["TestSomething"].End)
	assert.True(t, testStats["TestSomething"].Passed)

	assert.NotZero(t, testStats["TestPanic"].Start)
	assert.NotZero(t, testStats["TestPanic"].End)
	assert.False(t, testStats["TestPanic"].Passed)
}

// FailfastSyncSuite will test the behavior when running with the failfast flag
// It logs calls in the callOrder slice which we then use to assert the correct calls were made
type FailfastSyncSuite struct {
	SyncSuite
	callOrder []string
}

func (s *FailfastSyncSuite) call(method string) {
	s.callOrder = append(s.callOrder, method)
}

func TestFailfastSyncSuite(t *testing.T) {
	// This test suite is run twice. Once normally and once with the -failfast flag by TestFailfastSuiteFailFastOn
	// If you need to debug it run this test directly with the failfast flag set on/off as you need
	failFast := flag.Lookup("test.failfast").Value.(flag.Getter).Get().(bool)
	s := new(FailfastSyncSuite)
	ok := testing.RunTests(
		allTestsFilter,
		[]testing.InternalTest{{
			Name: t.Name() + "/FailfastSyncSuite",
			F: func(t *testing.T) {
				Run(t, s)
			},
		}},
	)
	assert.False(t, ok)
	var expect []string
	if failFast {
		expect = []string{"SetupTest", "Test A Fails", "TearDownTest"}
	} else {
		expect = []string{"SetupTest", "Test A Fails", "TearDownTest", "SetupTest", "Test B Passes", "TearDownTest"}
	}
	callOrderAssert(t, expect, s.callOrder)
}

func TestFailfastSyncSuiteFailFastOn(t *testing.T) {
	// To test this with failfast on (and isolated from other intended test failures in our test suite) we launch it in its own process
	cmd := exec.Command("go", "test", "-v", "-race", "-run", "TestFailfastSyncSuite", "-failfast")
	var out bytes.Buffer
	cmd.Stdout = &out
	t.Log("Running go test -v -race -run TestFailfastSyncSuite -failfast")
	err := cmd.Run()
	t.Log(out.String())
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func (s *FailfastSyncSuite) SetupTest() {
	s.call("SetupTest")
}

func (s *FailfastSyncSuite) TearDownTest() {
	s.call("TearDownTest")
}

func (s *FailfastSyncSuite) Test_A_Fails() {
	s.call("Test A Fails")
	s.T().Error("Test A meant to fail")
}

func (s *FailfastSyncSuite) Test_B_Passes() {
	s.call("Test B Passes")
	s.Require().True(true)
}

type unInitializedSyncSuite struct {
	SyncSuite
}

// TestUnInitializedSuites asserts the behavior of the suite methods when the
// suite is not initialized
func TestUnInitializedSyncSuites(t *testing.T) {
	t.Run("should panic on Require", func(t *testing.T) {
		suite := new(unInitializedSyncSuite)

		assert.Panics(t, func() {
			suite.Require().True(true)
		})
	})

	t.Run("should panic on Assert", func(t *testing.T) {
		suite := new(unInitializedSyncSuite)

		assert.Panics(t, func() {
			suite.Assert().True(true)
		})
	})
}

// SyncSuiteSignatureValidationTester tests valid and invalid method signatures.
type SyncSuiteSignatureValidationTester struct {
	SyncSuite

	executedTestCount int
}

// Valid test method — should run.
func (s *SyncSuiteSignatureValidationTester) TestValidSignature() {
	s.executedTestCount++
}

// Invalid: has return value.
func (s *SyncSuiteSignatureValidationTester) TestInvalidSignatureReturnValue() interface{} {
	s.executedTestCount++
	return nil
}

// Invalid: has input arg.
func (s *SyncSuiteSignatureValidationTester) TestInvalidSignatureArg(somearg string) {
	s.executedTestCount++
}

// Invalid: both input arg and return value.
func (s *SyncSuiteSignatureValidationTester) TestInvalidSignatureBoth(somearg string) interface{} {
	s.executedTestCount++
	return nil
}

// TestSuiteSignatureValidation ensures that invalid signature methods fail and valid method runs.
func TestSyncSuiteSignatureValidation(t *testing.T) {
	suiteTester := new(SyncSuiteSignatureValidationTester)

	ok := testing.RunTests(allTestsFilter, []testing.InternalTest{
		{
			Name: "signature validation",
			F: func(t *testing.T) {
				Run(t, suiteTester)
			},
		},
	})

	require.False(t, ok, "Suite should fail due to invalid method signatures")

	assert.Equal(t, 1, suiteTester.executedTestCount, "Only the valid test method should have been executed")
}

type syncSuiteTimeTest struct {
	SyncSuite
	ticker *time.Ticker
}

func (s *syncSuiteTimeTest) SetupTest() {
	// there will no test that have a timeout set to 24h
	s.ticker = time.NewTicker(time.Hour * 24)
}

func (s *syncSuiteTimeTest) TestTimeIsAdvanced() {
	// check if ticker is advanced and the test does not timeout
	<-s.ticker.C
}

// Check if suite test is called in a synctest bubble and aswell if the
// ticker is created inside the bubble (SetupTest called inside synctest.Test).
func TestSyncSuiteTimeTest(t *testing.T) {
	t.Setenv("GODEBUG", "asynctimerchan=0") // since our go.mod says `go 1.17`
	RunSync(t, new(syncSuiteTimeTest))
}

type syncSuiteSetTime struct {
	SyncSuite
}

func (s *syncSuiteSetTime) TestSetTime() {
	// time is not set because timestamp is before 1. Jan 2000
	s.SetTime(time.Date(1970, time.January, 1, 1, 0, 1, 0, time.UTC))
	s.Equal(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now().UTC())

	// time is advanced to given timestamp
	ts := time.Date(2001, time.January, 1, 1, 0, 0, 0, time.UTC)
	s.SetTime(ts)
	s.Equal(ts, time.Now().UTC())
}

func TestSyncSuiteSetTime(t *testing.T) {
	t.Setenv("GODEBUG", "asynctimerchan=0") // since our go.mod says `go 1.17`
	RunSync(t, new(syncSuiteSetTime))
}
