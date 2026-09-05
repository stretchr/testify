//go:build !go1.23 && !goexperiment.rangefunc

package assert

// seqToSlice would convert a sequence of elements to a slice of the respective type.
// However, since sequences are not supported given the build tags, it just returns x as-is.
func seqToSlice(x interface{}) interface{} {
	return x
}
