package suite

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CallOrderSuite struct {
	Suite
	callOrder []string
}

func (s *CallOrderSuite) call(method string) {
	//	s.Mutex.Lock()
	// defer s.Mutex.Unlock()

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
	assert.Equal(s.T(), "SetupSuite;SetupTest;Test A;TearDownTest;TearDownSuite", strings.Join(s.callOrder, ";"))
}
func (s *CallOrderSuite) SetupTest() {
	s.T().Parallel()
	s.call("SetupTest")
}

func (s *CallOrderSuite) TearDownTest() {
	s.call("TearDownTest")
}

func (s *CallOrderSuite) Test_A() {
	s.call("Test A")
}

//func (s *CallOrderSuite) Test_B() {
//	time.Sleep(time.Second)
//	s.call("Test B")
//}
