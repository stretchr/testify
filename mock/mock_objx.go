// This source file isolates the uses of the objx module to ease
// maintenance of downstream forks that remove that dependency.
// See https://github.com/stretchr/testify/issues/1752

package mock

import "github.com/stretchr/objx"

type testData = objx.Map

// TestData holds any data that might be useful for testing.  Testify ignores
// this data completely allowing you to do whatever you like with it.
func (m *Mock) TestData() objx.Map {
	if m.testData == nil {
		m.testData = make(objx.Map)
	}

	return m.testData
}
