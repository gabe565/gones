//go:build !wasm

package console

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func (c *Console) SaveStateNum(num uint8, createUndo bool) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	if num == AutoSaveNum {
		log.Info().Str("file", filepath.Base(path)).Msg("Auto-saving state")
	} else {
		log.Info().Str("file", filepath.Base(path)).Msg("Saving state")
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

	log.Info().Str("file", filepath.Base(path)).Msg("Loading state")

	if err := c.CreateUndoLoadState(); err != nil {
		return err
	}

	if err := c.LoadState(f); err != nil {
		if num == AutoSaveNum {
			log.Err(err).Msg("Load state failed. Moving state file and continuing.")
			if err := os.Rename(path, path+".failed"); err != nil {
				return err
			}
			return nil
		}

		return err
	}

	return nil
}
