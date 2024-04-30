package console

import (
	"encoding/base64"
	"path/filepath"
	"strings"
	"syscall/js"

	"github.com/rs/zerolog/log"
)

func (c *Console) SaveStateNum(num uint8, createUndo bool) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	log.Info().Str("file", filepath.Base(path)).Msg("Saving state to db")

	if createUndo && num != AutoSaveNum {
		vals, err := await(js.Global().Get("GonesClient").Call("dbGet", "states", path))
		if err == nil {
			data, err := base64.StdEncoding.DecodeString(vals[0].String())
			if err == nil {
				if err := c.CreateUndoSaveState(data); err != nil {
					return err
				}
			}
		}
	}

	var buf strings.Builder
	b64w := base64.NewEncoder(base64.StdEncoding, &buf)
	if err := c.SaveState(b64w); err != nil {
		return err
	}
	if err := b64w.Close(); err != nil {
		return err
	}

	_, err = await(js.Global().Get("GonesClient").Call("dbPut", "states", path, buf.String()))
	return err
}

func (c *Console) LoadStateNum(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	vals, err := await(js.Global().Get("GonesClient").Call("dbGet", "states", path))
	if err != nil {
		return err
	}
	data := vals[0]

	if data.IsNull() {
		return nil
	}

	log.Info().Str("file", filepath.Base(path)).Msg("Loading state from db")

	if err := c.CreateUndoLoadState(); err != nil {
		return err
	}

	r := strings.NewReader(data.String())
	b64r := base64.NewDecoder(base64.StdEncoding, r)
	if err := c.LoadState(b64r); err != nil {
		return err
	}

	return nil
}
