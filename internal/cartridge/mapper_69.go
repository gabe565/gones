package cartridge

import (
	"log/slog"

	"github.com/gabe565/gones/internal/util"
)

func NewMapper69(cartridge *Cartridge) *Mapper69 {
	prgCount := len(cartridge.prg) / 0x2000
	mapper := &Mapper69{
		cartridge: cartridge,
		PRGCount:  byte(prgCount),
		PRGBanks:  [5]int{0, 0, 0, 0, prgCount - 1},
	}
	return mapper
}

type Mapper69 struct {
	cartridge *Cartridge
	Command   byte

	PRGCount   byte   `msgpack:"alias:PrgCount"`
	PRGBanks   [5]int `msgpack:"alias:PrgBanks"`
	RAMSelect  bool   `msgpack:"alias:RamSelect"`
	RAMEnabled bool   `msgpack:"alias:RamEnabled"`

	CHRBanks [8]int `msgpack:"alias:ChrBanks"`

	IRQEnabled        bool   `msgpack:"alias:IrqEnable"`
	IRQCounterEnabled bool   `msgpack:"alias:IrqCounterEnable"`
	IRQCounter        uint16 `msgpack:"alias:IrqCounter"`
	IRQPending        bool   `msgpack:"alias:IrqPending"`
}

func (m *Mapper69) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper69) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper69) OnCPUStep(cycles uint) {
	if m.IRQCounterEnabled {
		prev := m.IRQCounter
		m.IRQCounter -= uint16(cycles)
		if m.IRQEnabled && m.IRQCounter > prev {
			m.IRQPending = true
		}
	}
}

func (m *Mapper69) IRQ() bool { return m.IRQPending }

func (m *Mapper69) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		// CHR
		bank := addr / 0x400
		offset := int(addr % 0x400)
		addr := m.CHRBanks[bank]*0x400 + offset
		return m.cartridge.CHR[addr%len(m.cartridge.CHR)]
	case 0x6000 <= addr && addr < 0x8000:
		// PRG/SRAM banks
		if m.RAMSelect {
			if m.RAMEnabled {
				return m.cartridge.SRAM[addr-0x6000]
			}
			// open bus
			return 0
		}
		fallthrough
	case 0x8000 <= addr:
		addr -= 0x6000
		bank := addr / 0x2000
		offset := int(addr % 0x2000)
		addr := m.PRGBanks[bank]*0x2000 + offset
		return m.cartridge.prg[addr%len(m.cartridge.prg)]
	default:
		slog.Error("Invalid mapper 69 read", "addr", util.HexAddr(addr))
		return 0
	}
}

func (m *Mapper69) WriteMem(addr uint16, data byte) {
	switch {
	case 0x6000 <= addr && addr < 0x8000:
		// SRAM register
		if m.RAMSelect && m.RAMEnabled {
			m.cartridge.SRAM[addr-0x6000] = data
		}
	case 0x8000 <= addr && addr < 0xA000:
		// Command register
		m.Command = data & 0xF
	case 0xA000 <= addr && addr < 0xC000:
		// Parameter register
		m.runCommand(data)
	default:
		slog.Error("Invalid mapper 69 write", "addr", util.HexAddr(addr))
	}
}

func (m *Mapper69) runCommand(data byte) {
	switch m.Command {
	case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7:
		// CHR bank switch
		bank := m.Command
		m.CHRBanks[bank] = int(data)
	case 0x8:
		// PRG bank switch (with RAM)
		m.RAMSelect = data>>6&1 == 1
		m.RAMEnabled = data>>7&1 == 1
		fallthrough
	case 0x9, 0xA, 0xB:
		// PRG bank switch
		bank := m.Command - 0x8
		m.PRGBanks[bank] = int(data & 0x1F)
	case 0xC:
		// Nametable Mirroring
		m.cartridge.Mirror = Mirror(data & 0x3)
	case 0xD:
		// IRQ control
		m.IRQEnabled = data&1 == 1
		m.IRQCounterEnabled = data>>7&1 == 1
		m.IRQPending = false
	case 0xE:
		// IRQ counter LO
		m.IRQCounter = m.IRQCounter&0xFF00 | uint16(data)
	case 0xF:
		// IRQ counter HI
		m.IRQCounter = uint16(data)<<8 | m.IRQCounter&0xFF
	}
}
