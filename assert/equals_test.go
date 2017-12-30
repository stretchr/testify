package assert

import (
	"testing"
	"time"
)

type equalT struct {
	main  int
	trash int
}

func (v1 equalT) Equal(v2 equalT) bool {
	return v1.main == v2.main
}

type invalidEqualArgT struct {
	v int
}

func (v1 invalidEqualArgT) Equal(i equalT) bool {
	return true
}

type invalidEqualResultT struct {
	v int
}

func (v1 invalidEqualResultT) Equal(i equalT) (bool, bool) {
	return true, true
}

type allT struct {
	b    bool
	i    int
	i8   int8
	i16  int16
	i32  int32
	i64  int64
	u    uint
	u8   uint8
	u16  uint16
	u32  uint32
	u64  uint64
	f32  float32
	f64  float64
	c64  complex64
	c128 complex128
	s    string
	ba   [2]byte
	rs   []rune
	m    map[int]int
	f    func() int
}

func TestDeepEqual(t *testing.T) {
	v1 := equalT{1, 1}
	v2 := equalT{1, 2}
	if !deepEqual(v1, v2) {
		t.Error("deepEqual should return true: equalT{1, 1} equals equalT{1, 2}")
	}
	if !deepEqual(&v1, &v2) {
		t.Error("deepEqual should return true: &equalT{1, 1} equals &equalT{1, 2}")
	}
	t1 := time.Now()
	if !deepEqual(t1, t1.In(time.UTC)) {
		t.Error("deepEqual should return true: 't1' equals 't1.In(UTC)'")
	}
	if deepEqual(t1, time.Time{}) {
		t.Error("deepEqual should return false: 't1' does not equal zero time")
	}
	if deepEqual(invalidEqualArgT{1}, invalidEqualArgT{2}) {
		t.Error("deepEqual should return false: 'invalidEqualArgT{1}' does not equal 'invalidEqualArgT{2}'")
	}
	if deepEqual(invalidEqualResultT{1}, invalidEqualResultT{2}) {
		t.Error("deepEqual should return false: 'invalidEqualResultT{1}' does not equal 'invalidEqualResultT{2}'")
	}

	a1 := allT{
		true,
		-1, -2, -3, -4, -5,
		1, 2, 3, 4, 5,
		0.1, 0.2,
		complex(6, 3), complex(4, 4),
		"string",
		[2]byte{'b', 'n'},
		[]rune{'a', 'e', 'i'},
		map[int]int{31: 5},
		nil,
	}
	a2 := allT{
		true,
		-1, -2, -3, -4, -5,
		1, 2, 3, 4, 5,
		0.1, 0.2,
		complex(6, 3), complex(4, 4),
		"string",
		[2]byte{'b', 'n'},
		[]rune{'a', 'e', 'i'},
		map[int]int{31: 5},
		nil,
	}
	if !deepEqual(a1, a2) {
		t.Error("deepEqual should return true: a1 equals a2")
	}
	if deepEqual(a1, allT{}) {
		t.Error("deepEqual should return false: a1 does not equal allT{}")
	}

	ch1 := make(chan int)
	if !deepEqual(ch1, ch1) {
		t.Error("deepEqual should return true: ch1 equals ch1")
	}
	if deepEqual(ch1, make(chan int)) {
		t.Error("deepEqual should return false: ch1 does not equal make(chan int)")
	}

	var nch1 <-chan int
	var nch2 <-chan int
	if !deepEqual(nch1, nch2) {
		t.Error("deepEqual should return true: null channel is equal to null channel of same type)")
	}
}
