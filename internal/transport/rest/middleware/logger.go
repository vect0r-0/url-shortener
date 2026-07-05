package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logger(logger *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			reqID := GetRequestID(r.Context())

			start := time.Now()

			sW := newStatusCode(w)

			next.ServeHTTP(sW, r)

			logger.Info("completed request",
				slog.String("request_id", reqID),
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.Int("status_code", sW.statusCode),
				slog.Duration("duration", time.Since(start)))

		})
	}
}
