package mock_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
)

type expectationsResultT struct {
	errors []string
}

func (r *expectationsResultT) Logf(string, ...interface{}) {}
func (r *expectationsResultT) Errorf(format string, args ...interface{}) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}
func (r *expectationsResultT) FailNow() { panic("unexpected FailNow") }

type expectationsWitness struct {
	checked bool
}

func (w *expectationsWitness) AssertExpectations(mock.TestingT) bool {
	w.checked = true
	return true
}

func TestAssertExpectationsForObjectsResult(t *testing.T) {
	for _, invalidFirst := range []bool{true, false} {
		t.Run(fmt.Sprintf("mixed objects invalidFirst=%t", invalidFirst), func(t *testing.T) {
			recorder := &expectationsResultT{}
			witness := &expectationsWitness{}
			objects := []interface{}{witness, mock.Mock{}}
			if invalidFirst {
				objects[0], objects[1] = objects[1], objects[0]
			}
			if mock.AssertExpectationsForObjects(recorder, objects...) {
				t.Error("returned true after reporting an invalid mock reference")
			}
			if !witness.checked || len(recorder.errors) != 1 {
				t.Fatalf("valid object checked=%t, errors=%v", witness.checked, recorder.errors)
			}
		})
	}
	t.Run("invalid mock reference", func(t *testing.T) {
		recorder := &expectationsResultT{}
		// This is the same by-value usage diagnosed by the existing regression
		// test for issue 1227. Reporting a failure must not return success.
		if mock.AssertExpectationsForObjects(recorder, mock.Mock{}) {
			t.Error("returned true after reporting an invalid mock reference")
		}
		if len(recorder.errors) != 1 || !strings.Contains(recorder.errors[0], "Invalid test object type mock.Mock") {
			t.Fatalf("unexpected errors: %v", recorder.errors)
		}
	})
	t.Run("valid mock", func(t *testing.T) {
		recorder := &expectationsResultT{}
		m := &mock.Mock{}
		m.On("Do").Return().Once()
		m.MethodCalled("Do")
		if !mock.AssertExpectationsForObjects(recorder, m) || len(recorder.errors) != 0 {
			t.Fatalf("valid expectations failed: %v", recorder.errors)
		}
	})
	t.Run("unmet expectation", func(t *testing.T) {
		recorder := &expectationsResultT{}
		m := &mock.Mock{}
		m.On("Do").Return().Once()
		if mock.AssertExpectationsForObjects(recorder, m) || len(recorder.errors) != 1 {
			t.Fatalf("unexpected assertion result: %v", recorder.errors)
		}
	})
	t.Run("empty object list", func(t *testing.T) {
		recorder := &expectationsResultT{}
		if !mock.AssertExpectationsForObjects(recorder) || len(recorder.errors) != 0 {
			t.Fatalf("empty object list failed: %v", recorder.errors)
		}
	})
}
