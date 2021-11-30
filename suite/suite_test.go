package suite_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// SuiteRequireTwice is intended to test the usage of suite.Require in two
// different tests
type SuiteRequireTwice struct{}

var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }

// TestSuiteRequireTwice checks for regressions of issue #149 where
// suite.requirements was not initialised in suite.SetT()
// A regression would result on these tests panicking rather than failing.
func TestSuiteRequireTwice(t *testing.T) {
	ok := testing.RunTests(
		allTestsFilter,
		[]testing.InternalTest{{
			Name: "TestSuiteRequireTwice",
			F: func(t *testing.T) {
				suite.Run(t, new(SuiteRequireTwice))
			},
		}},
	)
	assert.Equal(t, false, ok)
}

func (s *SuiteRequireTwice) TestRequireOne(t *suite.T) {
	r := t.Require()
	r.Equal(1, 2)
}

func (s *SuiteRequireTwice) TestRequireTwo(t *suite.T) {
	r := t.Require()
	r.Equal(1, 2)
}

type panickingSuite struct {
	panicInSetupSuite    bool
	panicInSetupTest     bool
	panicInBeforeTest    bool
	panicInTest          bool
	panicInAfterTest     bool
	panicInTearDownTest  bool
	panicInTearDownSuite bool
}

func (s *panickingSuite) SetupSuite(_ *suite.T) {
	if s.panicInSetupSuite {
		panic("oops in setup suite")
	}
}

func (s *panickingSuite) SetupTest(_ *suite.T) {
	if s.panicInSetupTest {
		panic("oops in setup test")
	}
}

func (s *panickingSuite) BeforeTest(_ *suite.T, _, _ string) {
	if s.panicInBeforeTest {
		panic("oops in before test")
	}
}

func (s *panickingSuite) Test(_ *suite.T) {
	if s.panicInTest {
		panic("oops in test")
	}
}

func (s *panickingSuite) AfterTest(_ *suite.T, _, _ string) {
	if s.panicInAfterTest {
		panic("oops in after test")
	}
}

func (s *panickingSuite) TearDownTest(_ *suite.T) {
	if s.panicInTearDownTest {
		panic("oops in tear down test")
	}
}

func (s *panickingSuite) TearDownSuite(_ *suite.T) {
	if s.panicInTearDownSuite {
		panic("oops in tear down suite")
	}
}

