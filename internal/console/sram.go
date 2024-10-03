//go:build !js

package console

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
)

func (c *Console) SaveSRAM() error {
	if !c.Cartridge.Battery {
		return nil
	}

	path, err := c.SRAMPath()
	if err != nil {
		return err
	}

	slog.Debug("Writing save to disk", "file", filepath.Base(path))

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

	path, err := c.SRAMPath()
	if err != nil {
		return err
	}

	sram, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	slog.Debug("Loading save from disk", "file", filepath.Base(path))

	c.Cartridge.SRAM = sram
	return nil
}
