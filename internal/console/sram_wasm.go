package console

import (
	"encoding/base64"
	"path/filepath"
	"syscall/js"

	log "github.com/sirupsen/logrus"
)

func (c *Console) SaveSRAM() error {
	path, err := c.Cartridge.SRAMPath()
	if err != nil {
		return err
	}

	log.WithField("file", filepath.Base(path)).Info("Writing save to db")

	data := base64.StdEncoding.EncodeToString(c.Cartridge.SRAM)

	_, err = await(js.Global().Get("GonesClient").Call("dbPut", "saves", path, data))
	return err
}

func (c *Console) LoadSRAM() error {
	path, err := c.Cartridge.SRAMPath()
	if err != nil {
		return err
	}

	vals, err := await(js.Global().Get("GonesClient").Call("dbGet", "saves", path))
	if err != nil {
		return err
	}
	data := vals[0]

	if data.IsNull() {
		return nil
	}

	log.WithField("file", filepath.Base(path)).Info("Loading save from db")

	if _, err := base64.StdEncoding.Decode(c.Cartridge.SRAM, []byte(data.String())); err != nil {
		return err
	}

	return nil
}
