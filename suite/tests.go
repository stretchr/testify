package suite

import "testing"

type test = struct {
	name string
	run  func(t *testing.T)
}

type tests []test

func (ts tests) run(t *testing.T) {
	if len(ts) == 0 {
		t.Log("warning: no tests to run")
		return
	}

	for _, test := range ts {
		t.Run(test.name, test.run)
	}
}
