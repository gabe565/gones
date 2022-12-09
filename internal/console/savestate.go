package console

import (
	"compress/gzip"
	"encoding/gob"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func (c *Console) SaveState(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	log.WithField("file", path).Info("Saving state")

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

	gzw := gzip.NewWriter(f)
	defer func() {
		_ = gzw.Close()
	}()

	if err := gob.NewEncoder(gzw).Encode(c); err != nil {
		return err
	}

	if err := gzw.Close(); err != nil {
		return err
	}

	return f.Close()
}

func (c *Console) LoadState(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
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

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer func() {
		_ = gzr.Close()
	}()

	if err := gob.NewDecoder(gzr).Decode(c); err != nil {
		return err
	}

	if err := gzr.Close(); err != nil {
		return err
	}

	return nil
}
