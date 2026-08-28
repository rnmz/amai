package config

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func initSqlTables(db *sqlx.DB) {
	tx := db.MustBegin()

	if _, filesTableErr := tx.Exec(`
	CREATE TABLE IF NOT EXISTS files (
		file_id UUID PRIMARY KEY,
		file_ext TEXT NOT NULL
	);
	`); filesTableErr != nil {
		slog.Error("[DB] Init table files failed", "error", filesTableErr)
		panic(fmt.Errorf("init table files failed: %v", filesTableErr))
	}

	if _, articlesTableErr := tx.Exec(`
	CREATE TABLE IF NOT EXISTS articles (
		id      UUID PRIMARY KEY,
		title   TEXT NOT NULL,
		poster  TEXT NOT NULL,
		created TIMESTAMP WITH TIME ZONE NOT NULL,
		updated TIMESTAMP WITH TIME ZONE NOT NULL,
		body    TEXT NOT NULL
	);
	`); articlesTableErr != nil {
		slog.Error("[DB] Init table articles failed", "error", articlesTableErr)
		panic(fmt.Errorf("init table articles failed: %v", articlesTableErr))
	}

	tx.Commit()
}

func InitDatabase() *sqlx.DB {
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		slog.Error("[DB] Error connecting to database", "error", err)
		return nil
	}
	db.SetMaxOpenConns(16)
	db.SetMaxIdleConns(16)

	if err := db.Ping(); err != nil {
		slog.Error("[DB] Error pinging database", "error", err)
		panic("database not available")
	}

	initSqlTables(db)

	return db
}

var closeOnce sync.Once

func CloseDatabase(db *sqlx.DB) {
	closeOnce.Do(func() {
		if db != nil {
			if err := db.Close(); err != nil {
				slog.Error("[DB] Error closing database", "error", err)
			}
		}
	})
}
