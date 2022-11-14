package ppu

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/gabe565/gones/internal/bitflags"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/interrupts"
	"github.com/gabe565/gones/internal/ppu/registers"
	log "github.com/sirupsen/logrus"
)

func New(cart *cartridge.Cartridge) *PPU {
	return &PPU{
		chr:         cart.Chr,
		mirroring:   cart.Mirror,
		interruptCh: make(chan interrupts.Interrupt, 1),
		Picture:     pixel.MakePictureData(pixel.R(0, 0, Width, Height)),
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

	scanline    uint16
	cycles      uint
	interruptCh chan interrupts.Interrupt

	readBuf byte
	Picture *pixel.PictureData
}

func (p *PPU) WriteAddr(data byte) {
	p.addr.Update(data)
}

func (p *PPU) WriteCtrl(data byte) {
	beforeNmi := p.ctrl.HasEnableNMI()
	p.ctrl = registers.Control(data)
	if !beforeNmi && p.ctrl.HasEnableNMI() && p.status.Has(registers.Vblank) {
		p.interruptCh <- interrupts.NMI
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
		p.WriteOam(data)
	}
}

func (p *PPU) ReadOam() byte {
	return p.oam[p.oamAddr]
}

func (p *PPU) WriteScroll(data byte) {
	p.scroll.Write(data)
}

func (p *PPU) ReadStatus() byte {
	defer func() {
		p.status.Remove(registers.Vblank)
	}()
	p.addr.ResetLatch()
	p.scroll.ResetLatch()
	return byte(p.status)
}

func (p *PPU) Write(data byte) {
	addr := p.addr.Get()
	switch {
	case addr < 0x2000:
		log.WithField("address", fmt.Sprintf("$%02X", addr)).
			Error("attempt to write to cartridge ROM")
	case 0x2000 <= addr && addr < 0x3000:
		addr := p.MirrorVramAddr(addr)
		p.vram[addr] = data
	case 0x3000 <= addr && addr < 0x3F00:
		log.WithField("address", fmt.Sprintf("$%02X", addr)).
			Error("bad PPU write")
	case 0x3F00 <= addr && addr < 0x4000:
		addr &= 0x3F1F
		switch addr {
		case 0x3F10, 0x3F14, 0x3F18, 0x3F1C:
			addr -= 0x10
		}
		addr -= 0x3F00
		p.palette[addr] = data
	default:
		log.WithField("address", fmt.Sprintf("%02X", addr)).
			Error("unexpected write to mirrored space")
	}
	p.addr.Increment(p.ctrl.VramAddr())
}

func (p *PPU) Read() byte {
	addr := p.addr.Get()
	p.addr.Increment(p.ctrl.VramAddr())

	switch {
	case addr < 0x2000:
		result := p.readBuf
		p.readBuf = p.chr[addr]
		return result
	case 0x2000 <= addr && addr < 0x3000:
		result := p.readBuf
		addr := p.MirrorVramAddr(addr)
		p.readBuf = p.vram[addr]
		return result
	case 0x3000 <= addr && addr < 0x3F00:
		log.WithField("address", fmt.Sprintf("$%02X", addr)).
			Error("bad PPU write")
		return 0
	case 0x3F00 <= addr && addr < 0x4000:
		addr &= 0x3F1F
		switch addr {
		case 0x3F10, 0x3F14, 0x3F18, 0x3F1C:
			addr -= 0x10
		}
		addr -= 0x3F00
		return p.palette[addr]
	default:
		log.WithField("address", fmt.Sprintf("%02X", addr)).
			Error("unexpected access to mirrored space")
		return 0
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

func (p *PPU) updateSpriteOverflow() {
	size := int(p.ctrl.SpriteSize())
	var count uint
	for i := 0; i < len(p.oam)/4; i += 1 {
		i := i * 4
		tileY := p.oam[i]
		row := int(p.scanline) - int(tileY)
		if row < 0 || row >= size {
			continue
		}
		count += 1
	}
	count &= 0b1111
	if count == 8 {
		p.status.Insert(registers.SpriteOverflow)
	}
}

func (p *PPU) Step() bool {
	p.cycles += 1

	switch {
	case p.cycles == 257:
		p.updateSpriteOverflow()
	case p.cycles > 340:
		if p.SpriteZeroHit(p.cycles) {
			p.status.Insert(registers.SpriteZeroHit)
		}

		p.cycles = 0
		p.scanline += 1

		switch {
		case p.scanline == 241:
			p.status.Insert(registers.Vblank)
			p.status.Remove(registers.SpriteZeroHit)
			if p.ctrl.HasEnableNMI() {
				p.interruptCh <- interrupts.NMI
			}
		case p.scanline >= 262:
			p.scanline = 0
			p.status.Remove(registers.Vblank | registers.SpriteOverflow | registers.SpriteZeroHit)
			return true
		}
	}
	return false
}

func (p *PPU) SpriteZeroHit(cycle uint) bool {
	x := p.oam[3]
	y := p.oam[0]
	return uint16(y) == p.scanline && uint(x) <= cycle && p.mask.Has(registers.SpriteEnable)
}

func (p *PPU) Reset() {
	p.cycles = 0
	p.scanline = 0
	p.WriteCtrl(0)
	p.WriteMask(0)
	p.WriteOamAddr(0)
	p.addr = registers.AddrRegister{}
	p.scroll = registers.Scroll{}
}

func (p *PPU) Interrupt() <-chan interrupts.Interrupt {
	return p.interruptCh
}
