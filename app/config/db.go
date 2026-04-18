package app

import (
	"fmt"
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
		panic(fmt.Errorf("init table files failed: %v", filesTableErr))
	}

	if _, postTableErr := tx.Exec(`
	CREATE TABLE IF NOT EXISTS posts (
		id      UUID PRIMARY KEY,
		title   TEXT NOT NULL,
		poster  TEXT NOT NULL,
		created TIMESTAMP WITH TIME ZONE NOT NULL,
		updated TIMESTAMP WITH TIME ZONE NOT NULL,
		body    TEXT NOT NULL
	);
	`); postTableErr != nil {
		panic(fmt.Errorf("init table posts failed: %v", postTableErr))
	}

	tx.Commit()
}

func InitDatabase() *sqlx.DB {
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Error connecting to database. Message %s", err.Error())
	}
	db.SetMaxOpenConns(16)
	db.SetMaxIdleConns(16)

	if err := db.Ping(); err != nil {
		fmt.Printf("Error pinging database. Message %s", err.Error())
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
				fmt.Printf("Error closing database: %s\n", err)
			}
		}
	})
}
