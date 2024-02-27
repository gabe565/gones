//go:build !wasm

package console

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func (c *Console) SaveStateNum(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	if num == AutoSaveNum {
		log.WithField("file", filepath.Base(path)).Info("Auto-saving state")
	} else {
		log.WithField("file", filepath.Base(path)).Info("Saving state")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o777); err != nil {
		return err
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

	log.WithField("file", filepath.Base(path)).Info("Loading state")

	if err := c.LoadState(f); err != nil {
		if num == AutoSaveNum {
			log.WithError(err).Error("Load state failed. Moving state file and continuing.")
			if err := os.Rename(path, path+".failed"); err != nil {
				return err
			}
			return nil
		}

		return err
	}

	return nil
}
