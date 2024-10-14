package cartridge

import (
	"log/slog"

	"gabe565.com/gones/internal/consts"
	"gabe565.com/gones/internal/log"
)

func NewMapper2(cartridge *Cartridge) *Mapper2 {
	prgBanks := uint(len(cartridge.PRG) / consts.PRGChunkSize)
	mapper := &Mapper2{
		cartridge: cartridge,
		PRGBanks:  prgBanks,
		PRGBank2:  prgBanks - 1,
	}
	return mapper
}

type Mapper2 struct {
	cartridge *Cartridge
	PRGBanks  uint `msgpack:"alias:PrgBanks"`
	PRGBank1  uint `msgpack:"alias:PrgBank1"`
	PRGBank2  uint `msgpack:"alias:PrgBank2"`
}

func (m *Mapper2) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper2) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper2) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return m.cartridge.CHR[addr]
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		return m.cartridge.SRAM[addr]
	case 0x8000 <= addr && addr < 0xC000:
		addr := uint(addr)
		addr -= 0x8000
		addr += m.PRGBank1 * consts.PRGChunkSize
		return m.cartridge.PRG[addr]
	case 0xC000 <= addr:
		addr := uint(addr)
		addr -= 0xC000
		addr += m.PRGBank2 * consts.PRGChunkSize
		return m.cartridge.PRG[addr]
	default:
		slog.Error("Invalid mapper 2 read", "addr", log.HexAddr(addr))
		return 0
	}
}

func (m *Mapper2) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		m.cartridge.CHR[addr] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		m.cartridge.SRAM[addr] = data
	case 0x8000 <= addr:
		data := uint(data)
		data %= m.PRGBanks
		m.PRGBank1 = data
	default:
		slog.Error("Invalid mapper 2 write", "addr", log.HexAddr(addr))
	}
}
