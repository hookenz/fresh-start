package main

import (
	"os"

	"embed"

	"github.com/hookenz/moneygo/api/server"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed dist/public*
var frontend embed.FS

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	s := server.New(":8080", frontend)
	s.Start()
}
