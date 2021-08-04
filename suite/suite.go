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

var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }
var matchMethod = flag.String("testify.m", "", "regular expression to select tests of the testify suite to run")

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

// setT sets the current *testing.T context.
func (suite *Suite) setT(t *testing.T) {
	if suite.t != nil {
		panic("suite.t already set, can't overwrite")
	}
	suite.t = t
	suite.Assertions = assert.New(t)
	suite.require = require.New(t)
}

func (suite *Suite) clearT() {
	suite.t = nil
	suite.Assertions = nil
	suite.require = nil
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

func failOnPanic(t *testing.T) {
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
// Deprecated: This method doesn't handle parallel sub-tests and will be removed in v2.
func (suite *Suite) Run(name string, subtest func()) bool {
	oldT := suite.T()
	defer func() {
		suite.clearT()
		suite.setT(oldT)
	}()
	return oldT.Run(name, func(t *testing.T) {
		suite.clearT()
		suite.setT(t)
		subtest()
	})
}

// Run takes a testing suite and runs all of the tests attached
// to it.
func Run(t *testing.T, suite TestingSuite) {
	defer failOnPanic(t)

	var suiteSetupDone bool

	var stats *SuiteInformation
	if _, ok := suite.(WithStats); ok {
		stats = newSuiteInformation()
	}

	tests := []testing.InternalTest{}
	methodFinder := reflect.TypeOf(suite)
	suiteName := methodFinder.Elem().Name()

	t.Run("All", func(t *testing.T) {
		defer failOnPanic(t)

		suite.setT(t)

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
					setupAllSuite.SetupSuite()
				}

				suiteSetupDone = true
			}

			test := testing.InternalTest{
				Name: method.Name,
				F: func(t *testing.T) {
					defer failOnPanic(t)
					childSuite := suite
					if c, ok := suite.(CopySuite); ok {
						childSuite = c.Copy()
						childSuite.clearT()
					}
					childSuite.setT(t)
					defer func() {
						if _, ok := suite.(CopySuite); !ok {
							defer suite.clearT()
						}
						if stats != nil {
							passed := !t.Failed()
							stats.end(method.Name, passed)
						}

						if tearDownTestSuite, ok := childSuite.(TearDownTestSuite); ok {
							tearDownTestSuite.TearDownTest()
						}

						if afterTestSuite, ok := childSuite.(AfterTest); ok {
							afterTestSuite.AfterTest(suiteName, method.Name)
						}
					}()

					if beforeTestSuite, ok := childSuite.(BeforeTest); ok {
						beforeTestSuite.BeforeTest(methodFinder.Elem().Name(), method.Name)
					}
					if setupTestSuite, ok := childSuite.(SetupTestSuite); ok {
						setupTestSuite.SetupTest()
					}

					if stats != nil {
						stats.start(method.Name)
					}

					method.Func.Call([]reflect.Value{reflect.ValueOf(childSuite)})
				},
			}
			tests = append(tests, test)
		}

		suite.clearT()
		runTests(t, tests)
	})
	if suiteSetupDone {
		suite.setT(t)
		if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
			tearDownAllSuite.TearDownSuite()
		}

		if suiteWithStats, measureStats := suite.(WithStats); measureStats {
			stats.End = time.Now()
			suiteWithStats.HandleStats(suiteName, stats)
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

func runTests(t *testing.T, tests []testing.InternalTest) {
	if len(tests) == 0 {
		t.Log("warning: no tests to run")
		return
	}

	for _, test := range tests {
		t.Run(test.Name, test.F)
	}
}
