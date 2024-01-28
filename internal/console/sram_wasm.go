package console

import (
	"encoding/base64"
	"path/filepath"
	"syscall/js"

	log "github.com/sirupsen/logrus"
)

func (c *Console) SaveSram() error {
	path, err := c.Cartridge.SramPath()
	if err != nil {
		return err
	}

	log.WithField("file", filepath.Base(path)).Info("Writing save to localstorage")

	data := base64.StdEncoding.EncodeToString(c.Cartridge.Sram)
	js.Global().Get("localStorage").Call("setItem", path, data)
	return nil
}

func (c *Console) LoadSram() error {
	path, err := c.Cartridge.SramPath()
	if err != nil {
		return err
	}

	log.WithField("file", filepath.Base(path)).Info("Loading save from localstorage")

	data := js.Global().Get("localStorage").Call("getItem", path)
	if data.IsNull() {
		return nil
	}

	if _, err := base64.StdEncoding.Decode(c.Cartridge.Sram, []byte(data.String())); err != nil {
		return err
	}

	return nil
}
