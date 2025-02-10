package rest

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status int
	header http.Header
	Error  error
	Body   any
}

func (r Response) WithBody(body any) Response {
	r.Body = body
	return r
}

func (r Response) WithError(err error) Response {
	r.Error = err
	return r
}

func (r Response) AddHeader(name, value string) Response {
	if r.header == nil {
		r.header = http.Header{}
	}
	r.header.Add(name, value)
	return r
}

func (r Response) SetHeader(name, value string) Response {
	if r.header == nil {
		r.header = http.Header{}
	}
	r.header.Set(name, value)
	return r
}

func (r Response) Header() http.Header {
	return r.header
}

func ParseBody[T any](r *http.Request) (T, error) {
	var target, zero T
	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		return zero, err
	}
	return target, nil
}
