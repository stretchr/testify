//go:build go1.25

package suite

import (
	"reflect"
	"testing"
	"testing/synctest"
	"time"
)

// SyncSuite is a basic testing/synctest suite with methods for storing and
// retrieving the current *testing.T context. Each method is run inside a
// synctest bubble.
type SyncSuite struct {
	sharedSuite
}

// SetTime is a helper function to set the fake clock of the synctest
// bubble to the given time. Only timestamps after 1.Jan 2000 is allowed,
// because synctest time starts at 1. Jan 2000
func (s *SyncSuite) SetTime(ts time.Time) {
	time.Sleep(time.Until(ts))
}

// Wait is just a wrapper for synctest.Wait()
func (s *SyncSuite) Wait() {
	synctest.Wait()
}

// RunSync takes a testing/synctest suite and runs all of the tests attached
// to it. Each test is run inside its own synctest bubble.
func RunSync(t *testing.T, suite TestingSuite) {
	defer recoverAndFailOnPanic(t)

	suite.SetT(t)
	suite.SetS(suite)

	var stats *SuiteInformation
	if _, ok := suite.(WithStats); ok {
		stats = newSuiteInformation()
	}

	var tests tests
	methodFinder := reflect.TypeOf(suite)
	suiteName := methodFinder.Elem().Name()

	for i := 0; i < methodFinder.NumMethod(); i++ {
		method := methodFinder.Method(i)

		if !isTestMethod(method) {
			continue
		}

		if test, ok := checkMethodSignature(method); !ok {
			tests = append(tests, test)
			continue
		}

		test := test{
			name: method.Name,
			run: func(t *testing.T) {
				synctest.Test(t, func(t *testing.T) {
					parentT := suite.T()
					suite.SetT(t)
					defer recoverAndFailOnPanic(t)
					defer func() {
						t.Helper()

						r := recover()

						stats.end(method.Name, !t.Failed() && r == nil)

						if afterTestSuite, ok := suite.(AfterTest); ok {
							afterTestSuite.AfterTest(suiteName, method.Name)
						}

						if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
							tearDownTestSuite.TearDownTest()
						}

						suite.SetT(parentT)
						failOnPanic(t, r)
					}()

					if setupTestSuite, ok := suite.(SetupTestSuite); ok {
						setupTestSuite.SetupTest()
					}
					if beforeTestSuite, ok := suite.(BeforeTest); ok {
						beforeTestSuite.BeforeTest(methodFinder.Elem().Name(), method.Name)
					}

					stats.start(method.Name)

					method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
				})
			},
		}
		tests = append(tests, test)
	}

	if len(tests) == 0 {
		return
	}

	if stats != nil {
		stats.Start = time.Now()
	}

	defer func() {
		if suiteWithStats, measureStats := suite.(WithStats); measureStats {
			stats.End = time.Now()
			suiteWithStats.HandleStats(suiteName, stats)
		}
	}()

	tests.run(t)
}
