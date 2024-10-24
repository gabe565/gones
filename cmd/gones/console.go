//go:build !js

package gones

import (
	"log/slog"

	"gabe565.com/gones/internal/cartridge"
	"gabe565.com/gones/internal/config"
	"gabe565.com/gones/internal/console"
	"github.com/ncruces/zenity"
)

func loadCartridge(path string) (*cartridge.Cartridge, error) {
	if path == "" {
		var err error
		if path, err = zenity.SelectFile(
			zenity.Title("Choose a ROM file"),
			zenity.FileFilter{
				Name:     "NES ROM",
				Patterns: []string{"*.nes"},
				CaseFold: true,
			},
		); err != nil {
			return nil, err
		}
	}

	cart, err := cartridge.FromINESFile(path)
	if err != nil {
		return nil, err
	}
	slog.Info("Loaded cartridge", "", cart)

	return cart, nil
}

func newConsole(conf *config.Config, cart *cartridge.Cartridge) (*console.Console, error) {
	return console.New(conf, cart)
}
