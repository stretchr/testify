package suite

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

var (
	matchMethod   = flag.String("testify.m", "", "regular expression to select tests of the testify suite to run")
	matchMethodRE *regexp.Regexp
)

func isTestMethod(method reflect.Method) bool {
	if !strings.HasPrefix(method.Name, "Test") {
		return false
	}

	// compile once if needed
	if *matchMethod != "" && matchMethodRE == nil {
		var err error
		matchMethodRE, err = regexp.Compile(*matchMethod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testify: invalid regexp for -m: %s\n", err)
			os.Exit(1)
		}
	}

	// Apply -testify.m filter
	if matchMethodRE != nil && !matchMethodRE.MatchString(method.Name) {
		return false
	}

	return true
}

func checkMethodSignature(method reflect.Method) (test, bool) {
	if method.Type.NumIn() > 1 || method.Type.NumOut() > 0 {
		return test{
			name: method.Name,
			run: func(t *testing.T) {
				t.Errorf(
					"testify: suite method %q has invalid signature: expected no input or output parameters, method has %d input parameters and %d output parameters",
					method.Name, method.Type.NumIn()-1, method.Type.NumOut(),
				)
			},
		}, false
	}

	return test{}, true
}
