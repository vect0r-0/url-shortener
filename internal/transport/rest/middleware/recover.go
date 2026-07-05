package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recover(logger *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("recover from panic",
						slog.Any("error", rec),
						slog.String("request_id", GetRequestID(r.Context())),
						slog.String("stack", string(debug.Stack())))

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
