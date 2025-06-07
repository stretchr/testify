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
		want     []FieldDiff
	}{
		{
			name:     "Equal strings",
			expected: "test",
			actual:   "test",
			want:     []FieldDiff{},
		},
		{
			name:     "Different strings",
			expected: "test",
			actual:   "other",
			want: []FieldDiff{{
				Path:     "",
				Expected: "test",
				Actual:   "other",
			}},
		},
		{
			name:     "Equal integers",
			expected: 42,
			actual:   42,
			want:     []FieldDiff{},
		},
		{
			name:     "Different integers",
			expected: 42,
			actual:   24,
			want: []FieldDiff{{
				Path:     "",
				Expected: 42,
				Actual:   24,
			}},
		},
		{
			name:     "Equal booleans",
			expected: true,
			actual:   true,
			want:     []FieldDiff{},
		},
		{
			name:     "Different booleans",
			expected: true,
			actual:   false,
			want: []FieldDiff{{
				Path:     "",
				Expected: true,
				Actual:   false,
			}},
		},
		{
			name:     "Different types",
			expected: 42,
			actual:   "42",
			want: []FieldDiff{{
				Path:     "",
				Expected: reflect.Int,
				Actual:   reflect.String,
			}},
		},
		{
			name:     "nil values",
			expected: nil,
			actual:   nil,
			want:     []FieldDiff{},
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
		want     []FieldDiff
	}{
		{
			name:     "NaN vs NaN",
			expected: math.NaN(),
			actual:   math.NaN(),
			want: []FieldDiff{
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
			want:     []FieldDiff{},
		},
		{
			name:     "Inf vs -Inf",
			expected: math.Inf(1),
			actual:   math.Inf(-1),
			want: []FieldDiff{
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
			want: []FieldDiff{
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
		want     []FieldDiff
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
			want: []FieldDiff{},
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
			want: []FieldDiff{
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
			want: []FieldDiff{},
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
			want: []FieldDiff{
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
			want: []FieldDiff{
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
		want     []FieldDiff
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
			want: []FieldDiff{},
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
			want: []FieldDiff{
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
			want: []FieldDiff{
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
			want:     []FieldDiff{},
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
		want     []FieldDiff
	}{
		{
			name:     "Equal slices",
			expected: []int{1, 2, 3},
			actual:   []int{1, 2, 3},
			want:     []FieldDiff{},
		},
		{
			name:     "Different slice values",
			expected: []int{1, 2, 3},
			actual:   []int{1, 4, 3},
			want: []FieldDiff{
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
			want: []FieldDiff{
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
			want:     []FieldDiff{},
		},
		{
			name:     "Different array values",
			expected: []int{1, 2, 3},
			actual:   []int{1, 5, 3},
			want: []FieldDiff{
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
			want: []FieldDiff{
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
			want: []FieldDiff{
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
		want     []FieldDiff
	}{
		{
			name:     "Equal maps",
			expected: map[string]int{"a": 1, "b": 2},
			actual:   map[string]int{"a": 1, "b": 2},
			want:     []FieldDiff{},
		},
		{
			name:     "Different map values",
			expected: map[string]int{"a": 1, "b": 2},
			actual:   map[string]int{"a": 1, "b": 3},
			want: []FieldDiff{
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
			want: []FieldDiff{
				{
					Path:     "",
					Expected: map[string]int{"a": 1, "b": 2, "c": 3},
					Actual:   map[string]int{"a": 1, "b": 2},
				},
			},
		},
		{
			name:     "Extra keys",
			expected: map[string]int{"a": 1, "b": 2},
			actual:   map[string]int{"a": 1, "b": 2, "c": 3},
			want: []FieldDiff{
				{
					Path:     "",
					Expected: map[string]int{"a": 1, "b": 2},
					Actual:   map[string]int{"a": 1, "b": 2, "c": 3},
				},
			},
		},
		{
			name:     "nil vs empty map",
			expected: map[string]int(nil),
			actual:   map[string]int{},
			want: []FieldDiff{
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
			want: []FieldDiff{
				{
					Path:     "[x].[b]",
					Expected: 2,
					Actual:   3,
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
		want     []FieldDiff
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
			want: []FieldDiff{},
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
			want: []FieldDiff{
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

// Helper function to compare FieldDiff slices without relying on order
func diffsAreEqual(a, b []FieldDiff) bool {
	if len(a) != len(b) {
		return false
	}

	// Create maps for easier comparison
	mapA := make(map[string]FieldDiff)
	mapB := make(map[string]FieldDiff)

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