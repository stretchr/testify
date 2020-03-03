package suite

import "time"

type SuiteStats struct {
	testStats map[string]*TestStats
}

// Stats stores information about the execution of some test.
type TestStats struct {
	TestName           string
	StartTime, EndTime time.Time
	Passed             bool
}

func newSuiteStats() *SuiteStats {
	testStats := make(map[string]*TestStats)

	return &SuiteStats{
		testStats: testStats,
	}
}

func (s SuiteStats) start(testName string) {
	s.testStats[testName] = &TestStats{
		TestName:  testName,
		StartTime: time.Now(),
	}
}

func (s SuiteStats) end(testName string, passed bool) {
	s.testStats[testName].EndTime = time.Now()
	s.testStats[testName].Passed = passed
}

func (s SuiteStats) get(testName string) *TestStats {
	if stats, exists := s.testStats[testName]; exists {
		return stats
	}

	return nil
}
