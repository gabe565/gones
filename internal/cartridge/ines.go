package cartridge

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/database"
)

type iNESFileHeader struct {
	Magic    [4]byte
	PRGCount byte
	CHRCount byte
	Control  [3]byte
	_        [7]byte
}

func (i iNESFileHeader) Mapper() byte {
	return i.Control[1]&0xF0 | i.Control[0]>>4
}

func (i iNESFileHeader) Mirror() Mirror {
	if i.Control[0]&0x8 != 0 {
		return FourScreen
	}
	return Mirror(i.Control[0] & 1)
}

func (i iNESFileHeader) Battery() bool {
	return i.Control[0]&0x2 != 0
}

func (i iNESFileHeader) NESv2() bool {
	return i.Control[1]&0xC == 0x8
}

func (i iNESFileHeader) Submapper() byte {
	return i.Control[2] >> 4
}

var ErrInvalidRom = errors.New("invalid ROM")

func FromiNesFile(path string) (*Cartridge, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	cartridge, err := FromiNes(f)
	if cartridge.name == "" {
		cartridge.name = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	return cartridge, err
}

func FromiNes(r io.ReadSeeker) (*Cartridge, error) {
	var header iNESFileHeader
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	if header.Magic != [4]byte{'N', 'E', 'S', 0x1A} {
		return nil, ErrInvalidRom
	}

	cartridge := New()
	cartridge.Mapper = header.Mapper()
	cartridge.Mirror = header.Mirror()
	cartridge.Battery = header.Battery()

	if header.NESv2() {
		cartridge.Submapper = header.Submapper()
	}

	slog.Debug("Loaded iNES header",
		"battery", cartridge.Battery,
		"mapper", cartridge.Mapper,
		"mirror", cartridge.Mirror,
		"prg", header.PRGCount,
		"chr", header.CHRCount,
	)

	cartridge.prg = make([]byte, int(header.PRGCount)*consts.PRGChunkSize)
	if _, err := io.ReadFull(r, cartridge.prg); err != nil {
		return nil, err
	}

	cartridge.CHR = make([]byte, int(header.CHRCount)*consts.CHRChunkSize)
	if _, err := io.ReadFull(r, cartridge.CHR); err != nil {
		return nil, err
	}

	if header.CHRCount == 0 {
		cartridge.CHR = make([]byte, consts.CHRChunkSize)
	}

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return cartridge, err
	}
	md5 := md5.New()
	if _, err := io.Copy(md5, r); err != nil {
		return cartridge, err
	}
	cartridge.hash = hex.EncodeToString(md5.Sum(nil))
	if cartridge.hash != "" {
		cartridge.name, _ = database.FindNameByHash(cartridge.hash)
	}

	return cartridge, nil
}
