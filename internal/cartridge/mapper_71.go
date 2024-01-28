package cartridge

import (
	"github.com/gabe565/gones/internal/consts"
	log "github.com/sirupsen/logrus"
)

func NewMapper71(cartridge *Cartridge) Mapper {
	prgCount := uint(len(cartridge.prg) / consts.PrgChunkSize)
	mapper := &Mapper71{
		cartridge: cartridge,
		PrgCount:  prgCount,
		PrgLast:   prgCount - 1,
	}
	return mapper
}

type Mapper71 struct {
	cartridge *Cartridge
	PrgCount  uint
	PrgActive uint
	PrgLast   uint
}

func (m *Mapper71) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper71) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper71) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return m.cartridge.Chr[addr]
	case 0x8000 <= addr && addr < 0xC000:
		addr := uint(addr)
		addr -= 0x8000
		addr += m.PrgActive * consts.PrgChunkSize
		return m.cartridge.prg[addr]
	case 0xC000 <= addr:
		addr := uint(addr)
		addr -= 0xC000
		addr += m.PrgLast * consts.PrgChunkSize
		return m.cartridge.prg[addr]
	default:
		log.Warnf("invalid mapper 71 read from $%04X", addr)
		return 0
	}
}

func (m *Mapper71) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		m.cartridge.Chr[addr] = data
	case 0x8000 <= addr && addr < 0x9000:
		// Ignored for compatibility
		// https://www.nesdev.org/wiki/INES_Mapper_071#Mirroring_($8000-$9FFF)
	case 0x9000 <= addr && addr < 0xA000:
		m.cartridge.Mirror = Mirror(data >> 4 & 1)
	case 0xC000 <= addr:
		data := uint(data & 0xF)
		data %= m.PrgCount
		m.PrgActive = data
	default:
		log.Warnf("invalid mapper 71 write to $%04X", addr)
	}
}
