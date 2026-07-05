package rest

import (
	"log/slog"
	"net/http"

	"github.com/vect0r-0/url-shortener/internal/transport/rest/middleware"
)

type Server struct {
	mux    *http.ServeMux
	logger *slog.Logger
}

func New(logger *slog.Logger) *Server {
	mux := http.NewServeMux()
	return &Server{mux: mux, logger: logger}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.Chain(s.mux,
		middleware.Recover(s.logger),
		middleware.RequestID,
		middleware.Logger(s.logger),
	).ServeHTTP(w, r)
}
