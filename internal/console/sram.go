package console

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func (c *Console) SaveSram() error {
	path, err := c.cartridge.SramPath()
	if err != nil {
		return err
	}

	log.WithField("file", path).Info("Writing save to disk")

	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
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

	if _, err := f.Write(c.cartridge.Sram); err != nil {
		return err
	}

	return f.Close()
}

func (c *Console) LoadSram() error {
	path, err := c.cartridge.SramPath()
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

	log.WithField("file", path).Info("Loading save from disk")

	if _, err := f.Read(c.cartridge.Sram); err != nil {
		return err
	}

	return nil
}
