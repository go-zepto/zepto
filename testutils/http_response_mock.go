package testutils

import (
	"net/http"
	"testing"
)

type ResponseMock struct {
	t       *testing.T
	headers http.Header
	body    []byte
	status  int
}

func NewResponseMock(t *testing.T) *ResponseMock {
	return &ResponseMock{
		t:       t,
		headers: make(http.Header),
	}
}

func (r *ResponseMock) Header() http.Header {
	return r.headers
}

func (r *ResponseMock) Write(body []byte) (int, error) {
	r.body = body
	return len(body), nil
}

func (r *ResponseMock) WriteHeader(status int) {
	r.status = status
}

func (r *ResponseMock) Assert(status int, body string) {
	if r.status != status {
		r.t.Errorf("expected status %+v to equal %+v", r.status, status)
	}
	if string(r.body) != body {
		r.t.Errorf("expected body %+v to equal %+v", string(r.body), body)
	}
}
