package errors

import "net/http"

type HTTPError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func BadRequest(msg string) HTTPError { return HTTPError{Status: http.StatusBadRequest, Message: msg} }
func NotFound(msg string) HTTPError   { return HTTPError{Status: http.StatusNotFound, Message: msg} }
func BadGateway(msg string) HTTPError { return HTTPError{Status: http.StatusBadGateway, Message: msg} }
func Internal(msg string) HTTPError {
	return HTTPError{Status: http.StatusInternalServerError, Message: msg}
}
