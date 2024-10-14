package gones

import (
	"log/slog"

	"gabe565.com/gones/cmd/options"
	"gabe565.com/gones/internal/config"
)

type Command struct{}

func (c *Command) Execute() error {
	slog.Info("Loaded config")
	conf := config.NewDefault()
	return run(nil, conf, "")
}

func New(_ ...options.Option) *Command {
	return &Command{}
}
