package http

import (
	"net/http"
)

// TestResponseWriter is a http.ResponseWriter object that keeps track of all activity
// allowing you to make assertions about how it was used.
type TestResponseWriter struct {

	// WrittenHeaderInt is the last int written by the call to WriteHeader(int)
	WrittenHeaderInt int

	// Output is a string containing the written bytes using the Write([]byte) func.
	Output string

	// header is the internal storage of the http.Header object
	header http.Header
}

// Header gets the http.Header describing the headers that were set in this response.
func (rw *TestResponseWriter) Header() http.Header {

	if rw.header == nil {
		rw.header = make(http.Header)
	}

	return rw.header
}

// Write writes the specified bytes to Output.
func (rw *TestResponseWriter) Write(bytes []byte) (int, error) {

	// add these bytes to the output string
	rw.Output = rw.Output + string(bytes)

	// return normal values
	return 0, nil

}

// WriteHeader stores the HTTP status code in the WrittenHeaderInt.
func (rw *TestResponseWriter) WriteHeader(i int) {
	rw.WrittenHeaderInt = i
}
