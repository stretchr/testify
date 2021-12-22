package suite

import (
	"flag"
	"fmt"
	"log"
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
	f       *testing.F
}

// T retrieves the current *testing.T context.
func (suite *Suite) T() *testing.T {
	return suite.t
}

// F retrieves the current *testing.F context.
func (suite *Suite) F() *testing.F {
	return suite.f
}

// SetT sets the current *testing.T context.
func (suite *Suite) SetT(t *testing.T) {
	suite.t = t
	suite.Assertions = assert.New(t)
	suite.require = require.New(t)
}

// SetT sets the current *testing.T context.
func (suite *Suite) SetF(f *testing.F) {
	suite.f = f
	suite.Assertions = assert.New(f)
	suite.require = require.New(f)
}

// Require returns a require context for suite.
func (suite *Suite) Require() *require.Assertions {
	if suite.require == nil {

		var tT require.TestingT
		if suite.T() != nil {
			tT = suite.T()
		}
		if suite.F() != nil {
			tT = suite.F()
		}

		if tT == nil {
			panic("Neither T nor F was non-nil")
		}
		log.Printf("tT:%v\n", tT)
		suite.require = require.New(tT)
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
		var tT assert.TestingT
		if suite.T() != nil {
			tT = suite.T()
		}
		if suite.F() != nil {
			tT = suite.F()
		}
		suite.Assertions = assert.New(tT)
	}
	return suite.Assertions
}

func failOnPanic(t testing.TB) {
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
func Run(t testing.TB, suite TestingSuite) {
	defer failOnPanic(t)

	if tT, ok := t.(*testing.T); ok {
		suite.SetT(tT)
	}

	if fF, ok := t.(*testing.F); ok {
		suite.SetF(fF)

	}

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
				setupAllSuite.SetupSuite()
			}

			suiteSetupDone = true
		}

		test := testing.InternalTest{
			Name: method.Name,
			F: func(t *testing.T) {
				parentT := suite.T()
				suite.SetT(t)
				defer failOnPanic(t)
				defer func() {
					if stats != nil {
						passed := !t.Failed()
						stats.end(method.Name, passed)
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
					beforeTestSuite.BeforeTest(methodFinder.Elem().Name(), method.Name)
				}

				if stats != nil {
					stats.start(method.Name)
				}

				method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
			},
		}
		tests = append(tests, test)
	}
	if suiteSetupDone {
		defer func() {
			if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
				tearDownAllSuite.TearDownSuite()
			}

			if suiteWithStats, measureStats := suite.(WithStats); measureStats {
				stats.End = time.Now()
				suiteWithStats.HandleStats(suiteName, stats)
			}
		}()
	}

	runTests(t, tests)
}

// Filtering method according to set regular expression
// specified command-line argument -m
func methodFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^(Test|Fuzz)", name); !ok {
		return false, nil
	}
	return regexp.MatchString(*matchMethod, name)
}

func runTests(t testing.TB, tests []testing.InternalTest) {
	if len(tests) == 0 {
		t.Log("warning: no tests to run")
		return
	}

	r, ok := t.(runner)
	if !ok { // backwards compatibility with Go 1.6 and below
		// t.F doesn't have Run() so Fuzz tests also
		// take this branch
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
