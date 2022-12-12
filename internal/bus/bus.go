package bus

import (
	"github.com/gabe565/gones/internal/apu"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/controller"
	"github.com/gabe565/gones/internal/ppu"
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
}

// ReadMem reads a byte from memory.
func (b *Bus) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		addr &= 0x07FF
		return b.CpuVram[addr]
	case 0x2000 <= addr && addr <= 0x2007, addr == 0x4014:
		return b.ppu.ReadMem(addr)
	case 0x2008 <= addr && addr < 0x4000:
		addr &= 0x2007
		return b.ReadMem(addr)
	case 0x4000 <= addr && addr < 0x4016:
		return b.apu.ReadMem(addr)
	case addr == 0x4016:
		return b.controller1.Read()
	case addr == 0x4017:
		return b.controller2.Read()
	case addr <= 0x4018 && addr < 0x4020:
		// Disabled
	default:
		return b.mapper.ReadMem(addr)
	}
	return 0
}

// WriteMem writes a byte to memory.
func (b *Bus) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		addr &= 0x07FF
		b.CpuVram[addr] = data
	case 0x2000 <= addr && addr <= 0x2007, addr == 0x4014:
		b.ppu.WriteMem(addr, data)
	case 0x2008 <= addr && addr < 0x4000:
		addr &= 0x2007
		b.WriteMem(addr, data)
	case 0x4000 <= addr && addr <= 0x4013, addr == 0x4015, addr == 0x4017:
		b.apu.WriteMem(addr, data)
	case addr == 0x4016:
		b.controller1.Write(data)
		b.controller2.Write(data)
	case addr <= 0x4018 && addr < 0x4020:
		// Disabled
	default:
		b.mapper.WriteMem(addr, data)
	}
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
