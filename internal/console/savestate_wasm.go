package console

import (
	"compress/gzip"
	"encoding/base64"
	"path/filepath"
	"strings"
	"syscall/js"

	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
)

func (c *Console) SaveState(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	log.WithField("file", filepath.Base(path)).Info("Saving state to db")

	var buf strings.Builder

	b64w := base64.NewEncoder(base64.StdEncoding, &buf)

	gzw := gzip.NewWriter(b64w)
	defer func() {
		_ = gzw.Close()
	}()

	encoder := msgpack.NewEncoder(gzw)
	encoder.UseCompactFloats(true)
	encoder.UseCompactInts(true)
	encoder.SetSortMapKeys(true)

	if err := encoder.Encode(c); err != nil {
		return err
	}

	if err := gzw.Close(); err != nil {
		return err
	}

	if err := b64w.Close(); err != nil {
		return err
	}

	_, err = await(js.Global().Call("SaveToIndexedDb", "states", path, buf.String()))
	return err
}

func (c *Console) LoadState(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	vals, err := await(js.Global().Call("GetFromIndexedDb", "states", path))
	if err != nil {
		return err
	}
	data := vals[0]

	if data.IsNull() {
		return nil
	}

	log.WithField("file", filepath.Base(path)).Info("Loading state from db")

	r := strings.NewReader(data.String())

	b64r := base64.NewDecoder(base64.StdEncoding, r)

	gzr, err := gzip.NewReader(b64r)
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
