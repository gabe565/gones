package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Load(_ *cobra.Command) (*Config, error) {
	log.Info("Loaded config")
	conf := NewDefault()
	return &conf, nil
}
