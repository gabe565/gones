package gones

import (
	"log/slog"

	"github.com/gabe565/gones/cmd/options"
	"github.com/gabe565/gones/internal/config"
)

type Command struct{}

func (c *Command) Execute() error {
	slog.Info("Loaded config")
	conf := config.NewDefault()
	return run(nil, &conf, "")
}

func New(_ ...options.Option) *Command {
	return &Command{}
}
