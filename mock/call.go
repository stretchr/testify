package mock

import "time"

/*
	Call
*/

// Call represents a method call and is used for setting expectations,
// as well as recording activity.
type Call struct {
	Parent *Mock

	// The name of the method that was or will be called.
	Method string

	// Holds the arguments of the method.
	Arguments Arguments

	// Holds the arguments that should be returned when
	// this method is called.
	ReturnArguments Arguments

	// Holds the caller info for the On() call
	callerInfo []string

	// The number of times to return the return arguments when setting
	// expectations. 0 means to always return the value.
	Repeatability int

	// Amount of times this call has been called
	totalCalls int

	// Call to this method can be optional
	optional bool

	// Holds a channel that will be used to block the Return until it either
	// receives a message or is closed. nil means it returns immediately.
	WaitFor <-chan time.Time

	waitTime time.Duration

	// Holds a handler used to manipulate arguments content that are passed by
	// reference. It's useful when mocking methods such as unmarshalers or
	// decoders.
	RunFn func(Arguments)
}

func newCall(parent *Mock, methodName string, callerInfo []string, methodArguments ...interface{}) *Call {
	return &Call{
		Parent:          parent,
		Method:          methodName,
		Arguments:       methodArguments,
		ReturnArguments: make([]interface{}, 0),
		callerInfo:      callerInfo,
		Repeatability:   0,
		WaitFor:         nil,
		RunFn:           nil,
	}
}

func (c *Call) lock() {
	c.Parent.mutex.Lock()
}

func (c *Call) unlock() {
	c.Parent.mutex.Unlock()
}

// Return specifies the return arguments for the expectation.
//
//    Mock.On("DoSomething").Return(errors.New("failed"))
func (c *Call) Return(returnArguments ...interface{}) *Call {
	c.lock()
	defer c.unlock()

	c.ReturnArguments = returnArguments

	return c
}

// Once indicates that that the mock should only return the value once.
//
//    Mock.On("MyMethod", arg1, arg2).Return(returnArg1, returnArg2).Once()
func (c *Call) Once() *Call {
	return c.Times(1)
}

// Twice indicates that that the mock should only return the value twice.
//
//    Mock.On("MyMethod", arg1, arg2).Return(returnArg1, returnArg2).Twice()
func (c *Call) Twice() *Call {
	return c.Times(2)
}

// Times indicates that that the mock should only return the indicated number
// of times.
//
//    Mock.On("MyMethod", arg1, arg2).Return(returnArg1, returnArg2).Times(5)
func (c *Call) Times(i int) *Call {
	c.lock()
	defer c.unlock()
	c.Repeatability = i
	return c
}

// WaitUntil sets the channel that will block the mock's return until its closed
// or a message is received.
//
//    Mock.On("MyMethod", arg1, arg2).WaitUntil(time.After(time.Second))
func (c *Call) WaitUntil(w <-chan time.Time) *Call {
	c.lock()
	defer c.unlock()
	c.WaitFor = w
	return c
}

// After sets how long to block until the call returns
//
//    Mock.On("MyMethod", arg1, arg2).After(time.Second)
func (c *Call) After(d time.Duration) *Call {
	c.lock()
	defer c.unlock()
	c.waitTime = d
	return c
}

// Run sets a handler to be called before returning. It can be used when
// mocking a method such as unmarshalers that takes a pointer to a struct and
// sets properties in such struct
//
//    Mock.On("Unmarshal", AnythingOfType("*map[string]interface{}").Return().Run(func(args Arguments) {
//    	arg := args.Get(0).(*map[string]interface{})
//    	arg["foo"] = "bar"
//    })
func (c *Call) Run(fn func(args Arguments)) *Call {
	c.lock()
	defer c.unlock()
	c.RunFn = fn
	return c
}

// Maybe allows the method call to be optional. Not calling an optional method
// will not cause an error while asserting expectations
func (c *Call) Maybe() *Call {
	c.lock()
	defer c.unlock()
	c.optional = true
	return c
}

// On chains a new expectation description onto the mocked interface. This
// allows syntax like.
//
//    Mock.
//       On("MyMethod", 1).Return(nil).
//       On("MyOtherMethod", 'a', 'b', 'c').Return(errors.New("Some Error"))
func (c *Call) On(methodName string, arguments ...interface{}) *Call {
	return c.Parent.On(methodName, arguments...)
}
