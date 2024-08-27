package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const filename string = "moneygo.db"

func Open() error {
	var err error
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	err = CreateTables(db)
	if err != nil {
		return fmt.Errorf("error creating database: %w", err)
	}

	return nil
}

func CreateTables(db *sql.DB) error {
	err := CreateUser(db)
	return err
}
