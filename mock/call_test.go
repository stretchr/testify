package mock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCaller(t *testing.T) {
	callerInfo := []string{"callerinfo"}
	args := []interface{}{1, "2", 3.0}
	methodName := "CallMethod"
	c := newCall(new(Mock), methodName, callerInfo, args...)

	assert.Equal(t, callerInfo, c.callerInfo)
	assert.Equal(t, methodName, c.Method)
	assert.Equal(t, Arguments(args), c.Arguments)

	t.Run("Should set repeabtability", func(t *testing.T) {
		c.Times(2)
		assert.Equal(t, 2, c.Repeatability)
	})

	t.Run("Should set WaitUntil channel", func(t *testing.T) {
		wait := make(chan time.Time)
		c.WaitUntil(wait)
		assert.True(t, wait == c.WaitFor)
	})

	t.Run("Should set wait duration", func(t *testing.T) {
		duration := time.Second
		c.After(duration)

		assert.Equal(t, duration, c.waitTime)
	})
}
