package suite

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
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
// Setup and Teardown functions which need to modify test state can accept pointer receivers safely.
// The calling of these functions is synchronized between the tests.
// func (s *ParallelTestSuite) SetupTest() {
//     s.internalState++
// }
//
// func (s *ParallelTestSuite) TearDownTest() {
// 	s.internalState++
// }
//
// All test methods in the suite will be run in parallel.
// All tests should accept value receivers to prevent data races in the test suite.
// All tests must accept t *testing.T as their only argument.
// func (s ParallelTestSuite) Test(t *testing.T) {
//     // do something
// }
//
// Tests should be run this way.
// func TestParallelSuite(t *testing.T) {
// 	RunParallel(t, new(ParallelTestSuite))
// }
// RunParallel synchronizes calls to SetupTestSuite and TearDownTestSuite. Any data accessed outside
// of those functions must be manually synchronized to be goroutine safe.
func RunParallel(t *testing.T, suite interface{}) {
	if setupAllSuite, ok := suite.(SetupAllSuite); ok {
		setupAllSuite.SetupSuite()
	}
	defer func() {
		if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
			tearDownAllSuite.TearDownSuite()
		}
	}()

	// Locks for synchronizing access to SetupTest() and TearDownTest().
	setupMtx, teardownMtx := sync.Mutex{}, sync.Mutex{}

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
			test := testing.InternalTest{
				Name: method.Name,
				F: func(t *testing.T) {
					t.Parallel()
					setupTestSuite, ok := suite.(SetupTestSuite)
					if ok {
						setupMtx.Lock()
						setupTestSuite.SetupTest()
						// This will deadlock if SetupTest() panics.
						// We cannot defer unlocking here because that would seralize running the tests.
						setupMtx.Unlock()
					}

					defer func() {
						if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
							teardownMtx.Lock()
							tearDownTestSuite.TearDownTest()
							// This will deadlock if TearDownTest() panics.
							teardownMtx.Unlock()
						}
					}()

					method.Func.Call([]reflect.Value{reflect.ValueOf(suite), reflect.ValueOf(t)})
				},
			}
			tests = append(tests, test)
		}
	}

	if !testing.RunTests(func(_, _ string) (bool, error) { return true, nil },
		tests) {
		t.Fail()
	}
}
