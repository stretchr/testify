package assert_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type durationRecorder struct {
	errors []string
	fatal  bool
}

func (r *durationRecorder) Errorf(format string, args ...interface{}) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}
func (r *durationRecorder) FailNow() { r.fatal = true }

func TestWithinDurationFullRange(t *testing.T) {
	const max = time.Duration(1<<63 - 1)
	const min = time.Duration(-1 << 63)
	base := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name   string
		actual time.Time
		delta  time.Duration
		want   bool
	}{
		{"equal_zero", base, 0, true},
		{"after_zero", base.Add(1), 0, false},
		{"before_zero", base.Add(-1), 0, false},
		{"positive_boundary", base.Add(time.Second), time.Second, true},
		{"negative_boundary", base.Add(-time.Second), time.Second, true},
		{"positive_outside", base.Add(time.Second + 1), time.Second, false},
		{"negative_outside", base.Add(-time.Second - 1), time.Second, false},
		{"max_positive_boundary", base.Add(max), max, true},
		{"max_negative_boundary", base.Add(-max), max, true},
		{"max_positive_outside", base.Add(max).Add(1), max, false},
		{"max_negative_outside", base.Add(-max).Add(-1), max, false},
		{"centuries_after", base.AddDate(400, 0, 0), max, false},
		{"centuries_before", base.AddDate(-400, 0, 0), max, false},
		{"negative_delta", base, -1, false},
		{"min_delta_equal", base, min, false},
		{"min_delta_before", base.Add(-max).Add(-1), min, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, reverse := range []bool{false, true} {
				expected, actual := base, tc.actual
				if reverse {
					expected, actual = actual, expected
				}
				r := new(durationRecorder)
				got := assert.WithinDuration(r, expected, actual, tc.delta, "duration marker %s", "kept")
				if got != tc.want || (len(r.errors) == 0) != tc.want {
					t.Errorf("reverse=%v got=%v errors=%v want=%v", reverse, got, r.errors, tc.want)
				}
				if !tc.want && (len(r.errors) != 1 || !strings.Contains(r.errors[0], "duration marker kept")) {
					t.Errorf("lost failure message: %v", r.errors)
				}
			}
		})
	}
}

func TestWithinDurationClockAndLocation(t *testing.T) {
	monotonic := time.Now()
	cases := []struct {
		name             string
		expected, actual time.Time
		delta            time.Duration
		want             bool
	}{
		{"monotonic_equal", monotonic, monotonic, 0, true},
		{"monotonic_boundary", monotonic, monotonic.Add(time.Second), time.Second, true},
		{"monotonic_outside", monotonic, monotonic.Add(time.Second + 1), time.Second, false},
		{"mixed_clocks", monotonic, monotonic.Round(0), 0, true},
		{"different_locations", monotonic.UTC(), monotonic.In(time.FixedZone("offset", 3600)), 0, true},
		{"zero_time", time.Time{}, time.Time{}, 0, true},
		{"zero_to_current", time.Time{}, monotonic, time.Duration(1<<63 - 1), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := new(durationRecorder)
			if got := assert.WithinDuration(r, tc.expected, tc.actual, tc.delta); got != tc.want {
				t.Errorf("got=%v want=%v errors=%v", got, tc.want, r.errors)
			}
		})
	}
}

func TestWithinDurationPublicWrappers(t *testing.T) {
	const max = time.Duration(1<<63 - 1)
	base := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	wrappers := []struct {
		name  string
		run   func(*durationRecorder, time.Time, time.Time, time.Duration)
		fatal bool
	}{
		{"assert_f", func(r *durationRecorder, a, b time.Time, d time.Duration) {
			assert.WithinDurationf(r, a, b, d, "wrapper marker %d", 7)
		}, false},
		{"assert_method", func(r *durationRecorder, a, b time.Time, d time.Duration) {
			assert.New(r).WithinDuration(a, b, d, "wrapper marker %d", 7)
		}, false},
		{"assert_method_f", func(r *durationRecorder, a, b time.Time, d time.Duration) {
			assert.New(r).WithinDurationf(a, b, d, "wrapper marker %d", 7)
		}, false},
		{"require", func(r *durationRecorder, a, b time.Time, d time.Duration) {
			require.WithinDuration(r, a, b, d, "wrapper marker %d", 7)
		}, true},
		{"require_f", func(r *durationRecorder, a, b time.Time, d time.Duration) {
			require.WithinDurationf(r, a, b, d, "wrapper marker %d", 7)
		}, true},
		{"require_method", func(r *durationRecorder, a, b time.Time, d time.Duration) {
			require.New(r).WithinDuration(a, b, d, "wrapper marker %d", 7)
		}, true},
		{"require_method_f", func(r *durationRecorder, a, b time.Time, d time.Duration) {
			require.New(r).WithinDurationf(a, b, d, "wrapper marker %d", 7)
		}, true},
	}
	for _, w := range wrappers {
		t.Run(w.name, func(t *testing.T) {
			r := new(durationRecorder)
			w.run(r, base, base.Add(max), max)
			if len(r.errors) != 0 || r.fatal {
				t.Errorf("boundary failed: %+v", r)
			}
			r = new(durationRecorder)
			w.run(r, base, base.Add(max).Add(1), max)
			if len(r.errors) != 1 || r.fatal != w.fatal || !strings.Contains(r.errors[0], "wrapper marker 7") {
				t.Errorf("outside incorrectly accepted: %+v", r)
			}
		})
	}
}
