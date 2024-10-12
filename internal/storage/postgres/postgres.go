package postgres

import (
	"context"
	"fmt"
	"github.com/ankodd/url-shortener/internal/config"
	"github.com/ankodd/url-shortener/internal/storage"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

// Storage model which contains connection to database
type Storage struct {
	conn *pgx.Conn
}

// New initializing new Storage object
//
// # Connecting to database and creating table
//
// Returning *Storage
func New(pgConf *config.PostgreSQL) (*Storage, error) {
	const fn = "postgres.New"

	// Connect to database
	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     pgConf.Host,
		Port:     pgConf.Port,
		Database: pgConf.Database,
		User:     pgConf.User,
		Password: pgConf.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	// Test connection to database
	if err = conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	// Create table
	_, err = conn.Exec(`
		CREATE TABLE IF NOT EXISTS url(
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", fn, err)
	}

	return &Storage{
		conn: conn,
	}, nil
}

// SaveURL saving url to database
//
// Returning int64, error
func (s *Storage) SaveURL(URL, alias string) (int64, error) {
	const fn = "postgres.SaveURL"

	row := s.conn.QueryRow("INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id", URL, alias)

	var id int64

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%s, %w", fn, err)
		}

		return 0, fmt.Errorf(storage.ErrAliasAlreadyExists)
	}

	return id, nil
}

// GetUrl getting url from database by alias
//
// Returning string, error
func (s *Storage) GetUrl(alias string) (string, error) {
	const fn = "postgres.GetUrl"

	row := s.conn.QueryRow("SELECT url FROM url WHERE alias = $1", alias)

	var url string

	if err := row.Scan(&url); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s, %w", fn, err)
		}

		return "", errors.New(storage.ErrAliasNotFound)
	}

	return url, nil
}

func (s *Storage) Close() error {
	return s.conn.Close()
}
