package suite

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SuiteRequireTwice is intended to test the usage of suite.Require in two
// different tests
type SuiteRequireTwice struct{ Suite }

// TestSuiteRequireTwice checks for regressions of issue #149 where
// suite.requirements was not initialised in suite.SetT()
// A regression would result on these tests panicking rather than failing.
func TestSuiteRequireTwice(t *testing.T) {
	ok := testing.RunTests(
		allTestsFilter,
		[]testing.InternalTest{{
			Name: "TestSuiteRequireTwice",
			F: func(t *testing.T) {
				suite := new(SuiteRequireTwice)
				Run(t, suite)
			},
		}},
	)
	assert.Equal(t, false, ok)
}

func (s *SuiteRequireTwice) TestRequireOne() {
	r := s.Require()
	r.Equal(1, 2)
}

func (s *SuiteRequireTwice) TestRequireTwo() {
	r := s.Require()
	r.Equal(1, 2)
}

type panickingSuite struct {
	Suite
	panicInSetupSuite    bool
	panicInSetupTest     bool
	panicInBeforeTest    bool
	panicInTest          bool
	panicInAfterTest     bool
	panicInTearDownTest  bool
	panicInTearDownSuite bool
}

func (s *panickingSuite) SetupSuite() {
	if s.panicInSetupSuite {
		panic("oops in setup suite")
	}
}

func (s *panickingSuite) SetupTest() {
	if s.panicInSetupTest {
		panic("oops in setup test")
	}
}

func (s *panickingSuite) BeforeTest(_, _ string) {
	if s.panicInBeforeTest {
		panic("oops in before test")
	}
}

func (s *panickingSuite) Test() {
	if s.panicInTest {
		panic("oops in test")
	}
}

func (s *panickingSuite) AfterTest(_, _ string) {
	if s.panicInAfterTest {
		panic("oops in after test")
	}
}

func (s *panickingSuite) TearDownTest() {
	if s.panicInTearDownTest {
		panic("oops in tear down test")
	}
}

func (s *panickingSuite) TearDownSuite() {
	if s.panicInTearDownSuite {
		panic("oops in tear down suite")
	}
}

