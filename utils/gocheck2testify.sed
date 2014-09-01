# SED script for converting GoLang unit test code from GoCheck style to Testify style (for whatever reason)
# Version 0.1 8/31/2014 by polyglot-jones
# WARNING: Only some of these conversions have been tested!

# ######################################################################################
# IMPORTANT: This script assumes that all of your tests are contained WITHIN TEST SUITES
# ######################################################################################

# Overhead code
s|\.\s*"gopkg.in/check.v1"|"github.com/stretchr/testify/assert"\n\t"github.com/stretchr/testify/suite"\n\t"testing"|
s/type\s*(\w*)Suite\s*struct\s*\{/type \1Suite struct \{\n\tsuite.Suite/
s|var\s*_\s*=\s*Suite\(&(\w*)Suite\{\s*\}\s*\)|// The one testify function that launches our test suite\nfunc Test\1Suite\(t \*testing.T\) \{\n\tsuite.Run\(t, new\(\1Suite\)\)\n\}|

# The test methods (incl. setup & teardoen)
s/Suite\)\s*SetUpSuite\(c \*C\)/Suite\) SetupAllSuite\(\)/
s/Suite\)\s*SetUpTest\(c \*C\)/Suite\) SetupTestSuite\(\)/
s/Suite\)\s*TearDownTest\(c \*C\)/Suite\) TearDownTestSuite\(\)/
s/Suite\)\s*TearDownSuite\(c \*C\)/Suite\) TearDownAllSuite\(\)/
s/Suite\)\s*Test([^\(]*)\(c \*C\)/Suite\) Test\1\(\)/

# The assertions
s/c\.(Fail|FailNow|Fatal|Fatalf|Log|Logf|Error|Errorf|Skip)\(([^\)]*)/suite.T\(\).\1\(\2/
s/c\.Fatalf\(([^\)]*)/suite.T\(\).Fatalf\(\1/
s/c\.Assert\(([^,]*), ErrorMatches/assert.EqualError\(suite.T\(\), \1/
s/c\.Assert\(([^,]*), Matches/assert.True\(suite.T\(\), strings\.Matches\(\1\)/
s/c\.Assert\(([^,]*), Equals/assert.Equal\(suite.T\(\), \1/
s/c\.Assert\(([^,]*), Not\(Equals\)/assert.NotEqual\(suite.T\(\), \1/
s/c\.Assert\(([^,]*), IsNil/assert.Nil\(suite.T\(\), \1/
s/c\.Assert\(([^,]*), NotNil/assert.NotNil\(suite.T\(\), \1/
s/c\.Assert\(([^,]*), FitsTypeOf/assert.IsType\(suite.T\(\), \1/
s/c\.Assert\(([^,]*), HasLen/assert.Len\(suite.T\(\), \1/
s/c\.Assert\(([^,]*), Not\(Panics\)/assert.NotPanics \(suite.T\(\), \1/

# Catch-all for Checker names that are either the same in testify (Implements, Panics) or there is no direct equivalent(DeepEquals, PanicMatches, etc.) in which case manual intervention will be needed.
s/c\.Assert\(([^,]*), (\w*)/assert.\2\(suite.T\(\), \1/

# This one must follow all of the assertion conversions
s/(assert\..*)Commentf\((.*)\)\)/\1\2\)/

# Uncomment one of these two lines if your suite instance variables are named "s" or "this" rather than "suite"
# s/\bsuite\.T\b/s\.T/
# s/\bsuite\.T\b/this\.T/