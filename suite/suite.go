package suite

import (
	"reflect"
	"regexp"
	"testing"
)

// Suite is a basic testing suite with methods for storing and
// retrieving the current *testing.T context.
type Suite struct {
	t *testing.T
}

// T retrieves the current *testing.T context.
func (suite *Suite) T() *testing.T {
	return suite.t
}

// SetT sets the current *testing.T context.
func (suite *Suite) SetT(t *testing.T) {
	suite.t = t
}

// Run takes a testing suite and runs all of the tests attached
// to it.
func Run(t *testing.T, suite TestingSuite) {
	suite.SetT(t)

	if setupAllSuite, ok := suite.(SetupAllSuite); ok {
		setupAllSuite.SetupSuite()
	}

	methodFinder := reflect.TypeOf(suite)
	tests := []testing.InternalTest{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		if ok, _ := regexp.MatchString("^Test", method.Name); ok {
			test := testing.InternalTest{
				Name: method.Name,
				F: func(t *testing.T) {
					parentT := suite.T()
					suite.SetT(t)
					if setupTestSuite, ok := suite.(SetupTestSuite); ok {
						setupTestSuite.SetupTest()
					}
					method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
					if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
						tearDownTestSuite.TearDownTest()
					}
					suite.SetT(parentT)
				},
			}
			tests = append(tests, test)
		}
	}

	if !testing.RunTests(func(_, _ string) (bool, error) {return true, nil},
		tests) {
		t.Fail()
	}

	if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
		tearDownAllSuite.TearDownSuite()
	}
}
