package config

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func Load(_ *cobra.Command) (*Config, error) {
	slog.Info("Loaded config")
	conf := NewDefault()
	return &conf, nil
}
