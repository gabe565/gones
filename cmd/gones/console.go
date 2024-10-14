//go:build !js

package gones

import (
	"log/slog"

	"gabe565.com/gones/internal/cartridge"
	"gabe565.com/gones/internal/config"
	"gabe565.com/gones/internal/console"
	"github.com/ncruces/zenity"
)

func newConsole(conf *config.Config, path string) (*console.Console, error) {
	if path == "" {
		var err error
		path, err = zenity.SelectFile(
			zenity.Title("Choose a ROM file"),
			zenity.FileFilter{
				Name:     "NES ROM",
				Patterns: []string{"*.nes"},
				CaseFold: true,
			},
		)
		if err != nil {
			return nil, err
		}
	}

	cart, err := cartridge.FromINESFile(path)
	if err != nil {
		return nil, err
	}
	slog.Info("Loaded cartridge", "", cart)

	return console.New(conf, cart)
}
