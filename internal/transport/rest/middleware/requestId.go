package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestIDKey string

const (
	requestIDVal    = requestIDKey("RequestID")
	RequestIDHeader = "X-Request-ID"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := r.Header.Get(RequestIDHeader)

		if v == "" {
			v = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), requestIDVal, v)

		w.Header().Set(RequestIDHeader, v)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func GetRequestID(ctx context.Context) string {

	v, ok := ctx.Value(requestIDVal).(string)
	if !ok {
		return ""
	}

	return v
}
