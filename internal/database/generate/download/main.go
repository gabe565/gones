package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	action, err := NewDownloader("Nintendo - Nintendo Entertainment System")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create downloader")
	}

	if err := action.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run downloader")
	}
}
