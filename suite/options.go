package suite

type runOptions struct {
	testFilter func(testName string) (bool, error)
}

// RunOption sets optional run parameter to specific value.
type RunOption func(*runOptions)

// WithTestFilter replaces default filter function which passed via -m flag to custom function.
func WithTestFilter(testFilter func(testName string) (bool, error)) RunOption {
	return func(opts *runOptions) {
		opts.testFilter = testFilter
	}
}

// WithIgnoreMatch ignores -m flag.
func WithIgnoreMatch() RunOption {
	return func(opts *runOptions) {
		opts.testFilter = func(testName string) (bool, error) {
			return true, nil
		}
	}
}
