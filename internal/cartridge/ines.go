package cartridge

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/database"
	log "github.com/sirupsen/logrus"
)

type iNESFileHeader struct {
	Magic    [4]byte
	PrgCount byte
	ChrCount byte
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

	cartridge, err := FromiNes(f)

	return cartridge, err
}

func FromiNes(r io.ReadSeeker) (*Cartridge, error) {
	var header iNESFileHeader
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	if header.Magic != iNesMagic {
		return nil, ErrInvalidRom
	}

	cartridge := New()
	cartridge.Mapper = header.Mapper()
	cartridge.Mirror = header.Mirror()
	cartridge.Battery = header.Battery()

	if header.NESv2() {
		cartridge.Submapper = header.Submapper()
	}

	cartridge.prg = make([]byte, int(header.PrgCount)*consts.PrgChunkSize)
	if _, err := io.ReadFull(r, cartridge.prg); err != nil {
		return nil, err
	}

	cartridge.Chr = make([]byte, int(header.ChrCount)*consts.ChrChunkSize)
	if _, err := io.ReadFull(r, cartridge.Chr); err != nil {
		return nil, err
	}

	if header.ChrCount == 0 {
		cartridge.Chr = make([]byte, consts.ChrChunkSize)
	}

	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return cartridge, err
	}
	md5 := md5.New()
	if _, err := io.Copy(md5, r); err != nil {
		return cartridge, err
	}
	cartridge.hash = fmt.Sprintf("%x", md5.Sum(nil))
	if cartridge.hash != "" {
		cartridge.name, _ = database.FindNameByHash(cartridge.hash)
	}

	log.WithFields(log.Fields{
		"battery": cartridge.Battery,
		"mapper":  cartridge.Mapper,
		"mirror":  cartridge.Mirror,
		"prg":     header.PrgCount,
		"chr":     header.ChrCount,
	}).Debug("Cartridge header info")

	return cartridge, nil
}
