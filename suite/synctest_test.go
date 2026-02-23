//go:build go1.25

package suite

import (
	"testing"
	"testing/synctest"
)

type SyncTestSuite struct {
	Suite
}

func TestSyncTest(t *testing.T) {
	t.Setenv("GODEBUG", "asynctimerchan=0") // since our go.mod says `go 1.17`
	Run(t, new(SyncTestSuite))
}

func (s *SyncTestSuite) TestSyncTest() {
	s.SyncTest(func() {
		synctest.Wait()
	})
	s.Run("subtest", func() {
		s.SyncTest(func() {
			synctest.Wait()
		})
	})
}
