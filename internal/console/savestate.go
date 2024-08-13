//go:build !wasm

package console

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
)

func (c *Console) SaveStateNum(num uint8, createUndo bool) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	logger := slog.With("file", filepath.Base(path))
	if num == AutoSaveNum {
		logger.Info("Auto-saving state")
	} else {
		logger.Info("Saving state")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
		return err
	}

	if createUndo && num != AutoSaveNum {
		if oldState, err := os.ReadFile(path); err == nil {
			if err := c.CreateUndoSaveState(oldState); err != nil {
				return err
			}
		}
	}

	if err := os.Rename(path, path+".bak"); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err := c.SaveState(f); err != nil {
		return err
	}

	return f.Close()
}

func (c *Console) LoadStateNum(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	slog.Info("Loading state", "file", filepath.Base(path))

	if num != AutoSaveNum {
		if err := c.CreateUndoLoadState(); err != nil {
			return err
		}
	}

	if err := c.LoadState(f); err != nil {
		if num == AutoSaveNum {
			slog.Error("Load state failed. Moving state file and continuing.", "error", err)
			if err := os.Rename(path, path+".failed"); err != nil {
				return err
			}
			return nil
		}

		return err
	}

	return nil
}
