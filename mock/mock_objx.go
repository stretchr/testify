// This source file isolates the uses of the objx module to ease
// maintenance of downstream forks that remove that dependency.
// See https://github.com/stretchr/testify/issues/1752

//go:build !testify_no_objx && !testify_no_deps

package mock

import (
	"reflect"
)

type testData = TestData

// TestData holds any data that might be useful for testing.  Testify ignores
// this data completely allowing you to do whatever you like with it.
//
// Deprecated: do not use. The new TestData do not return an [github.com/stretchr/objx.Map] anymore.
// See https://github.com/stretchr/testify/issues/1852.
func (m *Mock) TestData() TestData {
	if m.testData == nil {
		m.testData = make(TestData)
	}

	return m.testData
}

// TestData replaces [github.com/stretchr/objx.Map].
type TestData map[string]interface{}

type reflectValue = reflect.Value

// TestDataValue replaces [github.com/stretchr/objx.Value] and exposes the same methods as [reflect.Value].
// Only a subset of objx.Value methods are available.
type TestDataValue struct {
	reflectValue
}

// Set replaces [github.com/stretchr/objx.Map.Set].
func (td TestData) Set(selector string, v interface{}) {
	td[selector] = v
}

// Get replaces [github.com/stretchr/objx.Map.Get].
func (td TestData) Get(selector string) *TestDataValue {
	v, ok := td[selector]
	if !ok {
		return nil
	}
	return &TestDataValue{reflectValue: reflect.ValueOf(&v).Elem()}
}

// MustInter replaces [github.com/stretchr/objx.Value.MustInter].
func (v *TestDataValue) MustInter() interface{} {
	if v == nil {
		return nil
	}
	// v.reflectValue contains an interface (ex: error), so dereference it
	return v.reflectValue.Elem()
}

// MustInter replaces [github.com/stretchr/objx.Value.Data].
func (v *TestDataValue) Data() interface{} {
	if v == nil {
		return nil
	}
	return v.reflectValue.Interface()
}
