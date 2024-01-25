package cartridge

import (
	"encoding/gob"

	log "github.com/sirupsen/logrus"
)

func NewMapper69(cartridge *Cartridge) Mapper {
	prgCount := len(cartridge.prg) / 0x2000
	mapper := &Mapper69{
		cartridge: cartridge,
		PrgCount:  byte(prgCount),
		PrgBanks:  [5]int{0, 0, 0, 0, prgCount - 1},
	}
	gob.Register(mapper)
	return mapper
}

type Mapper69 struct {
	cartridge *Cartridge
	cpu       CPU
	Command   byte

	PrgCount   byte
	PrgBanks   [5]int
	RamSelect  bool
	RamEnabled bool

	ChrBanks [8]int

	IrqEnable        bool
	IrqCounterEnable bool
	IrqCounter       uint16
}

func (m *Mapper69) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper69) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper69) SetCpu(c CPU) { m.cpu = c }

func (m *Mapper69) OnCPUStep() {
	if m.IrqCounterEnable {
		m.IrqCounter -= 1
		if m.IrqEnable && m.IrqCounter == 0xFFFF {
			m.cpu.AddIrq()
		}
	}
}

func (m *Mapper69) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		// CHR
		bank := addr / 0x400
		offset := int(addr % 0x400)
		addr := m.ChrBanks[bank]*0x400 + offset
		return m.cartridge.Chr[addr%len(m.cartridge.Chr)]
	case 0x6000 <= addr && addr < 0x8000:
		// PRG/SRAM banks
		if m.RamSelect {
			if m.RamEnabled {
				return m.cartridge.Sram[addr-0x6000]
			} else {
				// open bus
				return 0
			}
		}
		fallthrough
	case 0x8000 <= addr:
		addr -= 0x6000
		bank := addr / 0x2000
		offset := int(addr % 0x2000)
		addr := m.PrgBanks[bank]*0x2000 + offset
		return m.cartridge.prg[addr%len(m.cartridge.prg)]
	default:
		log.Warnf("invalid mapper 69 read from $%04X", addr)
		return 0
	}
}

func (m *Mapper69) WriteMem(addr uint16, data byte) {
	switch {
	case 0x6000 <= addr && addr < 0x8000:
		// SRAM register
		if m.RamSelect && m.RamEnabled {
			m.cartridge.Sram[addr-0x6000] = data
		}
	case 0x8000 <= addr && addr < 0xA000:
		// Command register
		m.Command = data & 0xF
	case 0xA000 <= addr && addr < 0xC000:
		// Parameter register
		m.runCommand(data)
	default:
		log.Warnf("invalid mapper 69 write to $%04X", addr)
	}
}

func (m *Mapper69) runCommand(data byte) {
	switch m.Command {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		// CHR bank switch
		bank := m.Command
		m.ChrBanks[bank] = int(data)
	case 0x8:
		// PRG bank switch (with RAM)
		m.RamSelect = data>>6&1 == 1
		m.RamEnabled = data>>7&1 == 1
		fallthrough
	case 0x9, 0xA, 0xB:
		// PRG bank switch
		bank := m.Command - 0x8
		m.PrgBanks[bank] = int(data & 0x1F)
	case 0xC:
		// Nametable Mirroring
		m.cartridge.Mirror = Mirror(data & 0x3)
	case 0xD:
		// IRQ control
		m.IrqEnable = data&1 == 1
		m.IrqCounterEnable = data>>7&1 == 1
	case 0xE:
		// IRQ counter LO
		m.IrqCounter = m.IrqCounter&0xFF00 | uint16(data)
	case 0xF:
		// IRQ counter HI
		m.IrqCounter = uint16(data)<<8 | m.IrqCounter&0xFF
	}
}
