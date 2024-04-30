//go:build !wasm

package console

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func (c *Console) SaveSRAM() error {
	if !c.Cartridge.Battery {
		return nil
	}

	path, err := c.Cartridge.SRAMPath()
	if err != nil {
		return err
	}

	log.Debug().Str("file", filepath.Base(path)).Msg("Writing save to disk")

	if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
		return err
	}

	if err := os.Rename(path, path+".bak"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return os.WriteFile(path, c.Cartridge.SRAM, 0o666)
}

func (c *Console) LoadSRAM() error {
	if !c.Cartridge.Battery {
		return nil
	}

	path, err := c.Cartridge.SRAMPath()
	if err != nil {
		return err
	}

	sram, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	log.Debug().Str("file", filepath.Base(path)).Msg("Loading save from disk")

	c.Cartridge.SRAM = sram
	return nil
}
