package console

import (
	"compress/gzip"
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

func (c *Console) SaveState(w io.Writer) error {
	gzw := gzip.NewWriter(w)
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

	return gzw.Close()
}

func (c *Console) LoadState(r io.Reader) error {
	gzr, err := gzip.NewReader(r)
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
