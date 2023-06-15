package console

import (
	"encoding/base64"
	"path/filepath"
	"strings"
	"syscall/js"

	log "github.com/sirupsen/logrus"
)

func (c *Console) SaveSram() error {
	path, err := c.Cartridge.SramPath()
	if err != nil {
		return err
	}

	log.WithField("file", filepath.Base(path)).Info("Writing save to disk")

	var buf strings.Builder

	b64w := base64.NewEncoder(base64.StdEncoding, &buf)

	if _, err := buf.Write(c.Cartridge.Sram); err != nil {
		return err
	}

	if err := b64w.Close(); err != nil {
		return err
	}

	js.Global().Get("localStorage").Call("setItem", path, buf.String())
	return nil
}

func (c *Console) LoadSram() error {
	path, err := c.Cartridge.SramPath()
	if err != nil {
		return err
	}

	data := js.Global().Get("localStorage").Call("getItem", path)
	if data.IsNull() {
		return nil
	}

	r := strings.NewReader(data.String())

	b64r := base64.NewDecoder(base64.StdEncoding, r)

	log.WithField("file", filepath.Base(path)).Info("Loading save from disk")

	if _, err := b64r.Read(c.Cartridge.Sram); err != nil {
		return err
	}

	return nil
}
