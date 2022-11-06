package bus

import log "github.com/sirupsen/logrus"

func New() Bus {
	return Bus{}
}

type Bus struct {
	cpuVram [0x800]byte
}

const (
	RamAddr     = 0x0000
	RamLastAddr = 0x1FFF
	PpuAddr     = 0x2000
	PpuLastAddr = 0x3FFF
)

func (b *Bus) MemRead(addr uint16) byte {
	if RamAddr <= addr && addr <= RamLastAddr {
		addr &= 0b111_1111_1111
		return b.cpuVram[addr]
	} else if PpuAddr <= addr && addr <= PpuLastAddr {
		addr &= 0b10_0000_0000_0111
		panic("PPU unsupported")
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
		addr &= 0b10_0000_0000_0111
		panic("PPU unsupported")
	} else {
		log.WithField("address", addr).Warn("Ignoring memory write")
	}
}
