package assert

import (
	"errors"
	"io"
	"math"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"
)

var (
	i     interface{}
	zeros = []interface{}{
		false,
		byte(0),
		complex64(0),
		complex128(0),
		float32(0),
		float64(0),
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		rune(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		uintptr(0),
		"",
		[0]interface{}{},
		[]interface{}(nil),
		struct{ x int }{},
		(*interface{})(nil),
		(func())(nil),
		nil,
		interface{}(nil),
		map[interface{}]interface{}(nil),
		(chan interface{})(nil),
		(<-chan interface{})(nil),
		(chan<- interface{})(nil),
	}
	nonZeros = []interface{}{
		true,
		byte(1),
		complex64(1),
		complex128(1),
		float32(1),
		float64(1),
		int(1),
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		rune(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		uintptr(1),
		"s",
		[1]interface{}{1},
		[]interface{}{},
		struct{ x int }{1},
		(*interface{})(&i),
		(func())(func() {}),
		interface{}(1),
		map[interface{}]interface{}{},
		(chan interface{})(make(chan interface{})),
		(<-chan interface{})(make(chan interface{})),
		(chan<- interface{})(make(chan interface{})),
	}
)

// AssertionTesterInterface defines an interface to be used for testing assertion methods
type AssertionTesterInterface interface {
	TestMethod()
}

// AssertionTesterConformingObject is an object that conforms to the AssertionTesterInterface interface
type AssertionTesterConformingObject struct {
}

func (a *AssertionTesterConformingObject) TestMethod() {
}

// AssertionTesterNonConformingObject is an object that does not conform to the AssertionTesterInterface interface
type AssertionTesterNonConformingObject struct {
}

func TestObjectsAreEqual(t *testing.T) {

	if !ObjectsAreEqual("Hello World", "Hello World") {
		t.Error("objectsAreEqual should return true")
	}
	if !ObjectsAreEqual(123, 123) {
		t.Error("objectsAreEqual should return true")
	}
	if !ObjectsAreEqual(123.5, 123.5) {
		t.Error("objectsAreEqual should return true")
	}
	if !ObjectsAreEqual([]byte("Hello World"), []byte("Hello World")) {
		t.Error("objectsAreEqual should return true")
	}
	if !ObjectsAreEqual(nil, nil) {
		t.Error("objectsAreEqual should return true")
	}
	if ObjectsAreEqual(map[int]int{5: 10}, map[int]int{10: 20}) {
		t.Error("objectsAreEqual should return false")
	}
	if ObjectsAreEqual('x', "x") {
		t.Error("objectsAreEqual should return false")
	}
	if ObjectsAreEqual("x", 'x') {
		t.Error("objectsAreEqual should return false")
	}
	if ObjectsAreEqual(0, 0.1) {
		t.Error("objectsAreEqual should return false")
	}
	if ObjectsAreEqual(0.1, 0) {
		t.Error("objectsAreEqual should return false")
	}
	if ObjectsAreEqual(uint32(10), int32(10)) {
		t.Error("objectsAreEqual should return false")
	}
	if !ObjectsAreEqualValues(uint32(10), int32(10)) {
		t.Error("ObjectsAreEqualValues should return true")
	}
	if ObjectsAreEqualValues(0, nil) {
		t.Fail()
	}
	if ObjectsAreEqualValues(nil, 0) {
		t.Fail()
	}

}

func TestImplements(t *testing.T) {

	mockT := new(testing.T)

	if !Implements(mockT, (*AssertionTesterInterface)(nil), new(AssertionTesterConformingObject)) {
		t.Error("Implements method should return true: AssertionTesterConformingObject implements AssertionTesterInterface")
	}
	if Implements(mockT, (*AssertionTesterInterface)(nil), new(AssertionTesterNonConformingObject)) {
		t.Error("Implements method should return false: AssertionTesterNonConformingObject does not implements AssertionTesterInterface")
	}

}

func TestIsType(t *testing.T) {

	mockT := new(testing.T)

	if !IsType(mockT, new(AssertionTesterConformingObject), new(AssertionTesterConformingObject)) {
		t.Error("IsType should return true: AssertionTesterConformingObject is the same type as AssertionTesterConformingObject")
	}
	if IsType(mockT, new(AssertionTesterConformingObject), new(AssertionTesterNonConformingObject)) {
		t.Error("IsType should return false: AssertionTesterConformingObject is not the same type as AssertionTesterNonConformingObject")
	}

}

func TestEqual(t *testing.T) {

	mockT := new(testing.T)

	if !Equal(mockT, "Hello World", "Hello World") {
		t.Error("Equal should return true")
	}
	if !Equal(mockT, 123, 123) {
		t.Error("Equal should return true")
	}
	if !Equal(mockT, 123.5, 123.5) {
		t.Error("Equal should return true")
	}
	if !Equal(mockT, []byte("Hello World"), []byte("Hello World")) {
		t.Error("Equal should return true")
	}
	if !Equal(mockT, nil, nil) {
		t.Error("Equal should return true")
	}
	if !Equal(mockT, int32(123), int32(123)) {
		t.Error("Equal should return true")
	}
	if !Equal(mockT, uint64(123), uint64(123)) {
		t.Error("Equal should return true")
	}

}

func TestFormatUnequalValues(t *testing.T) {
	expected, actual := formatUnequalValues("foo", "bar")
	Equal(t, `"foo"`, expected, "value should not include type")
	Equal(t, `"bar"`, actual, "value should not include type")

	expected, actual = formatUnequalValues(123, 123)
	Equal(t, `123`, expected, "value should not include type")
	Equal(t, `123`, actual, "value should not include type")

	expected, actual = formatUnequalValues(int64(123), int32(123))
	Equal(t, `int64(123)`, expected, "value should include type")
	Equal(t, `int32(123)`, actual, "value should include type")

	type testStructType struct {
		Val string
	}

	expected, actual = formatUnequalValues(&testStructType{Val: "test"}, &testStructType{Val: "test"})
	Equal(t, `&assert.testStructType{Val:"test"}`, expected, "value should not include type annotation")
	Equal(t, `&assert.testStructType{Val:"test"}`, actual, "value should not include type annotation")
}

func TestNotNil(t *testing.T) {

	mockT := new(testing.T)

	if !NotNil(mockT, new(AssertionTesterConformingObject)) {
		t.Error("NotNil should return true: object is not nil")
	}
	if NotNil(mockT, nil) {
		t.Error("NotNil should return false: object is nil")
	}
	if NotNil(mockT, (*struct{})(nil)) {
		t.Error("NotNil should return false: object is (*struct{})(nil)")
	}

}

func TestNil(t *testing.T) {

	mockT := new(testing.T)

	if !Nil(mockT, nil) {
		t.Error("Nil should return true: object is nil")
	}
	if !Nil(mockT, (*struct{})(nil)) {
		t.Error("Nil should return true: object is (*struct{})(nil)")
	}
	if Nil(mockT, new(AssertionTesterConformingObject)) {
		t.Error("Nil should return false: object is not nil")
	}

}

func TestTrue(t *testing.T) {

	mockT := new(testing.T)

	if !True(mockT, true) {
		t.Error("True should return true")
	}
	if True(mockT, false) {
		t.Error("True should return false")
	}

}

func TestFalse(t *testing.T) {

	mockT := new(testing.T)

	if !False(mockT, false) {
		t.Error("False should return true")
	}
	if False(mockT, true) {
		t.Error("False should return false")
	}

}

func TestExactly(t *testing.T) {

	mockT := new(testing.T)

	a := float32(1)
	b := float64(1)
	c := float32(1)
	d := float32(2)

	if Exactly(mockT, a, b) {
		t.Error("Exactly should return false")
	}
	if Exactly(mockT, a, d) {
		t.Error("Exactly should return false")
	}
	if !Exactly(mockT, a, c) {
		t.Error("Exactly should return true")
	}

	if Exactly(mockT, nil, a) {
		t.Error("Exactly should return false")
	}
	if Exactly(mockT, a, nil) {
		t.Error("Exactly should return false")
	}

}

func TestNotEqual(t *testing.T) {

	mockT := new(testing.T)

	if !NotEqual(mockT, "Hello World", "Hello World!") {
		t.Error("NotEqual should return true")
	}
	if !NotEqual(mockT, 123, 1234) {
		t.Error("NotEqual should return true")
	}
	if !NotEqual(mockT, 123.5, 123.55) {
		t.Error("NotEqual should return true")
	}
	if !NotEqual(mockT, []byte("Hello World"), []byte("Hello World!")) {
		t.Error("NotEqual should return true")
	}
	if !NotEqual(mockT, nil, new(AssertionTesterConformingObject)) {
		t.Error("NotEqual should return true")
	}
	funcA := func() int { return 23 }
	funcB := func() int { return 42 }
	if !NotEqual(mockT, funcA, funcB) {
		t.Error("NotEqual should return true")
	}

	if NotEqual(mockT, "Hello World", "Hello World") {
		t.Error("NotEqual should return false")
	}
	if NotEqual(mockT, 123, 123) {
		t.Error("NotEqual should return false")
	}
	if NotEqual(mockT, 123.5, 123.5) {
		t.Error("NotEqual should return false")
	}
	if NotEqual(mockT, []byte("Hello World"), []byte("Hello World")) {
		t.Error("NotEqual should return false")
	}
	if NotEqual(mockT, new(AssertionTesterConformingObject), new(AssertionTesterConformingObject)) {
		t.Error("NotEqual should return false")
	}
}

type A struct {
	Name, Value string
}

func TestContains(t *testing.T) {

	mockT := new(testing.T)
	list := []string{"Foo", "Bar"}
	complexList := []*A{
		{"b", "c"},
		{"d", "e"},
		{"g", "h"},
		{"j", "k"},
	}
	simpleMap := map[interface{}]interface{}{"Foo": "Bar"}

	if !Contains(mockT, "Hello World", "Hello") {
		t.Error("Contains should return true: \"Hello World\" contains \"Hello\"")
	}
	if Contains(mockT, "Hello World", "Salut") {
		t.Error("Contains should return false: \"Hello World\" does not contain \"Salut\"")
	}

	if !Contains(mockT, list, "Bar") {
		t.Error("Contains should return true: \"[\"Foo\", \"Bar\"]\" contains \"Bar\"")
	}
	if Contains(mockT, list, "Salut") {
		t.Error("Contains should return false: \"[\"Foo\", \"Bar\"]\" does not contain \"Salut\"")
	}
	if !Contains(mockT, complexList, &A{"g", "h"}) {
		t.Error("Contains should return true: complexList contains {\"g\", \"h\"}")
	}
	if Contains(mockT, complexList, &A{"g", "e"}) {
		t.Error("Contains should return false: complexList contains {\"g\", \"e\"}")
	}
	if Contains(mockT, complexList, &A{"g", "e"}) {
		t.Error("Contains should return false: complexList contains {\"g\", \"e\"}")
	}
	if !Contains(mockT, simpleMap, "Foo") {
		t.Error("Contains should return true: \"{\"Foo\": \"Bar\"}\" contains \"Foo\"")
	}
	if Contains(mockT, simpleMap, "Bar") {
		t.Error("Contains should return false: \"{\"Foo\": \"Bar\"}\" does not contains \"Bar\"")
	}
}

func TestNotContains(t *testing.T) {

	mockT := new(testing.T)
	list := []string{"Foo", "Bar"}
	simpleMap := map[interface{}]interface{}{"Foo": "Bar"}

	if !NotContains(mockT, "Hello World", "Hello!") {
		t.Error("NotContains should return true: \"Hello World\" does not contain \"Hello!\"")
	}
	if NotContains(mockT, "Hello World", "Hello") {
		t.Error("NotContains should return false: \"Hello World\" contains \"Hello\"")
	}

	if !NotContains(mockT, list, "Foo!") {
		t.Error("NotContains should return true: \"[\"Foo\", \"Bar\"]\" does not contain \"Foo!\"")
	}
	if NotContains(mockT, list, "Foo") {
		t.Error("NotContains should return false: \"[\"Foo\", \"Bar\"]\" contains \"Foo\"")
	}
	if NotContains(mockT, simpleMap, "Foo") {
		t.Error("Contains should return true: \"{\"Foo\": \"Bar\"}\" contains \"Foo\"")
	}
	if !NotContains(mockT, simpleMap, "Bar") {
		t.Error("Contains should return false: \"{\"Foo\": \"Bar\"}\" does not contains \"Bar\"")
	}
}

func Test_includeElement(t *testing.T) {

	list1 := []string{"Foo", "Bar"}
	list2 := []int{1, 2}
	simpleMap := map[interface{}]interface{}{"Foo": "Bar"}

	ok, found := includeElement("Hello World", "World")
	True(t, ok)
	True(t, found)

	ok, found = includeElement(list1, "Foo")
	True(t, ok)
	True(t, found)

	ok, found = includeElement(list1, "Bar")
	True(t, ok)
	True(t, found)

	ok, found = includeElement(list2, 1)
	True(t, ok)
	True(t, found)

	ok, found = includeElement(list2, 2)
	True(t, ok)
	True(t, found)

	ok, found = includeElement(list1, "Foo!")
	True(t, ok)
	False(t, found)

	ok, found = includeElement(list2, 3)
	True(t, ok)
	False(t, found)

	ok, found = includeElement(list2, "1")
	True(t, ok)
	False(t, found)

	ok, found = includeElement(simpleMap, "Foo")
	True(t, ok)
	True(t, found)

	ok, found = includeElement(simpleMap, "Bar")
	True(t, ok)
	False(t, found)

	ok, found = includeElement(1433, "1")
	False(t, ok)
	False(t, found)
}

func TestCondition(t *testing.T) {
	mockT := new(testing.T)

	if !Condition(mockT, func() bool { return true }, "Truth") {
		t.Error("Condition should return true")
	}

	if Condition(mockT, func() bool { return false }, "Lie") {
		t.Error("Condition should return false")
	}

}

func TestDidPanic(t *testing.T) {

	if funcDidPanic, _ := didPanic(func() {
		panic("Panic!")
	}); !funcDidPanic {
		t.Error("didPanic should return true")
	}

	if funcDidPanic, _ := didPanic(func() {
	}); funcDidPanic {
		t.Error("didPanic should return false")
	}

}

func TestPanics(t *testing.T) {

	mockT := new(testing.T)

	if !Panics(mockT, func() {
		panic("Panic!")
	}) {
		t.Error("Panics should return true")
	}

	if Panics(mockT, func() {
	}) {
		t.Error("Panics should return false")
	}

}

func TestNotPanics(t *testing.T) {

	mockT := new(testing.T)

	if !NotPanics(mockT, func() {
	}) {
		t.Error("NotPanics should return true")
	}

	if NotPanics(mockT, func() {
		panic("Panic!")
	}) {
		t.Error("NotPanics should return false")
	}

}

func TestNoError(t *testing.T) {

	mockT := new(testing.T)

	// start with a nil error
	var err error

	True(t, NoError(mockT, err), "NoError should return True for nil arg")

	// now set an error
	err = errors.New("some error")

	False(t, NoError(mockT, err), "NoError with error should return False")

	// returning an empty error interface
	err = func() error {
		var err *customError
		if err != nil {
			t.Fatal("err should be nil here")
		}
		return err
	}()

	if err == nil { // err is not nil here!
		t.Errorf("Error should be nil due to empty interface", err)
	}

	False(t, NoError(mockT, err), "NoError should fail with empty error interface")
}

type customError struct{}

func (*customError) Error() string { return "fail" }

func TestError(t *testing.T) {

	mockT := new(testing.T)

	// start with a nil error
	var err error

	False(t, Error(mockT, err), "Error should return False for nil arg")

	// now set an error
	err = errors.New("some error")

	True(t, Error(mockT, err), "Error with error should return True")

	// returning an empty error interface
	err = func() error {
		var err *customError
		if err != nil {
			t.Fatal("err should be nil here")
		}
		return err
	}()

	if err == nil { // err is not nil here!
		t.Errorf("Error should be nil due to empty interface", err)
	}

	True(t, Error(mockT, err), "Error should pass with empty error interface")
}

func TestEqualError(t *testing.T) {
	mockT := new(testing.T)

	// start with a nil error
	var err error
	False(t, EqualError(mockT, err, ""),
		"EqualError should return false for nil arg")

	// now set an error
	err = errors.New("some error")
	False(t, EqualError(mockT, err, "Not some error"),
		"EqualError should return false for different error string")
	True(t, EqualError(mockT, err, "some error"),
		"EqualError should return true")
}

func Test_isEmpty(t *testing.T) {

	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}

	True(t, isEmpty(""))
	True(t, isEmpty(nil))
	True(t, isEmpty([]string{}))
	True(t, isEmpty(0))
	True(t, isEmpty(int32(0)))
	True(t, isEmpty(int64(0)))
	True(t, isEmpty(false))
	True(t, isEmpty(map[string]string{}))
	True(t, isEmpty(new(time.Time)))
	True(t, isEmpty(time.Time{}))
	True(t, isEmpty(make(chan struct{})))
	False(t, isEmpty("something"))
	False(t, isEmpty(errors.New("something")))
	False(t, isEmpty([]string{"something"}))
	False(t, isEmpty(1))
	False(t, isEmpty(true))
	False(t, isEmpty(map[string]string{"Hello": "World"}))
	False(t, isEmpty(chWithValue))

}

func TestEmpty(t *testing.T) {

	mockT := new(testing.T)
	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}
	var tiP *time.Time
	var tiNP time.Time
	var s *string
	var f *os.File

	True(t, Empty(mockT, ""), "Empty string is empty")
	True(t, Empty(mockT, nil), "Nil is empty")
	True(t, Empty(mockT, []string{}), "Empty string array is empty")
	True(t, Empty(mockT, 0), "Zero int value is empty")
	True(t, Empty(mockT, false), "False value is empty")
	True(t, Empty(mockT, make(chan struct{})), "Channel without values is empty")
	True(t, Empty(mockT, s), "Nil string pointer is empty")
	True(t, Empty(mockT, f), "Nil os.File pointer is empty")
	True(t, Empty(mockT, tiP), "Nil time.Time pointer is empty")
	True(t, Empty(mockT, tiNP), "time.Time is empty")

	False(t, Empty(mockT, "something"), "Non Empty string is not empty")
	False(t, Empty(mockT, errors.New("something")), "Non nil object is not empty")
	False(t, Empty(mockT, []string{"something"}), "Non empty string array is not empty")
	False(t, Empty(mockT, 1), "Non-zero int value is not empty")
	False(t, Empty(mockT, true), "True value is not empty")
	False(t, Empty(mockT, chWithValue), "Channel with values is not empty")
}

