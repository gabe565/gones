package main

import (
	"log/slog"
	"os"

	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/log"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	log.Init(os.Stderr)

	f, err := os.Create("config_example.toml")
	if err != nil {
		slog.Error("Failed to create config example", "error", err)
		os.Exit(1)
	}

	encoder := toml.NewEncoder(f)
	conf := config.NewDefault()
	if err := encoder.Encode(conf); err != nil {
		slog.Error("Failed to write config example", "error", err)
		os.Exit(1)
	}
}
