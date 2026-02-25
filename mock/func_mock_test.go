package mock

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ExampleError = errors.New("test")

type TestFuncTypeNoArgumentsNoOuts func()

type TestFuncTypeArgumentsNoOuts func(a, b string)

type TestFuncTypeNoArgumentsOuts func() (string, error)

type TestFuncTypeArgumentsOuts func(a, b string) (string, error)

func TestFuncMockFor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var typ TestFuncTypeNoArgumentsNoOuts
		res, err := FuncMockFor(typ)
		assert.NoError(t, err)
		assert.IsType(t, &FuncMock{}, res)
	})
	t.Run("fail", func(t *testing.T) {
		t.Run("not a func", func(t *testing.T) {
			var typ interface{}
			res, err := FuncMockFor(typ)
			assert.Error(t, err)
			assert.Nil(t, res)
		})
	})
}

func TestFuncMock(t *testing.T) {
	t.Run("with no arguments and with no outputs", func(t *testing.T) {
		var typ TestFuncTypeNoArgumentsNoOuts
		funcMock, err := FuncMockFor(typ)
		assert.NoError(t, err)
		assert.IsType(t, &FuncMock{}, funcMock)
		defer funcMock.AssertNumberOfCalls(t, 1)

		funcMock.On().Return()
		fn := funcMock.Build().(TestFuncTypeNoArgumentsNoOuts)
		assert.NotPanics(t, func() {
			fn()
		})
	})
	t.Run("with arguments and with no outputs", func(t *testing.T) {
		var typ TestFuncTypeArgumentsNoOuts
		funcMock, err := FuncMockFor(typ)
		assert.NoError(t, err)
		assert.IsType(t, &FuncMock{}, funcMock)
		funcMock.On("a", "b").Return()
		defer funcMock.AssertNumberOfCalls(t, 1)

		funcMock.On().Return()
		fn := funcMock.Build().(TestFuncTypeArgumentsNoOuts)
		assert.NotPanics(t, func() {
			fn("a", "b")
		})
	})
	t.Run("with no arguments and with outputs", func(t *testing.T) {
		t.Run("with no error", func(t *testing.T) {
			var typ TestFuncTypeNoArgumentsOuts
			funcMock, err := FuncMockFor(typ)
			assert.NoError(t, err)
			assert.IsType(t, &FuncMock{}, funcMock)
			funcMock.On().Return("test", nil)
			defer funcMock.AssertNumberOfCalls(t, 1)

			funcMock.On().Return()
			fn := funcMock.Build().(TestFuncTypeNoArgumentsOuts)
			assert.NotPanics(t, func() {
				res, err := fn()
				assert.ErrorIs(t, err, nil)
				assert.Equal(t, "test", res)
			})
		})

		t.Run("with error", func(t *testing.T) {
			var typ TestFuncTypeNoArgumentsOuts
			funcMock, err := FuncMockFor(typ)
			assert.NoError(t, err)
			assert.IsType(t, &FuncMock{}, funcMock)
			funcMock.On().Return("test", ExampleError)
			defer funcMock.AssertNumberOfCalls(t, 1)

			funcMock.On().Return()
			fn := funcMock.Build().(TestFuncTypeNoArgumentsOuts)
			assert.NotPanics(t, func() {
				res, err := fn()
				assert.ErrorIs(t, err, ExampleError)
				assert.Equal(t, "test", res)
			})
		})
	})
	t.Run("with arguments and with outputs", func(t *testing.T) {
		t.Run("with no error", func(t *testing.T) {
			var typ TestFuncTypeArgumentsOuts
			funcMock, err := FuncMockFor(typ)
			assert.NoError(t, err)
			assert.IsType(t, &FuncMock{}, funcMock)
			funcMock.On().Return("test", ExampleError)
			defer funcMock.AssertNumberOfCalls(t, 1)

			funcMock.On("1", "2").Return("1 2", nil)
			fn := funcMock.Build().(TestFuncTypeArgumentsOuts)
			assert.NotPanics(t, func() {
				res, err := fn("1", "2")
				assert.NoError(t, err)
				assert.Equal(t, "1 2", res)
			})
		})
		t.Run("with error", func(t *testing.T) {
			var typ TestFuncTypeArgumentsOuts
			funcMock, err := FuncMockFor(typ)
			assert.NoError(t, err)
			assert.IsType(t, &FuncMock{}, funcMock)
			funcMock.On().Return("test", ExampleError)
			defer funcMock.AssertNumberOfCalls(t, 1)

			funcMock.On("1", "2").Return("1 2", ExampleError)
			fn := funcMock.Build().(TestFuncTypeArgumentsOuts)
			assert.NotPanics(t, func() {
				res, err := fn("1", "2")
				assert.ErrorIs(t, err, ExampleError)
				assert.Equal(t, "1 2", res)
			})
		})
	})
}