func TestNotEmpty(t *testing.T) {

	mockT := new(testing.T)
	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}

	False(t, NotEmpty(mockT, ""), "Empty string is empty")
	False(t, NotEmpty(mockT, nil), "Nil is empty")
	False(t, NotEmpty(mockT, []string{}), "Empty string array is empty")
	False(t, NotEmpty(mockT, 0), "Zero int value is empty")
	False(t, NotEmpty(mockT, false), "False value is empty")
	False(t, NotEmpty(mockT, make(chan struct{})), "Channel without values is empty")

	True(t, NotEmpty(mockT, "something"), "Non Empty string is not empty")
	True(t, NotEmpty(mockT, errors.New("something")), "Non nil object is not empty")
	True(t, NotEmpty(mockT, []string{"something"}), "Non empty string array is not empty")
	True(t, NotEmpty(mockT, 1), "Non-zero int value is not empty")
	True(t, NotEmpty(mockT, true), "True value is not empty")
	True(t, NotEmpty(mockT, chWithValue), "Channel with values is not empty")
}

func Test_getLen(t *testing.T) {
	falseCases := []interface{}{
		nil,
		0,
		true,
		false,
		'A',
		struct{}{},
	}
	for _, v := range falseCases {
		ok, l := getLen(v)
		False(t, ok, "Expected getLen fail to get length of %#v", v)
		Equal(t, 0, l, "getLen should return 0 for %#v", v)
	}

	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	trueCases := []struct {
		v interface{}
		l int
	}{
		{[]int{1, 2, 3}, 3},
		{[...]int{1, 2, 3}, 3},
		{"ABC", 3},
		{map[int]int{1: 2, 2: 4, 3: 6}, 3},
		{ch, 3},

		{[]int{}, 0},
		{map[int]int{}, 0},
		{make(chan int), 0},

		{[]int(nil), 0},
		{map[int]int(nil), 0},
		{(chan int)(nil), 0},
	}

	for _, c := range trueCases {
		ok, l := getLen(c.v)
		True(t, ok, "Expected getLen success to get length of %#v", c.v)
		Equal(t, c.l, l)
	}
}

