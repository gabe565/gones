package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	action, err := NewDownloader("Nintendo - Nintendo Entertainment System")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create downloader")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := action.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to run downloader") //nolint:gocritic
	}
}
