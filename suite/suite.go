package suite

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }
var matchMethod = flag.String("testify.m", "", "regular expression to select tests of the testify suite to run")

// Suite is a basic testing suite with methods for storing and
// retrieving the current *testing.T context.
type Suite struct {
	*assert.Assertions

	mu      sync.RWMutex
	require *require.Assertions
	t       *testing.T

	// Parent suite to have access to the implemented methods of parent struct
	s TestingSuite

	isParallelTest map[string]struct{}
}

// T retrieves the current *testing.T context.
func (suite *Suite) T() *testing.T {
	suite.mu.RLock()
	defer suite.mu.RUnlock()

	if suite.isInSuiteParallelMethod() {
		panic("Avoid T() in parallel tests. Use the passed-in one instead.")
	}

	return suite.t
}

// SetT sets the current *testing.T context.
func (suite *Suite) SetT(t *testing.T) {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	suite.t = t
	suite.Assertions = assert.New(t)
	suite.require = require.New(t)
}

// SetS needs to set the current test suite as parent
// to get access to the parent methods
func (suite *Suite) SetS(s TestingSuite) {
	suite.s = s
}

// Require returns a require context for suite.
func (suite *Suite) Require() *require.Assertions {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	if suite.require == nil {
		panic("'Require' must not be called before 'Run' or 'SetT'")
	}

	if suite.isInSuiteParallelMethod() {
		panic("Avoid Require() in parallel tests. Use require(t, ...) instead.")
	}

	return suite.require
}

// Assert returns an assert context for suite.  Normally, you can call
// `suite.NoError(expected, actual)`, but for situations where the embedded
// methods are overridden (for example, you might want to override
// assert.Assertions with require.Assertions), this method is provided so you
// can call `suite.Assert().NoError()`.
func (suite *Suite) Assert() *assert.Assertions {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	if suite.Assertions == nil {
		panic("'Assert' must not be called before 'Run' or 'SetT'")
	}

	if suite.isInSuiteParallelMethod() {
		panic("Avoid Assert() in parallel tests. Use assert(t, ...) instead.")
	}

	return suite.Assertions
}

func (suite *Suite) isInSuiteParallelMethod() bool {
	for i := 1; ; i++ {
		pc, _, _, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// Example rawFuncName:
		// github.com/foo/bar/tests/e2e.(*E2ETestSuite).MyTest
		rawFuncName := runtime.FuncForPC(pc).Name()
		splittedFuncName := strings.Split(rawFuncName, ".")
		funcName := splittedFuncName[len(splittedFuncName)-1]

		if _, isParallel := suite.isParallelTest[funcName]; isParallel {
			return true
		}
	}

	return false
}

func recoverAndFailOnPanic(t *testing.T) {
	t.Helper()
	r := recover()
	failOnPanic(t, r)
}

func failOnPanic(t *testing.T, r interface{}) {
	t.Helper()
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

	return oldT.Run(name, func(t *testing.T) {
		suite.SetT(t)
		defer suite.SetT(oldT)

		defer recoverAndFailOnPanic(t)

		if setupSubTest, ok := suite.s.(SetupSubTest); ok {
			setupSubTest.SetupSubTest()
		}

		if tearDownSubTest, ok := suite.s.(TearDownSubTest); ok {
			defer tearDownSubTest.TearDownSubTest()
		}

		subtest()
	})
}

