package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"github.com/michaeljs1990/sqlitestore"
)

const maxAge = 3600

const createUserSQL = `CREATE TABLE IF NOT EXISTS user (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT UNIQUE NOT NULL CHECK (name LIKE '%_@__%.__%'),
	password TEXT NOT NULL
);`

const selectUserSQL = `SELECT id, name, password from USER WHERE name = '%'`

type SqliteStore struct {
	filename     string
	db           *sql.DB
	sessionStore *sqlitestore.SqliteStore
}

func NewSqliteStore(filename string) *SqliteStore {
	return &SqliteStore{
		filename: filename,
	}
}

func (s *SqliteStore) Open() error {
	if s.db != nil {
		return nil
	}

	var err error
	s.db, err = sql.Open("sqlite3", s.filename)
	if err != nil {
		return fmt.Errorf("error opening sqlite database file: %w", err)
	}

	err = s.createTables()
	if err != nil {
		return fmt.Errorf("error creating sqlite database: %w", err)
	}

	err = s.createSessionStore()
	if err != nil {
		return fmt.Errorf("error creating session store: %w", err)
	}

	return nil
}

func (s *SqliteStore) SelectUser(name string) (UserRecord, error) {
	row := s.db.QueryRow(selectUserSQL, name)
	user := UserRecord{}
	err := row.Scan(&user.ID, &user.Name, &user.Password)
	return user, err
}

func (s *SqliteStore) createTables() error {
	err := s.createTableUser()
	return err
}

func (s *SqliteStore) createTableUser() error {
	_, err := s.db.Exec(createUserSQL)
	if err != nil {
		return fmt.Errorf("error creating table user: %w", err)
	}

	return nil
}

func (s *SqliteStore) createSessionStore() error {
	var err error

	if s.db == nil {
		// this would only happen if we call this function before the store is open. We shouldn't do that.
		log.Fatal("db not open")
	}

	s.sessionStore, err = sqlitestore.NewSqliteStoreFromConnection(s.db, "sessions", "/", maxAge, []byte(os.Getenv("SESSION_KEY")))
	return err
}

func (s *SqliteStore) SessionStore() sessions.Store {
	return s.sessionStore
}
