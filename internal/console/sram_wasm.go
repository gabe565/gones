package console

import (
	"encoding/base64"
	"path/filepath"
	"syscall/js"

	"github.com/rs/zerolog/log"
)

func (c *Console) SaveSRAM() error {
	path, err := c.Cartridge.SRAMPath()
	if err != nil {
		return err
	}

	log.Info().Str("file", filepath.Base(path)).Msg("Writing save to db")

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

	log.Info().Str("file", filepath.Base(path)).Msg("Loading save from db")

	if _, err := base64.StdEncoding.Decode(c.Cartridge.SRAM, []byte(data.String())); err != nil {
		return err
	}

	return nil
}
