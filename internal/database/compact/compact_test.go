package compact

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"
)

func TestKey(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		hash    string
		want    uint64
		wantErr assert.ErrorAssertionFunc
	}{
		{"full md5", "85f0ddddfe4ab67c42aba48498f42fdc", 0x85f0ddddfe4ab67c, assert.NoError},
		{"exactly 8 bytes", "85f0ddddfe4ab67c", 0x85f0ddddfe4ab67c, assert.NoError},
		{"uppercase", "85F0DDDDFE4AB67C", 0x85f0ddddfe4ab67c, assert.NoError},
		{"too short", "85f0dddd", 0, assert.Error},
		{"empty hash", "", 0, assert.Error},
		{"not hex", "zzzzzzzzzzzzzzzz", 0, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Key(tt.hash)
			if !tt.wantErr(t, err) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFind(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		hashes []uint64
		names  []string
	}{
		{"no entries", nil, nil},
		{"single", []uint64{1}, []string{"Metroid (USA)"}},
		{
			"many",
			[]uint64{0x85f0ddddfe4ab67c, 0x397d10e475266ad2, 0},
			[]string{"Super Mario Bros. 3 (USA)", "Metroid (USA)", "Zero Hash"},
		},
		{"empty name", []uint64{1, 2}, []string{"", "After Empty"}},
		{"unicode", []uint64{1, 2}, []string{"Bishōjo Senshi Sailor Moon (Japan)", "Next"}},
		{"comma and quotes", []uint64{1}, []string{`Wow, "Quoted" (USA)`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data, err := Marshal(tt.hashes, tt.names)
			require.NoError(t, err)

			// every entry is findable at its own index
			for i, hash := range tt.hashes {
				got, ok, err := Find(bytes.NewReader(data), hash)
				require.NoError(t, err)
				require.True(t, ok, "hash at index %d not found", i)
				assert.Equal(t, tt.names[i], got)
			}

			// a hash that is not present misses cleanly
			got, ok, err := Find(bytes.NewReader(data), 0xdeadbeefcafebabe)
			require.NoError(t, err)
			assert.False(t, ok)
			assert.Empty(t, got)
		})
	}
}

// TestFindDoesNotReadNamesOnMiss checks that a lookup which misses stops after
// the hash block. The database has to be much larger than the decoder's read
// buffer for the difference to be observable.
func TestFindDoesNotReadNamesOnMiss(t *testing.T) {
	t.Parallel()

	const count = 20000
	hashes := make([]uint64, count)
	names := make([]string, count)
	for i := range count {
		hashes[i] = uint64(i) + 1
		names[i] = fmt.Sprintf("A Reasonably Long Game Title Number %d (USA)", i)
	}

	data, err := Marshal(hashes, names)
	require.NoError(t, err)

	miss := &countingReader{data: data}
	_, ok, err := Find(miss, 0xdeadbeefcafebabe)
	require.NoError(t, err)
	require.False(t, ok)

	hit := &countingReader{data: data}
	name, ok, err := Find(hit, hashes[count-1])
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, names[count-1], name)

	assert.Less(t, miss.read, len(data)/2,
		"a miss should stop after the hash block")
	assert.Less(t, miss.read, hit.read,
		"a miss should read less than a hit")
}

type countingReader struct {
	data []byte
	pos  int
	read int
}

func (c *countingReader) Read(p []byte) (int, error) {
	if c.pos >= len(c.data) {
		return 0, io.EOF
	}
	n := copy(p, c.data[c.pos:])
	c.pos += n
	c.read += n
	return n, nil
}

func TestMarshalInvalid(t *testing.T) {
	t.Parallel()

	t.Run("length mismatch", func(t *testing.T) {
		t.Parallel()
		_, err := Marshal([]uint64{1, 2}, []string{"only one"})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalid)
	})

	t.Run("name with newline", func(t *testing.T) {
		t.Parallel()
		_, err := Marshal([]uint64{1}, []string{"two\nlines"})
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalid)
	})
}

func TestFindInvalid(t *testing.T) {
	t.Parallel()

	notAnArray, err := msgpack.Marshal("just a string")
	require.NoError(t, err)

	wrongFieldCount, err := msgpack.Marshal([]any{1, 2, 3})
	require.NoError(t, err)

	raggedHashes, err := msgpack.Marshal([]any{[]byte{1, 2, 3}, "One"})
	require.NoError(t, err)

	// Two hashes but a single name, so resolving the second one has no name.
	twoHashes := make([]byte, 2*HashLen)
	binary.BigEndian.PutUint64(twoHashes, 1)
	binary.BigEndian.PutUint64(twoHashes[HashLen:], 2)
	tooFewNames, err := msgpack.Marshal([]any{twoHashes, "only one"})
	require.NoError(t, err)

	tests := []struct {
		name string
		data []byte
		key  uint64
	}{
		{"empty data", nil, 0},
		{"not an array", notAnArray, 0},
		{"wrong field count", wrongFieldCount, 0},
		{"ragged hash block", raggedHashes, 0},
		{"fewer names than hashes", tooFewNames, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, _, err := Find(bytes.NewReader(tt.data), tt.key)
			require.Error(t, err)
		})
	}
}
