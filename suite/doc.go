// Package suite contains logic for creating testing suite structs
// and running the methods on those structs as tests.  The most useful
// piece of this package is that you can create setup/teardown methods
// on your testing suites, which will run before/after the whole suite
// or individual tests (depending on which interface(s) you
// implement).
//
// A testing suite is usually built by defining a Suite struct
// that includes all fields that tests need.
//
// After that, you can implement any of the interfaces in
// suite/interfaces.go to add setup/teardown functionality to your
// suite, and add any methods that start with "Test" to add tests.
// Test methods must match signature: `func(*suite.T)`. The suite.T
// object passed may be used to run sub-tests, verify assertions
// and control test execution.
// Methods that do not match any suite interfaces and do not begin
// with "Test" will not be run by testify, and can safely be used as
// helper methods.
//
// Once you've built your testing suite, you need to run the suite
// (using suite.Run from testify) inside any function that matches the
// identity that "go test" is already looking for (i.e.
// func(*testing.T)).
//
// Regular expression to select test suites specified command-line
// argument "-run". Regular expression to select the methods
// of test suites specified command-line argument "-m".
// Suite object has assertion methods.
//
// A crude example:
//     // Basic imports
//     import (
//         "testing"
//         "github.com/stretchr/testify/assert"
//         "github.com/stretchr/testify/suite"
//     )
//
//     // Define the suite, which is simply a struct with all
//     // fields that tests need.
//     type ExampleTestSuite struct {
//         VariableThatShouldStartAtFive int
//     }
//
//     // Make sure that VariableThatShouldStartAtFive is set to five
//     // before each test
//     func (suite *ExampleTestSuite) SetupTest(t *suite.T) {
//         suite.VariableThatShouldStartAtFive = 5
//     }
//
//     // All methods that begin with "Test" are run as tests within a
//     // suite.
//     func (suite *ExampleTestSuite) TestExample(t *suite.T) {
//         assert.Equal(t, 5, suite.VariableThatShouldStartAtFive)
//         t.Equal(5, suite.VariableThatShouldStartAtFive)
//     }
//
//     // In order for 'go test' to run this suite, we need to create
//     // a normal test function and pass our suite to suite.Run
//     func TestExampleTestSuite(t *testing.T) {
//         suite.Run(t, new(ExampleTestSuite))
//     }
package suite
