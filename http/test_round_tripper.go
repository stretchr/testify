package http

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// Deprecated: Use https://pkg.go.dev/net/http/httptest instead.
type TestRoundTripper struct {
	mock.Mock
}

// Deprecated: Use https://pkg.go.dev/net/http/httptest instead.
func (t *TestRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := t.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
