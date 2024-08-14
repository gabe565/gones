package console

import (
	"encoding/base64"
	"log/slog"
	"path/filepath"
	"syscall/js"
)

func (c *Console) SaveSRAM() error {
	if !c.Cartridge.Battery {
		return nil
	}

	path, err := c.Cartridge.SRAMPath()
	if err != nil {
		return err
	}

	slog.Info("Writing save to db", "file", filepath.Base(path))

	data := base64.StdEncoding.EncodeToString(c.Cartridge.SRAM)

	_, err = await(js.Global().Get("GonesClient").Call("dbPut", "saves", path, data))
	return err
}

func (c *Console) LoadSRAM() error {
	if !c.Cartridge.Battery {
		return nil
	}

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

	slog.Info("Loading save from db", "file", filepath.Base(path))

	if _, err := base64.StdEncoding.Decode(c.Cartridge.SRAM, []byte(data.String())); err != nil {
		return err
	}

	return nil
}
