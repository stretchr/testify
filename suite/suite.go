package suite

import (
	"testing"
	"reflect"
	"regexp"
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

	if beforeAllSuite, ok := suite.(BeforeAllSuite); ok {
		beforeAllSuite.BeforeSuite()
	}

	if afterAllSuite, ok := suite.(AfterAllSuite); ok {
		defer afterAllSuite.AfterSuite()
	}

	methodFinder := reflect.TypeOf(suite)
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		if ok, _ := regexp.MatchString("^Test", method.Name); ok {
			if beforeTestSuite, ok := suite.(BeforeTestSuite); ok {
				beforeTestSuite.BeforeTest()
			}
			method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
			if afterTestSuite, ok := suite.(AfterTestSuite); ok {
				afterTestSuite.AfterTest()
			}
		}
	}
}
