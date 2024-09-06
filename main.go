package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/hookenz/moneygo/api/db"
	"github.com/hookenz/moneygo/api/server"
	"golang.org/x/term"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed assets
var assets embed.FS

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var addUser bool
	flag.BoolVar(&addUser, "add-user", false, "add a new user")
	flag.Parse()

	store := db.NewSqliteStore("database.db")

	err := store.Open()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	if addUser {
		err = PromptAddUser(store)
		if err != nil {
			fmt.Printf("Error adding user: %v", err)
		}

		return
	}

	s := server.New(":9000", store, assets)
	s.Start()
}

func PromptAddUser(store db.Database) error {
	fmt.Println("Adding a new user to the database")
	fmt.Printf("Enter a username: ")

	var username string
	fmt.Scanln(&username)

	if username == "" {
		return fmt.Errorf("username cannot be blank")
	}

	fmt.Printf("Password: ")

	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return fmt.Errorf("password cannot be empty")
	}

	return store.InsertUser(username, string(password))
}
