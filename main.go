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

	err := db.Open()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	s := server.New(":9000", assets)
	s.Start()
}
