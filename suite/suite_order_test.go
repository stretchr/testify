package suite

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type CallOrderSuite struct {
	Suite
	callOrder []string
}

func (s *CallOrderSuite) call(method string) {
	time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
	s.callOrder = append(s.callOrder, method)
}

func TestSuiteCallOrder(t *testing.T) {
	Run(t, new(CallOrderSuite))
}
func (s *CallOrderSuite) SetupSuite() {
	s.call("SetupSuite")
}

func (s *CallOrderSuite) TearDownSuite() {
	s.call("TearDownSuite")
	assert.Equal(s.T(), "SetupSuite;SetupTest;Test A;TearDownTest;SetupTest;Test B;TearDownTest;TearDownSuite", strings.Join(s.callOrder, ";"))
}
func (s *CallOrderSuite) SetupTest() {
	s.call("SetupTest")
}

func (s *CallOrderSuite) TearDownTest() {
	s.call("TearDownTest")
}

func (s *CallOrderSuite) Test_A() {
	s.call("Test A")
}

func (s *CallOrderSuite) Test_B() {
	s.call("Test B")
}
