package assert

import "testing"

func TestPanicsWithValueUncomparable(t *testing.T) {
	pointerValue := 1
	otherPointerValue := 1
	cases := []struct {
		name             string
		expected, actual interface{}
		want             bool
	}{
		{"same pointer", &pointerValue, &pointerValue, true},
		{"distinct pointers", &pointerValue, &otherPointerValue, false},
		{"comparable struct with distinct pointers", struct{ V *int }{&pointerValue}, struct{ V *int }{&otherPointerValue}, false},
		{"interface containing slice", struct{ V interface{} }{[]int{1}}, struct{ V interface{} }{[]int{1}}, true},
		{"equal strings", "value", "value", true},
		{"different strings", "value", "other", false},
		{"different types", int(1), int64(1), false},
		{"equal slices", []int{1, 2}, []int{1, 2}, true},
		{"different slices", []int{1, 2}, []int{1, 3}, false},
		{"equal maps", map[string]int{"a": 1}, map[string]int{"a": 1}, true},
		{"different maps", map[string]int{"a": 1}, map[string]int{"a": 2}, false},
		{"struct containing slice", struct{ Values []int }{[]int{1}}, struct{ Values []int }{[]int{1}}, true},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := PanicsWithValue(new(mockTestingT), tt.expected, func() { panic(tt.actual) })
			if got != tt.want {
				t.Errorf("PanicsWithValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
