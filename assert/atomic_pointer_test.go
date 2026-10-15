//go:build go1.21
// +build go1.21

package assert_test

import (
	"math"
	"runtime"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// TestAtomicPointerEqualityIssue1403 exercises the real exported assertions.
// Go 1.21 allows this file to use generics while the module retains go 1.17.
func TestAtomicPointerEqualityIssue1403(t *testing.T) {
	one, anotherOne, two := 1, 1, 2
	onePointer, anotherOnePointer := &one, &anotherOne
	var ordinaryNil *int
	var atomicNil *atomic.Pointer[int]
	var left, equalValue, unequalValue, sameReferent atomic.Pointer[int]
	var zeroLeft, zeroRight atomic.Pointer[int]
	var otherType atomic.Pointer[string]
	left.Store(&one)
	equalValue.Store(&anotherOne)
	unequalValue.Store(&two)
	sameReferent.Store(&one)

	// Prevent aliasing from turning the regression into an identity comparison.
	if &left == &equalValue || left.Load() == equalValue.Load() {
		t.Fatal("regression requires distinct wrappers and distinct referents")
	}
	if *left.Load() != 1 || *equalValue.Load() != 1 {
		t.Fatal("regression requires equal loaded values")
	}

	cases := []struct {
		name      string
		expected  interface{}
		actual    interface{}
		wantEqual bool
	}{
		{"distinct_atomic_equal_loaded_values", &left, &equalValue, true},
		{"atomic_same_wrapper", &left, &left, true},
		{"atomic_same_referent", &left, &sameReferent, true},
		{"atomic_unequal_loaded_values", &left, &unequalValue, false},
		{"atomic_zero_loaded_pointers", &zeroLeft, &zeroRight, true},
		{"typed_nil_atomic_wrappers", atomicNil, atomicNil, true},
		{"typed_nil_atomic_vs_non_nil_wrapper", atomicNil, &zeroLeft, false},
		{"different_atomic_pointer_types", &zeroLeft, &otherType, false},
		{"ordinary_equal_pointer_values", &one, &anotherOne, true},
		{"ordinary_equal_double_pointer_values", &onePointer, &anotherOnePointer, true},
		{"ordinary_unequal_pointer_values", &one, &two, false},
		{"ordinary_typed_nil_pointers", ordinaryNil, ordinaryNil, true},
		{"ordinary_typed_nil_vs_non_nil", ordinaryNil, &one, false},
		{"ordinary_different_pointer_types", ordinaryNil, (*string)(nil), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantEqual {
				assert.Equal(t, tc.expected, tc.actual)
			} else {
				assert.NotEqual(t, tc.expected, tc.actual)
			}
		})
	}
}

// These lookalikes must keep ordinary equality and never invoke user methods.
type Pointer[T any] struct{ Value *T }

func (*Pointer[T]) Load() *T { panic("user Load must not be called") }

type atomicPointerWrapper struct{ atomic.Pointer[int] }

func (*atomicPointerWrapper) Load() *int { panic("embedded Load must not be called") }

type definedAtomicPointer *atomic.Pointer[int]

func TestAtomicPointerEqualityBoundaries(t *testing.T) {
	one, anotherOne := 1, 1
	nan, anotherNaN := math.NaN(), math.NaN()
	var left, equal atomic.Pointer[int]
	left.Store(&one)
	equal.Store(&anotherOne)
	var nanLeft, nanSame, nanDistinct atomic.Pointer[float64]
	nanLeft.Store(&nan)
	nanSame.Store(&nan)
	nanDistinct.Store(&anotherNaN)
	type node struct {
		value int
		next  *node
	}
	first, second, different := &node{value: 1}, &node{value: 1}, &node{value: 2}
	first.next, second.next, different.next = first, second, different
	var cycleLeft, cycleEqual, cycleUnequal atomic.Pointer[node]
	cycleLeft.Store(first)
	cycleEqual.Store(second)
	cycleUnequal.Store(different)
	var valueLeft, valueRight atomic.Value
	valueLeft.Store(&one)
	valueRight.Store(&anotherOne)
	wrappedLeft, wrappedRight := &atomicPointerWrapper{}, &atomicPointerWrapper{}
	wrappedLeft.Store(&one)
	wrappedRight.Store(&anotherOne)
	cases := []struct {
		name        string
		left, right interface{}
		equal       bool
	}{
		{"same_wrapper_nan", &nanLeft, &nanLeft, true},
		{"same_loaded_pointer_nan", &nanLeft, &nanSame, true},
		{"distinct_loaded_nan", &nanLeft, &nanDistinct, false},
		{"equal_cycles", &cycleLeft, &cycleEqual, true},
		{"unequal_cycles", &cycleLeft, &cycleUnequal, false},
		{"user_load_is_not_called", &Pointer[int]{&one}, &Pointer[int]{&anotherOne}, true},
		{"embedded_atomic_not_unwrapped", wrappedLeft, wrappedRight, false},
		{"defined_pointer_not_unwrapped", definedAtomicPointer(&left), definedAtomicPointer(&equal), false},
		{"atomic_value_unchanged", &valueLeft, &valueRight, true},
		{"unsafe_pointers_not_dereferenced", unsafe.Pointer(&one), unsafe.Pointer(&anotherOne), false},
		{"nil_and_empty_bytes", []byte(nil), []byte{}, false},
		{"nil_interface_and_typed_nil", nil, (*atomic.Pointer[int])(nil), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.equal {
				assert.Equal(t, tc.left, tc.right)
			} else {
				assert.NotEqual(t, tc.left, tc.right)
			}
		})
	}
	assert.True(t, assert.ObjectsAreEqual(&left, &equal))
	assert.True(t, assert.ObjectsAreEqualValues(&left, &equal))
	assert.Exactly(t, &left, &equal)
}

func TestAtomicPointerEqualityConcurrentStores(t *testing.T) {
	one, anotherOne, equalOne := 1, 1, 1
	var left, right atomic.Pointer[int]
	left.Store(&one)
	right.Store(&equalOne)
	stop, done, started := make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		defer close(done)
		close(started)
		for {
			select {
			case <-stop:
				return
			default:
				left.Store(&one)
				left.Store(&anotherOne)
				runtime.Gosched()
			}
		}
	}()
	defer func() { close(stop); <-done }()
	<-started
	for i := 0; i < 1000; i++ {
		if !assert.ObjectsAreEqual(&left, &right) {
			t.Fatal("each atomically loaded pointer must refer to an equal, immutable value")
		}
		runtime.Gosched()
	}
}
