//go:build !wasm

package cmd

import (
	"errors"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/console"
)

func newConsole(path string) (*console.Console, error) {
	if path == "" {
		return nil, errors.New("No ROM provided")
	}

	cart, err := cartridge.FromiNesFile(path)
	if err != nil {
		return nil, err
	}

	return console.New(cart)
}
