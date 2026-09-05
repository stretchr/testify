package suite

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SuiteRequire is a basic testing suite with methods for storing and
// retrieving the current *testing.T context.
// the difference with Suite is SuiteRequire uses require as default not assert.
type SuiteRequire struct {
	*require.Assertions

	mu     sync.RWMutex
	assert *assert.Assertions
	t      *testing.T

	// Parent suite to have access to the implemented methods of parent struct
	s TestingSuite
}

// T retrieves the current *testing.T context.
func (suite *SuiteRequire) T() *testing.T {
	suite.mu.RLock()
	defer suite.mu.RUnlock()
	return suite.t
}

// SetT sets the current *testing.T context.
func (suite *SuiteRequire) SetT(t *testing.T) {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	suite.t = t
	suite.Assertions = require.New(t)
	suite.assert = assert.New(t)
}

// SetS needs to set the current test suite as parent
// to get access to the parent methods
func (suite *SuiteRequire) SetS(s TestingSuite) {
	suite.s = s
}

// Require returns a require context for suite.
func (suite *SuiteRequire) Require() *require.Assertions {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	if suite.Assertions == nil {
		suite.Assertions = require.New(suite.T())
	}
	return suite.Assertions
}

// Assert returns an assert context for suite.  Normally, you can call
// `suite.NoError(expected, actual)`, but for situations where the embedded
// methods are overridden (for example, you might want to override
// assert.Assertions with require.Assertions), this method is provided so you
// can call `suite.Assert().NoError()`.
func (suite *SuiteRequire) Assert() *assert.Assertions {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	if suite.assert == nil {
		suite.assert = assert.New(suite.T())
	}
	return suite.assert
}

// Run provides suite functionality around golang subtests.  It should be
// called in place of t.Run(name, func(t *testing.T)) in test suite code.
// The passed-in func will be executed as a subtest with a fresh instance of t.
// Provides compatibility with go test pkg -run TestSuite/TestName/SubTestName.
func (suite *SuiteRequire) Run(name string, subtest func()) bool {
	oldT := suite.T()

	if setupSubTest, ok := suite.s.(SetupSubTest); ok {
		setupSubTest.SetupSubTest()
	}

	defer func() {
		suite.SetT(oldT)
		if tearDownSubTest, ok := suite.s.(TearDownSubTest); ok {
			tearDownSubTest.TearDownSubTest()
		}
	}()

	return oldT.Run(name, func(t *testing.T) {
		suite.SetT(t)
		subtest()
	})
}
