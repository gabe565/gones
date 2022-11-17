package console

import (
	"encoding/gob"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func (c *Console) SaveState(num uint8) error {
	path, err := c.cartridge.StatePath(num)
	if err != nil {
		return err
	}

	log.WithField("file", path).Info("Saving state")

	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err := gob.NewEncoder(f).Encode(c); err != nil {
		return err
	}

	return f.Close()
}

func (c *Console) LoadState(num uint8) error {
	path, err := c.cartridge.StatePath(num)
	if err != nil {
		return err
	}

	log.WithField("file", path).Info("Loading state")

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err := gob.NewDecoder(f).Decode(c); err != nil {
		return err
	}

	return nil
}
