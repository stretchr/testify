package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Diff with a FunctionalOptions expectation against a nil actual argument
// must report a difference instead of panicking on the nil type.
func Test_Arguments_Diff_FunctionalOptions_NilActual(t *testing.T) {
	t.Parallel()

	args := Arguments{FunctionalOptions(OpNum(1))}
	diff, count := args.Diff([]interface{}{nil})

	require.Equal(t, 1, count)
	assert.Contains(t, diff, "type <nil>")
}

// Diff with an empty FunctionalOptions expectation against a nil actual
// argument must report a difference instead of panicking in assertOpts.
func Test_Arguments_Diff_FunctionalOptions_EmptyExpected_NilActual(t *testing.T) {
	t.Parallel()

	args := Arguments{FunctionalOptions()}
	diff, count := args.Diff([]interface{}{nil})

	require.Equal(t, 1, count)
	assert.Contains(t, diff, "type <nil>")
}

// Diff with a nil option value in the expectation must report a
// difference instead of panicking when rendering the expected type.
func Test_Arguments_Diff_FunctionalOptions_NilExpectedValue(t *testing.T) {
	t.Parallel()

	args := Arguments{FunctionalOptions(nil)}
	diff, count := args.Diff([]interface{}{OpNum(1)})

	require.Equal(t, 1, count)
	assert.Contains(t, diff, "FAIL")
}

// A FunctionalOptions holding a nil value must render without panicking.
func Test_FunctionalOptions_NilValue_String(t *testing.T) {
	t.Parallel()

	assert.Contains(t, FunctionalOptions(nil).String(), "<nil>")
}

// Full mock flow: expectation uses FunctionalOptions, the call passes nil.
// The call must fail as an unexpected call, not crash the test run.
func Test_Mock_FunctionalOptions_NilActual(t *testing.T) {
	t.Parallel()

	mockT := &MockTestingT{}
	m := &Mock{}
	m.Test(mockT)
	m.On("Fn", FunctionalOptions(OpNum(1))).Return()

	assert.PanicsWithValue(t, mockTestingTFailNowCalled, func() {
		m.MethodCalled("Fn", nil)
	})
	assert.Equal(t, 1, mockT.errorfCount)
}

// Diff where the actual options slice holds a nil func element must
// report a difference instead of panicking in funcName/isFuncSame.
func Test_Arguments_Diff_FunctionalOptions_NilFuncElement(t *testing.T) {
	t.Parallel()

	args := Arguments{FunctionalOptions(OpNum(1))}
	diff, count := args.Diff([]interface{}{[]OptionFn{nil}})

	require.Equal(t, 1, count)
	assert.Contains(t, diff, "FAIL")
}

// Diff where a nil sits among the expected option values must report a
// difference instead of panicking in runtimeFunc.
func Test_Arguments_Diff_FunctionalOptions_NilInExpectedValues(t *testing.T) {
	t.Parallel()

	args := Arguments{FunctionalOptions(OpNum(1), nil)}
	diff, count := args.Diff([]interface{}{[]OptionFn{OpNum(1), OpNum(2)}})

	require.Equal(t, 1, count)
	assert.Contains(t, diff, "FAIL")
}
