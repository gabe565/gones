package rdb

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"
)

// build assembles a valid rdb file from the given entries.
func build(t *testing.T, entries ...map[string]any) []byte {
	t.Helper()

	var body bytes.Buffer
	for _, entry := range entries {
		b, err := msgpack.Marshal(entry)
		require.NoError(t, err)
		body.Write(b)
	}

	meta, err := msgpack.Marshal(map[string]any{"count": len(entries)})
	require.NoError(t, err)

	var buf bytes.Buffer
	buf.WriteString(magic)
	require.NoError(t, binary.Write(&buf, binary.BigEndian, uint64(headerSize+body.Len())))
	buf.Write(body.Bytes())
	buf.Write(meta)
	return buf.Bytes()
}

func collect(data []byte) ([]Entry, error) {
	var entries []Entry
	for entry, err := range All(data) {
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func TestAll(t *testing.T) {
	t.Parallel()

	mario := []byte{
		0x85, 0xf0, 0xdd, 0xdd, 0xfe, 0x4a, 0xb6, 0x7c,
		0x42, 0xab, 0xa4, 0x84, 0x98, 0xf4, 0x2f, 0xdc,
	}

	t.Run("decodes entries and ignores unknown fields", func(t *testing.T) {
		t.Parallel()
		data := build(t,
			map[string]any{
				"name":     "Super Mario Bros. 3 (USA)",
				"md5":      mario,
				"rom_name": "Super Mario Bros. 3 (USA).nes",
				"size":     393232,
				"genre":    "Platformer",
			},
			map[string]any{"name": "Metroid (USA)", "md5": mario},
		)

		entries, err := collect(data)
		require.NoError(t, err)
		require.Len(t, entries, 2)
		assert.Equal(t, "Super Mario Bros. 3 (USA)", entries[0].Name)
		assert.Equal(t, mario, entries[0].MD5)
		assert.Equal(t, "Metroid (USA)", entries[1].Name)
	})

	t.Run("entry without md5", func(t *testing.T) {
		t.Parallel()
		entries, err := collect(build(t, map[string]any{"name": "No Hash"}))
		require.NoError(t, err)
		require.Len(t, entries, 1)
		assert.Empty(t, entries[0].MD5)
	})

	t.Run("stops early when consumer breaks", func(t *testing.T) {
		t.Parallel()
		data := build(t,
			map[string]any{"name": "First", "md5": mario},
			map[string]any{"name": "Second", "md5": mario},
		)

		var count int
		for range All(data) {
			count++
			break
		}
		assert.Equal(t, 1, count)
	})

	t.Run("empty database", func(t *testing.T) {
		t.Parallel()
		entries, err := collect(build(t))
		require.NoError(t, err)
		assert.Empty(t, entries)
	})
}

func TestAllInvalid(t *testing.T) {
	t.Parallel()

	valid := build(t, map[string]any{"name": "Game", "md5": []byte{0x01}})

	badOffset := bytes.Clone(valid)
	binary.BigEndian.PutUint64(badOffset[8:headerSize], uint64(len(valid)+1))

	shortOffset := bytes.Clone(valid)
	binary.BigEndian.PutUint64(shortOffset[8:headerSize], 4)

	tests := []struct {
		name string
		data []byte
	}{
		{"bad magic", append([]byte("NOTANRDB"), valid[8:]...)},
		{"truncated header", valid[:8]},
		{"empty", nil},
		{"offset past end", badOffset},
		{"offset inside header", shortOffset},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := collect(tt.data)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrInvalid)
		})
	}
}
