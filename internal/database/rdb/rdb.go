// Package rdb reads libretro RetroArch database (.rdb) files.
//
// The format is an 8-byte magic, a big-endian uint64 offset to a trailing
// metadata map, then a stream of concatenated MessagePack maps, one per game.
package rdb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"iter"

	"github.com/vmihailenco/msgpack/v5"
)

const (
	magic      = "RARCHDB\x00"
	headerSize = 16
)

var ErrInvalid = errors.New("invalid rdb")

// Entry is a single game record. Fields that gones does not use are skipped.
type Entry struct {
	Name string `msgpack:"name"`
	MD5  []byte `msgpack:"md5"`
}

func All(data []byte) iter.Seq2[Entry, error] {
	return func(yield func(Entry, error) bool) {
		if len(data) < headerSize || !bytes.HasPrefix(data, []byte(magic)) {
			yield(Entry{}, fmt.Errorf("%w: bad magic", ErrInvalid))
			return
		}

		// Entries run from the end of the header up to the metadata map.
		end := binary.BigEndian.Uint64(data[8:headerSize])
		if end < headerSize || end > uint64(len(data)) {
			yield(Entry{}, fmt.Errorf("%w: metadata offset %d out of range", ErrInvalid, end))
			return
		}

		dec := msgpack.NewDecoder(bytes.NewReader(data[headerSize:end]))
		for {
			var entry Entry
			if err := dec.Decode(&entry); err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				yield(Entry{}, err)
				return
			}

			if !yield(entry, nil) {
				return
			}
		}
	}
}
