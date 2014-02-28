package suite

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"regexp"
	"testing"
)

// Suite is a basic testing suite with methods for storing and
// retrieving the current *testing.T context.
type Suite struct {
	*assert.Assertions
}

// SetT sets the current *testing.T context.
func (suite *Suite) SetT(t *testing.T) {
	suite.Assertions = assert.New(t)
}

// Run takes a testing suite and runs all of the tests attached
// to it.
func Run(t *testing.T, suite TestingSuite) {
	suite.SetT(t)

	if setupAllSuite, ok := suite.(SetupAllSuite); ok {
		setupAllSuite.SetupSuite()
	}

	methodFinder := reflect.TypeOf(suite)
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		if ok, _ := regexp.MatchString("^Test", method.Name); ok {
			if setupTestSuite, ok := suite.(SetupTestSuite); ok {
				setupTestSuite.SetupTest()
			}
			method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
			if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
				tearDownTestSuite.TearDownTest()
			}
		}
	}

	if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
		tearDownAllSuite.TearDownSuite()
	}
}
