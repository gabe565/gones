package console

import (
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func (c *Console) SaveSram() error {
	path, err := c.Cartridge.SramPath()
	if err != nil {
		return err
	}

	log.WithField("file", filepath.Base(path)).Info("Writing save to disk")

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

	if _, err := f.Write(c.Cartridge.Sram); err != nil {
		return err
	}

	return f.Close()
}

func (c *Console) LoadSram() error {
	path, err := c.Cartridge.SramPath()
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

	log.WithField("file", filepath.Base(path)).Info("Loading save from disk")

	if _, err := f.Read(c.Cartridge.Sram); err != nil {
		return err
	}

	return nil
}
