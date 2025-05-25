package suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mySuite struct {
	Suite
}

func (s *mySuite) SetupTest() {
	s.T().Skip("Just because!")
}
func (s *mySuite) HandleStats(_ string, _ *SuiteInformation) {}

func (s *mySuite) TestSomething() {
	panic("Should not get here.")
}

func TestSuiteWithStatsAndSkip(t *testing.T) {
	assert.NotPanics(
		t,
		func() {
			Run(t, &mySuite{})
		},
	)
}
