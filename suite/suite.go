package suite

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var matchMethod = flag.String("testify.m", "", "regular expression to select tests of the testify suite to run")

// T is a basic testing suite with methods for accessing
// appropriate *testing.T objects in tests.
type T struct {
	*assert.Assertions
	require  *require.Assertions
	testingT testingT
}

func (t *T) Cleanup(f func()) {
	t.testingT.Cleanup(f)
}

func (t *T) Error(args ...interface{}) {
	t.testingT.Error(args...)
}

func (t *T) Errorf(format string, args ...interface{}) {
	t.testingT.Errorf(format, args...)
}

func (t *T) Fail() {
	t.testingT.Fail()
}

func (t *T) FailNow() {
	t.testingT.FailNow()
}

func (t *T) Failed() bool {
	return t.testingT.Failed()
}

func (t *T) Fatal(args ...interface{}) {
	t.testingT.Fatal(args...)
}

func (t *T) Fatalf(format string, args ...interface{}) {
	t.testingT.Fatalf(format, args...)
}

func (t *T) Helper() {
	t.testingT.Helper()
}

func (t *T) Log(args ...interface{}) {
	t.testingT.Log(args...)
}

func (t *T) Logf(format string, args ...interface{}) {
	t.testingT.Logf(format, args...)
}

func (t *T) Name() string {
	return t.testingT.Name()
}
func (t *T) Skip(args ...interface{}) {
	t.testingT.Skip(args...)
}

func (t *T) SkipNow() {
	t.testingT.SkipNow()
}

func (t *T) Skipf(format string, args ...interface{}) {
	t.testingT.Skipf(format, args...)
}
func (t *T) Skipped() bool {
	return t.testingT.Skipped()
}

func (t *T) TempDir() string {
	if testingT, ok := t.testingT.(testingT115); ok {
		return testingT.TempDir()
	}
	panic("*testing.T does not support TempDir()")
}

func (t *T) Deadline() (deadline time.Time, ok bool) {
	if testingT, ok := t.testingT.(testingT115); ok {
		return testingT.Deadline()
	}
	panic("*testing.T does not support Deadline()")
}

func (t *T) Setenv(key, value string) {
	if testingT, ok := t.testingT.(testingT117); ok {
		testingT.Setenv(key, value)
		return
	}
	panic("*testing.T does not support Setenv()")
}

func (t *T) Parallel() {
	t.testingT.Parallel()
}

// setT sets the current *testing.T context.
func (t *T) setT(testingT *testing.T) {
	if t.testingT != nil {
		panic("T.testingT already set, can't overwrite")
	}
	t.testingT = testingT
	t.Assertions = assert.New(testingT)
	t.require = require.New(testingT)
}

// Require returns a require context for suite.
func (t *T) Require() *require.Assertions {
	if t.testingT == nil {
		panic("T.testingT not set, can't get Require object")
	}
	return t.require
}

func failOnPanic(t *T) {
	r := recover()
	if r != nil {
		t.Errorf("test panicked: %v\n%s", r, debug.Stack())
		t.FailNow()
	}
}

// Run provides suite functionality around golang subtests.  It should be
// called in place of t.Run(name, func(t *testing.T)) in test suite code.
// The passed-in func will be executed as a subtest with a fresh instance of t.
// Provides compatibility with go test pkg -run TestSuite/TestName/SubTestName.
func (t *T) Run(name string, subtest func(t *T)) bool {
	return t.testingT.Run(name, func(testingT *testing.T) {
		t := &T{}
		t.setT(testingT)
		subtest(t)
	})
}

// Run takes a testing suite and runs all of the tests attached
// to it.
func Run(testingT *testing.T, suite interface{}) {
	t := &T{}
	t.setT(testingT)

	defer failOnPanic(t)

	var suiteSetupDone bool

	var stats *SuiteInformation
	if _, ok := suite.(WithStats); ok {
		stats = newSuiteInformation()
	}

	tests := []testing.InternalTest{}
	methodFinder := reflect.TypeOf(suite)
	suiteName := methodFinder.Elem().Name()

	for i := 0; i < methodFinder.NumMethod(); i++ {
		method := methodFinder.Method(i)

		ok, err := methodFilter(method.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testify: invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}

		if !ok {
			continue
		}

		if !suiteSetupDone {
			if stats != nil {
				stats.Start = time.Now()
			}

			if setupAllSuite, ok := suite.(SetupAllSuite); ok {
				setupAllSuite.SetupSuite(t)
			}

			suiteSetupDone = true
		}

		test := testing.InternalTest{
			Name: method.Name,
			F: func(testingT *testing.T) {
				t := &T{}
				t.setT(testingT)

				defer failOnPanic(t)

				defer func() {
					if stats != nil {
						passed := !t.Failed()
						stats.end(method.Name, passed)
					}

					if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
						tearDownTestSuite.TearDownTest(t)
					}

					if afterTestSuite, ok := suite.(AfterTest); ok {
						afterTestSuite.AfterTest(t, suiteName, method.Name)
					}
				}()

				if beforeTestSuite, ok := suite.(BeforeTest); ok {
					beforeTestSuite.BeforeTest(t, methodFinder.Elem().Name(), method.Name)
				}
				if setupTestSuite, ok := suite.(SetupTestSuite); ok {
					setupTestSuite.SetupTest(t)
				}

				if stats != nil {
					stats.start(method.Name)
				}

				method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(t)})
			},
		}
		tests = append(tests, test)
	}

	if len(tests) == 0 {
		testingT.Log("warning: no tests to run")
		return
	}

	// run sub-tests in a group so tearDownSuite is called in the right order
	testingT.Run("All", func(testingT *testing.T) {
		for _, test := range tests {
			testingT.Run(test.Name, test.F)
		}
	})
	if suiteSetupDone {
		if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
			tearDownAllSuite.TearDownSuite(t)
		}

		if suiteWithStats, measureStats := suite.(WithStats); measureStats {
			stats.End = time.Now()
			suiteWithStats.HandleStats(t, suiteName, stats)
		}
	}
}

// Filtering method according to set regular expression
// specified command-line argument -m
func methodFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^Test", name); !ok {
		return false, nil
	}
	return regexp.MatchString(*matchMethod, name)
}
