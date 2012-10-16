package http

import (
	"net/http"
)

type TestResponseWriter struct {
	WrittenHeaderInt int
	Output           string
	header           http.Header
}

func (rw *TestResponseWriter) Header() http.Header {

	if rw.header == nil {
		rw.header = make(http.Header)
	}

	return rw.header
}
func (rw *TestResponseWriter) Write(bytes []byte) (int, error) {

	// add these bytes to the output string
	rw.Output = rw.Output + string(bytes)

	// return normal values
	return 0, nil

}
func (rw *TestResponseWriter) WriteHeader(i int) {
	rw.WrittenHeaderInt = i
}
