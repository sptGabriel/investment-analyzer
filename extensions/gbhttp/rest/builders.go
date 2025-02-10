package rest

import "net/http"

func OK(body any) Response {
	return NewResponse(http.StatusOK).WithBody(body)
}

func InternalServerError(err error) Response {
	return NewResponse(http.StatusInternalServerError).WithError(err)
}

func BadRequest(err error, body any) Response {
	return NewResponse(http.StatusBadRequest).WithError(err).WithBody(body)
}

func NewResponse(status int) Response {
	return Response{
		Status: status,
		header: http.Header{},
	}
}
