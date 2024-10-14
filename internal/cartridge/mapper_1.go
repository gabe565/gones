package cartridge

import (
	"log/slog"

	"gabe565.com/gones/internal/consts"
	"gabe565.com/gones/internal/log"
)

func NewMapper1(cartridge *Cartridge) *Mapper1 {
	mapper := &Mapper1{
		cartridge:     cartridge,
		ShiftRegister: 0x10,
	}
	mapper.PRGOffsets[1] = mapper.prgBankOffset(-1)
	return mapper
}

type Mapper1 struct {
	cartridge     *Cartridge
	ShiftRegister byte
	Control       byte
	PRGMode       byte   `msgpack:"alias:PrgMode"`
	CHRMode       bool   `msgpack:"alias:ChrMode"`
	PRGBank       byte   `msgpack:"alias:PrgBank"`
	CHRBank0      byte   `msgpack:"alias:ChrBank0"`
	CHRBank1      byte   `msgpack:"alias:ChrBank1"`
	PRGOffsets    [2]int `msgpack:"alias:PrgOffsets"`
	CHROffsets    [2]int `msgpack:"alias:ChrOffsets"`
}

func (m *Mapper1) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper1) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper1) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		bank := addr / 0x1000
		offset := int(addr % 0x1000)
		return m.cartridge.CHR[m.CHROffsets[bank]+offset]
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		return m.cartridge.SRAM[addr]
	case 0x8000 <= addr:
		addr -= 0x8000
		bank := addr / consts.PRGChunkSize
		offset := int(addr % consts.PRGChunkSize)
		return m.cartridge.prg[m.PRGOffsets[bank]+offset]
	default:
		slog.Error("Invalid mapper 1 read", "addr", log.HexAddr(addr))
		return 0
	}
}

func (m *Mapper1) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		bank := addr / 0x1000
		offset := int(addr % 0x1000)
		m.cartridge.CHR[m.CHROffsets[bank]+offset] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		m.cartridge.SRAM[addr] = data
	case 0x8000 <= addr:
		if data>>7&1 == 1 {
			m.ShiftRegister = 0x10
			m.writeControl(m.Control | 0x0C)
		} else {
			complete := m.ShiftRegister&1 == 1
			m.ShiftRegister >>= 1
			m.ShiftRegister |= data & 1 << 4
			if complete {
				data := m.ShiftRegister
				switch {
				case addr < 0xA000:
					m.writeControl(data)
				case 0xA000 <= addr && addr < 0xC000:
					m.CHRBank0 = data
					m.updateOffsets()
				case 0xC000 <= addr && addr < 0xE000:
					m.CHRBank1 = data
					m.updateOffsets()
				case 0xE000 <= addr:
					m.PRGBank = data & 0xF
					m.updateOffsets()
				}
				m.ShiftRegister = 0x10
			}
		}
	default:
		slog.Error("Invalid mapper 1 write", "addr", log.HexAddr(addr))
	}
}

func (m *Mapper1) writeControl(data byte) {
	m.Control = data
	m.CHRMode = data>>4&1 == 1
	m.PRGMode = data >> 2 & 3
	switch data & 3 {
	case 0:
		m.cartridge.Mirror = SingleLower
	case 1:
		m.cartridge.Mirror = SingleUpper
	case 2:
		m.cartridge.Mirror = Vertical
	case 3:
		m.cartridge.Mirror = Horizontal
	}
	m.updateOffsets()
}

func (m *Mapper1) prgBankOffset(i int) int {
	if i >= 0x80 {
		i -= 0x100
	}
	i %= len(m.cartridge.prg) / consts.PRGChunkSize
	offset := i * consts.PRGChunkSize
	if offset < 0 {
		offset += len(m.cartridge.prg)
	}
	return offset
}

func (m *Mapper1) chrBankOffset(i int) int {
	if i >= 0x80 {
		i -= 0x100
	}
	i %= len(m.cartridge.CHR) / 0x1000
	offset := i * 0x1000
	if offset < 0 {
		offset += len(m.cartridge.CHR)
	}
	return offset
}

func (m *Mapper1) updateOffsets() {
	switch m.PRGMode {
	case 0, 1:
		m.PRGOffsets[0] = m.prgBankOffset(int(m.PRGBank & 0xFE))
		m.PRGOffsets[1] = m.prgBankOffset(int(m.PRGBank | 0x01))
	case 2:
		m.PRGOffsets[0] = 0
		m.PRGOffsets[1] = m.prgBankOffset(int(m.PRGBank))
	case 3:
		m.PRGOffsets[0] = m.prgBankOffset(int(m.PRGBank))
		m.PRGOffsets[1] = m.prgBankOffset(-1)
	}

	if m.CHRMode {
		m.CHROffsets[0] = m.chrBankOffset(int(m.CHRBank0))
		m.CHROffsets[1] = m.chrBankOffset(int(m.CHRBank1))
	} else {
		m.CHROffsets[0] = m.chrBankOffset(int(m.CHRBank0 & 0xFE))
		m.CHROffsets[1] = m.chrBankOffset(int(m.CHRBank0 | 0x01))
	}
}
