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
	Control  [2]byte
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
	cartridge.Mapper = getMapper(header.Control)
	cartridge.Mirror = getMirror(header.Control[0])
	cartridge.Battery = hasBattery(header.Control[0])

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

func getMapper(data [2]byte) byte {
	mapper1 := data[0] >> 4
	mapper2 := data[1] >> 4
	return mapper2<<4 | mapper1
}

func getMirror(data byte) Mirror {
	mirror1 := data & 1
	mirror2 := (data >> 3) & 1
	return Mirror(mirror2<<1 | mirror1)
}

func hasBattery(data byte) bool {
	return (data>>1)&1 == 1
}
