package cartridge

import (
	"github.com/gabe565/gones/internal/ppu/registers"
	log "github.com/sirupsen/logrus"
)

func NewMapper4(cartridge *Cartridge) Mapper {
	mapper := &Mapper4{cartridge: cartridge}
	mapper.PrgOffsets[0] = mapper.prgBankOffset(0)
	mapper.PrgOffsets[1] = mapper.prgBankOffset(1)
	mapper.PrgOffsets[2] = mapper.prgBankOffset(-2)
	mapper.PrgOffsets[3] = mapper.prgBankOffset(-1)
	return mapper
}

type Mapper4 struct {
	cartridge  *Cartridge
	Register   byte
	Registers  [8]byte
	PrgMode    bool
	ChrMode    bool
	PrgOffsets [4]int
	ChrOffsets [8]int
	Reload     byte
	Counter    byte
	IrqEnable  bool
	IrqPending bool
	PrevA12    bool
}

func (m *Mapper4) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper4) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper4) OnScanline() {
	if m.Counter == 0 {
		m.Counter = m.Reload
	} else {
		m.Counter -= 1
		if m.Counter == 0 && m.IrqEnable {
			m.IrqPending = true
		}
	}
}

func (m *Mapper4) Irq() bool { return m.IrqPending }

func (m *Mapper4) OnVramAddr(addr registers.Address) {
	curr := addr.FineY&1 == 1
	switch m.cartridge.Submapper {
	case SubmapperMcAcc:
		if m.PrevA12 && !curr {
			m.OnScanline()
		}
	default:
		if !m.PrevA12 && curr {
			m.OnScanline()
		}
	}
	m.PrevA12 = curr
}

func (m *Mapper4) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		bank := addr / 0x400
		offset := int(addr % 0x400)
		return m.cartridge.Chr[m.ChrOffsets[bank]+offset]
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		return m.cartridge.Sram[addr]
	case 0x8000 <= addr:
		addr -= 0x8000
		bank := addr / 0x2000
		offset := int(addr % 0x2000)
		return m.cartridge.prg[m.PrgOffsets[bank]+offset]
	default:
		log.Warnf("invalid mapper 4 read from $%04X", addr)
		return 0
	}
}

func (m *Mapper4) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		bank := addr / 0x400
		offset := int(addr % 0x400)
		m.cartridge.Chr[m.ChrOffsets[bank]+offset] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		m.cartridge.Sram[addr] = data
	case 0x8000 <= addr && addr < 0xA000:
		if addr%2 == 0 {
			// Bank select
			m.PrgMode = data&0x40 == 0x40
			m.ChrMode = data&0x80 == 0x80
			m.Register = data & 7
			m.updateOffsets()
		} else {
			// Bank data
			m.Registers[m.Register] = data
			m.updateOffsets()
		}
	case 0xA000 <= addr && addr < 0xC000:
		if addr%2 == 0 {
			// Mirror
			switch data & 1 {
			case 0:
				m.cartridge.Mirror = Vertical
			case 1:
				m.cartridge.Mirror = Horizontal
			}
		}
	case 0xC000 <= addr && addr < 0xE000:
		if addr%2 == 0 {
			// IRQ Latch
			m.Reload = data
		} else {
			// IRQ Reload
			m.Counter = 0
		}
	case 0xE000 <= addr:
		m.IrqEnable = addr%2 == 1
		if !m.IrqEnable {
			m.IrqPending = false
		}
	default:
		log.Warnf("invalid mapper 4 write to $%04X", addr)
	}
}

func (m *Mapper4) prgBankOffset(i int) int {
	if i >= 0x80 {
		i -= 0x100
	}
	i %= len(m.cartridge.prg) / 0x2000
	offset := i * 0x2000
	if offset < 0 {
		offset += len(m.cartridge.prg)
	}
	return offset
}

func (m *Mapper4) chrBankOffset(i int) int {
	if i >= 0x80 {
		i -= 0x100
	}
	i %= len(m.cartridge.Chr) / 0x0400
	offset := i * 0x0400
	if offset < 0 {
		offset += len(m.cartridge.Chr)
	}
	return offset
}

func (m *Mapper4) updateOffsets() {
	if m.PrgMode {
		m.PrgOffsets[0] = m.prgBankOffset(-2)
		m.PrgOffsets[1] = m.prgBankOffset(int(m.Registers[7]))
		m.PrgOffsets[2] = m.prgBankOffset(int(m.Registers[6]))
		m.PrgOffsets[3] = m.prgBankOffset(-1)
	} else {
		m.PrgOffsets[0] = m.prgBankOffset(int(m.Registers[6]))
		m.PrgOffsets[1] = m.prgBankOffset(int(m.Registers[7]))
		m.PrgOffsets[2] = m.prgBankOffset(-2)
		m.PrgOffsets[3] = m.prgBankOffset(-1)
	}

	if m.ChrMode {
		m.ChrOffsets[0] = m.chrBankOffset(int(m.Registers[2]))
		m.ChrOffsets[1] = m.chrBankOffset(int(m.Registers[3]))
		m.ChrOffsets[2] = m.chrBankOffset(int(m.Registers[4]))
		m.ChrOffsets[3] = m.chrBankOffset(int(m.Registers[5]))
		m.ChrOffsets[4] = m.chrBankOffset(int(m.Registers[0] & 0xFE))
		m.ChrOffsets[5] = m.chrBankOffset(int(m.Registers[0] | 1))
		m.ChrOffsets[6] = m.chrBankOffset(int(m.Registers[1] & 0xFE))
		m.ChrOffsets[7] = m.chrBankOffset(int(m.Registers[1] | 1))
	} else {
		m.ChrOffsets[0] = m.chrBankOffset(int(m.Registers[0] & 0xFE))
		m.ChrOffsets[1] = m.chrBankOffset(int(m.Registers[0] | 1))
		m.ChrOffsets[2] = m.chrBankOffset(int(m.Registers[1] & 0xFE))
		m.ChrOffsets[3] = m.chrBankOffset(int(m.Registers[1] | 1))
		m.ChrOffsets[4] = m.chrBankOffset(int(m.Registers[2]))
		m.ChrOffsets[5] = m.chrBankOffset(int(m.Registers[3]))
		m.ChrOffsets[6] = m.chrBankOffset(int(m.Registers[4]))
		m.ChrOffsets[7] = m.chrBankOffset(int(m.Registers[5]))
	}
}