func TestLen(t *testing.T) {
	mockT := new(testing.T)

	False(t, Len(mockT, nil, 0), "nil does not have length")
	False(t, Len(mockT, 0, 0), "int does not have length")
	False(t, Len(mockT, true, 0), "true does not have length")
	False(t, Len(mockT, false, 0), "false does not have length")
	False(t, Len(mockT, 'A', 0), "Rune does not have length")
	False(t, Len(mockT, struct{}{}, 0), "Struct does not have length")

	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3

	cases := []struct {
		v interface{}
		l int
	}{
		{[]int{1, 2, 3}, 3},
		{[...]int{1, 2, 3}, 3},
		{"ABC", 3},
		{map[int]int{1: 2, 2: 4, 3: 6}, 3},
		{ch, 3},

		{[]int{}, 0},
		{map[int]int{}, 0},
		{make(chan int), 0},

		{[]int(nil), 0},
		{map[int]int(nil), 0},
		{(chan int)(nil), 0},
	}

	for _, c := range cases {
		True(t, Len(mockT, c.v, c.l), "%#v have %d items", c.v, c.l)
	}

	cases = []struct {
		v interface{}
		l int
	}{
		{[]int{1, 2, 3}, 4},
		{[...]int{1, 2, 3}, 2},
		{"ABC", 2},
		{map[int]int{1: 2, 2: 4, 3: 6}, 4},
		{ch, 2},

		{[]int{}, 1},
		{map[int]int{}, 1},
		{make(chan int), 1},

		{[]int(nil), 1},
		{map[int]int(nil), 1},
		{(chan int)(nil), 1},
	}

	for _, c := range cases {
		False(t, Len(mockT, c.v, c.l), "%#v have %d items", c.v, c.l)
	}
}

