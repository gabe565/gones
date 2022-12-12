package cartridge

import (
	"encoding/gob"
	"github.com/gabe565/gones/internal/interrupts"
	log "github.com/sirupsen/logrus"
)

func NewMapper4(cartridge *Cartridge) Mapper {
	mapper := &Mapper4{cartridge: cartridge}
	mapper.prgOffsets[0] = mapper.prgBankOffset(0)
	mapper.prgOffsets[1] = mapper.prgBankOffset(1)
	mapper.prgOffsets[2] = mapper.prgBankOffset(-2)
	mapper.prgOffsets[3] = mapper.prgBankOffset(-1)
	gob.Register(mapper)
	return mapper
}

type Mapper4 struct {
	cartridge  *Cartridge
	cpu        CPU
	register   byte
	registers  [8]byte
	prgMode    bool
	chrMode    bool
	prgOffsets [4]int
	chrOffsets [8]int
	reload     byte
	counter    byte
	irqEnable  bool
}

func (m *Mapper4) Step(renderEnabled bool, scanline uint16, cycle uint) {
	switch {
	case cycle != 260:
		return
	case 240 <= scanline && scanline < 261:
		return
	case !renderEnabled:
		return
	case m.counter == 0:
		m.counter = m.reload
	default:
		m.counter -= 1
		if m.counter == 0 && m.irqEnable {
			m.cpu.AddInterrupt(&interrupts.IRQ)
		}
	}
}

func (m *Mapper4) Cartridge() *Cartridge { return m.cartridge }

func (m *Mapper4) SetCpu(c CPU) { m.cpu = c }

func (m *Mapper4) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		bank := addr / 0x400
		offset := int(addr % 0x400)
		return m.cartridge.Chr[m.chrOffsets[bank]+offset]
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		return m.cartridge.Sram[addr]
	case 0x8000 <= addr:
		addr -= 0x8000
		bank := addr / 0x2000
		offset := int(addr % 0x2000)
		return m.cartridge.prg[m.prgOffsets[bank]+offset]
	default:
		log.Fatalf("invalid mapper 4 read from $%04X", addr)
		return 0
	}
}

func (m *Mapper4) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		bank := addr / 0x400
		offset := int(addr % 0x400)
		m.cartridge.Chr[m.chrOffsets[bank]+offset] = data
	case 0x6000 <= addr && addr < 0x8000:
		addr -= 0x6000
		m.cartridge.Sram[addr] = data
	case 0x8000 <= addr:
		switch {
		case addr < 0xA000:
			if addr%2 == 0 {
				// Bank select
				m.prgMode = data&0x40 == 0x40
				m.chrMode = data&0x80 == 0x80
				m.register = data & 7
				m.updateOffsets()
			} else {
				// Bank data
				m.registers[m.register] = data
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
				m.reload = data
			} else {
				// IRQ Reload
				m.counter = 0
			}
		case 0xE000 <= addr:
			m.irqEnable = addr%2 == 1
		}
	default:
		log.Fatalf("invalid mapper 4 write to $%04X", addr)
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
	if m.prgMode {
		m.prgOffsets[0] = m.prgBankOffset(-2)
		m.prgOffsets[1] = m.prgBankOffset(int(m.registers[7]))
		m.prgOffsets[2] = m.prgBankOffset(int(m.registers[6]))
		m.prgOffsets[3] = m.prgBankOffset(-1)
	} else {
		m.prgOffsets[0] = m.prgBankOffset(int(m.registers[6]))
		m.prgOffsets[1] = m.prgBankOffset(int(m.registers[7]))
		m.prgOffsets[2] = m.prgBankOffset(-2)
		m.prgOffsets[3] = m.prgBankOffset(-1)
	}

	if m.chrMode {
		m.chrOffsets[0] = m.chrBankOffset(int(m.registers[2]))
		m.chrOffsets[1] = m.chrBankOffset(int(m.registers[3]))
		m.chrOffsets[2] = m.chrBankOffset(int(m.registers[4]))
		m.chrOffsets[3] = m.chrBankOffset(int(m.registers[5]))
		m.chrOffsets[4] = m.chrBankOffset(int(m.registers[0] & 0xFE))
		m.chrOffsets[5] = m.chrBankOffset(int(m.registers[0] | 1))
		m.chrOffsets[6] = m.chrBankOffset(int(m.registers[1] & 0xFE))
		m.chrOffsets[7] = m.chrBankOffset(int(m.registers[1] | 1))
	} else {
		m.chrOffsets[0] = m.chrBankOffset(int(m.registers[0] & 0xFE))
		m.chrOffsets[1] = m.chrBankOffset(int(m.registers[0] | 1))
		m.chrOffsets[2] = m.chrBankOffset(int(m.registers[1] & 0xFE))
		m.chrOffsets[3] = m.chrBankOffset(int(m.registers[1] | 1))
		m.chrOffsets[4] = m.chrBankOffset(int(m.registers[2]))
		m.chrOffsets[5] = m.chrBankOffset(int(m.registers[3]))
		m.chrOffsets[6] = m.chrBankOffset(int(m.registers[4]))
		m.chrOffsets[7] = m.chrBankOffset(int(m.registers[5]))
	}
}
