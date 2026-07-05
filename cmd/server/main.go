package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/vect0r-0/url-shortener/internal/config"
	"github.com/vect0r-0/url-shortener/internal/storage/postgres"
	"github.com/vect0r-0/url-shortener/internal/transport/rest"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	storage, err := postgres.New(cfg.DB)

	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := storage.Close(); err != nil {
			log.Error("failed to close storage", "error", err)
		}
	}()

	log.Info("running database migrations")

	if err := storage.Migrate(); err != nil {
		log.Error("failed to migrate storage", "error", err)
		os.Exit(1)
	}
	log.Info("database is up to date")

	h := &http.Server{
		Handler: rest.New(log),
		Addr:    cfg.HttpServer.Address,
	}

	log.Info("starting http server", slog.String("addr", cfg.HttpServer.Address))
	if err := h.ListenAndServe(); err != nil {
		log.Error("failed to start http server", "error", err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	}

	return log
}
