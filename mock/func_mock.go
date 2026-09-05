package mock

import (
	"errors"
	"reflect"
	"testing"
)

var ErrFuncMockNotFunc = errors.New("not a function")

func FuncMockFor(fn interface{}) (*FuncMock, error) {
	typ := reflect.TypeOf(fn)
	if typ == nil || typ.Kind() != reflect.Func {
		return nil, ErrFuncMockNotFunc
	}

	return &FuncMock{
		mock: Mock{},
		typ:  typ,
	}, nil
}

type FuncMock struct {
	mock Mock
	typ  reflect.Type
}

func (m *FuncMock) Build() interface{} {
	return reflect.MakeFunc(m.typ, func(args []reflect.Value) []reflect.Value {
		argsAsInterface := make([]interface{}, len(args))
		for i, arg := range args {
			argsAsInterface[i] = arg.Interface()
		}
		outs := m.mock.MethodCalled("func", argsAsInterface...)
		res := make([]reflect.Value, m.typ.NumOut())
		for i := 0; i < m.typ.NumOut(); i++ {
			val := outs.Get(i)
			if val == nil {
				res[i] = reflect.Zero(m.typ.Out(i))
				continue
			}
			res[i] = reflect.ValueOf(val)
		}
		return res
	}).Interface()
}

func (m *FuncMock) On(args ...interface{}) *Call {
	return m.mock.On("func", args...)
}

func (m *FuncMock) AssertExpectations(t *testing.T) {
	m.mock.AssertExpectations(t)
}

func (m *FuncMock) AssertNotCalled(t *testing.T, arguments ...interface{}) {
	m.mock.AssertNotCalled(t, "func", arguments...)
}

func (m *FuncMock) AssertCalled(t *testing.T, arguments ...interface{}) {
	m.mock.AssertCalled(t, "func", arguments...)
}

func (m *FuncMock) AssertNumberOfCalls(t *testing.T, expectedCalls int) {
	m.mock.AssertNumberOfCalls(t, "func", expectedCalls)
}
