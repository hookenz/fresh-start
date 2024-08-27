package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const createUserSQL = `CREATE TABLE IF NOT EXISTS user (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT UNIQUE NOT NULL CHECK (name LIKE '%_@__%.__%'),
	password TEXT NOT NULL
);`

func CreateUser(db *sql.DB) error {
	_, err := db.Exec(createUserSQL)
	if err != nil {
		return fmt.Errorf("error creating table user: %w", err)
	}

	return nil
}
