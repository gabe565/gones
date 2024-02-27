package console

import (
	"encoding/base64"
	"path/filepath"
	"strings"
	"syscall/js"

	log "github.com/sirupsen/logrus"
)

func (c *Console) SaveStateNum(num uint8) error {
	path, err := c.Cartridge.StatePath(num)
	if err != nil {
		return err
	}

	log.WithField("file", filepath.Base(path)).Info("Saving state to db")

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

	log.WithField("file", filepath.Base(path)).Info("Loading state from db")

	r := strings.NewReader(data.String())
	b64r := base64.NewDecoder(base64.StdEncoding, r)
	if err := c.LoadState(b64r); err != nil {
		return err
	}

	return nil
}
