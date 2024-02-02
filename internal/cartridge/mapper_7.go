package cartridge

import (
	"github.com/gabe565/gones/internal/consts"
	log "github.com/sirupsen/logrus"
)

func NewMapper7(cartridge *Cartridge) Mapper {
	mapper := &Mapper7{cartridge: cartridge}
	return mapper
}

type Mapper7 struct {
	cartridge *Cartridge
	PrgBank   uint
}

func (m *Mapper7) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper7) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper7) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return m.cartridge.Chr[addr&0x1FFF]
	case 0x6000 <= addr && addr < 0x8000:
		addr := uint(addr)
		addr -= 0x6000
		return m.cartridge.Sram[addr]
	case 0x8000 <= addr:
		addr := uint(addr)
		addr -= 0x8000
		addr += m.PrgBank * 2 * consts.PrgChunkSize
		addr %= uint(len(m.cartridge.prg))
		return m.cartridge.prg[addr]
	default:
		log.Warnf("invalid mapper 7 read from $%04X", addr)
		return 0
	}
}

func (m *Mapper7) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		m.cartridge.Chr[addr%0x1FFF] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr := uint(addr)
		addr -= 0x6000
		m.cartridge.Sram[addr] = data
	case 0x8000 <= addr:
		switch data >> 4 & 1 {
		case 0:
			m.cartridge.Mirror = SingleLower
		case 1:
			m.cartridge.Mirror = SingleUpper
		}
		m.PrgBank = uint(data & 7)
	default:
		log.Warnf("invalid mapper 7 write to $%04X", addr)
	}
}
