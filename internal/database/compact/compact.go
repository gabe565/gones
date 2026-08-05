// Package compact implements the database format embedded into release builds.
//
// A database is a two element msgpack array:
//
//	[0] bin  concatenated big endian HashLen byte MD5 prefixes
//	[1] str  newline separated names, parallel to the hashes
//
// Hashes are truncated to 64 bits, which halves the largest part of the file.
// Hashes are random, so they are incompressible; storing them as hex text
// instead costs twice the space and still cannot be compressed away. The
// database contains no collisions at this length, and an unrecognized ROM has
// a ~1.4e-15 chance of matching an entry by accident.
//
// Keeping the hashes and names in separate blocks means a lookup can scan the
// hashes and return without ever decoding the names, and lets gzip exploit the
// redundancy between similar titles. Entries are written in name order so that
// similar titles stay adjacent.
package compact

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

// HashLen is the number of MD5 bytes retained per entry.
const HashLen = 8

var ErrInvalid = errors.New("invalid database")

// fields is the number of elements in the msgpack array.
const fields = 2

const sep = "\n"

// Key returns the lookup key for a hex encoded MD5.
func Key(hash string) (uint64, error) {
	if len(hash) < HashLen*2 {
		return 0, fmt.Errorf("%w: short hash %q", ErrInvalid, hash)
	}

	b, err := hex.DecodeString(hash[:HashLen*2])
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	return binary.BigEndian.Uint64(b), nil
}

// Marshal encodes parallel hash and name slices. Names must not contain a newline.
func Marshal(hashes []uint64, names []string) ([]byte, error) {
	if len(hashes) != len(names) {
		return nil, fmt.Errorf("%w: %d hashes for %d names", ErrInvalid, len(hashes), len(names))
	}

	blob := make([]byte, 0, len(hashes)*HashLen)
	for _, hash := range hashes {
		blob = binary.BigEndian.AppendUint64(blob, hash)
	}

	for _, name := range names {
		if strings.Contains(name, sep) {
			return nil, fmt.Errorf("%w: name contains separator: %q", ErrInvalid, name)
		}
	}

	// Encoded field by field to mirror Find, which has to decode the same way
	// in order to stream.
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	if err := enc.EncodeArrayLen(fields); err != nil {
		return nil, err
	}
	if err := enc.EncodeBytes(blob); err != nil {
		return nil, err
	}
	if err := enc.EncodeString(strings.Join(names, sep)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Find streams r and stops as soon as key is located. Names are only decoded
// once a hash matches, so a lookup that misses never reads them.
func Find(r io.Reader, key uint64) (string, bool, error) {
	dec := msgpack.NewDecoder(r)

	n, err := dec.DecodeArrayLen()
	if err != nil {
		return "", false, err
	}
	if n != fields {
		return "", false, fmt.Errorf("%w: expected 2 fields, got %d", ErrInvalid, n)
	}

	hashes, err := dec.DecodeBytes()
	if err != nil {
		return "", false, err
	}
	if len(hashes)%HashLen != 0 {
		return "", false, fmt.Errorf("%w: %d hash bytes is not a multiple of %d",
			ErrInvalid, len(hashes), HashLen)
	}

	idx := -1
	for i := 0; i < len(hashes); i += HashLen {
		if binary.BigEndian.Uint64(hashes[i:]) == key {
			idx = i / HashLen
			break
		}
	}
	if idx == -1 {
		return "", false, nil
	}

	names, err := dec.DecodeBytes()
	if err != nil {
		return "", false, err
	}

	name, ok := nth(names, idx)
	if !ok {
		return "", false, fmt.Errorf("%w: no name at index %d", ErrInvalid, idx)
	}
	return string(name), true, nil
}

// nth returns the i-th newline separated field.
func nth(b []byte, i int) ([]byte, bool) {
	var curr int
	for name := range bytes.SplitSeq(b, []byte(sep)) {
		if curr == i {
			return name, true
		}
		curr++
	}
	return nil, false
}
