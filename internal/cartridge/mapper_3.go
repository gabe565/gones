package cartridge

import (
	"github.com/gabe565/gones/internal/consts"
	log "github.com/sirupsen/logrus"
)

func NewMapper3(cartridge *Cartridge) *Mapper3 {
	prgBanks := uint(len(cartridge.prg) / consts.PRGChunkSize)
	mapper := &Mapper3{
		cartridge: cartridge,
		PRGBank2:  prgBanks - 1,
	}
	return mapper
}

type Mapper3 struct {
	cartridge *Cartridge
	CHRBank   uint `msgpack:"alias:ChrBank"`
	PRGBank1  uint `msgpack:"alias:PrgBank1"`
	PRGBank2  uint `msgpack:"alias:PrgBank2"`
}

func (m *Mapper3) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper3) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper3) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		addr := uint(addr)
		addr += m.CHRBank * 0x2000
		return m.cartridge.CHR[addr]
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		return m.cartridge.SRAM[addr]
	case 0x8000 <= addr && addr < 0xC000:
		addr := uint(addr)
		addr -= 0x8000
		addr += m.PRGBank1 * consts.PRGChunkSize
		return m.cartridge.prg[addr]
	case 0xC000 <= addr:
		addr := uint(addr)
		addr -= 0xC000
		addr += m.PRGBank2 * consts.PRGChunkSize
		return m.cartridge.prg[addr]
	default:
		log.Warnf("invalid mapper 3 read from $%04X", addr)
		return 0
	}
}

func (m *Mapper3) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		addr := uint(addr)
		addr += m.CHRBank * 0x2000
		m.cartridge.CHR[addr] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		m.cartridge.SRAM[addr] = data
	case 0x8000 <= addr:
		m.CHRBank = uint(data & 3)
	default:
		log.Warnf("invalid mapper 3 write to $%04X", addr)
	}
}
