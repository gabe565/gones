package cartridge

import (
	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/util"
	"github.com/rs/zerolog/log"
)

func NewMapper7(cartridge *Cartridge) *Mapper7 {
	mapper := &Mapper7{cartridge: cartridge}
	return mapper
}

type Mapper7 struct {
	cartridge *Cartridge
	PRGBank   uint `msgpack:"alias:PrgBank"`
}

func (m *Mapper7) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper7) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper7) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return m.cartridge.CHR[addr&0x1FFF]
	case 0x6000 <= addr && addr < 0x8000:
		addr := uint(addr)
		addr -= 0x6000
		return m.cartridge.SRAM[addr]
	case 0x8000 <= addr:
		addr := uint(addr)
		addr -= 0x8000
		addr += m.PRGBank * 2 * consts.PRGChunkSize
		addr %= uint(len(m.cartridge.prg))
		return m.cartridge.prg[addr]
	default:
		log.Error().Str("addr", util.EncodeHexAddr(addr)).Msg("Invalid mapper 7 read")
		return 0
	}
}

func (m *Mapper7) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		m.cartridge.CHR[addr%0x1FFF] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr := uint(addr)
		addr -= 0x6000
		m.cartridge.SRAM[addr] = data
	case 0x8000 <= addr:
		switch data >> 4 & 1 {
		case 0:
			m.cartridge.Mirror = SingleLower
		case 1:
			m.cartridge.Mirror = SingleUpper
		}
		m.PRGBank = uint(data & 7)
	default:
		log.Error().Str("addr", util.EncodeHexAddr(addr)).Msg("Invalid mapper 7 write")
	}
}