func TestWithinDuration(t *testing.T) {

	mockT := new(testing.T)
	a := time.Now()
	b := a.Add(10 * time.Second)

	True(t, WithinDuration(mockT, a, b, 10*time.Second), "A 10s difference is within a 10s time difference")
	True(t, WithinDuration(mockT, b, a, 10*time.Second), "A 10s difference is within a 10s time difference")

	False(t, WithinDuration(mockT, a, b, 9*time.Second), "A 10s difference is not within a 9s time difference")
	False(t, WithinDuration(mockT, b, a, 9*time.Second), "A 10s difference is not within a 9s time difference")

	False(t, WithinDuration(mockT, a, b, -9*time.Second), "A 10s difference is not within a 9s time difference")
	False(t, WithinDuration(mockT, b, a, -9*time.Second), "A 10s difference is not within a 9s time difference")

	False(t, WithinDuration(mockT, a, b, -11*time.Second), "A 10s difference is not within a 9s time difference")
	False(t, WithinDuration(mockT, b, a, -11*time.Second), "A 10s difference is not within a 9s time difference")
}

func TestInDelta(t *testing.T) {
	mockT := new(testing.T)

	True(t, InDelta(mockT, 1.001, 1, 0.01), "|1.001 - 1| <= 0.01")
	True(t, InDelta(mockT, 1, 1.001, 0.01), "|1 - 1.001| <= 0.01")
	True(t, InDelta(mockT, 1, 2, 1), "|1 - 2| <= 1")
	False(t, InDelta(mockT, 1, 2, 0.5), "Expected |1 - 2| <= 0.5 to fail")
	False(t, InDelta(mockT, 2, 1, 0.5), "Expected |2 - 1| <= 0.5 to fail")
	False(t, InDelta(mockT, "", nil, 1), "Expected non numerals to fail")
	False(t, InDelta(mockT, 42, math.NaN(), 0.01), "Expected NaN for actual to fail")
	False(t, InDelta(mockT, math.NaN(), 42, 0.01), "Expected NaN for expected to fail")

	cases := []struct {
		a, b  interface{}
		delta float64
	}{
		{uint8(2), uint8(1), 1},
		{uint16(2), uint16(1), 1},
		{uint32(2), uint32(1), 1},
		{uint64(2), uint64(1), 1},

		{int(2), int(1), 1},
		{int8(2), int8(1), 1},
		{int16(2), int16(1), 1},
		{int32(2), int32(1), 1},
		{int64(2), int64(1), 1},

		{float32(2), float32(1), 1},
		{float64(2), float64(1), 1},
	}

	for _, tc := range cases {
		True(t, InDelta(mockT, tc.a, tc.b, tc.delta), "Expected |%V - %V| <= %v", tc.a, tc.b, tc.delta)
	}
}

