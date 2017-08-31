package suite

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }
var matchMethod = flag.String("testify.m", "", "regular expression to select tests of the testify suite to run")

type BaseSuite struct {
	*assert.Assertions
	require *require.Assertions
	tb      testing.TB
}

// TB retrieves the current testing.TB context.
func (suite *BaseSuite) TB() testing.TB {
	return suite.tb
}

// SetTB sets the current testing.TB context.
func (suite *BaseSuite) SetTB(tb testing.TB) {
	suite.tb = tb
	suite.Assertions = assert.New(tb)
	suite.require = require.New(tb)
}

// Suite is a basic testing suite with methods for storing and
// retrieving the current *testing.T context.
type Suite struct {
	BaseSuite
	t *testing.T
}

// BenchmarkSuite is a basic testing suite with methods for storing and
// retrieving the current *testing.T context.
type BenchmarkSuite struct {
	BaseSuite
	b *testing.B
}

// T retrieves the current *testing.T context.
func (suite *Suite) T() *testing.T {
	return suite.t
}

// B retrieves the current *testing.B context.
func (suite *BenchmarkSuite) B() *testing.B {
	return suite.b
}

// SetT sets the current *testing.T context.
func (suite *Suite) SetT(t *testing.T) {
	suite.t = t
	suite.SetTB(t)
}

// SetB sets the current *testing.B context.
func (suite *BenchmarkSuite) SetB(b *testing.B) {
	suite.b = b
	suite.SetTB(b)
}

// Require returns a require context for suite.
func (suite *BaseSuite) Require() *require.Assertions {
	if suite.require == nil {
		suite.require = require.New(suite.TB())
	}
	return suite.require
}

// Assert returns an assert context for suite.  Normally, you can call
// `suite.NoError(expected, actual)`, but for situations where the embedded
// methods are overridden (for example, you might want to override
// assert.Assertions with require.Assertions), this method is provided so you
// can call `suite.Assert().NoError()`.
func (suite *BaseSuite) Assert() *assert.Assertions {
	if suite.Assertions == nil {
		suite.Assertions = assert.New(suite.TB())
	}
	return suite.Assertions
}

// Run takes a testing suite and runs all of the tests attached
// to it.
func Run(t *testing.T, suite TestingSuite) {
	suite.SetT(t)

	if setupAllSuite, ok := suite.(SetupAllSuite); ok {
		setupAllSuite.SetupSuite()
	}
	defer func() {
		if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
			tearDownAllSuite.TearDownSuite()
		}
	}()

	methodFinder := reflect.TypeOf(suite)
	suiteName := methodFinder.Elem().Name()

	tests := []testing.InternalTest{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := methodFilter(method.Name, "Test")
		if err != nil {
			fmt.Fprintf(os.Stderr, "testify: invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}
		if ok {
			test := testing.InternalTest{
				Name: method.Name,
				F: func(t *testing.T) {
					parentT := suite.T()
					suite.SetT(t)
					beforeHandlers(suite, suiteName, method.Name)

					defer func() {
						afterHandlers(suite, suiteName, method.Name)
						suite.SetT(parentT)
					}()

					method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
				},
			}
			tests = append(tests, test)
		}
	}
	runTests(t, tests)
}

func beforeHandlers(suite interface{}, suiteName, methodName string) {
	if setupTestSuite, ok := suite.(SetupTestSuite); ok {
		setupTestSuite.SetupTest()
	}
	if beforeTestSuite, ok := suite.(BeforeTest); ok {
		beforeTestSuite.BeforeTest(suiteName, methodName)
	}
}

func afterHandlers(suite interface{}, suiteName, methodName string) {
	if afterTestSuite, ok := suite.(AfterTest); ok {
		afterTestSuite.AfterTest(suiteName, methodName)
	}
	if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
		tearDownTestSuite.TearDownTest()
	}
}

func runTests(t testing.TB, tests []testing.InternalTest) {
	r, ok := t.(testRunner)
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

// RunBenchmark takes a benchmarking suite and runs all of the benchmarks
// attached to it.
func RunBenchmark(b *testing.B, suite BenchmarkingSuite) {
	suite.SetB(b)

	if setupAllSuite, ok := suite.(SetupAllSuite); ok {
		setupAllSuite.SetupSuite()
	}
	defer func() {
		if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
			tearDownAllSuite.TearDownSuite()
		}
	}()

	methodFinder := reflect.TypeOf(suite)
	suiteName := methodFinder.Elem().Name()

	benchmarks := []testing.InternalBenchmark{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := methodFilter(method.Name, "Benchmark")
		if err != nil {
			fmt.Fprintf(os.Stderr, "testify: invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}
		if ok {
			test := testing.InternalBenchmark{
				Name: method.Name,
				F: func(b *testing.B) {
					parentB := suite.B()
					suite.SetB(b)
					beforeHandlers(suite, suiteName, method.Name)

					defer func() {
						afterHandlers(suite, suiteName, method.Name)
						suite.SetB(parentB)
					}()

					method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
				},
			}
			benchmarks = append(benchmarks, test)
		}
	}
	runBenchmarks(b, benchmarks)
}

func runBenchmarks(t testing.TB, benchmarks []testing.InternalBenchmark) {
	r, ok := t.(benchmarkRunner)

	if !ok { // backwards compatibility with Go 1.6 and below
		testing.RunBenchmarks(allTestsFilter, benchmarks)
	}

	for _, benchmark := range benchmarks {
		r.Run(strings.TrimPrefix(benchmark.Name, "Benchmark"), benchmark.F)
	}
}

// Filtering method according to set regular expression
// specified command-line argument -m
func methodFilter(name, prefix string) (bool, error) {
	if !strings.HasPrefix(name, prefix) {
		return false, nil
	}
	return regexp.MatchString(*matchMethod, name)
}

type testRunner interface {
	Run(name string, f func(t *testing.T)) bool
}

type benchmarkRunner interface {
	Run(name string, f func(b *testing.B)) bool
}
