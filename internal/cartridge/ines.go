package cartridge

import (
	"encoding/binary"
	"errors"
	"github.com/gabe565/gones/internal/consts"
	"io"
	"os"
)

type iNESFileHeader struct {
	Magic    [4]byte
	PrgCount byte
	ChrCount byte
	Control1 byte
	Control2 byte
	RAMCount byte
	_        [7]byte
}

var iNesMagic = [4]byte{'N', 'E', 'S', 0x1A}

var ErrInvalidRom = errors.New("invalid ROM")

func FromiNesFile(path string) (*Cartridge, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	return FromiNes(f)
}

func FromiNes(r io.Reader) (*Cartridge, error) {
	var header iNESFileHeader
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	if header.Magic != iNesMagic {
		return nil, ErrInvalidRom
	}

	cartridge := New()

	mapper1 := header.Control1 >> 4
	mapper2 := header.Control2 >> 4
	cartridge.Mapper = mapper1 | mapper2<<4

	mirror1 := header.Control1 & 1
	mirror2 := (header.Control1 >> 3) & 1
	cartridge.Mirror = Mirror(mirror2<<1 | mirror1)

	cartridge.Battery = (header.Control1>>1)&1 == 1

	cartridge.Prg = make([]byte, int(header.PrgCount)*consts.PrgChunkSize)
	if _, err := io.ReadFull(r, cartridge.Prg); err != nil {
		return nil, err
	}

	cartridge.Chr = make([]byte, int(header.ChrCount)*consts.ChrChunkSize)
	if _, err := io.ReadFull(r, cartridge.Chr); err != nil {
		return nil, err
	}

	if header.ChrCount == 0 {
		cartridge.Chr = make([]byte, consts.ChrChunkSize)
	}

	return cartridge, nil
}
