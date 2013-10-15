// The suite package contains logic for creating testing suite structs
// and running the methods on those structs as tests.  The most useful
// piece of this package is that you can create setup/teardown methods
// on your testing suites, which will run before/after the whole suite
// or individual tests (depending on which interface(s) you
// implement).
//
// Once you've built your testing suite, you need to run the suite
// inside any function that matches the identity that "go test" is
// already looking for (i.e. func(*testing.T)).
//
// A crude example:
//     // Basic imports
//     import (
//         "testing"
//         "github.com/stretchr/testify/assert"
//         "github.com/stretchr/testify/suite"
//     )
//
//     // Define the suite, and absorb the built-in basic suite
//     // functionality from testify - including a T() method which
//     // returns the current testing context
//     type ExampleTestSuite struct {
//         suite.Suite
//         VariableThatShouldStartAtFive int
//     }
//
//     // Make sure that VariableThatShouldStartAtFive is set to five
//     // before each test
//     func (suite *ExampleTestSuite) SetupTest() {
//         suite.VariableThatShouldStartAtFive = 5
//     }
//
//     // All methods that begin with "Test" are run as tests within a
//     // suite.
//     func (suite *ExampleTestSuite) TestExample() {
//         assert.Equal(suite.T(), suite.VariableThatShouldStartAtFive, 5)
//     }
//
//     // In order for 'go test' to run this suite, we need to create
//     // a normal test function and pass our suite to suite.Run
//     func TestExampleTestSuite(t *testing.T) {
//         suite.Run(t, new(ExampleTestSuite))
//     }
package suite
