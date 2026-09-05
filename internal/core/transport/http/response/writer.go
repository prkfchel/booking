package core_http_response

import "net/http"

const (
	StatusCodeUninitialized = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: rw,
		StatusCode:     StatusCodeUninitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.StatusCode = statusCode
}

func (rw *ResponseWriter) GetStatusCode() int {
	if rw.StatusCode == StatusCodeUninitialized {
		return http.StatusOK
	}

	return rw.StatusCode
}
