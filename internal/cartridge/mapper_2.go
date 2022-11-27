package cartridge

import (
	"encoding/gob"
	"github.com/gabe565/gones/internal/consts"
	log "github.com/sirupsen/logrus"
)

func NewMapper2(cartridge *Cartridge) Mapper {
	prgBanks := uint16(len(cartridge.prg) / consts.PrgChunkSize)
	mapper := &Mapper2{
		cartridge: cartridge,
		PrgBanks:  uint(prgBanks),
		PrgBank2:  uint(prgBanks - 1),
	}
	gob.Register(mapper)
	return mapper
}

type Mapper2 struct {
	cartridge *Cartridge
	PrgBanks  uint
	PrgBank1  uint
	PrgBank2  uint
}

func (m *Mapper2) Step() {}

func (m *Mapper2) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper2) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return m.cartridge.Chr[addr]
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		return m.cartridge.Sram[addr]
	case 0x8000 <= addr && addr < 0xC000:
		addr := uint(addr)
		addr -= 0x8000
		addr += m.PrgBank1 * consts.PrgChunkSize
		return m.cartridge.prg[addr]
	case 0xC000 <= addr:
		addr := uint(addr)
		addr -= 0xC000
		addr += m.PrgBank2 * consts.PrgChunkSize
		return m.cartridge.prg[addr]
	default:
		log.Fatalf("invalid mapper 2 read from $%04X", addr)
		return 0
	}
}

func (m *Mapper2) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		m.cartridge.Chr[addr] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr := uint(addr)
		addr -= 0x6000
		m.cartridge.Sram[addr] = data
	case 0x8000 <= addr:
		addr := uint(addr)
		addr %= m.PrgBanks
		m.PrgBank1 = addr
	default:
		log.Fatalf("invalid mapper 2 write to $%04X", addr)
	}
}
