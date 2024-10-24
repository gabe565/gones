package gones

import (
	"gabe565.com/gones/cmd/options"
	"gabe565.com/gones/internal/config"
)

type Command struct{}

func (c *Command) Execute() error {
	conf := config.NewDefault()

	cart, err := loadCartridge()
	if err != nil {
		return err
	}

	return run(nil, conf, cart)
}

func New(_ ...options.Option) *Command {
	return &Command{}
}
