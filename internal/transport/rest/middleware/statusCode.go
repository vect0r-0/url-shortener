package middleware

import "net/http"

type statusCode struct {
	http.ResponseWriter
	statusCode int
}

func newStatusCode(w http.ResponseWriter) *statusCode {
	return &statusCode{ResponseWriter: w, statusCode: http.StatusOK}
}

func (s *statusCode) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}
