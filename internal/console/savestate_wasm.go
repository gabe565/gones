package console

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/gob"
	"path/filepath"
	"strings"
	"syscall/js"

	"github.com/gabe565/gones/internal/cartridge"
	log "github.com/sirupsen/logrus"
)

func (c *Console) SaveState(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	log.WithField("file", filepath.Base(path)).Info("Saving state")

	var buf strings.Builder

	b64w := base64.NewEncoder(base64.StdEncoding, &buf)

	gzw := gzip.NewWriter(b64w)
	defer func() {
		_ = gzw.Close()
	}()

	if err := gob.NewEncoder(gzw).Encode(c); err != nil {
		return err
	}

	if err := gzw.Close(); err != nil {
		return err
	}

	if err := b64w.Close(); err != nil {
		return err
	}

	js.Global().Get("localStorage").Call("setItem", path, buf.String())
	return nil
}

func (c *Console) LoadState(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	log.WithField("file", filepath.Base(path)).Info("Loading state")

	data := js.Global().Get("localStorage").Call("getItem", path)
	if data.IsNull() {
		return nil
	}

	r := strings.NewReader(data.String())

	b64r := base64.NewDecoder(base64.StdEncoding, r)

	gzr, err := gzip.NewReader(b64r)
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
