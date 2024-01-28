//go:build !wasm

package console

import (
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
)

func (c *Console) SaveState(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	log.WithField("file", filepath.Base(path)).Info("Saving state")

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

	gzw := gzip.NewWriter(f)
	defer func() {
		_ = gzw.Close()
	}()

	encoder := msgpack.NewEncoder(gzw)
	encoder.UseCompactFloats(true)
	encoder.UseCompactInts(true)

	if err := encoder.Encode(c); err != nil {
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

	log.WithField("file", filepath.Base(path)).Info("Loading state")

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

	if err := msgpack.NewDecoder(gzr).Decode(c); err != nil {
		return err
	}

	if err := gzr.Close(); err != nil {
		return err
	}

	c.PPU.UpdatePalette(c.PPU.Mask.Get())

	return nil
}
