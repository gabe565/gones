package console

import (
	"compress/gzip"
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"

	"github.com/gabe565/gones/internal/cartridge"
	log "github.com/sirupsen/logrus"
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

	if err := gob.NewDecoder(gzr).Decode(c); err != nil {
		return err
	}

	if err := gzr.Close(); err != nil {
		return err
	}

	// gob doesn't handle interfaces in the same way as regular pointers
	// c.Mapper will be a new instance and needs to be setup
	c.Mapper.SetCartridge(c.Cartridge)
	if mapper, ok := c.Mapper.(cartridge.MapperInterrupts); ok {
		mapper.SetCpu(c.CPU)
	}
	c.Bus.SetMapper(c.Mapper)
	c.PPU.SetMapper(c.Mapper)

	return nil
}