func TestInDeltaSlice(t *testing.T) {
	mockT := new(testing.T)

	True(t, InDeltaSlice(mockT,
		[]float64{1.001, 0.999},
		[]float64{1, 1},
		0.1), "{1.001, 0.009} is element-wise close to {1, 1} in delta=0.1")

	True(t, InDeltaSlice(mockT,
		[]float64{1, 2},
		[]float64{0, 3},
		1), "{1, 2} is element-wise close to {0, 3} in delta=1")

	False(t, InDeltaSlice(mockT,
		[]float64{1, 2},
		[]float64{0, 3},
		0.1), "{1, 2} is not element-wise close to {0, 3} in delta=0.1")

	False(t, InDeltaSlice(mockT, "", nil, 1), "Expected non numeral slices to fail")
}

func TestInEpsilon(t *testing.T) {
	mockT := new(testing.T)

	cases := []struct {
		a, b    interface{}
		epsilon float64
	}{
		{uint8(2), uint16(2), .001},
		{2.1, 2.2, 0.1},
		{2.2, 2.1, 0.1},
		{-2.1, -2.2, 0.1},
		{-2.2, -2.1, 0.1},
		{uint64(100), uint8(101), 0.01},
		{0.1, -0.1, 2},
		{0.1, 0, 2},
	}

	for _, tc := range cases {
		True(t, InEpsilon(t, tc.a, tc.b, tc.epsilon, "Expected %V and %V to have a relative difference of %v", tc.a, tc.b, tc.epsilon), "test: %q", tc)
	}

	cases = []struct {
		a, b    interface{}
		epsilon float64
	}{
		{uint8(2), int16(-2), .001},
		{uint64(100), uint8(102), 0.01},
		{2.1, 2.2, 0.001},
		{2.2, 2.1, 0.001},
		{2.1, -2.2, 1},
		{2.1, "bla-bla", 0},
		{0.1, -0.1, 1.99},
		{0, 0.1, 2}, // expected must be different to zero
	}

	for _, tc := range cases {
		False(t, InEpsilon(mockT, tc.a, tc.b, tc.epsilon, "Expected %V and %V to have a relative difference of %v", tc.a, tc.b, tc.epsilon))
	}

}

