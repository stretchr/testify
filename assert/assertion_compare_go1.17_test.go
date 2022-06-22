//go:build go1.17
// +build go1.17

// TODO: once support for Go 1.16 is dropped, this file can be
//       merged/removed with assertion_compare_can_convert.go and
//       assertion_compare_legacy.go

package assert

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestCompare17(t *testing.T) {
	type customTime time.Time
	type customBytes []byte
	for _, currCase := range []struct {
		less    interface{}
		greater interface{}
		cType   string
	}{
		{less: time.Now(), greater: time.Now().Add(time.Hour), cType: "time.Time"},
		{less: customTime(time.Now()), greater: customTime(time.Now().Add(time.Hour)), cType: "time.Time"},
		{less: []byte{1, 1}, greater: []byte{1, 2}, cType: "[]byte"},
		{less: customBytes([]byte{1, 1}), greater: customBytes([]byte{1, 2}), cType: "[]byte"},
	} {
		resLess, isComparable := compare(currCase.less, currCase.greater, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object should be comparable for type " + currCase.cType)
		}

		if resLess != compareLess {
			t.Errorf("object less (%v) should be less than greater (%v) for type "+currCase.cType,
				currCase.less, currCase.greater)
		}

		resGreater, isComparable := compare(currCase.greater, currCase.less, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object are comparable for type " + currCase.cType)
		}

		if resGreater != compareGreater {
			t.Errorf("object greater should be greater than less for type " + currCase.cType)
		}

		resEqual, isComparable := compare(currCase.less, currCase.less, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object are comparable for type " + currCase.cType)
		}

		if resEqual != 0 {
			t.Errorf("objects should be equal for type " + currCase.cType)
		}
	}
}

func TestGreater17(t *testing.T) {
	mockT := new(testing.T)

	if !Greater(mockT, 2, 1) {
		t.Error("Greater should return true")
	}

	if Greater(mockT, 1, 1) {
		t.Error("Greater should return false")
	}

	if Greater(mockT, 1, 2) {
		t.Error("Greater should return false")
	}

	// Check error report
	for _, currCase := range []struct {
		less    interface{}
		greater interface{}
		msg     string
	}{
		{less: []byte{1, 1}, greater: []byte{1, 2}, msg: `"[1 1]" is not greater than "[1 2]"`},
		{less: time.Time{}, greater: time.Time{}.Add(time.Hour), msg: `"0001-01-01 00:00:00 +0000 UTC" is not greater than "0001-01-01 01:00:00 +0000 UTC"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		False(t, Greater(out, currCase.less, currCase.greater))
		Contains(t, out.buf.String(), currCase.msg)
		Contains(t, out.helpers, "github.com/stretchr/testify/assert.Greater")
	}
}

func TestGreaterOrEqual17(t *testing.T) {
	mockT := new(testing.T)

	if !GreaterOrEqual(mockT, 2, 1) {
		t.Error("GreaterOrEqual should return true")
	}

	if !GreaterOrEqual(mockT, 1, 1) {
		t.Error("GreaterOrEqual should return true")
	}

	if GreaterOrEqual(mockT, 1, 2) {
		t.Error("GreaterOrEqual should return false")
	}

	// Check error report
	for _, currCase := range []struct {
		less    interface{}
		greater interface{}
		msg     string
	}{
		{less: []byte{1, 1}, greater: []byte{1, 2}, msg: `"[1 1]" is not greater than or equal to "[1 2]"`},
		{less: time.Time{}, greater: time.Time{}.Add(time.Hour), msg: `"0001-01-01 00:00:00 +0000 UTC" is not greater than or equal to "0001-01-01 01:00:00 +0000 UTC"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		False(t, GreaterOrEqual(out, currCase.less, currCase.greater))
		Contains(t, out.buf.String(), currCase.msg)
		Contains(t, out.helpers, "github.com/stretchr/testify/assert.GreaterOrEqual")
	}
}

func TestLess17(t *testing.T) {
	mockT := new(testing.T)

	if !Less(mockT, 1, 2) {
		t.Error("Less should return true")
	}

	if Less(mockT, 1, 1) {
		t.Error("Less should return false")
	}

	if Less(mockT, 2, 1) {
		t.Error("Less should return false")
	}

	// Check error report
	for _, currCase := range []struct {
		less    interface{}
		greater interface{}
		msg     string
	}{
		{less: []byte{1, 1}, greater: []byte{1, 2}, msg: `"[1 2]" is not less than "[1 1]"`},
		{less: time.Time{}, greater: time.Time{}.Add(time.Hour), msg: `"0001-01-01 01:00:00 +0000 UTC" is not less than "0001-01-01 00:00:00 +0000 UTC"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		False(t, Less(out, currCase.greater, currCase.less))
		Contains(t, out.buf.String(), currCase.msg)
		Contains(t, out.helpers, "github.com/stretchr/testify/assert.Less")
	}
}

func TestLessOrEqual17(t *testing.T) {
	mockT := new(testing.T)

	if !LessOrEqual(mockT, 1, 2) {
		t.Error("LessOrEqual should return true")
	}

	if !LessOrEqual(mockT, 1, 1) {
		t.Error("LessOrEqual should return true")
	}

	if LessOrEqual(mockT, 2, 1) {
		t.Error("LessOrEqual should return false")
	}

	// Check error report
	for _, currCase := range []struct {
		less    interface{}
		greater interface{}
		msg     string
	}{
		{less: []byte{1, 1}, greater: []byte{1, 2}, msg: `"[1 2]" is not less than or equal to "[1 1]"`},
		{less: time.Time{}, greater: time.Time{}.Add(time.Hour), msg: `"0001-01-01 01:00:00 +0000 UTC" is not less than or equal to "0001-01-01 00:00:00 +0000 UTC"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		False(t, LessOrEqual(out, currCase.greater, currCase.less))
		Contains(t, out.buf.String(), currCase.msg)
		Contains(t, out.helpers, "github.com/stretchr/testify/assert.LessOrEqual")
	}
}
