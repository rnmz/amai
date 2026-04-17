package app

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)
 
func InitDatabase() *sqlx.DB {
	dsn := os.Getenv("DSN")
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Error connecting to database. Message %s", err.Error())
	}
	db.SetMaxOpenConns(16)
	db.SetMaxIdleConns(16)

	initCommands, _ := os.ReadFile(filepath.Join(".", "init.sql"))
	db.MustExec(string(initCommands))

	if err := db.Ping(); err != nil {
		fmt.Printf("Error pinging database. Message %s", err.Error())
		panic("database not available")
	}
	return db
}

func CloseDatabase(db *sqlx.DB) {
	sync.OnceFunc(func() {
		if db != nil {
			if err := db.Close(); err != nil {
				fmt.Printf("Error closing database. Message %s", err.Error())
			}
		}
	})
}
