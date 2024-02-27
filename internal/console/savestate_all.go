package console

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"slices"

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

var ErrNoPreviousState = errors.New("no previous state available")

func (c *Console) CreateUndoLoadState() error {
	var buf bytes.Buffer
	if err := c.SaveState(&buf); err != nil {
		return err
	}

	if len(c.undoLoadStates) >= c.config.State.UndoStateCount {
		c.undoLoadStates = slices.Delete(c.undoLoadStates, 0, 1)
	}
	c.undoLoadStates = append(c.undoLoadStates, buf.Bytes())

	return nil
}

func (c *Console) UndoLoadState() error {
	if len(c.undoLoadStates) == 0 {
		return ErrNoPreviousState
	}

	prev := c.undoLoadStates[len(c.undoLoadStates)-1]
	if err := c.LoadState(bytes.NewReader(prev)); err != nil {
		return err
	}

	c.undoLoadStates = slices.Delete(c.undoLoadStates, len(c.undoLoadStates)-1, len(c.undoLoadStates))
	return nil
}
