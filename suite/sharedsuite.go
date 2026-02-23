package suite

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sharedSuite struct {
	*assert.Assertions

	mu      sync.RWMutex
	require *require.Assertions
	t       *testing.T

	// Parent suite to have access to the implemented methods of parent struct
	s TestingSuite
}

// T retrieves the current *testing.T context.
func (suite *sharedSuite) T() *testing.T {
	suite.mu.RLock()
	defer suite.mu.RUnlock()
	return suite.t
}

// SetT sets the current *testing.T context.
func (suite *sharedSuite) SetT(t *testing.T) {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	suite.t = t
	suite.Assertions = assert.New(t)
	suite.require = require.New(t)
}

// SetS needs to set the current test suite as parent
// to get access to the parent methods
func (suite *sharedSuite) SetS(s TestingSuite) {
	suite.s = s
}

// Require returns a require context for suite.
func (suite *sharedSuite) Require() *require.Assertions {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	if suite.require == nil {
		panic("'Require' must not be called before 'Run' or 'SetT'")
	}
	return suite.require
}

// Assert returns an assert context for suite. Normally, you can call:
//
//	suite.NoError(err)
//
// But for situations where the embedded methods are overridden (for example,
// you might want to override assert.Assertions with require.Assertions), this
// method is provided so you can call:
//
//	suite.Assert().NoError(err)
func (suite *sharedSuite) Assert() *assert.Assertions {
	suite.mu.Lock()
	defer suite.mu.Unlock()
	if suite.Assertions == nil {
		panic("'Assert' must not be called before 'Run' or 'SetT'")
	}
	return suite.Assertions
}
