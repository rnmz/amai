package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS files (
	file_id  TEXT PRIMARY KEY,
	file_ext TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS articles (
	id      TEXT PRIMARY KEY,
	title   TEXT NOT NULL,
	poster  TEXT NOT NULL,
	created DATETIME NOT NULL,
	updated DATETIME NOT NULL,
	body    TEXT NOT NULL
);`

func InitDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", os.Getenv("DSN"))
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("init schema: %w", err)
	}

	return db, nil
}

func CloseDatabase(db *sqlx.DB) {
	if db == nil {
		return
	}
	if err := db.Close(); err != nil {
		slog.Error("[DB] Error closing database", "error", err)
	}
}
