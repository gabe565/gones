package bus

import (
	"log/slog"

	"github.com/gabe565/gones/internal/apu"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/controller"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/gabe565/gones/internal/util"
)

func New(conf *config.Config, mapper cartridge.Mapper, ppu *ppu.PPU, apu *apu.APU) *Bus {
	return &Bus{
		mapper:      mapper,
		apu:         apu,
		ppu:         ppu,
		controller1: controller.NewController(conf, controller.Player1),
		controller2: controller.NewController(conf, controller.Player2),
	}
}

type Bus struct {
	CPUVRAM     [0x800]byte `msgpack:"alias:CpuVram"`
	mapper      cartridge.Mapper
	apu         *apu.APU
	ppu         *ppu.PPU
	controller1 controller.Controller
	controller2 controller.Controller
	OpenBus     byte
}

// ReadMem reads a byte from memory.
func (b *Bus) ReadMem(addr uint16) byte {
	switch {
	case addr < 0x2000:
		addr &= 0x07FF
		b.OpenBus = b.CPUVRAM[addr]
	case 0x2000 <= addr && addr <= 0x2007, addr == 0x4014:
		return b.ppu.ReadMem(addr)
	case 0x2008 <= addr && addr < 0x4000:
		addr &= 0x2007
		return b.ppu.ReadMem(addr)
	case 0x4000 <= addr && addr < 0x4016:
		b.OpenBus = b.apu.ReadMem(addr)
	case addr == 0x4016:
		b.OpenBus &^= 0xF
		b.OpenBus |= b.controller1.Read()
	case addr == 0x4017:
		b.OpenBus &^= 0xF
		b.OpenBus |= b.controller2.Read()
	case addr <= 0x4018 && addr < 0x4020:
		// Disabled test registers
	case 0x4020 <= addr:
		b.OpenBus = b.mapper.ReadMem(addr)
	default:
		slog.Error("Invalid Bus read", "addr", util.EncodeHexAddr(addr))
		return 0
	}
	return b.OpenBus
}

// ReadMemSafe reads a byte from memory, but immediately returns 0xFF for any reads with side effects.
func (b *Bus) ReadMemSafe(addr uint16) byte {
	switch {
	case 0x2001 <= addr && addr < 0x4000,
		0x4004 <= addr && addr <= 0x4007,
		0x4015 <= addr && addr <= 0x4017:
		return 0xFF
	default:
		return b.ReadMem(addr)
	}
}

// WriteMem writes a byte to memory.
func (b *Bus) WriteMem(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		addr &= 0x07FF
		b.CPUVRAM[addr] = data
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
		slog.Error("Invalid Bus write", "addr", util.EncodeHexAddr(addr))
	}
	b.OpenBus = data
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

func (b *Bus) SetMapper(m cartridge.Mapper) {
	b.mapper = m
}
