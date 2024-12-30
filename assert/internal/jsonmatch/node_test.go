package jsonmatch

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"number": 1,
		"float": 1.1,
		"string": "hello",
		"bool": true,
		"null": null,
		"array": [1, 2, 3],
		"object": {
			"key": "value"
		}
	}`

	n := Node{}

	err := json.Unmarshal([]byte(jsonStr), &n)
	if err != nil {
		t.Fatal(err)
	}

	expected := Node{
		kind: Object,
		value: map[string]Node{
			"number": {kind: Number, value: 1.0},
			"float":  {kind: Number, value: 1.1},
			"string": {kind: String, value: "hello"},
			"bool":   {kind: Bool, value: true},
			"null":   {kind: Null, value: nil},
			"array": {
				kind: Array,
				value: []Node{
					{kind: Number, value: 1.0},
					{kind: Number, value: 2.0},
					{kind: Number, value: 3.0},
				},
			},
			"object": {
				kind: Object,
				value: map[string]Node{
					"key": {kind: String, value: "value"},
				},
			},
		},
	}

	if !n.Matches(expected) {
		t.Fatalf("\n\nexpected %v\n\n actual %v", expected, n)
	}
}