func TestInEpsilonSlice(t *testing.T) {
	mockT := new(testing.T)

	True(t, InEpsilonSlice(mockT,
		[]float64{2.2, 2.0},
		[]float64{2.1, 2.1},
		0.06), "{2.2, 2.0} is element-wise close to {2.1, 2.1} in espilon=0.06")

	False(t, InEpsilonSlice(mockT,
		[]float64{2.2, 2.0},
		[]float64{2.1, 2.1},
		0.04), "{2.2, 2.0} is not element-wise close to {2.1, 2.1} in espilon=0.04")

	False(t, InEpsilonSlice(mockT, "", nil, 1), "Expected non numeral slices to fail")
}

func TestRegexp(t *testing.T) {
	mockT := new(testing.T)

	cases := []struct {
		rx, str string
	}{
		{"^start", "start of the line"},
		{"end$", "in the end"},
		{"[0-9]{3}[.-]?[0-9]{2}[.-]?[0-9]{2}", "My phone number is 650.12.34"},
	}

	for _, tc := range cases {
		True(t, Regexp(mockT, tc.rx, tc.str))
		True(t, Regexp(mockT, regexp.MustCompile(tc.rx), tc.str))
		False(t, NotRegexp(mockT, tc.rx, tc.str))
		False(t, NotRegexp(mockT, regexp.MustCompile(tc.rx), tc.str))
	}

	cases = []struct {
		rx, str string
	}{
		{"^asdfastart", "Not the start of the line"},
		{"end$", "in the end."},
		{"[0-9]{3}[.-]?[0-9]{2}[.-]?[0-9]{2}", "My phone number is 650.12a.34"},
	}

	for _, tc := range cases {
		False(t, Regexp(mockT, tc.rx, tc.str), "Expected \"%s\" to not match \"%s\"", tc.rx, tc.str)
		False(t, Regexp(mockT, regexp.MustCompile(tc.rx), tc.str))
		True(t, NotRegexp(mockT, tc.rx, tc.str))
		True(t, NotRegexp(mockT, regexp.MustCompile(tc.rx), tc.str))
	}
}

