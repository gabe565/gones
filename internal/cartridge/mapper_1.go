package cartridge

import (
	"encoding/gob"

	"github.com/gabe565/gones/internal/consts"
	log "github.com/sirupsen/logrus"
)

func NewMapper1(cartridge *Cartridge) Mapper {
	mapper := &Mapper1{
		cartridge:     cartridge,
		ShiftRegister: 0x10,
	}
	mapper.PrgOffsets[1] = mapper.prgBankOffset(-1)
	gob.Register(mapper)
	return mapper
}

type Mapper1 struct {
	cartridge     *Cartridge
	ShiftRegister byte
	Control       byte
	PrgMode       byte
	ChrMode       bool
	PrgBank       byte
	ChrBank0      byte
	ChrBank1      byte
	PrgOffsets    [2]int
	ChrOffsets    [2]int
}

func (m *Mapper1) Step(_ bool, _ uint16, _ uint) {}

func (m *Mapper1) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper1) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper1) SetCpu(_ CPU) {}

func (m *Mapper1) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		bank := addr / 0x1000
		offset := int(addr % 0x1000)
		return m.cartridge.Chr[m.ChrOffsets[bank]+offset]
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		return m.cartridge.Sram[addr]
	case 0x8000 <= addr:
		addr -= 0x8000
		bank := addr / consts.PrgChunkSize
		offset := int(addr % consts.PrgChunkSize)
		return m.cartridge.prg[m.PrgOffsets[bank]+offset]
	default:
		log.Fatalf("invalid mapper 1 read from $%04X", addr)
		return 0
	}
}

func (m *Mapper1) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		bank := addr / 0x1000
		offset := int(addr % 0x1000)
		m.cartridge.Chr[m.ChrOffsets[bank]+offset] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		m.cartridge.Sram[addr] = data
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
					m.ChrBank0 = data
					m.updateOffsets()
				case 0xC000 <= addr && addr < 0xE000:
					m.ChrBank1 = data
					m.updateOffsets()
				case 0xE000 <= addr:
					m.PrgBank = data & 0xF
					m.updateOffsets()
				}
				m.ShiftRegister = 0x10
			}
		}
	default:
		log.Fatalf("invalid mapper 1 write to $%04X", addr)
	}
}

func (m *Mapper1) writeControl(data byte) {
	m.Control = data
	m.ChrMode = data>>4&1 == 1
	m.PrgMode = data >> 2 & 3
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
	i %= len(m.cartridge.prg) / consts.PrgChunkSize
	offset := i * consts.PrgChunkSize
	if offset < 0 {
		offset += len(m.cartridge.prg)
	}
	return offset
}

func (m *Mapper1) chrBankOffset(i int) int {
	if i >= 0x80 {
		i -= 0x100
	}
	i %= len(m.cartridge.Chr) / 0x1000
	offset := i * 0x1000
	if offset < 0 {
		offset += len(m.cartridge.Chr)
	}
	return offset
}

func (m *Mapper1) updateOffsets() {
	switch m.PrgMode {
	case 0, 1:
		m.PrgOffsets[0] = m.prgBankOffset(int(m.PrgBank & 0xFE))
		m.PrgOffsets[1] = m.prgBankOffset(int(m.PrgBank | 0x01))
	case 2:
		m.PrgOffsets[0] = 0
		m.PrgOffsets[1] = m.prgBankOffset(int(m.PrgBank))
	case 3:
		m.PrgOffsets[0] = m.prgBankOffset(int(m.PrgBank))
		m.PrgOffsets[1] = m.prgBankOffset(-1)
	}

	if m.ChrMode {
		m.ChrOffsets[0] = m.chrBankOffset(int(m.ChrBank0))
		m.ChrOffsets[1] = m.chrBankOffset(int(m.ChrBank1))
	} else {
		m.ChrOffsets[0] = m.chrBankOffset(int(m.ChrBank0 & 0xFE))
		m.ChrOffsets[1] = m.chrBankOffset(int(m.ChrBank0 | 0x01))
	}
}
