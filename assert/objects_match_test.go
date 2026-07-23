package assert

import (
	"testing"
)

func TestObjectsMatchStructsWithUnorderedSlices(t *testing.T) {
	type Test struct {
		Names []string
	}
	c1 := Test{Names: []string{"Joe", "Rick"}}
	c2 := Test{Names: []string{"Rick", "Joe"}}
	c3 := Test{Names: []string{"Joe", "Bob"}}

	mockT := new(testing.T)
	if !ObjectsMatch(mockT, c1, c2) {
		t.Fatal("expected match for same elements different order")
	}
	if ObjectsMatch(mockT, c1, c3) {
		t.Fatal("expected mismatch for different elements")
	}
}

func TestObjectsMatchNested(t *testing.T) {
	type Inner struct {
		Tags []string
	}
	type Outer struct {
		Items []Inner
	}
	a := Outer{Items: []Inner{{Tags: []string{"x", "y"}}, {Tags: []string{"a"}}}}
	b := Outer{Items: []Inner{{Tags: []string{"a"}}, {Tags: []string{"y", "x"}}}}
	mockT := new(testing.T)
	if !ObjectsMatch(mockT, a, b) {
		t.Fatal("expected nested unordered match")
	}
}

func TestJsonContentsMatch(t *testing.T) {
	expected := `{"participants":["Joe","Rick"],"event":"Birthday party"}`
	actual := `{"event":"Birthday party","participants":["Rick","Joe"]}`
	mockT := new(testing.T)
	if !JsonContentsMatch(mockT, expected, actual) {
		t.Fatal("expected JSON contents match")
	}
	if JsonContentsMatch(mockT, expected, `{"event":"other"}`) {
		t.Fatal("expected mismatch")
	}
}

func TestObjectsMatchMaps(t *testing.T) {
	a := map[string][]int{"k": {1, 2, 3}}
	b := map[string][]int{"k": {3, 1, 2}}
	mockT := new(testing.T)
	if !ObjectsMatch(mockT, a, b) {
		t.Fatal("map slice order should not matter")
	}
}
