package user

import (
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/hookenz/moneygo/api/db"
)

type UserView struct {
	Name      string
	SessionID string
}

func Authenticate(db db.Database, username, password string) (UserView, error) {
	user := UserView{}
	record, err := db.SelectUser(username)
	if err != nil {
		return user, err
	}

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return user, err
	}

	match, err := argon2id.ComparePasswordAndHash(record.Password, hash)
	if err != nil {
		return user, fmt.Errorf("authentication failure")
	}

	if !match {
		return user, fmt.Errorf("authentication failure")
	}

	// Create a session

	return user, nil
}
