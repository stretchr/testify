package suite

import "time"

// Suite stats stores stats for the whole suite execution.
type SuiteStats struct {
	StartTime, EndTime time.Time
	Passed             bool
	TestStats          map[string]*TestStats
}

// TestStats stores information about the execution of each test.
type TestStats struct {
	TestName           string
	StartTime, EndTime time.Time
	Passed             bool
}

func newSuiteStats() *SuiteStats {
	testStats := make(map[string]*TestStats)

	return &SuiteStats{
		TestStats: testStats,
		Passed:    true,
	}
}

func (s SuiteStats) start(testName string) {
	s.TestStats[testName] = &TestStats{
		TestName:  testName,
		StartTime: time.Now(),
	}
}

func (s SuiteStats) end(testName string, passed bool) {
	s.TestStats[testName].EndTime = time.Now()
	s.TestStats[testName].Passed = passed
}
