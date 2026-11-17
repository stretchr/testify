package assert

import (
	"net/http"
	"strings"
	"testing"
)

func TestHTTPBodyAssertionsRejectInvalidRequests(t *testing.T) {
	for _, request := range []struct {
		name, method, target string
	}{
		{"URL", "GET", "://bad"},
		{"method", "invalid method", "/"},
	} {
		for _, check := range []struct {
			name   string
			assert func(TestingT, http.HandlerFunc, string, string, string) bool
		}{
			{"contains", func(tt TestingT, h http.HandlerFunc, m, u, s string) bool {
				return HTTPBodyContains(tt, h, m, u, nil, s, "request context")
			}},
			{"not contains", func(tt TestingT, h http.HandlerFunc, m, u, s string) bool {
				return HTTPBodyNotContains(tt, h, m, u, nil, s, "request context")
			}},
		} {
			for _, substring := range []string{"", "absent"} {
				t.Run(request.name+"/"+check.name+"/"+substring, func(t *testing.T) {
					called := false
					handler := func(http.ResponseWriter, *http.Request) { called = true }
					tt := new(mockTestingT)
					if check.assert(tt, handler, request.method, request.target, substring) {
						t.Error("invalid request assertion returned true")
					}
					if !tt.Failed() || !strings.Contains(tt.errorString(), "Failed to build test request") || !strings.Contains(tt.errorString(), "request context") {
						t.Errorf("missing request failure or caller context: %s", tt.errorString())
					}
					if called {
						t.Error("handler ran for an invalid request")
					}
				})
			}
		}
	}
}

func TestHTTPBodyAssertionsValidEmptyBody(t *testing.T) {
	handler := func(http.ResponseWriter, *http.Request) {}
	for _, check := range []struct {
		name string
		run  func(TestingT) bool
		want bool
	}{
		{"contains empty", func(tt TestingT) bool { return HTTPBodyContains(tt, handler, "GET", "/", nil, "") }, true},
		{"contains absent", func(tt TestingT) bool { return HTTPBodyContains(tt, handler, "GET", "/", nil, "absent") }, false},
		{"not contains empty", func(tt TestingT) bool { return HTTPBodyNotContains(tt, handler, "GET", "/", nil, "") }, false},
		{"not contains absent", func(tt TestingT) bool { return HTTPBodyNotContains(tt, handler, "GET", "/", nil, "absent") }, true},
	} {
		t.Run(check.name, func(t *testing.T) {
			tt := new(mockTestingT)
			if got := check.run(tt); got != check.want || tt.Failed() == check.want {
				t.Errorf("got %v, failed=%v; want %v", got, tt.Failed(), check.want)
			}
		})
	}
}

func TestHTTPBodyInvalidRequestCompatibility(t *testing.T) {
	for _, method := range []string{"GET", "invalid method"} {
		called := false
		got := HTTPBody(func(http.ResponseWriter, *http.Request) { called = true }, method, "://bad", nil)
		if got != "" || called {
			t.Errorf("HTTPBody returned %q, called=%v", got, called)
		}
	}
}
