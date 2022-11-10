package ppu

import (
	"fmt"
	"github.com/gabe565/gones/internal/bitflags"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/interrupts"
	"github.com/gabe565/gones/internal/ppu/registers"
	log "github.com/sirupsen/logrus"
)

func New(cart *cartridge.Cartridge) *PPU {
	return &PPU{
		chr:       cart.Chr,
		mirroring: cart.Mirror,
	}
}

type PPU struct {
	chr       []byte
	mirroring cartridge.Mirror
	ctrl      registers.Control
	mask      bitflags.Flags
	scroll    registers.Scroll
	status    bitflags.Flags
	addr      registers.AddrRegister
	vram      [0x800]byte

	oamAddr byte
	oam     [0x100]byte
	palette [0x20]byte

	scanline  uint16
	cycles    uint
	interrupt *interrupts.Interrupt

	readBuf byte
}

func (p *PPU) WriteAddr(data byte) {
	p.addr.Update(data)
}

func (p *PPU) WriteCtrl(data byte) {
	beforeNmi := p.ctrl.GenerateVblankNmi()
	p.ctrl = registers.Control(data)
	if !beforeNmi && p.ctrl.GenerateVblankNmi() && p.status.Has(registers.VblankStarted) {
		p.interrupt = &interrupts.NMI
	}
}

func (p *PPU) WriteMask(data byte) {
	p.mask = bitflags.Flags(data)
}

func (p *PPU) WriteOamAddr(data byte) {
	p.oamAddr = data
}

func (p *PPU) WriteOam(data byte) {
	p.oam[p.oamAddr] = data
	p.oamAddr += 1
}

func (p *PPU) WriteOamDma(data [0x100]byte) {
	for _, data := range data {
		p.oam[p.oamAddr] = data
		p.oamAddr += 1
	}
}

func (p *PPU) ReadOam() byte {
	return p.oam[p.oamAddr]
}

func (p *PPU) WriteScroll(data byte) {
	p.scroll.Write(data)
}

func (p *PPU) ReadStatus() byte {
	data := p.status
	p.status.Remove(registers.VblankStarted)
	p.addr.ResetLatch()
	p.scroll.ResetLatch()
	return byte(data)
}

func (p *PPU) Write(data byte) {
	addr := p.addr.Get()
	switch {
	case addr < 0x2000:
		log.WithField("address", fmt.Sprintf("$%02X", addr)).Warn("attempt to write to CHR ROM")
	case addr < 0x3000:
		addr := p.MirrorVramAddr(addr)
		p.vram[addr] = data
	case addr < 0x3F00:
		log.WithField("address", fmt.Sprintf("$%02X", addr)).Error("bad PPU write")
	case addr < 0x4000:
		switch addr {
		case 0x3F10, 0x3F14, 0x3F18, 0x3F1C:
			addr -= 0x10
		}
		addr -= 0x3F00
		p.palette[addr] = data
	default:
		panic(fmt.Sprintf("unexpected write to mirrored space: $%02X", addr))
	}
	p.addr.Increment(p.ctrl.VramAddrIncrement())
}

func (p *PPU) Read() byte {
	addr := p.addr.Get()
	p.addr.Increment(p.ctrl.VramAddrIncrement())

	switch {
	case addr < 0x2000:
		result := p.readBuf
		p.readBuf = p.chr[addr]
		return result
	case addr < 0x3000:
		result := p.readBuf
		addr := p.MirrorVramAddr(addr)
		p.readBuf = p.vram[addr]
		return result
	case addr < 0x3F00:
		log.WithField("address", fmt.Sprintf("$%02X", addr)).Error("bad PPU write")
		return 0
	case addr < 0x4000:
		switch addr {
		case 0x3F10, 0x3F14, 0x3F18, 0x3F1C:
			addr -= 0x10
		}
		addr -= 0x3F00
		return p.palette[addr]
	default:
		panic(fmt.Sprintf("unexpected access to mirrored space: $%02X", addr))
	}
}

func (p *PPU) MirrorVramAddr(addr uint16) uint16 {
	addr &= 0x2FFF
	addr -= 0x2000
	nameTable := addr / 0x400

	switch p.mirroring {
	case cartridge.Vertical:
		switch nameTable {
		case 2, 3:
			return addr - 0x800
		}
	case cartridge.Horizontal:
		switch nameTable {
		case 1, 2:
			return addr - 0x400
		case 3:
			return addr - 0x800
		}
	}
	return addr
}

func (p *PPU) Tick(cycles uint) bool {
	p.cycles += cycles
	if p.cycles >= 341 {
		p.cycles -= 341
		p.scanline += 1

		if p.scanline == 241 {
			p.status.Insert(registers.VblankStarted)
			p.status.Remove(registers.SpriteZeroHit)
			if p.ctrl.GenerateVblankNmi() {
				p.interrupt = &interrupts.NMI
			}
		}

		if p.scanline >= 262 {
			p.scanline = 0
			p.interrupt = nil
			p.status.Remove(registers.SpriteZeroHit | registers.VblankStarted)
			return true
		}
	}
	return false
}

func (p *PPU) ReadInterrupt() *interrupts.Interrupt {
	return p.interrupt
}
