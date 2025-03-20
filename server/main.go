package main

import (
	"alertwest-interview-q1/lib"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	db := lib.NewDB()
	server := NewServer(db)

	// Start components
	db.Run()
	server.Start(":8080")
}
