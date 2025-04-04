package cartridge

import (
	"log/slog"

	"gabe565.com/gones/internal/consts"
	"gabe565.com/gones/internal/log"
)

func NewMapper71(cartridge *Cartridge) *Mapper71 {
	prgCount := uint(len(cartridge.PRG) / consts.PRGChunkSize)
	mapper := &Mapper71{
		cartridge: cartridge,
		PRGCount:  prgCount,
		PRGLast:   prgCount - 1,
	}
	return mapper
}

type Mapper71 struct {
	cartridge *Cartridge
	PRGCount  uint `msgpack:"alias:PrgCount"`
	PRGActive uint `msgpack:"alias:PrgActive"`
	PRGLast   uint `msgpack:"alias:PrgLast"`
}

func (m *Mapper71) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper71) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper71) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return m.cartridge.CHR[addr]
	case 0x8000 <= addr && addr < 0xC000:
		addr := uint(addr)
		addr -= 0x8000
		addr += m.PRGActive * consts.PRGChunkSize
		return m.cartridge.PRG[addr]
	case 0xC000 <= addr:
		addr := uint(addr)
		addr -= 0xC000
		addr += m.PRGLast * consts.PRGChunkSize
		return m.cartridge.PRG[addr]
	default:
		slog.Error("Invalid mapper 71 read", "addr", log.HexAddr(addr))
		return 0
	}
}

func (m *Mapper71) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		m.cartridge.CHR[addr] = data
	case 0x8000 <= addr && addr < 0x9000:
		// Ignored for compatibility
		// https://www.nesdev.org/wiki/INES_Mapper_071#Mirroring_($8000-$9FFF)
	case 0x9000 <= addr && addr < 0xA000:
		m.cartridge.Mirror = Mirror(data >> 4 & 1)
	case 0xC000 <= addr:
		data := uint(data & 0xF)
		data %= m.PRGCount
		m.PRGActive = data
	default:
		slog.Error("Invalid mapper 71 write", "addr", log.HexAddr(addr))
	}
}
