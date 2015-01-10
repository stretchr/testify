package suite

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testPattern   = "^Test"
	defaultFormat = "{method}"
)

var (
	matchMethod = flag.String("m", "", "Run only those test methods matching the regular expression.")
	format      = flag.String("nameformat", "",
		"The format to use for test name output.  Use {function} for the test function name, "+
			"{suite} for the suite type name, and {method} for the test method name.  Overrides "+
			"formatting set by NamedRun.")
)

type testIdentifier struct {
	function string
	suite    string
	method   string
}

func (t testIdentifier) Name(format string) string {
	name := strings.Replace(format, "{function}", t.function, -1)
	name = strings.Replace(name, "{suite}", t.suite, -1)
	name = strings.Replace(name, "{method}", t.method, -1)
	return name
}

// Suite is a basic testing suite with methods for storing and
// retrieving the current *testing.T context.
type Suite struct {
	*assert.Assertions
	t *testing.T
}

// T retrieves the current *testing.T context.
func (suite *Suite) T() *testing.T {
	return suite.t
}

// SetT sets the current *testing.T context.
func (suite *Suite) SetT(t *testing.T) {
	suite.t = t
	suite.Assertions = assert.New(t)
}

// Run takes a testing suite and runs all of the tests attached
// to it.  For legacy reasons, Run uses a default test name format
// of "{method}", which only prints out the test method name.  See
// NamedRun for a version of Run that accepts other format strings
// for the test name in testing output.
func Run(t *testing.T, suite TestingSuite) {
	run("", t, suite)
}

// NamedRun performs a test run the same as Run, but with a different
// test name format.  Pass a pattern to use as your test names for
// nameFormat.  The following patterns will be replaced by details
// about the current test:
//
// * "{function}": The name of the function that called NamedRun.
// * "{suite}": The name of the testing suite's underlying type.
// * "{method}": The name of the test method.
//
// Be aware that the test flag "nameformat" will override nameFormat
// globally.
func NamedRun(nameFormat string, t *testing.T, suite TestingSuite) {
	run(nameFormat, t, suite)
}

// run is used to ensure that Run and NamedRun (and any other suite
// running functions) are at the same calling stack depth.
func run(nameFormat string, t *testing.T, suite TestingSuite) {
	if nameFormat == "" {
		nameFormat = defaultFormat
	}
	if *format != "" {
		nameFormat = *format
	}
	id := testIdentifier{}
	if callPC, _, _, ok := runtime.Caller(2); ok {
		id.function = runtime.FuncForPC(callPC).Name()
		id.function = id.function[strings.LastIndex(id.function, ".")+1:]
	}

	suite.SetT(t)

	if setupAllSuite, ok := suite.(SetupAllSuite); ok {
		setupAllSuite.SetupSuite()
	}
	defer func() {
		if tearDownAllSuite, ok := suite.(TearDownAllSuite); ok {
			tearDownAllSuite.TearDownSuite()
		}
	}()

	suiteType := reflect.TypeOf(suite)
	methodFinder := suiteType
	if suiteType.Kind() == reflect.Ptr {
		suiteType = suiteType.Elem()
	}
	id.suite = suiteType.Name()

	tests := []testing.InternalTest{}
	for index := 0; index < methodFinder.NumMethod(); index++ {
		method := methodFinder.Method(index)
		ok, err := methodFilter(method.Name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testify: invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}
		if ok {
			id.method = method.Name
			test := testing.InternalTest{
				Name: id.Name(nameFormat),
				F: func(t *testing.T) {
					parentT := suite.T()
					suite.SetT(t)
					if setupTestSuite, ok := suite.(SetupTestSuite); ok {
						setupTestSuite.SetupTest()
					}
					defer func() {
						if tearDownTestSuite, ok := suite.(TearDownTestSuite); ok {
							tearDownTestSuite.TearDownTest()
						}
						suite.SetT(parentT)
					}()
					method.Func.Call([]reflect.Value{reflect.ValueOf(suite)})
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

// methodFilter filters method names similar to `go test -run regexp`.  Note
// that methods must match testPattern prior to matching the passed in regexp.
func methodFilter(name string) (bool, error) {
	if ok, _ := regexp.MatchString(testPattern, name); !ok {
		return false, nil
	}
	return regexp.MatchString(*matchMethod, name)
}
