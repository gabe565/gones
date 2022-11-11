package bus

import (
	"fmt"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/interrupts"
	"github.com/gabe565/gones/internal/joypad"
	"github.com/gabe565/gones/internal/ppu"
	log "github.com/sirupsen/logrus"
)

func New(cart *cartridge.Cartridge, ppu *ppu.PPU) *Bus {
	return &Bus{
		cartridge:   cart,
		ppu:         ppu,
		Joypad1:     &joypad.Joypad{},
		RenderStart: make(chan struct{}),
		RenderDone:  make(chan struct{}),
	}
}

type Bus struct {
	cpuVram     [0x800]byte
	cartridge   *cartridge.Cartridge
	ppu         *ppu.PPU
	Joypad1     *joypad.Joypad
	cycles      uint
	RenderStart chan struct{}
	RenderDone  chan struct{}
}

func (b *Bus) MemRead(addr uint16) byte {
	switch {
	case addr < 0x2000:
		addr &= 0x07FF
		return b.cpuVram[addr]
	case addr == 0x2000, addr == 0x2001, addr == 0x2003, addr == 0x2005, addr == 0x2006, addr == 0x4014:
		return 0
	case addr == 0x2002:
		return b.ppu.ReadStatus()
	case addr == 0x2004:
		return b.ppu.ReadOam()
	case addr == 0x2007:
		return b.ppu.Read()
	case 0x2008 <= addr && addr < 0x4000:
		addr &= 0x2007
		return b.MemRead(addr)
	case 0x4000 <= addr && addr < 0x4016:
		// APU
		return 0
	case addr == 0x4016:
		return b.Joypad1.Read()
	case addr == 0x4017:
		// Joypad 2
		return 0
	case addr <= 0x4018 && addr < 0x4020:
		// Disabled APU
		return 0
	default:
		addr -= 0x8000
		if len(b.cartridge.Prg) == 0x4000 {
			addr %= 0x4000
		}
		return b.cartridge.Prg[addr]
	}
}

func (b *Bus) MemWrite(addr uint16, data byte) {
	switch {
	case addr < 0x2000:
		addr &= 0x07FF
		b.cpuVram[addr] = data
	case addr == 0x2000:
		b.ppu.WriteCtrl(data)
	case addr == 0x2001:
		b.ppu.WriteMask(data)
	case addr == 0x2002:
		log.WithField("address", fmt.Sprintf("%02X", addr)).
			Error("attempt to write to PPU status register")
	case addr == 0x2003:
		b.ppu.WriteOamAddr(data)
	case addr == 0x2004:
		b.ppu.WriteOam(data)
	case addr == 0x2005:
		b.ppu.WriteScroll(data)
	case addr == 0x2006:
		b.ppu.WriteAddr(data)
	case addr == 0x2007:
		b.ppu.Write(data)
	case 0x2008 <= addr && addr < 0x4000:
		addr &= 0x2007
		b.MemWrite(addr, data)
	case addr == 0x4014:
		var buf [256]byte
		hi := uint16(data) << 8
		for k := range buf {
			buf[k] = b.MemRead(hi + uint16(k))
		}
		b.ppu.WriteOamDma(buf)
	case 0x4000 <= addr && addr < 0x4013, addr == 0x4015:
		// APU
	case addr == 0x4016:
		b.Joypad1.Write(data)
	case addr == 0x4017:
		// Joypad 2
	case addr <= 0x4018 && addr < 0x4020:
		// Disabled
	default:
		log.WithField("address", fmt.Sprintf("%02X", addr)).
			Error("attempt to write to cartridge ROM")
	}
}

func (b *Bus) Tick(cycles uint) {
	b.cycles += cycles

	if b.ppu.Tick(cycles * 3) {
		b.RenderStart <- struct{}{}
		<-b.RenderDone
	}
}

func (b *Bus) ReadInterrupt() *interrupts.Interrupt {
	return b.ppu.ReadInterrupt()
}
