package assert

import (
	"bytes"
	"errors"
	"io"
	http "net/http"
)

type builder struct {
	code           int
	body           io.ReadCloser
	expectedBody   bytes.Buffer
	requestHeader  http.Header
	responseHeader http.Header
	err            error
}

type HttpOption func(*builder) error

func WithCode(code int) HttpOption {
	return func(b *builder) error {
		if code < 100 || code > 511 {
			return errors.New("Given HTTP code is outside range of possible values assignement")
		}

		b.code = code
		return nil
	}
}

func WithErr(err error) HttpOption {
	return func(b *builder) error {
		b.err = err
		return nil
	}
}

func WithBody(body io.ReadCloser) HttpOption {
	return func(b *builder) error {
		b.body = body
		return nil
	}
}

func WithExpectedBody(expectedBody bytes.Buffer) HttpOption {
	return func(b *builder) error {
		b.expectedBody = expectedBody
		return nil
	}
}

func WithRequestHeader(requestHeader http.Header) HttpOption {
	return func(b *builder) error {
		b.requestHeader = requestHeader
		return nil
	}
}

func WithResponseHeader(responseHeader http.Header) HttpOption {
	return func(b *builder) error {
		b.responseHeader = responseHeader
		return nil
	}
}
