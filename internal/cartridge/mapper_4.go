package cartridge

import (
	"log/slog"

	"github.com/gabe565/gones/internal/ppu/registers"
	"github.com/gabe565/gones/internal/util"
)

func NewMapper4(cartridge *Cartridge) *Mapper4 {
	mapper := &Mapper4{cartridge: cartridge}
	mapper.PRGOffsets[0] = mapper.prgBankOffset(0)
	mapper.PRGOffsets[1] = mapper.prgBankOffset(1)
	mapper.PRGOffsets[2] = mapper.prgBankOffset(-2)
	mapper.PRGOffsets[3] = mapper.prgBankOffset(-1)
	return mapper
}

type Mapper4 struct {
	cartridge  *Cartridge
	Register   byte
	Registers  [8]byte
	PRGMode    bool   `msgpack:"alias:PrgMode"`
	CHRMode    bool   `msgpack:"alias:ChrMode"`
	PRGOffsets [4]int `msgpack:"alias:PrgOffsets"`
	CHROffsets [8]int `msgpack:"alias:ChrOffsets"`
	Reload     byte
	Counter    byte
	IRQEnabled bool `msgpack:"alias:IrqEnable"`
	IRQPending bool `msgpack:"alias:IrqPending"`
	PrevA12    bool
}

func (m *Mapper4) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper4) SetCartridge(c *Cartridge) { m.cartridge = c }

func (m *Mapper4) OnScanline() {
	if m.Counter == 0 {
		m.Counter = m.Reload
	} else {
		m.Counter--
		if m.Counter == 0 && m.IRQEnabled {
			m.IRQPending = true
		}
	}
}

func (m *Mapper4) IRQ() bool { return m.IRQPending }

func (m *Mapper4) OnVRAMAddr(addr registers.Address) {
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
		return m.cartridge.CHR[m.CHROffsets[bank]+offset]
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		return m.cartridge.SRAM[addr]
	case 0x8000 <= addr:
		addr -= 0x8000
		bank := addr / 0x2000
		offset := int(addr % 0x2000)
		return m.cartridge.prg[m.PRGOffsets[bank]+offset]
	default:
		slog.Error("Invalid mapper 4 read", "addr", util.EncodeHexAddr(addr))
		return 0
	}
}

func (m *Mapper4) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		bank := addr / 0x400
		offset := int(addr % 0x400)
		m.cartridge.CHR[m.CHROffsets[bank]+offset] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		m.cartridge.SRAM[addr] = data
	case 0x8000 <= addr && addr < 0xA000:
		if addr%2 == 0 {
			// Bank select
			m.PRGMode = data&0x40 == 0x40
			m.CHRMode = data&0x80 == 0x80
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
		m.IRQEnabled = addr%2 == 1
		if !m.IRQEnabled {
			m.IRQPending = false
		}
	default:
		slog.Error("Invalid mapper 4 write", "addr", util.EncodeHexAddr(addr))
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
	i %= len(m.cartridge.CHR) / 0x0400
	offset := i * 0x0400
	if offset < 0 {
		offset += len(m.cartridge.CHR)
	}
	return offset
}

func (m *Mapper4) updateOffsets() {
	if m.PRGMode {
		m.PRGOffsets[0] = m.prgBankOffset(-2)
		m.PRGOffsets[1] = m.prgBankOffset(int(m.Registers[7]))
		m.PRGOffsets[2] = m.prgBankOffset(int(m.Registers[6]))
		m.PRGOffsets[3] = m.prgBankOffset(-1)
	} else {
		m.PRGOffsets[0] = m.prgBankOffset(int(m.Registers[6]))
		m.PRGOffsets[1] = m.prgBankOffset(int(m.Registers[7]))
		m.PRGOffsets[2] = m.prgBankOffset(-2)
		m.PRGOffsets[3] = m.prgBankOffset(-1)
	}

	if m.CHRMode {
		m.CHROffsets[0] = m.chrBankOffset(int(m.Registers[2]))
		m.CHROffsets[1] = m.chrBankOffset(int(m.Registers[3]))
		m.CHROffsets[2] = m.chrBankOffset(int(m.Registers[4]))
		m.CHROffsets[3] = m.chrBankOffset(int(m.Registers[5]))
		m.CHROffsets[4] = m.chrBankOffset(int(m.Registers[0] & 0xFE))
		m.CHROffsets[5] = m.chrBankOffset(int(m.Registers[0] | 1))
		m.CHROffsets[6] = m.chrBankOffset(int(m.Registers[1] & 0xFE))
		m.CHROffsets[7] = m.chrBankOffset(int(m.Registers[1] | 1))
	} else {
		m.CHROffsets[0] = m.chrBankOffset(int(m.Registers[0] & 0xFE))
		m.CHROffsets[1] = m.chrBankOffset(int(m.Registers[0] | 1))
		m.CHROffsets[2] = m.chrBankOffset(int(m.Registers[1] & 0xFE))
		m.CHROffsets[3] = m.chrBankOffset(int(m.Registers[1] | 1))
		m.CHROffsets[4] = m.chrBankOffset(int(m.Registers[2]))
		m.CHROffsets[5] = m.chrBankOffset(int(m.Registers[3]))
		m.CHROffsets[6] = m.chrBankOffset(int(m.Registers[4]))
		m.CHROffsets[7] = m.chrBankOffset(int(m.Registers[5]))
	}
}
