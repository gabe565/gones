package bus

import (
	"github.com/gabe565/gones/internal/apu"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/controller"
	"github.com/gabe565/gones/internal/ppu"
	log "github.com/sirupsen/logrus"
)

func New(mapper cartridge.Mapper, ppu *ppu.PPU, apu *apu.APU) *Bus {
	return &Bus{
		mapper: mapper,
		apu:    apu,
		ppu:    ppu,
		controller1: controller.Controller{
			Enabled: true,
			Keymap:  controller.Player1Keymap,
		},
		controller2: controller.Controller{
			Enabled: true,
			Keymap:  controller.Player2Keymap,
		},
	}
}

type Bus struct {
	CpuVram     [0x800]byte
	mapper      cartridge.Mapper
	apu         *apu.APU
	ppu         *ppu.PPU
	controller1 controller.Controller
	controller2 controller.Controller
	openBus     byte
}

// ReadMem reads a byte from memory.
func (b *Bus) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		addr &= 0x07FF
		b.openBus = b.CpuVram[addr]
	case 0x2000 <= addr && addr <= 0x2007, addr == 0x4014:
		return b.ppu.ReadMem(addr)
	case 0x2008 <= addr && addr < 0x4000:
		addr &= 0x2007
		return b.ppu.ReadMem(addr)
	case 0x4000 <= addr && addr < 0x4016:
		b.openBus = b.apu.ReadMem(addr)
	case addr == 0x4016:
		b.openBus &^= 0xF
		b.openBus |= b.controller1.Read()
	case addr == 0x4017:
		b.openBus &^= 0xF
		b.openBus |= b.controller2.Read()
	case addr <= 0x4018 && addr < 0x4020:
		// Disabled test registers
	case 0x4020 <= addr:
		b.openBus = b.mapper.ReadMem(addr)
	default:
		log.Errorf("invalid Bus read from $%02X", addr)
		return 0
	}
	return b.openBus
}

// WriteMem writes a byte to memory.
func (b *Bus) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		addr &= 0x07FF
		b.CpuVram[addr] = data
	case 0x2000 <= addr && addr <= 0x2007, addr == 0x4014:
		b.ppu.WriteMem(addr, data)
		return
	case 0x2008 <= addr && addr < 0x4000:
		addr &= 0x2007
		b.ppu.WriteMem(addr, data)
		return
	case 0x4000 <= addr && addr <= 0x4013, addr == 0x4015, addr == 0x4017:
		b.apu.WriteMem(addr, data)
	case addr == 0x4016:
		b.controller1.Write(data)
		b.controller2.Write(data)
	case addr <= 0x4018 && addr < 0x4020:
		// Disabled test registers
	case 0x4020 <= addr:
		b.mapper.WriteMem(addr, data)
	default:
		log.Errorf("invalid Bus write to $%02X", addr)
	}
	b.openBus = data
}

// ReadMem16 reads two bytes from memory.
func (b *Bus) ReadMem16(addr uint16) uint16 {
	lo := uint16(b.ReadMem(addr))
	hi := uint16(b.ReadMem(addr + 1))
	return hi<<8 | lo
}

// WriteMem16 writes two bytes to memory.
func (b *Bus) WriteMem16(addr uint16, data uint16) {
	hi := byte(data >> 8)
	lo := byte(data & 0xFF)
	b.WriteMem(addr, lo)
	b.WriteMem(addr+1, hi)
}

func (b *Bus) UpdateInput() {
	b.controller1.UpdateInput()
	b.controller2.UpdateInput()
}