func testAutogeneratedFunction() {
	defer func() {
		if err := recover(); err == nil {
			panic("did not panic")
		}
		CallerInfo()
	}()
	t := struct {
		io.Closer
	}{}
	var c io.Closer
	c = t
	c.Close()
}

func TestCallerInfoWithAutogeneratedFunctions(t *testing.T) {
	NotPanics(t, func() {
		testAutogeneratedFunction()
	})
}

func TestZero(t *testing.T) {
	mockT := new(testing.T)

	for _, test := range zeros {
		True(t, Zero(mockT, test, "%#v is not the %v zero value", test, reflect.TypeOf(test)))
	}

	for _, test := range nonZeros {
		False(t, Zero(mockT, test, "%#v is not the %v zero value", test, reflect.TypeOf(test)))
	}
}

func TestNotZero(t *testing.T) {
	mockT := new(testing.T)

	for _, test := range zeros {
		False(t, NotZero(mockT, test, "%#v is not the %v zero value", test, reflect.TypeOf(test)))
	}

	for _, test := range nonZeros {
		True(t, NotZero(mockT, test, "%#v is not the %v zero value", test, reflect.TypeOf(test)))
	}
}

func TestJSONEq_EqualSONString(t *testing.T) {
	mockT := new(testing.T)
	True(t, JSONEq(mockT, `{"hello": "world", "foo": "bar"}`, `{"hello": "world", "foo": "bar"}`))
}

