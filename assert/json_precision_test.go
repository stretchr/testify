package assert

import (
	"strings"
	"testing"
)

func TestJSONEqNumericPrecision(t *testing.T) {
	cases := []struct {
		name, expected, actual string
	}{
		{"large integers", `{"id":9007199254740992}`, `{"id":9007199254740993}`},
		{"negative integers", `[-9007199254740992]`, `[-9007199254740993]`},
		{"int64 limits", `9223372036854775806`, `9223372036854775807`},
		{"decimal digits", `0.1`, `0.10000000000000001`},
		{"exponent digits", `1e0`, `1.0000000000000001e0`},
		{"underflow", `1e-400`, `0`},
		{"negative underflow", `-1e-400`, `0`},
		{"rounded subnormal", `4e-324`, `5e-324`},
		{"nested values", `{"items":[{"amount":0.1}]}`, `{"items":[{"amount":0.10000000000000001}]}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			if JSONEq(mockT, tc.expected, tc.actual) {
				t.Errorf("different JSON numbers compared equal: %s and %s", tc.expected, tc.actual)
			}
			if !mockT.Failed() {
				t.Error("missing assertion failure message")
			}
		})
	}
}

func TestJSONEqNumericEquivalence(t *testing.T) {
	cases := []struct {
		name, expected, actual string
	}{
		{"integer and decimal", `1`, `1.0`},
		{"integer and exponent", `1000`, `1e3`},
		{"fraction zeroes", `0.00000100`, `1e-6`},
		{"negative value", `-120.00100`, `-120001e-3`},
		{"signed zero", `-0`, `0.0`},
		{"zero exponent", `-0e99999999999999999999999999999`, `0`},
		{"tiny exponent", `1e-99999999999999999999999999999`, `10e-100000000000000000000000000000`},
		{"large integer", `9007199254740993`, `9007199254740993.0`},
		{"nested key order", `{"a":1,"b":[2.00]}`, `{"b":[2e0],"a":1.0}`},
		{"duplicate key", `{"value":1,"value":2}`, `{"value":2.0}`},
		{"same bytes", `{"value":9007199254740993}`, `{"value":9007199254740993}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			if !JSONEq(mockT, tc.expected, tc.actual) {
				t.Errorf("equivalent JSON numbers failed comparison: %s", mockT.errorString())
			}
		})
	}
}

func TestJSONEqPreservesValidation(t *testing.T) {
	cases := []struct {
		name, expected, actual, diagnostic string
	}{
		{"expected overflow", `1e1000`, `0`, "Expected value ('1e1000') is not valid json."},
		{"actual overflow", `0`, `1e1000`, "Input ('1e1000') needs to be valid json."},
		{"same overflow", `1e1000`, `1e1000`, "Expected value ('1e1000') is not valid json."},
		{"expected invalid", `01`, `1`, "Expected value ('01') is not valid json."},
		{"actual invalid", `1`, `1 2`, "Input ('1 2') needs to be valid json."},
		{"same invalid", `1 2`, `1 2`, "Expected value ('1 2') is not valid json."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockT := &mockTestingT{}
			if JSONEq(mockT, tc.expected, tc.actual) {
				t.Fatal("invalid input compared equal")
			}
			if !strings.Contains(mockT.errorString(), tc.diagnostic) {
				t.Errorf("missing %q in failure message: %s", tc.diagnostic, mockT.errorString())
			}
		})
	}
}
