package suite

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/mitchellh/copystructure"
)

// RunParallel takes a testing suite and runs all of the tests attached
// to it. This function is required for test suites that want to execute concurrently
// via t.Parallel().
//
// Example:
// type ParallelTestSuite struct {
//    internalState int
// }
//
// func (s *ParallelTestSuite) SetupTest() {
//     s.internalState++
// }
//
// func (s *ParallelTestSuite) TearDownTest() {
// 	s.internalState++
// }
//
// All test methods in the suite will be run in parallel.
// All test methods must accept t *testing.T as their only argument.
// func (s *ParallelTestSuite) Test(t *testing.T) {
//     // do something
// }
//
// Tests should be run this way.
// func TestParallelSuite(t *testing.T) {
// 	RunParallel(t, new(ParallelTestSuite))
// }
// RunParallel creates a copy of the test suite for each test method to prevent sharing data.
func RunParallel(t *testing.T, suite interface{}) {
	if setupAllSuite, ok := suite.(SetupAllSuite); ok {
		setupAllSuite.SetupSuite()
	}
	defer func() {
		if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
			tearDownAllSuite.TearDownSuite()
		}
	}()

	methodFinder := reflect.TypeOf(suite)
	tests := []testing.InternalTest{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := methodFilter(method.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testify: invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}

		if ok {
			var err error
			suite, err = copystructure.Copy(suite)
			if err != nil {
				t.Fatal(err)
			}

			test := testing.InternalTest{
				Name: method.Name,
				F:    makeInnerTest(t, suite, method),
			}
			tests = append(tests, test)
		}
	}

	if !testing.RunTests(func(_, _ string) (bool, error) { return true, nil },
		tests) {
		t.Fail()
	}
}

// makeInnerTest is it's own function to prevent closing in the suite value.
// This allows us to pass in different copies of the suite value which is
// would not be possible if makeInnerTest's logic was defined as an anonymous
// function in RunParallel.
func makeInnerTest(t *testing.T, suite interface{}, method reflect.Method) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		setupTestSuite, ok := suite.(SetupTestSuite)
		if ok {
			setupTestSuite.SetupTest()
		}

		defer func() {
			if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
				tearDownTestSuite.TearDownTest()
			}
		}()

		method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(t)})
	}
}
