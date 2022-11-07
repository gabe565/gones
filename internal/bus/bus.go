package bus

import (
	"github.com/gabe565/gones/internal/cartridge"
	log "github.com/sirupsen/logrus"
)

func New(cart *cartridge.Cartridge) *Bus {
	return &Bus{
		cartridge: cart,
	}
}

type Bus struct {
	cpuVram   [0x800]byte
	cartridge *cartridge.Cartridge
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
	if RamAddr <= addr && addr <= RamLastAddr {
		addr &= 0b111_1111_1111
		return b.cpuVram[addr]
	} else if PpuAddr <= addr && addr <= PpuLastAddr {
		// addr &= 0b10_0000_0000_0111
		log.Error("PPU unsupported")
		return 0
	} else if PrgRomAddr <= addr && addr <= 0xFFFF {
		addr -= PrgRomAddr
		if len(b.cartridge.Prg) == PrgRomMirror {
			addr %= PrgRomMirror
		}
		return b.cartridge.Prg[addr]
	} else {
		log.WithField("address", addr).Warn("Ignoring memory read")
		return 0
	}
}

func (b *Bus) MemWrite(addr uint16, data byte) {
	if RamAddr <= addr && addr <= RamLastAddr {
		addr &= 0b111_1111_1111
		b.cpuVram[addr] = data
	} else if PpuAddr <= addr && addr <= PpuLastAddr {
		// addr &= 0b10_0000_0000_0111
		log.Error("PPU unsupported")
	} else if PrgRomAddr <= addr && addr <= 0xFFFF {
		panic("Attempt to write to cartridge ROM")
	} else {
		log.WithField("address", addr).Warn("Ignoring memory write")
	}
}
