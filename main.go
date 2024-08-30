package main

import (
	"embed"
	"os"

	"github.com/hookenz/moneygo/api/db"
	"github.com/hookenz/moneygo/api/server"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed assets
var assets embed.FS

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	store := db.NewSqliteStore("database.db")

	err := store.Open()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	s := server.New(":9000", store, assets)
	s.Start()
}
