package jsonmatch

import (
	"encoding/json"
	"errors"
)

type (
	Node struct {
		kind  Kind
		value interface{}
	}

	Kind int
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

func (n *Node) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	newNode, err := NewNode(v)
	if err != nil {
		return err
	}

	n.kind = newNode.kind
	n.value = newNode.value

	return nil
}

func NewNode(v interface{}) (*Node, error) {
	n := &Node{}

	switch v := v.(type) {
	case map[string]interface{}:
		n.kind = Object
		val := make(map[string]Node)

		for k, childVal := range v {
			child, err := NewNode(childVal)
			if err != nil {
				return nil, err
			}
			val[k] = *child
		}

		n.value = val
	case []interface{}:
		n.kind = Array
		val := make([]Node, 0)
		for _, childVal := range v {
			child, err := NewNode(childVal)
			if err != nil {
				return nil, err
			}
			val = append(val, *child)
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
	default:
		return nil, ErrInvalidType
	}

	return n, nil
}

func (n Node) Matches(other Node) bool {
	if n.kind != other.kind {
		return false
	}

	switch n.kind {
	case Object:
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
			if !val[k].Matches(otherVal[k]) {
				return false
			}
		}
	case Array:
		for i := range n.value.([]Node) {
			val := n.value.([]Node)[i]
			otherVal := other.value.([]Node)[i]
			if !val.Matches(otherVal) {
				return false
			}
		}
	default:
		match := n.value == other.value
		if !match {
			return false
		}
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