func TestJSONEq_EquivalentButNotEqual(t *testing.T) {
	mockT := new(testing.T)
	True(t, JSONEq(mockT, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`))
}

func TestJSONEq_HashOfArraysAndHashes(t *testing.T) {
	mockT := new(testing.T)
	True(t, JSONEq(mockT, "{\r\n\t\"numeric\": 1.5,\r\n\t\"array\": [{\"foo\": \"bar\"}, 1, \"string\", [\"nested\", \"array\", 5.5]],\r\n\t\"hash\": {\"nested\": \"hash\", \"nested_slice\": [\"this\", \"is\", \"nested\"]},\r\n\t\"string\": \"foo\"\r\n}",
		"{\r\n\t\"numeric\": 1.5,\r\n\t\"hash\": {\"nested\": \"hash\", \"nested_slice\": [\"this\", \"is\", \"nested\"]},\r\n\t\"string\": \"foo\",\r\n\t\"array\": [{\"foo\": \"bar\"}, 1, \"string\", [\"nested\", \"array\", 5.5]]\r\n}"))
}

func TestJSONEq_Array(t *testing.T) {
	mockT := new(testing.T)
	True(t, JSONEq(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `["foo", {"nested": "hash", "hello": "world"}]`))
}

func TestJSONEq_HashAndArrayNotEquivalent(t *testing.T) {
	mockT := new(testing.T)
	False(t, JSONEq(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `{"foo": "bar", {"nested": "hash", "hello": "world"}}`))
}

func TestJSONEq_HashesNotEquivalent(t *testing.T) {
	mockT := new(testing.T)
	False(t, JSONEq(mockT, `{"foo": "bar"}`, `{"foo": "bar", "hello": "world"}`))
}

func TestJSONEq_ActualIsNotJSON(t *testing.T) {
	mockT := new(testing.T)
	False(t, JSONEq(mockT, `{"foo": "bar"}`, "Not JSON"))
}

func TestJSONEq_ExpectedIsNotJSON(t *testing.T) {
	mockT := new(testing.T)
	False(t, JSONEq(mockT, "Not JSON", `{"foo": "bar", "hello": "world"}`))
}

func TestJSONEq_ExpectedAndActualNotJSON(t *testing.T) {
	mockT := new(testing.T)
	False(t, JSONEq(mockT, "Not JSON", "Not JSON"))
}

func TestJSONEq_ArraysOfDifferentOrder(t *testing.T) {
	mockT := new(testing.T)
	False(t, JSONEq(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `[{ "hello": "world", "nested": "hash"}, "foo"]`))
}

func TestDiff(t *testing.T) {
	expected := `

Diff:
--- Expected
+++ Actual
@@ -1,3 +1,3 @@
 (struct { foo string }) {
- foo: (string) (len=5) "hello"
+ foo: (string) (len=3) "bar"
 }
`
	actual := diff(
		struct{ foo string }{"hello"},
		struct{ foo string }{"bar"},
	)
	Equal(t, expected, actual)

	expected = `

Diff:
--- Expected
+++ Actual
@@ -2,5 +2,5 @@
  (int) 1,
- (int) 2,
  (int) 3,
- (int) 4
+ (int) 5,
+ (int) 7
 }
`
	actual = diff(
		[]int{1, 2, 3, 4},
		[]int{1, 3, 5, 7},
	)
	Equal(t, expected, actual)

	expected = `

Diff:
--- Expected
+++ Actual
@@ -2,4 +2,4 @@
  (int) 1,
- (int) 2,
- (int) 3
+ (int) 3,
+ (int) 5
 }
`
	actual = diff(
		[]int{1, 2, 3, 4}[0:3],
		[]int{1, 3, 5, 7}[0:3],
	)
	Equal(t, expected, actual)

	expected = `

Diff:
--- Expected
+++ Actual
@@ -1,6 +1,6 @@
 (map[string]int) (len=4) {
- (string) (len=4) "four": (int) 4,
+ (string) (len=4) "five": (int) 5,
  (string) (len=3) "one": (int) 1,
- (string) (len=5) "three": (int) 3,
- (string) (len=3) "two": (int) 2
+ (string) (len=5) "seven": (int) 7,
+ (string) (len=5) "three": (int) 3
 }
`

	actual = diff(
		map[string]int{"one": 1, "two": 2, "three": 3, "four": 4},
		map[string]int{"one": 1, "three": 3, "five": 5, "seven": 7},
	)
	Equal(t, expected, actual)
}

func TestDiffEmptyCases(t *testing.T) {
	Equal(t, "", diff(nil, nil))
	Equal(t, "", diff(struct{ foo string }{}, nil))
	Equal(t, "", diff(nil, struct{ foo string }{}))
	Equal(t, "", diff(1, 2))
	Equal(t, "", diff(1, 2))
	Equal(t, "", diff([]int{1}, []bool{true}))
}

// Ensure there are no data races
func TestDiffRace(t *testing.T) {
	t.Parallel()

	expected := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}

	actual := map[string]string{
		"d": "D",
		"e": "E",
		"f": "F",
	}

	// run diffs in parallel simulating tests with t.Parallel()
	numRoutines := 10
	rChans := make([]chan string, numRoutines)
	for idx := range rChans {
		rChans[idx] = make(chan string)
		go func(ch chan string) {
			defer close(ch)
			ch <- diff(expected, actual)
		}(rChans[idx])
	}

	for _, ch := range rChans {
		for msg := range ch {
			NotZero(t, msg) // dummy assert
		}
	}
}

type mockTestingT struct {
}

func (m *mockTestingT) Errorf(format string, args ...interface{}) {}

func TestFailNowWithPlainTestingT(t *testing.T) {
	mockT := &mockTestingT{}

	Panics(t, func() {
		FailNow(mockT, "failed")
	}, "should panic since mockT is missing FailNow()")
}

type mockFailNowTestingT struct {
}

func (m *mockFailNowTestingT) Errorf(format string, args ...interface{}) {}

func (m *mockFailNowTestingT) FailNow() {}

func TestFailNowWithFullTestingT(t *testing.T) {
	mockT := &mockFailNowTestingT{}

	NotPanics(t, func() {
		FailNow(mockT, "failed")
	}, "should call mockT.FailNow() rather than panicking")
}
