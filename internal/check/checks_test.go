package check_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/internal/check"
)

func Test_IsEmpty(t *testing.T) {

	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}

	assert.True(t, check.IsEmpty(""))
	assert.True(t, check.IsEmpty(nil))
	assert.True(t, check.IsEmpty([]string{}))
	assert.True(t, check.IsEmpty(0))
	assert.True(t, check.IsEmpty(int32(0)))
	assert.True(t, check.IsEmpty(int64(0)))
	assert.True(t, check.IsEmpty(false))
	assert.True(t, check.IsEmpty(map[string]string{}))
	assert.True(t, check.IsEmpty(new(time.Time)))
	assert.True(t, check.IsEmpty(time.Time{}))
	assert.True(t, check.IsEmpty(make(chan struct{})))
	assert.True(t, check.IsEmpty([1]int{}))
	assert.False(t, check.IsEmpty("something"))
	assert.False(t, check.IsEmpty(errors.New("something")))
	assert.False(t, check.IsEmpty([]string{"something"}))
	assert.False(t, check.IsEmpty(1))
	assert.False(t, check.IsEmpty(true))
	assert.False(t, check.IsEmpty(map[string]string{"Hello": "World"}))
	assert.False(t, check.IsEmpty(chWithValue))
	assert.False(t, check.IsEmpty([1]int{42}))
}

func TestObjectsAreEqual(t *testing.T) {
	cases := []struct {
		expected interface{}
		actual   interface{}
		result   bool
	}{
		// cases that are expected to be equal
		{"Hello World", "Hello World", true},
		{123, 123, true},
		{123.5, 123.5, true},
		{[]byte("Hello World"), []byte("Hello World"), true},
		{nil, nil, true},

		// cases that are expected not to be equal
		{map[int]int{5: 10}, map[int]int{10: 20}, false},
		{'x', "x", false},
		{"x", 'x', false},
		{0, 0.1, false},
		{0.1, 0, false},
		{time.Now, time.Now, false},
		{func() {}, func() {}, false},
		{uint32(10), int32(10), false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("ObjectsAreEqual(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := check.ObjectsAreEqual(c.expected, c.actual)

			if res != c.result {
				t.Errorf("ObjectsAreEqual(%#v, %#v) should return %#v", c.expected, c.actual, c.result)
			}

		})
	}

	// Cases where type differ but values are equal
	if !check.ObjectsAreEqualValues(uint32(10), int32(10)) {
		t.Error("ObjectsAreEqualValues should return true")
	}
	if check.ObjectsAreEqualValues(0, nil) {
		t.Fail()
	}
	if check.ObjectsAreEqualValues(nil, 0) {
		t.Fail()
	}

}

func TestDiffLists(t *testing.T) {
	tests := []struct {
		name   string
		listA  interface{}
		listB  interface{}
		extraA []interface{}
		extraB []interface{}
	}{
		{
			name:   "equal empty",
			listA:  []string{},
			listB:  []string{},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "equal same order",
			listA:  []string{"hello", "world"},
			listB:  []string{"hello", "world"},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "equal different order",
			listA:  []string{"hello", "world"},
			listB:  []string{"world", "hello"},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "extra A",
			listA:  []string{"hello", "hello", "world"},
			listB:  []string{"hello", "world"},
			extraA: []interface{}{"hello"},
			extraB: nil,
		},
		{
			name:   "extra A twice",
			listA:  []string{"hello", "hello", "hello", "world"},
			listB:  []string{"hello", "world"},
			extraA: []interface{}{"hello", "hello"},
			extraB: nil,
		},
		{
			name:   "extra B",
			listA:  []string{"hello", "world"},
			listB:  []string{"hello", "hello", "world"},
			extraA: nil,
			extraB: []interface{}{"hello"},
		},
		{
			name:   "extra B twice",
			listA:  []string{"hello", "world"},
			listB:  []string{"hello", "hello", "world", "hello"},
			extraA: nil,
			extraB: []interface{}{"hello", "hello"},
		},
		{
			name:   "integers 1",
			listA:  []int{1, 2, 3, 4, 5},
			listB:  []int{5, 4, 3, 2, 1},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "integers 2",
			listA:  []int{1, 2, 1, 2, 1},
			listB:  []int{2, 1, 2, 1, 2},
			extraA: []interface{}{1},
			extraB: []interface{}{2},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			actualExtraA, actualExtraB := check.DiffLists(test.listA, test.listB)
			assert.Equal(t, test.extraA, actualExtraA, "extra A does not match for listA=%v listB=%v",
				test.listA, test.listB)
			assert.Equal(t, test.extraB, actualExtraB, "extra B does not match for listA=%v listB=%v",
				test.listA, test.listB)
		})
	}
}
