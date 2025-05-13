package suite_test

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type parallelSuiteData struct {
	calls          []string
	parallelSuiteT map[string]*testing.T
}

type parallelSuite struct {
	suite.Suite
	mutex sync.Mutex
	data  *parallelSuiteData
}

func (s *parallelSuite) recordCall(method string) {
	time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data.calls = append(s.data.calls, method)
}

func TestSuiteParallel(t *testing.T) {
	data := parallelSuiteData{
		calls:          []string{},
		parallelSuiteT: map[string]*testing.T{},
	}
	s := &parallelSuite{data: &data}
	suite.Run(t, s)
}

func (s *parallelSuite) SetupSuite() {
	s.recordCall("SetupSuite")
}

func (s *parallelSuite) TearDownSuite() {
	t := s.T()

	s.recordCall("TearDownSuite")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	testACalls := []string{}
	testBCalls := []string{}
	for _, call := range s.data.calls {
		if strings.Contains(call, "Test_A") {
			testACalls = append(testACalls, call)
		} else if strings.Contains(call, "Test_B") {
			testBCalls = append(testBCalls, call)
		}
	}
	assert.Equal(
		t,
		[]string{
			fmt.Sprintf("BeforeTest %s/parallel/ParallelTest_A ParallelTest_A", t.Name()),
			"Test_A",
			fmt.Sprintf("AfterTest %s/parallel/ParallelTest_A ParallelTest_A", t.Name()),
		},
		testACalls,
	)
	assert.Equal(
		t,
		[]string{
			fmt.Sprintf("BeforeTest %s/parallel/ParallelTest_B ParallelTest_B", t.Name()),
			"Test_B",
			fmt.Sprintf("AfterTest %s/parallel/ParallelTest_B ParallelTest_B", t.Name()),
		},
		testBCalls,
	)

	require.NotEmpty(t, s.data.calls)

	assert.Equal(t, "SetupSuite", s.data.calls[0])
	assert.Equal(t, "TearDownSuite", s.data.calls[len(s.data.calls)-1])

	assert.NotEqual(t, s, s.data.parallelSuiteT["Test_A"])
	assert.NotEqual(t, s, s.data.parallelSuiteT["Test_B"])
	assert.NotEqual(t, s.data.parallelSuiteT["Test_A"], s.data.parallelSuiteT["Test_B"])
}

func (s *parallelSuite) BeforeTest(t *testing.T, _, testName string) {
	s.recordCall(fmt.Sprintf("BeforeTest %s %s", t.Name(), testName))
}

func (s *parallelSuite) AfterTest(t *testing.T, _, testName string) {
	s.recordCall(fmt.Sprintf("AfterTest %s %s", t.Name(), testName))
}

func (s *parallelSuite) ParallelTest_A(t *testing.T) {
	s.recordCall("Test_A")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data.parallelSuiteT["Test_A"] = t
}

func (s *parallelSuite) ParallelTest_B(t *testing.T) {
	s.recordCall("Test_B")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data.parallelSuiteT["Test_B"] = t
}
