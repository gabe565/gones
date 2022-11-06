package cartridge

import (
	"encoding/binary"
	"errors"
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

var iNesMagic = [4]byte{0x4E, 0x45, 0x53, 0x1A}

var (
	ErrInvalidRom = errors.New("invalid ROM")
	ErrNES2       = errors.New("NES2.0 format is not supported")
)

func FromiNes(path string) (*Cartridge, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var header iNESFileHeader
	if err := binary.Read(f, binary.LittleEndian, &header); err != nil {
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

	cartridge.Battery = (header.Control1 >> 1) & 1

	cartridge.Prg = make([]byte, int(header.PrgCount)*16384)
	if _, err := io.ReadFull(f, cartridge.Prg); err != nil {
		return nil, err
	}

	cartridge.Chr = make([]byte, int(header.ChrCount)*8192)
	if _, err := io.ReadFull(f, cartridge.Chr); err != nil {
		return nil, err
	}

	if header.ChrCount == 0 {
		cartridge.Chr = make([]byte, 8192)
	}

	return cartridge, nil
}
