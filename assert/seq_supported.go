//go:build go1.23 || goexperiment.rangefunc

package assert

import "reflect"

var (
	boolType = reflect.TypeOf(true)
)

// seqToSlice checks if x is a sequence, and converts it to a slice of the
// same element type. Otherwise, x is returned as-is.
func seqToSlice(x interface{}) interface{} {
	if x == nil {
		return nil
	}

	xv := reflect.ValueOf(x)
	xt := xv.Type()
	// We're looking for a function with exactly one input parameter and no return values.
	if xt.Kind() != reflect.Func || xt.NumIn() != 1 || xt.NumOut() != 0 {
		return x
	}

	// The input parameter should be of type func(T) bool
	paramType := xt.In(0)
	if paramType.Kind() != reflect.Func || paramType.NumIn() != 1 || paramType.NumOut() != 1 || paramType.Out(0) != boolType {
		return x
	}

	elemType := paramType.In(0)
	resultType := reflect.SliceOf(elemType)
	result := reflect.MakeSlice(resultType, 0, 0)

	yieldFunc := reflect.MakeFunc(paramType, func(args []reflect.Value) []reflect.Value {
		result = reflect.Append(result, args[0])
		return []reflect.Value{reflect.ValueOf(true)}
	})

	// Call the function with the yield function as the argument
	xv.Call([]reflect.Value{yieldFunc})

	return result.Interface()
}
