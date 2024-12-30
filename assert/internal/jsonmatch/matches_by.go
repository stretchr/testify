package jsonmatch

import (
	"errors"
)

type (
	// reprepsents a node with a value in a json object
	Node struct {
		kind  Kind
		value interface{}
	}

	// an enum of the kinds of nodes in a json object
	Kind int

	// a map of matchers with a function to determine if a value matches
	// the key of the matchers corresponds to the value passed into
	// the function
	Matchers map[string]func(interface{}) bool
)

const (
	Object Kind = iota
	Array
	String
	Number
	Bool
	Null
)

var ErrInvalidType = errors.New("invalid data type in json object")

func MatchesBy(expected, actual interface{}, matchers Matchers) bool {
	expectedNode := New(expected)
	actualNode := New(actual)

	return expectedNode.matchesBy(actualNode, matchers)
}

func New(v interface{}) Node {
	n := Node{}

	switch v := v.(type) {
	case map[string]interface{}:
		n.kind = Object
		val := make(map[string]Node)

		for k, childVal := range v {
			val[k] = New(childVal)
		}

		n.value = val
	case []interface{}:
		n.kind = Array
		val := make([]Node, 0)
		for _, childVal := range v {
			val = append(val, New(childVal))
		}
		n.value = val
	case string:
		n.kind = String
		n.value = v
	case float64:
		n.kind = Number
		n.value = v
	case bool:
		n.kind = Bool
		n.value = v
	case nil:
		n.kind = Null
		n.value = nil
	}

	return n
}

func (n Node) matchesBy(other Node, matchers Matchers) bool {
	switch n.kind {
	case Object:
		if other.kind != Object {
			return false
		}

		val := n.value.(map[string]Node)
		otherVal := other.value.(map[string]Node)

		nodeKeys := keys(val)
		otherKeys := keys(otherVal)
		if len(nodeKeys) != len(otherKeys) {
			return false
		}

		for i := range nodeKeys {
			k := nodeKeys[i]
			if _, ok := otherVal[k]; !ok {
				return false
			}
			if !val[k].matchesBy(otherVal[k], matchers) {
				return false
			}
		}
	case Array:
		if other.kind != Array {
			return false
		}

		for i := range n.value.([]Node) {
			val := n.value.([]Node)[i]
			otherVal := other.value.([]Node)[i]
			if !val.matchesBy(otherVal, matchers) {
				return false
			}
		}
	default:
		for val, f := range matchers {
			if val == n.value {
				return f(other.value)
			}
		}

		return n.value == other.value
	}

	return true
}

func keys(m map[string]Node) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
