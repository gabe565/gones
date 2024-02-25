//go:build !wasm

package gones

import (
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/console"
	"github.com/ncruces/zenity"
	log "github.com/sirupsen/logrus"
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

	cart, err := cartridge.FromiNesFile(path)
	if err != nil {
		return nil, err
	}
	log.WithField("title", cart.Name()).Info("Loaded cartridge")

	return console.New(conf, cart)
}
