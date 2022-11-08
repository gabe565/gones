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
	if addr <= RamLastAddr {
		addr &= 0b111_1111_1111
		return b.cpuVram[addr]
	} else if addr <= PpuLastAddr {
		// addr &= 0b10_0000_0000_0111
		log.Error("PPU unsupported")
		return 0
	} else {
		addr -= PrgRomAddr
		if len(b.cartridge.Prg) == PrgRomMirror {
			addr %= PrgRomMirror
		}
		return b.cartridge.Prg[addr]
	}
}

func (b *Bus) MemWrite(addr uint16, data byte) {
	if addr <= RamLastAddr {
		addr &= 0b111_1111_1111
		b.cpuVram[addr] = data
	} else if addr <= PpuLastAddr {
		// addr &= 0b10_0000_0000_0111
		log.Error("PPU unsupported")
	} else {
		panic("Attempt to write to cartridge ROM")
	}
}
