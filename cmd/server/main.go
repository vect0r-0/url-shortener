package main

import (
	"log/slog"
	"os"

	"github.com/vect0r-0/url-shortener/internal/config"
	"github.com/vect0r-0/url-shortener/internal/storage/postgres"
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
