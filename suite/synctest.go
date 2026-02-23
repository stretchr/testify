//go:build go1.25

package suite

import (
	"testing"
	"testing/synctest"
)

// SyncTest executes f in a new [synctest] bubble.
func (suite *Suite) SyncTest(f func()) {
	oldT := suite.T()
	synctest.Test(oldT, func(t *testing.T) {
		suite.SetT(t)
		defer suite.SetT(oldT)
		f()
	})
}
