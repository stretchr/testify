package assert

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestFindDifferences_BasicTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		want     []fieldDiff
	}{
		{
			name:     "Equal strings",
			expected: "test",
			actual:   "test",
			want:     []fieldDiff{},
		},
		{
			name:     "Different strings",
			expected: "test",
			actual:   "other",
			want: []fieldDiff{{
				Path:     "",
				Expected: "test",
				Actual:   "other",
			}},
		},
		{
			name:     "Equal integers",
			expected: 42,
			actual:   42,
			want:     []fieldDiff{},
		},
		{
			name:     "Different integers",
			expected: 42,
			actual:   24,
			want: []fieldDiff{{
				Path:     "",
				Expected: 42,
				Actual:   24,
			}},
		},
		{
			name:     "Equal booleans",
			expected: true,
			actual:   true,
			want:     []fieldDiff{},
		},
		{
			name:     "Different booleans",
			expected: true,
			actual:   false,
			want: []fieldDiff{{
				Path:     "",
				Expected: true,
				Actual:   false,
			}},
		},
		{
			name:     "Different types",
			expected: 42,
			actual:   "42",
			want: []fieldDiff{{
				Path:     "",
				Expected: reflect.Int,
				Actual:   reflect.String,
			}},
		},
		{
			name:     "nil values",
			expected: nil,
			actual:   nil,
			want:     []fieldDiff{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := findDifferences(tt.expected, tt.actual)
			if !diffsAreEqual(got, tt.want) {
				t.Errorf("findDifferences() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindDifferences_SpecialValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		want     []fieldDiff
	}{
		{
			name:     "NaN vs NaN",
			expected: math.NaN(),
			actual:   math.NaN(),
			want: []fieldDiff{
				{
					Path:     "",
					Expected: math.NaN(),
					Actual:   math.NaN(),
				},
			},
		},
		{
			name:     "Inf vs Inf",
			expected: math.Inf(1),
			actual:   math.Inf(1),
			want:     []fieldDiff{},
		},
		{
			name:     "Inf vs -Inf",
			expected: math.Inf(1),
			actual:   math.Inf(-1),
			want: []fieldDiff{
				{
					Path:     "",
					Expected: math.Inf(1),
					Actual:   math.Inf(-1),
				},
			},
		},
		{
			name:     "Float vs NaN",
			expected: 1.0,
			actual:   math.NaN(),
			want: []fieldDiff{
				{
					Path:     "",
					Expected: 1.0,
					Actual:   math.NaN(),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := findDifferences(tt.expected, tt.actual)

			//  NaN we need a special check
			if tt.name == "NaN vs NaN" {
				if len(got) != 1 {
					t.Errorf("findDifferences() = %v, should return 1 diff for NaN comparison", got)
				}
				return
			}

			if !diffsAreEqual(got, tt.want) {
				t.Errorf("findDifferences() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindDifferences_Structs(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	type Employee struct {
		Person
		Department string
		Salary     float64
	}

	type Address struct {
		Street  string
		City    string
		ZipCode string
	}

	type ComplexPerson struct {
		Name    string
		Age     int
		Address Address
	}

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		want     []fieldDiff
	}{
		{
			name: "Equal simple structs",
			expected: Person{
				Name: "John",
				Age:  30,
			},
			actual: Person{
				Name: "John",
				Age:  30,
			},
			want: []fieldDiff{},
		},
		{
			name: "Different simple structs",
			expected: Person{
				Name: "John",
				Age:  30,
			},
			actual: Person{
				Name: "Jane",
				Age:  25,
			},
			want: []fieldDiff{
				{
					Path:     "Name",
					Expected: "John",
					Actual:   "Jane",
				},
				{
					Path:     "Age",
					Expected: 30,
					Actual:   25,
				},
			},
		},
		{
			name: "Equal nested structs",
			expected: ComplexPerson{
				Name: "John",
				Age:  30,
				Address: Address{
					Street:  "123 Main St",
					City:    "New York",
					ZipCode: "10001",
				},
			},
			actual: ComplexPerson{
				Name: "John",
				Age:  30,
				Address: Address{
					Street:  "123 Main St",
					City:    "New York",
					ZipCode: "10001",
				},
			},
			want: []fieldDiff{},
		},
		{
			name: "Different nested structs",
			expected: ComplexPerson{
				Name: "John",
				Age:  30,
				Address: Address{
					Street:  "123 Main St",
					City:    "New York",
					ZipCode: "10001",
				},
			},
			actual: ComplexPerson{
				Name: "John",
				Age:  30,
				Address: Address{
					Street:  "456 Oak Ave",
					City:    "Boston",
					ZipCode: "10001",
				},
			},
			want: []fieldDiff{
				{
					Path:     "Address.Street",
					Expected: "123 Main St",
					Actual:   "456 Oak Ave",
				},
				{
					Path:     "Address.City",
					Expected: "New York",
					Actual:   "Boston",
				},
			},
		},
		{
			name: "Embedded structs",
			expected: Employee{
				Person: Person{
					Name: "John",
					Age:  30,
				},
				Department: "Engineering",
				Salary:     100000,
			},
			actual: Employee{
				Person: Person{
					Name: "John",
					Age:  35,
				},
				Department: "Engineering",
				Salary:     90000,
			},
			want: []fieldDiff{
				{
					Path:     "Person.Age",
					Expected: 30,
					Actual:   35,
				},
				{
					Path:     "Salary",
					Expected: float64(100000),
					Actual:   float64(90000),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := findDifferences(tt.expected, tt.actual)
			if !diffsAreEqual(got, tt.want) {
				t.Errorf("findDifferences() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindDifferences_PointersAndNil(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name    string
		Address *string
	}

	addressA := "123 Main St"
	addressB := "456 Oak Ave"

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		want     []fieldDiff
	}{
		{
			name: "Both pointers equal",
			expected: Person{
				Name:    "John",
				Address: &addressA,
			},
			actual: Person{
				Name:    "John",
				Address: &addressA,
			},
			want: []fieldDiff{},
		},
		{
			name: "Pointers to different values",
			expected: Person{
				Name:    "John",
				Address: &addressA,
			},
			actual: Person{
				Name:    "John",
				Address: &addressB,
			},
			want: []fieldDiff{
				{
					Path:     "Address",
					Expected: "123 Main St",
					Actual:   "456 Oak Ave",
				},
			},
		},
		{
			name: "One pointer nil",
			expected: Person{
				Name:    "John",
				Address: &addressA,
			},
			actual: Person{
				Name:    "John",
				Address: nil,
			},
			want: []fieldDiff{
				{
					Path:     "Address",
					Expected: &addressA,
					Actual:   nil,
				},
			},
		},
		{
			name:     "Equal pointers to nil",
			expected: (*string)(nil),
			actual:   (*string)(nil),
			want:     []fieldDiff{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := findDifferences(tt.expected, tt.actual)
			if !diffsAreEqual(got, tt.want) {
				t.Errorf("findDifferences() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindDifferences_SlicesAndArrays(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		want     []fieldDiff
	}{
		{
			name:     "Equal slices",
			expected: []int{1, 2, 3},
			actual:   []int{1, 2, 3},
			want:     []fieldDiff{},
		},
		{
			name:     "Different slice values",
			expected: []int{1, 2, 3},
			actual:   []int{1, 4, 3},
			want: []fieldDiff{
				{
					Path:     "[1]",
					Expected: 2,
					Actual:   4,
				},
			},
		},
		{
			name:     "Different slice lengths",
			expected: []int{1, 2, 3},
			actual:   []int{1, 2},
			want: []fieldDiff{
				{
					Path:     "",
					Expected: []int{1, 2, 3},
					Actual:   []int{1, 2},
				},
			},
		},
		{
			name:     "Equal arrays",
			expected: []int{1, 2, 3},
			actual:   []int{1, 2, 3},
			want:     []fieldDiff{},
		},
		{
			name:     "Different array values",
			expected: []int{1, 2, 3},
			actual:   []int{1, 5, 3},
			want: []fieldDiff{
				{
					Path:     "[1]",
					Expected: 2,
					Actual:   5,
				},
			},
		},
		{
			name:     "nil vs empty slice",
			expected: []int(nil),
			actual:   []int{},
			want: []fieldDiff{
				{
					Path:     "",
					Expected: []int(nil),
					Actual:   []int{},
				},
			},
		},
		{
			name:     "Nested slices",
			expected: [][]int{{1, 2}, {3, 4}},
			actual:   [][]int{{1, 2}, {3, 5}},
			want: []fieldDiff{
				{
					Path:     "[1].[1]",
					Expected: 4,
					Actual:   5,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := findDifferences(tt.expected, tt.actual)
			if !diffsAreEqual(got, tt.want) {
				t.Errorf("findDifferences() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindDifferences_Maps(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		want     []fieldDiff
	}{
		{
			name:     "Equal maps",
			expected: map[string]int{"a": 1, "b": 2},
			actual:   map[string]int{"a": 1, "b": 2},
			want:     []fieldDiff{},
		},
		{
			name:     "Different map values",
			expected: map[string]int{"a": 1, "b": 2},
			actual:   map[string]int{"a": 1, "b": 3},
			want: []fieldDiff{
				{
					Path:     "[b]",
					Expected: 2,
					Actual:   3,
				},
			},
		},
		{
			name:     "Missing keys",
			expected: map[string]int{"a": 1, "b": 2, "c": 3},
			actual:   map[string]int{"a": 1, "b": 2},
			want: []fieldDiff{
				{
					Path:     "[c]",
					Expected: 3,
					Actual:   "<missing>",
				},
			},
		},
		{
			name:     "Extra keys",
			expected: map[string]int{"a": 1, "b": 2},
			actual:   map[string]int{"a": 1, "b": 2, "c": 3},
			want: []fieldDiff{
				{
					Path:     "[c]",
					Expected: "<missing>",
					Actual:   3,
				},
			},
		},
		{
			name:     "nil vs empty map",
			expected: map[string]int(nil),
			actual:   map[string]int{},
			want: []fieldDiff{
				{
					Path:     "",
					Expected: map[string]int(nil),
					Actual:   map[string]int{},
				},
			},
		},
		{
			name:     "Nested maps",
			expected: map[string]map[string]int{"x": {"a": 1, "b": 2}},
			actual:   map[string]map[string]int{"x": {"a": 1, "b": 3}},
			want: []fieldDiff{
				{
					Path:     "[x].[b]",
					Expected: 2,
					Actual:   3,
				},
			},
		},
		{
			name: "Multiple different keys and values",
			expected: map[string]interface{}{
				"name":  "Alice",
				"age":   30,
				"email": "alice@example.com",
			},
			actual: map[string]interface{}{
				"name":    "Bob",
				"age":     25,
				"address": "123 Main St",
			},
			want: []fieldDiff{
				{
					Path:     "[name]",
					Expected: "Alice",
					Actual:   "Bob",
				},
				{
					Path:     "[age]",
					Expected: 30,
					Actual:   25,
				},
				{
					Path:     "[email]",
					Expected: "alice@example.com",
					Actual:   "<missing>",
				},
				{
					Path:     "[address]",
					Expected: "<missing>",
					Actual:   "123 Main St",
				},
			},
		},
		{
			name: "Complex nested maps",
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"profile": map[string]interface{}{
						"settings": map[string]interface{}{
							"theme":      "dark",
							"fontSize":   12,
							"showAvatar": true,
						},
					},
				},
			},
			actual: map[string]interface{}{
				"user": map[string]interface{}{
					"profile": map[string]interface{}{
						"settings": map[string]interface{}{
							"theme":         "light",
							"fontSize":      14,
							"showAvatar":    true,
							"notifications": true,
						},
					},
				},
			},
			want: []fieldDiff{
				{
					Path:     "[user].[profile].[settings].[theme]",
					Expected: "dark",
					Actual:   "light",
				},
				{
					Path:     "[user].[profile].[settings].[fontSize]",
					Expected: 12,
					Actual:   14,
				},
				{
					Path:     "[user].[profile].[settings].[notifications]",
					Expected: "<missing>",
					Actual:   true,
				},
			},
		},
		{
			name: "Maps with numeric keys",
			expected: map[int]string{
				1: "one",
				2: "two",
				3: "three",
			},
			actual: map[int]string{
				1: "one",
				2: "deux",
				4: "four",
			},
			want: []fieldDiff{
				{
					Path:     "[2]",
					Expected: "two",
					Actual:   "deux",
				},
				{
					Path:     "[3]",
					Expected: "three",
					Actual:   "<missing>",
				},
				{
					Path:     "[4]",
					Expected: "<missing>",
					Actual:   "four",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := findDifferences(tt.expected, tt.actual)
			if !diffsAreEqual(got, tt.want) {
				t.Errorf("findDifferences() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindDifferences_ComplexNesting(t *testing.T) {
	t.Parallel()

	type Address struct {
		Street string
		City   string
	}

	type Person struct {
		Name      string
		Age       int
		Addresses []Address
		Tags      map[string]string
	}

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
		want     []fieldDiff
	}{
		{
			name: "Complex equal structures",
			expected: Person{
				Name: "John",
				Age:  30,
				Addresses: []Address{
					{Street: "123 Main St", City: "New York"},
					{Street: "456 Oak Ave", City: "Boston"},
				},
				Tags: map[string]string{
					"department": "Engineering",
					"level":      "Senior",
				},
			},
			actual: Person{
				Name: "John",
				Age:  30,
				Addresses: []Address{
					{Street: "123 Main St", City: "New York"},
					{Street: "456 Oak Ave", City: "Boston"},
				},
				Tags: map[string]string{
					"department": "Engineering",
					"level":      "Senior",
				},
			},
			want: []fieldDiff{},
		},
		{
			name: "Complex different structures",
			expected: Person{
				Name: "John",
				Age:  30,
				Addresses: []Address{
					{Street: "123 Main St", City: "New York"},
					{Street: "456 Oak Ave", City: "Boston"},
				},
				Tags: map[string]string{
					"department": "Engineering",
					"level":      "Senior",
				},
			},
			actual: Person{
				Name: "John",
				Age:  35,
				Addresses: []Address{
					{Street: "123 Main St", City: "New York"},
					{Street: "456 Oak Ave", City: "Chicago"},
				},
				Tags: map[string]string{
					"department": "Engineering",
					"level":      "Lead",
				},
			},
			want: []fieldDiff{
				{
					Path:     "Age",
					Expected: 30,
					Actual:   35,
				},
				{
					Path:     "Addresses.[1].City",
					Expected: "Boston",
					Actual:   "Chicago",
				},
				{
					Path:     "Tags.[level]",
					Expected: "Senior",
					Actual:   "Lead",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := findDifferences(tt.expected, tt.actual)
			if !diffsAreEqual(got, tt.want) {
				t.Errorf("findDifferences() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare fieldDiff slices without relying on order
func diffsAreEqual(a, b []fieldDiff) bool {
	if len(a) != len(b) {
		return false
	}

	// Create maps for easier comparison
	mapA := make(map[string]fieldDiff)
	mapB := make(map[string]fieldDiff)

	for _, diff := range a {
		mapA[diff.Path] = diff
	}

	for _, diff := range b {
		mapB[diff.Path] = diff
	}

	// Check if all diffs in a are in b
	for path, diffA := range mapA {
		diffB, ok := mapB[path]
		if !ok {
			return false
		}

		// For simplicity, we'll just compare expected and actual as strings
		// This isn't perfect but works for basic testing
		if fmt.Sprintf("%v", diffA.Expected) != fmt.Sprintf("%v", diffB.Expected) ||
			fmt.Sprintf("%v", diffA.Actual) != fmt.Sprintf("%v", diffB.Actual) {
			return false
		}
	}

	return true
}