func TestSuiteRecoverPanic(t *testing.T) {
	ok := true
	panickingTests := []testing.InternalTest{
		{
			Name: "TestPanicInSetupSuite",
			F:    func(t *testing.T) { suite.Run(t, &panickingSuite{panicInSetupSuite: true}) },
		},
		{
			Name: "TestPanicInSetupTest",
			F:    func(t *testing.T) { suite.Run(t, &panickingSuite{panicInSetupTest: true}) },
		},
		{
			Name: "TestPanicInBeforeTest",
			F:    func(t *testing.T) { suite.Run(t, &panickingSuite{panicInBeforeTest: true}) },
		},
		{
			Name: "TestPanicInTest",
			F:    func(t *testing.T) { suite.Run(t, &panickingSuite{panicInTest: true}) },
		},
		{
			Name: "TestPanicInAfterTest",
			F:    func(t *testing.T) { suite.Run(t, &panickingSuite{panicInAfterTest: true}) },
		},
		{
			Name: "TestPanicInTearDownTest",
			F:    func(t *testing.T) { suite.Run(t, &panickingSuite{panicInTearDownTest: true}) },
		},
		{
			Name: "TestPanicInTearDownSuite",
			F:    func(t *testing.T) { suite.Run(t, &panickingSuite{panicInTearDownSuite: true}) },
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

	SuiteT *suite.T
	TestT  map[string]*suite.T
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (s *SuiteTester) SetupSuite(t *suite.T) {
	s.SetupSuiteRunCount++
	s.SuiteT = t
}

func (s *SuiteTester) BeforeTest(t *suite.T, suiteName, testName string) {
	s.SuiteNameBefore = append(s.SuiteNameBefore, suiteName)
	s.TestNameBefore = append(s.TestNameBefore, testName)
	s.TimeBefore = append(s.TimeBefore, time.Now())
	s.TestT[testName] = t
}

func (s *SuiteTester) AfterTest(t *suite.T, suiteName, testName string) {
	s.SuiteNameAfter = append(s.SuiteNameAfter, suiteName)
	s.TestNameAfter = append(s.TestNameAfter, testName)
	s.TimeAfter = append(s.TimeAfter, time.Now())
	// T should be from the sub-test
	assert.True(t, s.TestT[testName] == t)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (s *SuiteTester) TearDownSuite(t *suite.T) {
	s.TearDownSuiteRunCount++
	// T should be from the suite
	assert.True(t, s.SuiteT == t)
}

// The SetupTest method will be run before every test in the suite.
func (s *SuiteTester) SetupTest(t *suite.T) {
	s.SetupTestRunCount++
	var subSuites []*suite.T
	for _, value := range s.TestT {
		subSuites = append(subSuites, value)
	}
	// T should be from one of the tests, not suite
	assert.Contains(t, subSuites, t)
	assert.False(t, s.SuiteT == t)
}

// The TearDownTest method will be run after every test in the suite.
func (s *SuiteTester) TearDownTest(t *suite.T) {
	s.TearDownTestRunCount++
	var subSuites []*suite.T
	for _, value := range s.TestT {
		subSuites = append(subSuites, value)
	}
	// T should be from one of the tests, not suite
	assert.Contains(t, subSuites, t)
	assert.False(t, s.SuiteT == t)
}

// Every method in a testing suite that begins with "Test" will be run
// as a test.  TestOne is an example of a test.  For the purposes of
// this example, we've included assertions in the tests, since most
// tests will issue assertions.
func (s *SuiteTester) TestOne(t *suite.T) {
	beforeCount := s.TestOneRunCount
	s.TestOneRunCount++
	assert.Equal(t, s.TestOneRunCount, beforeCount+1)
	t.Equal(s.TestOneRunCount, beforeCount+1)
	// T should be from the right test
	assert.True(t, t == s.TestT["TestOne"])
}

// TestTwo is another example of a test.
func (s *SuiteTester) TestTwo(t *suite.T) {
	beforeCount := s.TestTwoRunCount
	s.TestTwoRunCount++
	assert.NotEqual(t, s.TestTwoRunCount, beforeCount)
	t.NotEqual(s.TestTwoRunCount, beforeCount)
	// T should be from the right test
	assert.True(t, t == s.TestT["TestTwo"])
}

func (s *SuiteTester) TestSkip(t *suite.T) {
	t.Skip()
	// T should be from the right test
	assert.True(t, t == s.TestT["TestSkip"])
}

// NonTestMethod does not begin with "Test", so it will not be run by
// testify as a test in the suite.  This is useful for creating helper
// methods for your tests.
func (s *SuiteTester) NonTestMethod() {
	s.NonTestMethodRunCount++
}

func (s *SuiteTester) TestSubtest(t *suite.T) {
	s.TestSubtestRunCount++

	// T should be from the right test
	assert.True(t, t == s.TestT["TestSubtest"])

	for _, spec := range []struct {
		testName string
	}{
		{"first"},
		{"second"},
	} {
		suiteT := t
		t.Run(spec.testName, func(t *suite.T) {
			// We should get a different *testing.T for subtests, so that
			// go test recognizes them as proper subtests for output formatting
			// and running individual subtests
			subTestT := t
			t.NotEqual(subTestT, suiteT)
		})
	}
}

type SuiteSkipTester struct {
	// Keep counts of how many times each method is run.
	SetupSuiteRunCount    int
	TearDownSuiteRunCount int
}

func (s *SuiteSkipTester) SetupSuite(t *suite.T) {
	s.SetupSuiteRunCount++
	t.Skip()
}

func (s *SuiteSkipTester) TestNothing(_ *suite.T) {
	// SetupSuite is only called when at least one test satisfies
	// test filter. For this suite to be set up (and then tore down)
	// it is necessary to add at least one test method.
}

func (s *SuiteSkipTester) TearDownSuite(_ *suite.T) {
	s.TearDownSuiteRunCount++
}

// TestRunSuite will be run by the 'go test' command, so within it, we
// can run our suite using the Run(*testing.T, TestingSuite) function.
func TestRunSuite(t *testing.T) {
	suiteTester := &SuiteTester{
		TestT: make(map[string]*suite.T),
		//SuiteT: t,
	}
	suite.Run(t, suiteTester)

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
	suite.Run(t, suiteSkipTester)

	// The suite was only run once, so SetupSuite method should have been
	// run only once. Since SetupSuite called Skip(), Teardown isn't called.
	assert.Equal(t, suiteSkipTester.SetupSuiteRunCount, 1)
	assert.Equal(t, suiteSkipTester.TearDownSuiteRunCount, 0)

}

// This suite has no Test... methods. It's setup and teardown must be skipped.
type SuiteSetupSkipTester struct {
	setUp    bool
	toreDown bool
}

func (s *SuiteSetupSkipTester) SetupSuite(_ *suite.T) {
	s.setUp = true
}

func (s *SuiteSetupSkipTester) NonTestMethod() {

}

func (s *SuiteSetupSkipTester) TearDownSuite(_ *suite.T) {
	s.toreDown = true
}

func TestSkippingSuiteSetup(t *testing.T) {
	suiteTester := new(SuiteSetupSkipTester)
	suite.Run(t, suiteTester)
	assert.False(t, suiteTester.setUp)
	assert.False(t, suiteTester.toreDown)
}

type SuiteLoggingTester struct{}

func (s *SuiteLoggingTester) TestLoggingPass(t *suite.T) {
	t.Log("TESTLOGPASS")
}

func (s *SuiteLoggingTester) TestLoggingFail(t *suite.T) {
	t.Log("TESTLOGFAIL")
	assert.NotNil(t, nil) // expected to fail
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
			suite.Run(subT, suiteLoggingTester)
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

type CallOrderSuite struct {
	callOrder []string
}

func (s *CallOrderSuite) call(method string) {
	time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
	s.callOrder = append(s.callOrder, method)
}

func TestSuiteCallOrder(t *testing.T) {
	suite.Run(t, new(CallOrderSuite))
}
func (s *CallOrderSuite) SetupSuite(_ *suite.T) {
	s.call("SetupSuite")
}

func (s *CallOrderSuite) TearDownSuite(t *suite.T) {
	s.call("TearDownSuite")
	assert.Equal(t, "SetupSuite;SetupTest;Test A;TearDownTest;SetupTest;Test B;TearDownTest;TearDownSuite", strings.Join(s.callOrder, ";"))
}

func (s *CallOrderSuite) SetupTest(_ *suite.T) {
	s.call("SetupTest")
}

func (s *CallOrderSuite) TearDownTest(_ *suite.T) {
	s.call("TearDownTest")
}

func (s *CallOrderSuite) Test_A(_ *suite.T) {
	s.call("Test A")
}

func (s *CallOrderSuite) Test_B(_ *suite.T) {
	s.call("Test B")
}

type suiteWithStats struct {
	wasCalled bool
	stats     *suite.SuiteInformation
}

func (s *suiteWithStats) HandleStats(_ *suite.T, _ string, stats *suite.SuiteInformation) {
	s.wasCalled = true
	s.stats = stats
}

func (s *suiteWithStats) TestSomething(t *suite.T) {
	t.Equal(1, 1)
}

func TestSuiteWithStats(t *testing.T) {
	suiteWithStats := new(suiteWithStats)
	suite.Run(t, suiteWithStats)

	assert.True(t, suiteWithStats.wasCalled)
	assert.NotZero(t, suiteWithStats.stats.Start)
	assert.NotZero(t, suiteWithStats.stats.End)
	assert.True(t, suiteWithStats.stats.Passed())

	testStats := suiteWithStats.stats.TestStats["TestSomething"]
	assert.NotZero(t, testStats.Start)
	assert.NotZero(t, testStats.End)
	assert.True(t, testStats.Passed)
}

// FailfastSuite will test the behavior when running with the failfast flag
// It logs calls in the callOrder slice which we then use to assert the correct calls were made
type FailfastSuite struct {
	callOrder []string
}

func (s *FailfastSuite) call(method string) {
	s.callOrder = append(s.callOrder, method)
}

func TestFailfastSuite(t *testing.T) {
	// This test suite is run twice. Once normally and once with the -failfast flag by TestFailfastSuiteFailFastOn
	// If you need to debug it run this test directly with the failfast flag set on/off as you need
	failFast := flag.Lookup("test.failfast").Value.(flag.Getter).Get().(bool)
	s := new(FailfastSuite)
	ok := testing.RunTests(
		allTestsFilter,
		[]testing.InternalTest{{
			Name: "TestFailfastSuite",
			F: func(t *testing.T) {
				suite.Run(t, s)
			},
		}},
	)
	assert.Equal(t, false, ok)
	if failFast {
		// Test A Fails and because we are running with failfast Test B never runs and we proceed straight to TearDownSuite
		assert.Equal(t, "SetupSuite;SetupTest;Test A Fails;TearDownTest;TearDownSuite", strings.Join(s.callOrder, ";"))
	} else {
		// Test A Fails and because we are running without failfast we continue and run Test B and then proceed to TearDownSuite
		assert.Equal(t, "SetupSuite;SetupTest;Test A Fails;TearDownTest;SetupTest;Test B Passes;TearDownTest;TearDownSuite", strings.Join(s.callOrder, ";"))
	}
}
func TestFailfastSuiteFailFastOn(t *testing.T) {
	// To test this with failfast on (and isolated from other intended test failures in our test suite) we launch it in its own process
	cmd := exec.Command("go", "test", "-v", "-race", "-run", "TestFailfastSuite", "-failfast")
	var out bytes.Buffer
	cmd.Stdout = &out
	t.Log("Running go test -v -race -run TestFailfastSuite -failfast")
	err := cmd.Run()
	t.Log(out.String())
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
func (s *FailfastSuite) SetupSuite(_ *suite.T) {
	s.call("SetupSuite")
}

func (s *FailfastSuite) TearDownSuite(_ *suite.T) {
	s.call("TearDownSuite")
}
func (s *FailfastSuite) SetupTest(_ *suite.T) {
	s.call("SetupTest")
}

func (s *FailfastSuite) TearDownTest(_ *suite.T) {
	s.call("TearDownTest")
}

func (s *FailfastSuite) Test_A_Fails(t *suite.T) {
	s.call("Test A Fails")
	t.Error("Test A meant to fail")
}

func (s *FailfastSuite) Test_B_Passes(t *suite.T) {
	s.call("Test B Passes")
	t.Require().True(true)
}

type parallelSuiteData struct {
	calls          []string
	callsIndex     map[string]int
	parallelSuiteT map[string]*suite.T
}

type parallelSuite struct {
	mutex *sync.Mutex
	data  *parallelSuiteData
}

func (s *parallelSuite) call(method string) {
	time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data.calls = append(s.data.calls, method)
	s.data.callsIndex[method] = len(s.data.calls) - 1
}

func TestSuiteParallel(t *testing.T) {
	data := parallelSuiteData{
		calls:          []string{},
		callsIndex:     make(map[string]int, 8),
		parallelSuiteT: make(map[string]*suite.T, 2),
	}
	s := &parallelSuite{mutex: &sync.Mutex{}, data: &data}
	suite.Run(t, s)
}

func (s *parallelSuite) SetupSuite(_ *suite.T) {
	s.call("SetupSuite")
}

func (s *parallelSuite) TearDownSuite(t *suite.T) {
	s.call("TearDownSuite")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	// first 3 calls and last call is known ordering
	assert.Equal(t, []string{"SetupSuite", "BeforeTest Test_A", "BeforeTest Test_B"}, s.data.calls[:3])
	assert.Equal(t, "TearDownSuite", s.data.calls[len(s.data.calls)-1])
	// should have these calls
	assert.Subset(t, s.data.calls[3:], []string{"Test_A", "AfterTest Test_A", "Test_B", "AfterTest Test_B"})
	// there won't be any other ordering guarantees between tests A and B since they are run in parallel,
	// but verify that AfterTest is run after the test
	assert.Greater(t, s.data.callsIndex["AfterTest Test_A"], s.data.callsIndex["Test_A"])
	assert.Greater(t, s.data.callsIndex["AfterTest Test_B"], s.data.callsIndex["Test_B"])
	// verify that copies of s are created correctly
	assert.NotEqual(t, s, s.data.parallelSuiteT["Test_A"])
	assert.NotEqual(t, s, s.data.parallelSuiteT["Test_B"])
	assert.NotEqual(t, s.data.parallelSuiteT["Test_A"], s.data.parallelSuiteT["Test_B"])
}

func (s *parallelSuite) BeforeTest(_ *suite.T, _, testName string) {
	s.call(fmt.Sprintf("BeforeTest %s", testName))
}

func (s *parallelSuite) AfterTest(_ *suite.T, _, testName string) {
	s.call(fmt.Sprintf("AfterTest %s", testName))
}

func (s *parallelSuite) Test_A(t *suite.T) {
	t.Parallel()
	s.call("Test_A")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data.parallelSuiteT["Test_A"] = t
}

func (s *parallelSuite) Test_B(t *suite.T) {
	t.Parallel()
	s.call("Test_B")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data.parallelSuiteT["Test_B"] = t
}
