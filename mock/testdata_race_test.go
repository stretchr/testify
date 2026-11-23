package mock

import (
	"sync"
	"testing"

	"github.com/stretchr/objx"
)

func TestMockTestDataConcurrentInitialization(t *testing.T) {
	const callers = 32
	for attempt := 0; attempt < 20; attempt++ {
		var mocked Mock
		maps := make([]objx.Map, callers)
		start := make(chan struct{})
		var ready sync.WaitGroup
		var done sync.WaitGroup
		ready.Add(callers)
		done.Add(callers)
		for i := 0; i < callers; i++ {
			go func(index int) {
				defer done.Done()
				ready.Done()
				<-start
				// Map writes remain the caller's responsibility. Only retrieve
				// the shared map concurrently in this regression test.
				maps[index] = mocked.TestData()
				if maps[index] == nil {
					t.Error("TestData returned nil")
				}
			}(i)
		}
		ready.Wait()
		close(start)
		done.Wait()
		// Writes and reads here are sequential, after every accessor returns.
		maps[0]["shared"] = attempt
		for i, values := range maps {
			if values["shared"] != attempt {
				t.Errorf("caller %d did not receive the shared map", i)
			}
		}
	}
}
