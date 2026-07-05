package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/vect0r-0/url-shortener/internal/config"
	"github.com/vect0r-0/url-shortener/internal/storage"
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

func (s *Storage) SaveURL(ctx context.Context, id uuid.UUID, urlToSave string, alias string) error {
	const op = "storage.postgres.SaveURL"

	_, err := s.db.ExecContext(ctx,
		"INSERT INTO url(id, url, alias) VALUES ($1, $2, $3)",
		id, urlToSave, alias,
	)

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrURLExist)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	var url string

	if err := s.db.QueryRowContext(ctx,
		"SELECT url FROM url WHERE alias = $1",
		alias).Scan(&url); err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
	}

	return url, nil
}

func (s *Storage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.postgres.DeleteURL"

	if _, err := s.db.ExecContext(ctx,
		"DELETE FROM url WHERE alias = $1", alias); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAllURLS(ctx context.Context) ([]string, error) {
	const op = "storage.postgres.GetAllURLS"

	var urls []string

	rows, err := s.db.QueryContext(ctx, "SELECT url FROM url")

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	for rows.Next() {
		var url string

		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		urls = append(urls, url)
	}

	return urls, nil
}
