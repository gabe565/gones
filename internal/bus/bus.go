package bus

import (
	"fmt"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/interrupts"
	"github.com/gabe565/gones/internal/ppu"
)

func New(cart *cartridge.Cartridge) *Bus {
	return &Bus{
		cartridge: cart,
		ppu:       ppu.New(cart),
	}
}

type Bus struct {
	cpuVram   [0x800]byte
	cartridge *cartridge.Cartridge
	ppu       *ppu.PPU
	cycles    uint
}

const (
	RamAddr      = 0x0000
	RamLastAddr  = 0x1FFF
	PpuAddr      = 0x2000
	PpuLastAddr  = 0x3FFF
	PrgRomAddr   = 0x8000
	PrgRomMirror = 0x4000
)

func (b *Bus) MemRead(addr uint16) byte {
	switch {
	case addr <= RamLastAddr:
		addr &= 0b111_1111_1111
		return b.cpuVram[addr]
	case addr == 0x2000, addr == 0x2001, addr == 0x2003, addr == 0x2005, addr == 0x2006, addr == 0x4014:
		panic(fmt.Sprintf("attempt to read from write-only PPU address $%02X", addr))
	case addr == 0x2007:
		return b.ppu.Read()
	case 0x2008 <= addr && addr <= PpuLastAddr:
		addr &= 0b0010_0000_0000_0111
		return b.MemRead(addr)
	default:
		addr -= PrgRomAddr
		if len(b.cartridge.Prg) == PrgRomMirror {
			addr %= PrgRomMirror
		}
		return b.cartridge.Prg[addr]
	}
}

func (b *Bus) MemWrite(addr uint16, data byte) {
	switch {
	case addr <= RamLastAddr:
		addr &= 0b111_1111_1111
		b.cpuVram[addr] = data
	case addr == 0x2000:
		b.ppu.WriteCtrl(data)
	case addr == 0x2006:
		b.ppu.WriteAddr(data)
	case addr == 0x2007:
		b.ppu.Write(data)
	case 0x2008 <= addr && addr <= PpuLastAddr:
		addr &= 0b10_0000_0000_0111
		b.MemWrite(addr, data)
	default:
		panic("Attempt to write to cartridge ROM")
	}
}

func (b *Bus) Tick(cycles uint) {
	b.cycles += cycles
	b.ppu.Tick(cycles * 3)
}

func (b *Bus) ReadInterrupt() *interrupts.Interrupt {
	return b.ppu.ReadInterrupt()
}
