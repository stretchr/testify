package suite

import "time"

// SuiteInformation stats stores stats for the whole suite execution.
type SuiteInformation struct {
	Start, End time.Time
	Passed     bool
	TestStats  map[string]*TestInformation
}

// TestInformation stores information about the execution of each test.
type TestInformation struct {
	TestName   string
	Start, End time.Time
	Passed     bool
}

func newSuiteInformation() *SuiteInformation {
	testStats := make(map[string]*TestInformation)

	return &SuiteInformation{
		TestStats: testStats,
		Passed:    true,
	}
}

func (s SuiteInformation) start(testName string) {
	s.TestStats[testName] = &TestInformation{
		TestName: testName,
		Start:    time.Now(),
	}
}

func (s SuiteInformation) end(testName string, passed bool) {
	s.TestStats[testName].End = time.Now()
	s.TestStats[testName].Passed = passed
}
