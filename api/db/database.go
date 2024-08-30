package db

import (
	"github.com/gorilla/sessions"
)

type Database interface {
	Open() error

	SelectUser(name string) (UserRecord, error)
	SessionStore() sessions.Store
}

type UserRecord struct {
	ID       string
	Name     string
	Password string
}
