package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

/*
InitDB opens a SQLite database connection and returns *sqlx.DB.
It also creates the urls table if it does not exist.
*/
func InitDB(dbPath string) *sqlx.DB {
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original TEXT NOT NULL,
		short TEXT NOT NULL UNIQUE,
		created_at DATETIME NOT NULL
	);
	`

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("Failed to create urls table: %v", err)
	}

	fmt.Println("Database initialized successfully")
	return db
}
