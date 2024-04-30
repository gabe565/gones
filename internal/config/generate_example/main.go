package main

import (
	"os"

	"github.com/gabe565/gones/internal/config"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	f, err := os.Create("config_example.toml")
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create config example")
	}

	encoder := toml.NewEncoder(f)
	conf := config.NewDefault()
	if err := encoder.Encode(conf); err != nil {
		log.Panic().Err(err).Msg("Failed to write config example")
	}
}
