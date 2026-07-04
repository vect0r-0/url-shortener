package postgres

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/vect0r-0/url-shortener/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(storageConnect config.DB) (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", storageConnect.User, storageConnect.Password,
		storageConnect.Host, storageConnect.Port, storageConnect.DBName)

	db, err := sql.Open("pgx", psqlInfo)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) Migrate() error {
	const op = "storage.postgres.Migrate"

	if err := goose.Up(s.db, "migrations"); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
