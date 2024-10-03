package console

import (
	"encoding/base64"
	"log/slog"
	"path/filepath"
	"strings"
	"syscall/js"
)

func (c *Console) SaveStateNum(num uint8, createUndo bool) error {
	path, err := c.StatePath(num)
	if err != nil {
		return err
	}

	slog.Info("Saving state to db", "file", filepath.Base(path))

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
	path, err := c.StatePath(num)
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

	slog.Info("Loading state from db", "file", filepath.Base(path))

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
