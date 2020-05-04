package suite

import (
	"flag"
	"log"
	"reflect"
	"regexp"
	"runtime/debug"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }
var matchMethod = flag.String("testify.m", "", "regular expression to select tests of the testify suite to run")
var unitTestMatcher = regexp.MustCompile("^Test")

// Suite is a basic testing suite with methods for storing and
// retrieving the current *testing.T context.
type Suite struct {
	*assert.Assertions
	require *require.Assertions
	t       *testing.T
}

// T retrieves the current *testing.T context.
func (suite *Suite) T() *testing.T {
	return suite.t
}

// SetT sets the current *testing.T context.
func (suite *Suite) SetT(t *testing.T) {
	suite.t = t
	suite.Assertions = assert.New(t)
	suite.require = require.New(t)
}

// Require returns a require context for suite.
func (suite *Suite) Require() *require.Assertions {
	if suite.require == nil {
		suite.require = require.New(suite.T())
	}
	return suite.require
}

// Assert returns an assert context for suite.  Normally, you can call
// `suite.NoError(expected, actual)`, but for situations where the embedded
// methods are overridden (for example, you might want to override
// assert.Assertions with require.Assertions), this method is provided so you
// can call `suite.Assert().NoError()`.
func (suite *Suite) Assert() *assert.Assertions {
	if suite.Assertions == nil {
		suite.Assertions = assert.New(suite.T())
	}
	return suite.Assertions
}

// Run provides suite functionality around golang subtests.  It should be
// called in place of t.Run(name, func(t *testing.T)) in test suite code.
// The passed-in func will be executed as a subtest with a fresh instance of t.
// Provides compatibility with go test pkg -run TestSuite/TestName/SubTestName.
func (suite *Suite) Run(name string, subtest func()) bool {
	oldT := suite.T()
	defer suite.SetT(oldT)
	return oldT.Run(name, func(t *testing.T) {
		suite.SetT(t)
		subtest()
	})
}

// Run takes a testing suite and runs all of the tests attached
// to it.
func Run(t *testing.T, suite TestingSuite) {
	defer failOnPanic(t)

	testsSync := &sync.WaitGroup{}
	suite.SetT(t)

	var suiteSetupDone bool

	var stats *SuiteInformation
	if _, ok := suite.(WithStats); ok {
		stats = newSuiteInformation()
	}

	tests := []testing.InternalTest{}
	suiteInstance := reflect.TypeOf(suite)

	for i := 0; i < suiteInstance.NumMethod(); i++ {
		method := suiteInstance.Method(i)

		ok, err := isUnitTest(method.Name)
		if err != nil {
			log.Fatalf("testify: invalid regexp for -m: %s\n", err)
		}

		if !ok {
			continue
		}

		suiteName := suiteInstance.Elem().Name()

		if !suiteSetupDone {
			if stats != nil {
				stats.Start = time.Now()
			}

			defer func() {
				if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
					testsSync.Wait()
					tearDownAllSuite.TearDownSuite()
				}

				if suiteWithStats, measureStats := suite.(WithStats); measureStats {
					stats.End = time.Now()
					suiteWithStats.HandleStats(suiteName, stats)
				}
			}()

			if setupAllSuite, ok := suite.(SetupAllSuite); ok {
				setupAllSuite.SetupSuite()
			}
			suiteSetupDone = true
		}

		test := testing.InternalTest{
			Name: method.Name,
			F: func(t *testing.T) {
				defer testsSync.Done()
				parentT := suite.T()
				suite.SetT(t)
				defer failOnPanic(t)
				defer func() {
					if stats != nil {
						stats.end(method.Name, !t.Failed())
					}

					if afterTestSuite, ok := suite.(AfterTest); ok {
						afterTestSuite.AfterTest(suiteName, method.Name)
					}

					if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
						tearDownTestSuite.TearDownTest()
					}

					suite.SetT(parentT)
				}()

				if setupTestSuite, ok := suite.(SetupTestSuite); ok {
					setupTestSuite.SetupTest()
				}
				if beforeTestSuite, ok := suite.(BeforeTest); ok {
					beforeTestSuite.BeforeTest(suiteInstance.Elem().Name(), method.Name)
				}

				if stats != nil {
					stats.start(method.Name)
				}

				method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
			},
		}
		tests = append(tests, test)
		testsSync.Add(1)
	}
	runTests(t, tests)
}

func failOnPanic(t *testing.T) {
	r := recover()
	if r != nil {
		t.Errorf("test panicked: %v\n%s", r, debug.Stack())
		t.FailNow()
	}
}

// isUnitTest indicates whether a function conforms to the Go standard naming
// for unit tests or the custom matcher specified by testify.m
func isUnitTest(funcName string) (bool, error) {
	if ok := unitTestMatcher.MatchString(funcName); ok {
		return true, nil
	}
	if len(*matchMethod) == 0 {
		return false, nil
	}
	return regexp.MatchString(*matchMethod, funcName)
}

func runTests(t testing.TB, tests []testing.InternalTest) {
	if len(tests) == 0 {
		t.Log("warning: no tests to run")
		return
	}

	r, ok := t.(runner)
	if !ok { // backwards compatibility with Go 1.6 and below
		if !testing.RunTests(allTestsFilter, tests) {
			t.Fail()
		}
		return
	}

	for _, test := range tests {
		r.Run(test.Name, test.F)
	}
}

type runner interface {
	Run(name string, f func(t *testing.T)) bool
}
