package pq

import (
	"database/sql"
	"fmt"
	"log/slog"

	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3" // init sqlite3 driver
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string, log *slog.Logger) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		log.Error("Unable to")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS urls(
		id 		integer PRIMARY KEY,
		alias	text NOT NULL UNIQUE,
		url 	text NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);
	`)
	if err != nil {
		log.Error("Unable to create statement")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveUrl"

	_stmt := `INSERT INTO urls(url, alias) VALUES(?, ?)`

	stmt, err := s.db.Prepare(_stmt)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last inserted id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(urlToGet string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	// select statement
	stmt := `SELECT url FROM urls WHERE alias=?`

	// founded row or no rows
	row := s.db.QueryRow(stmt, urlToGet)

	var url string

	// result or ErrNoRows
	switch err := row.Scan(&url); err {
	case sql.ErrNoRows:
		return "", storage.ErrURLNotFound
	case nil:
		return url, nil
	default:
		return "", fmt.Errorf("%s: default proked: %w", op, err)
	}
}

func (s *Storage) DeleteUrl(urlToDelete string) (bool, error) {
	const op = "storage.sqlite.DeleteUrl"

	_stmt := `DELETE FROM urls WHERE alias = ?`

	stmt, err := s.db.Prepare(_stmt)
	if err != nil {
		return false, fmt.Errorf("%s: failed to prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(urlToDelete)
	if err != nil {
		return false, fmt.Errorf("%s: failed to exec statement: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: failed to get affected rows: %w", op, err)
	}

	if rows != 0 {
		return true, nil
	} else {
		return false, storage.ErrURLNotFound
	}
}