func TestSuiteRecoverPanic(t *testing.T) {
	ok := true
	panickingTests := []testing.InternalTest{
		{
			Name: "TestPanicInSetupSuite",
			F:    func(t *testing.T) { Run(t, &panickingSuite{panicInSetupSuite: true}) },
		},
		{
			Name: "TestPanicInSetupTest",
			F:    func(t *testing.T) { Run(t, &panickingSuite{panicInSetupTest: true}) },
		},
		{
			Name: "TestPanicInBeforeTest",
			F:    func(t *testing.T) { Run(t, &panickingSuite{panicInBeforeTest: true}) },
		},
		{
			Name: "TestPanicInTest",
			F:    func(t *testing.T) { Run(t, &panickingSuite{panicInTest: true}) },
		},
		{
			Name: "TestPanicInAfterTest",
			F:    func(t *testing.T) { Run(t, &panickingSuite{panicInAfterTest: true}) },
		},
		{
			Name: "TestPanicInTearDownTest",
			F:    func(t *testing.T) { Run(t, &panickingSuite{panicInTearDownTest: true}) },
		},
		{
			Name: "TestPanicInTearDownSuite",
			F:    func(t *testing.T) { Run(t, &panickingSuite{panicInTearDownSuite: true}) },
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
type SuiteTester struct {
	// Include our basic suite logic.
	Suite

	// Keep counts of how many times each method is run.
	SetupSuiteRunCount    int
	TearDownSuiteRunCount int
	SetupTestRunCount     int
	TearDownTestRunCount  int
	TestOneRunCount       int
	TestTwoRunCount       int
	TestSubtestRunCount   int
	NonTestMethodRunCount int

	SuiteNameBefore []string
	TestNameBefore  []string

	SuiteNameAfter []string
	TestNameAfter  []string

	TimeBefore []time.Time
	TimeAfter  []time.Time
}

type SuiteSkipTester struct {
	// Include our basic suite logic.
	Suite

	// Keep counts of how many times each method is run.
	SetupSuiteRunCount    int
	TearDownSuiteRunCount int
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *SuiteTester) SetupSuite() {
	suite.SetupSuiteRunCount++
}

func (suite *SuiteTester) BeforeTest(suiteName, testName string) {
	suite.SuiteNameBefore = append(suite.SuiteNameBefore, suiteName)
	suite.TestNameBefore = append(suite.TestNameBefore, testName)
	suite.TimeBefore = append(suite.TimeBefore, time.Now())
}

func (suite *SuiteTester) AfterTest(suiteName, testName string) {
	suite.SuiteNameAfter = append(suite.SuiteNameAfter, suiteName)
	suite.TestNameAfter = append(suite.TestNameAfter, testName)
	suite.TimeAfter = append(suite.TimeAfter, time.Now())
}

func (suite *SuiteSkipTester) SetupSuite() {
	suite.SetupSuiteRunCount++
	suite.T().Skip()
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *SuiteTester) TearDownSuite() {
	suite.TearDownSuiteRunCount++
}

func (suite *SuiteSkipTester) TearDownSuite() {
	suite.TearDownSuiteRunCount++
}

// The SetupTest method will be run before every test in the suite.
func (suite *SuiteTester) SetupTest() {
	suite.SetupTestRunCount++
}

// The TearDownTest method will be run after every test in the suite.
func (suite *SuiteTester) TearDownTest() {
	suite.TearDownTestRunCount++
}

// Every method in a testing suite that begins with "Test" will be run
// as a test.  TestOne is an example of a test.  For the purposes of
// this example, we've included assertions in the tests, since most
// tests will issue assertions.
func (suite *SuiteTester) TestOne() {
	beforeCount := suite.TestOneRunCount
	suite.TestOneRunCount++
	assert.Equal(suite.T(), suite.TestOneRunCount, beforeCount+1)
	suite.Equal(suite.TestOneRunCount, beforeCount+1)
}

// TestTwo is another example of a test.
func (suite *SuiteTester) TestTwo() {
	beforeCount := suite.TestTwoRunCount
	suite.TestTwoRunCount++
	assert.NotEqual(suite.T(), suite.TestTwoRunCount, beforeCount)
	suite.NotEqual(suite.TestTwoRunCount, beforeCount)
}

func (suite *SuiteTester) TestSkip() {
	suite.T().Skip()
}

// NonTestMethod does not begin with "Test", so it will not be run by
// testify as a test in the suite.  This is useful for creating helper
// methods for your tests.
func (suite *SuiteTester) NonTestMethod() {
	suite.NonTestMethodRunCount++
}

func (suite *SuiteTester) TestSubtest() {
	suite.TestSubtestRunCount++

	for _, t := range []struct {
		testName string
	}{
		{"first"},
		{"second"},
	} {
		suiteT := suite.T()
		suite.Run(t.testName, func() {
			// We should get a different *testing.T for subtests, so that
			// go test recognizes them as proper subtests for output formatting
			// and running individual subtests
			subTestT := suite.T()
			suite.NotEqual(subTestT, suiteT)
		})
		suite.Equal(suiteT, suite.T())
	}
}

// TestRunSuite will be run by the 'go test' command, so within it, we
// can run our suite using the Run(*testing.T, TestingSuite) function.
func TestRunSuite(t *testing.T) {
	suiteTester := new(SuiteTester)
	Run(t, suiteTester)

	// Normally, the test would end here.  The following are simply
	// some assertions to ensure that the Run function is working as
	// intended - they are not part of the example.

	// The suite was only run once, so the SetupSuite and TearDownSuite
	// methods should have each been run only once.
	assert.Equal(t, suiteTester.SetupSuiteRunCount, 1)
	assert.Equal(t, suiteTester.TearDownSuiteRunCount, 1)

	assert.Equal(t, len(suiteTester.SuiteNameAfter), 4)
	assert.Equal(t, len(suiteTester.SuiteNameBefore), 4)
	assert.Equal(t, len(suiteTester.TestNameAfter), 4)
	assert.Equal(t, len(suiteTester.TestNameBefore), 4)

	assert.Contains(t, suiteTester.TestNameAfter, "TestOne")
	assert.Contains(t, suiteTester.TestNameAfter, "TestTwo")
	assert.Contains(t, suiteTester.TestNameAfter, "TestSkip")
	assert.Contains(t, suiteTester.TestNameAfter, "TestSubtest")

	assert.Contains(t, suiteTester.TestNameBefore, "TestOne")
	assert.Contains(t, suiteTester.TestNameBefore, "TestTwo")
	assert.Contains(t, suiteTester.TestNameBefore, "TestSkip")
	assert.Contains(t, suiteTester.TestNameBefore, "TestSubtest")

	for _, suiteName := range suiteTester.SuiteNameAfter {
		assert.Equal(t, "SuiteTester", suiteName)
	}

	for _, suiteName := range suiteTester.SuiteNameBefore {
		assert.Equal(t, "SuiteTester", suiteName)
	}

	for _, when := range suiteTester.TimeAfter {
		assert.False(t, when.IsZero())
	}

	for _, when := range suiteTester.TimeBefore {
		assert.False(t, when.IsZero())
	}

	// There are four test methods (TestOne, TestTwo, TestSkip, and TestSubtest), so
	// the SetupTest and TearDownTest methods (which should be run once for
	// each test) should have been run four times.
	assert.Equal(t, suiteTester.SetupTestRunCount, 4)
	assert.Equal(t, suiteTester.TearDownTestRunCount, 4)

	// Each test should have been run once.
	assert.Equal(t, suiteTester.TestOneRunCount, 1)
	assert.Equal(t, suiteTester.TestTwoRunCount, 1)
	assert.Equal(t, suiteTester.TestSubtestRunCount, 1)

	// Methods that don't match the test method identifier shouldn't
	// have been run at all.
	assert.Equal(t, suiteTester.NonTestMethodRunCount, 0)

	suiteSkipTester := new(SuiteSkipTester)
	Run(t, suiteSkipTester)

	// The suite was only run once, so the SetupSuite and TearDownSuite
	// methods should have each been run only once, even though SetupSuite
	// called Skip()
	assert.Equal(t, suiteSkipTester.SetupSuiteRunCount, 1)
	assert.Equal(t, suiteSkipTester.TearDownSuiteRunCount, 1)

}

func TestSuiteGetters(t *testing.T) {
	suite := new(SuiteTester)
	suite.SetT(t)
	assert.NotNil(t, suite.Assert())
	assert.Equal(t, suite.Assertions, suite.Assert())
	assert.NotNil(t, suite.Require())
	assert.Equal(t, suite.require, suite.Require())
}

type SuiteLoggingTester struct {
	Suite
}

func (s *SuiteLoggingTester) TestLoggingPass() {
	s.T().Log("TESTLOGPASS")
}

func (s *SuiteLoggingTester) TestLoggingFail() {
	s.T().Log("TESTLOGFAIL")
	assert.NotNil(s.T(), nil) // expected to fail
}

type StdoutCapture struct {
	oldStdout *os.File
	readPipe  *os.File
}

func (sc *StdoutCapture) StartCapture() {
	sc.oldStdout = os.Stdout
	sc.readPipe, os.Stdout, _ = os.Pipe()
}

func (sc *StdoutCapture) StopCapture() (string, error) {
	if sc.oldStdout == nil || sc.readPipe == nil {
		return "", errors.New("StartCapture not called before StopCapture")
	}
	os.Stdout.Close()
	os.Stdout = sc.oldStdout
	bytes, err := ioutil.ReadAll(sc.readPipe)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func TestSuiteLogging(t *testing.T) {
	suiteLoggingTester := new(SuiteLoggingTester)
	capture := StdoutCapture{}
	internalTest := testing.InternalTest{
		Name: "SomeTest",
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
