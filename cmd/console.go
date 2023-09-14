//go:build !wasm

package cmd

import (
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/console"
	"github.com/ncruces/zenity"
)

func newConsole(path string) (*console.Console, error) {
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

	return console.New(cart)
}
