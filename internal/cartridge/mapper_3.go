package cartridge

import (
	"encoding/gob"

	"github.com/gabe565/gones/internal/consts"
	log "github.com/sirupsen/logrus"
)

func NewMapper3(cartridge *Cartridge) Mapper {
	prgBanks := uint(len(cartridge.prg) / consts.PrgChunkSize)
	mapper := &Mapper3{
		cartridge: cartridge,
		PrgBank2:  prgBanks - 1,
	}
	gob.Register(mapper)
	return mapper
}

type Mapper3 struct {
	cartridge *Cartridge
	ChrBank   uint
	PrgBank1  uint
	PrgBank2  uint
}

func (m *Mapper3) Step(_ bool, _ uint16, _ uint) {}

func (m *Mapper3) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper3) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper3) SetCpu(_ CPU) {}

func (m *Mapper3) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		addr := uint(addr)
		addr += m.ChrBank * 0x2000
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
		log.Fatalf("invalid mapper 3 read from $%04X", addr)
		return 0
	}
}

func (m *Mapper3) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		addr := uint(addr)
		addr += m.ChrBank * 0x2000
		m.cartridge.Chr[addr] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		m.cartridge.Sram[addr] = data
	case 0x8000 <= addr:
		m.ChrBank = uint(data & 3)
	default:
		log.Fatalf("invalid mapper 3 write to $%04X", addr)
	}
}