// Run takes a testing suite and runs all of the tests attached
// to it.
func Run(t *testing.T, suite TestingSuite) {
	defer recoverAndFailOnPanic(t)

	suite.SetT(t)
	suite.SetS(suite)

	var suiteSetupDone bool

	var stats *SuiteInformation
	if _, ok := suite.(WithStats); ok {
		stats = newSuiteInformation()
	}

	tests := []testing.InternalTest{}
	parallelTests := []testing.InternalTest{}

	methodFinder := reflect.TypeOf(suite)
	suiteName := methodFinder.Elem().Name()

	setupTestSuite, hasSetupTestWithoutT := suite.(SetupTestSuite)
	setupTestParallelSuite, hasSetupTestWithT := suite.(SetupParallelTestSuite)

	beforeTestSuite, hasBeforeTestWithoutT := suite.(BeforeTest)
	beforeTestParallelSuite, hasBeforeTestWithT := suite.(BeforeParallelTest)

	afterTestSuite, hasAfterTestWithoutT := suite.(AfterTest)
	afterTestParallelSuite, hasAfterTestWithT := suite.(AfterParallelTest)

	tearDownTestSuite, hasTearDownTestWithoutT := suite.(TearDownTestSuite)
	tearDownTestParallelSuite, hasTearDownTestWithT := suite.(TearDownParallelTestSuite)

	testifySuiteVal := GetEmbeddedValue(suite, reflect.TypeOf(Suite{}))
	if !testifySuiteVal.IsValid() {
		panic("nononoo")
	}

	var isParallelTestPtr *map[string]struct{}

	if testifySuiteVal.IsValid() {
		isParalleTestVal := testifySuiteVal.FieldByName("isParallelTest")
		if !isParalleTestVal.IsValid() {
			panic("Should be able to get isParallelTest!")
		}

		// We need unsafe here to circumvent Go’s prevention of accessing
		// unexported values.
		ptr := unsafe.Pointer(isParalleTestVal.UnsafeAddr())
		isParallelTestPtr = (*map[string]struct{})(ptr)

		if *isParallelTestPtr == nil {
			*isParallelTestPtr = map[string]struct{}{}
		}
	}

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

		isParallel := strings.HasPrefix(method.Name, "Parallel")

		if isParallel {
			if isParallelTestPtr != nil {
				(*isParallelTestPtr)[method.Name] = struct{}{}
			}

			var faultyMethods []string

			if hasSetupTestWithoutT {
				faultyMethods = append(faultyMethods, "SetupTest")
			}
			if hasBeforeTestWithoutT {
				faultyMethods = append(faultyMethods, "BeforeTest")
			}
			if hasAfterTestWithoutT {
				faultyMethods = append(faultyMethods, "AfterTest")
			}
			if hasTearDownTestWithoutT {
				faultyMethods = append(faultyMethods, "TearDownTest")
			}

			if len(faultyMethods) > 0 {
				joined := strings.Join(faultyMethods, " and ")
				t.Errorf("Suite contains a parallel test (%#q), so %s must accept a %T.", method.Name, joined, t)
				t.FailNow()
			}
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

				if isParallel {
					t.Parallel()
				} else {
					suite.SetT(t)
				}

				defer recoverAndFailOnPanic(t)
				defer func() {
					t.Helper()

					r := recover()

					if stats != nil {
						passed := !t.Failed() && r == nil
						stats.end(method.Name, passed)
					}

					if hasAfterTestWithoutT {
						afterTestSuite.AfterTest(suiteName, method.Name)
					} else if hasAfterTestWithT {
						afterTestParallelSuite.AfterTest(t, suiteName, method.Name)
					}

					if hasTearDownTestWithoutT {
						tearDownTestSuite.TearDownTest()
					} else if hasTearDownTestWithT {
						tearDownTestParallelSuite.TearDownTest(t)
					}

					if !isParallel {
						suite.SetT(parentT)
					}

					failOnPanic(t, r)
				}()

				if hasSetupTestWithoutT {
					setupTestSuite.SetupTest()
				} else if hasSetupTestWithT {
					setupTestParallelSuite.SetupTest(t)
				}

				if hasBeforeTestWithoutT {
					beforeTestSuite.BeforeTest(methodFinder.Elem().Name(), method.Name)
				} else if hasBeforeTestWithT {
					beforeTestParallelSuite.BeforeTest(t, methodFinder.Elem().Name(), method.Name)
				}

				if stats != nil {
					stats.start(method.Name)
				}

				methodArgs := []reflect.Value{
					reflect.ValueOf(suite),
				}

				if isParallel {
					verifyParallelMethod(t, method.Name, method.Type)
					methodArgs = append(methodArgs, reflect.ValueOf(t))
				} else {
					verifySequentialMethod(t, method.Name, method.Type)
				}

				method.Func.Call(methodArgs)
			},
		}

		if isParallel {
			parallelTests = append(parallelTests, test)
		} else {
			tests = append(tests, test)
		}
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

	if len(parallelTests) > 0 {
		runTests(
			t,
			[]testing.InternalTest{
				{
					Name: "parallel",
					F: func(t *testing.T) {
						runTests(t, parallelTests)
					},
				},
			},
		)
	}
}

// Filtering method according to set regular expression
// specified command-line argument -m
func methodFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString("^(?:Parallel)?Test", name); !ok {
		return false, nil
	}
	return regexp.MatchString(*matchMethod, name)
}

func verifyParallelMethod(t *testing.T, name string, rt reflect.Type) {
	if rt.NumIn() != 2 || rt.In(1) != reflect.TypeOf(t) {
		t.Errorf("%#q method is parallel, so it must accept a %T (and only that)", name, &testing.T{})
		t.FailNow()
	}
}

func verifySequentialMethod(t *testing.T, name string, rt reflect.Type) {
	if rt.NumIn() != 1 {
		t.Errorf("%#q method is sequential, so it must accept no arguments", name)
		t.FailNow()
	}
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
