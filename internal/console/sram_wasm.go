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

	log.WithField("file", filepath.Base(path)).Info("Writing save to db")

	data := base64.StdEncoding.EncodeToString(c.Cartridge.Sram)

	_, err = await(js.Global().Call("SaveToIndexedDb", "saves", path, data))
	return err
}

func (c *Console) LoadSram() error {
	path, err := c.Cartridge.SramPath()
	if err != nil {
		return err
	}

	vals, err := await(js.Global().Call("GetFromIndexedDb", "saves", path))
	if err != nil {
		return err
	}
	data := vals[0]

	if data.IsNull() {
		return nil
	}

	log.WithField("file", filepath.Base(path)).Info("Loading save from db")

	if _, err := base64.StdEncoding.Decode(c.Cartridge.Sram, []byte(data.String())); err != nil {
		return err
	}

	return nil
}
