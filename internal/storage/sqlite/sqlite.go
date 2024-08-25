package pq

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3" // init sqlite3 driver
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
