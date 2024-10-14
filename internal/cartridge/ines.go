package cartridge

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"gabe565.com/gones/internal/consts"
	"gabe565.com/gones/internal/database"
)

type INESFileHeader struct {
	Magic    [4]byte
	PRGCount byte
	CHRCount byte
	Control  [3]byte
	_        [7]byte
}

func (i INESFileHeader) Mapper() uint8 {
	return i.Control[1]&0xF0 | i.Control[0]>>4
}

func (i INESFileHeader) Mirror() Mirror {
	if i.Control[0]&0x8 != 0 {
		return FourScreen
	}
	return Mirror(i.Control[0] & 1)
}

func (i INESFileHeader) Battery() bool {
	return i.Control[0]&0x2 != 0
}

func (i INESFileHeader) NESv2() bool {
	return i.Control[1]&0xC == 0x8
}

func (i INESFileHeader) Submapper() uint8 {
	if i.NESv2() {
		return i.Control[2] >> 4
	}
	return 0
}

var ErrInvalidROM = errors.New("invalid ROM file")

func FromINESFile(path string) (*Cartridge, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	cartridge, err := FromINES(f)
	if err != nil {
		return nil, err
	}

	if cartridge.name == "" {
		cartridge.name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}
	return cartridge, nil
}

func FromINES(r io.Reader) (*Cartridge, error) {
	hasher := md5.New()
	tr := io.TeeReader(r, hasher)

	var header INESFileHeader
	if err := binary.Read(tr, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	if header.Magic != [4]byte{'N', 'E', 'S', 0x1A} {
		return nil, fmt.Errorf("%w: %s", ErrInvalidROM, "missing NES header")
	}

	cartridge := New()
	cartridge.Header = header
	cartridge.Mirror = header.Mirror()
	cartridge.Battery = header.Battery()

	slog.Debug("Loaded iNES header",
		"battery", cartridge.Battery,
		"mapper", header.Mapper(),
		"mirror", cartridge.Mirror,
		"prg", header.PRGCount,
		"chr", header.CHRCount,
	)

	cartridge.PRG = make([]byte, int(header.PRGCount)*consts.PRGChunkSize)
	if _, err := io.ReadFull(tr, cartridge.PRG); err != nil {
		return nil, err
	}

	if header.CHRCount == 0 {
		cartridge.CHR = make([]byte, consts.CHRChunkSize)
	} else {
		cartridge.CHR = make([]byte, int(header.CHRCount)*consts.CHRChunkSize)
		if _, err := io.ReadFull(tr, cartridge.CHR); err != nil {
			return nil, err
		}
	}

	// Ensure all bytes are written to hasher
	if _, err := io.Copy(hasher, r); err != nil {
		return nil, err
	}

	cartridge.hash = hex.EncodeToString(hasher.Sum(nil))
	cartridge.name, _ = database.FindNameByHash(cartridge.hash)
	return cartridge, nil
}
