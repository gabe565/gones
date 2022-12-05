package ppu

import (
	"fmt"
	"github.com/gabe565/gones/internal/bitflags"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/interrupts"
	"github.com/gabe565/gones/internal/ppu/registers"
	log "github.com/sirupsen/logrus"
	"image"
)

func New(mapper cartridge.Mapper) *PPU {
	return &PPU{
		mapper:      mapper,
		interruptCh: make(chan interrupts.Interrupt, 1),
		image:       image.NewRGBA(image.Rect(0, 0, Width, TrimmedHeight)),
	}
}

type PPU struct {
	mapper cartridge.Mapper

	Ctrl      registers.Control
	Mask      bitflags.Flags
	Status    bitflags.Flags
	Addr      registers.AddrRegister
	TmpAddr   registers.AddrRegister
	AddrLatch bool
	FineX     byte
	Vram      [0x800]byte

	OamAddr byte
	Oam     [0x100]byte
	Palette [0x20]byte

	Scanline    uint16
	Cycles      uint
	interruptCh chan interrupts.Interrupt

	ReadBuf byte
	image   *image.RGBA
}

func (p *PPU) WriteAddr(data byte) {
	if p.AddrLatch {
		p.TmpAddr.WriteLo(data)
		p.Addr = p.TmpAddr
	} else {
		p.TmpAddr.WriteHi(data)
	}
	p.AddrLatch = !p.AddrLatch
}

func (p *PPU) WriteCtrl(data byte) {
	beforeNmi := p.Ctrl.HasEnableNMI()
	p.Ctrl = registers.Control(data)
	p.TmpAddr.NametableX = bitflags.Flags(p.Ctrl).Intersects(registers.CtrlNametableX)
	p.TmpAddr.NametableY = bitflags.Flags(p.Ctrl).Intersects(registers.CtrlNametableY)
	if !beforeNmi && p.Ctrl.HasEnableNMI() && p.Status.Intersects(registers.Vblank) {
		p.interruptCh <- interrupts.NMI
	}
}

func (p *PPU) WriteMask(data byte) {
	p.Mask = bitflags.Flags(data)
}

func (p *PPU) WriteOamAddr(data byte) {
	p.OamAddr = data
}

func (p *PPU) WriteOam(data byte) {
	p.Oam[p.OamAddr] = data
	p.OamAddr += 1
}

func (p *PPU) WriteOamDma(data [0x100]byte) {
	for _, data := range data {
		p.WriteOam(data)
	}
}

func (p *PPU) ReadOam() byte {
	return p.Oam[p.OamAddr]
}

func (p *PPU) WriteScroll(data byte) {
	if p.AddrLatch {
		p.TmpAddr.WriteScrollY(data)
	} else {
		p.TmpAddr.WriteScrollX(data)
		p.FineX = data & 7
	}
	p.AddrLatch = !p.AddrLatch
}

func (p *PPU) ReadStatus() byte {
	defer func() {
		p.Status.Remove(registers.Vblank)
	}()
	p.AddrLatch = false
	return byte(p.Status)
}

func (p *PPU) Write(data byte) {
	addr := p.Addr.Get()
	switch {
	case addr < 0x2000:
		p.mapper.WriteMem(addr, data)
	case 0x2000 <= addr && addr < 0x3000:
		addr := p.MirrorVramAddr(addr)
		p.Vram[addr] = data
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
		p.Palette[addr] = data
	default:
		log.WithField("address", fmt.Sprintf("%02X", addr)).
			Error("unexpected write to mirrored space")
	}
	p.Addr.Increment(p.Ctrl.VramAddr())
}

func (p *PPU) Read() byte {
	addr := p.Addr.Get()
	p.Addr.Increment(p.Ctrl.VramAddr())

	val := p.ReadAddr(addr)
	if addr < 0x3000 {
		val, p.ReadBuf = p.ReadBuf, val
	}
	return val
}

func (p *PPU) ReadAddr(addr uint16) byte {
	switch {
	case addr < 0x2000:
		return p.mapper.ReadMem(addr)
	case 0x2000 <= addr && addr < 0x3000:
		addr := p.MirrorVramAddr(addr)
		return p.Vram[addr]
	case 0x3000 <= addr && addr < 0x3F00:
		log.WithField("address", fmt.Sprintf("$%02X", addr)).
			Error("bad PPU read")
		return 0
	case 0x3F00 <= addr && addr < 0x4000:
		addr &= 0x3F1F
		switch addr {
		case 0x3F10, 0x3F14, 0x3F18, 0x3F1C:
			addr -= 0x10
		}
		addr -= 0x3F00
		return p.Palette[addr]
	default:
		log.WithField("address", fmt.Sprintf("%02X", addr)).
			Error("unexpected access to mirrored space")
		return 0
	}
}

var MirrorLookup = [...][4]uint16{
	cartridge.Horizontal:  {0, 0, 1, 1},
	cartridge.Vertical:    {0, 1, 0, 1},
	cartridge.SingleLower: {0, 0, 0, 0},
	cartridge.SingleUpper: {1, 1, 1, 1},
	cartridge.FourScreen:  {0, 1, 2, 3},
}

func (p *PPU) MirrorVramAddr(addr uint16) uint16 {
	addr &= 0xFFF
	nameTable := addr / 0x400
	offset := addr % 0x400
	return offset + 0x400*MirrorLookup[p.mapper.Cartridge().Mirror][nameTable]
}

func (p *PPU) updateSpriteOverflow() {
	size := int(p.Ctrl.SpriteSize())
	var count uint
	for i := 0; i < len(p.Oam)/4; i += 1 {
		i := i * 4
		tileY := p.Oam[i]
		row := int(p.Scanline) - int(tileY)
		if row < 0 || row >= size {
			continue
		}
		count += 1
	}
	count &= 0b1111
	if count == 8 {
		p.Status.Insert(registers.SpriteOverflow)
	}
}

func (p *PPU) Step() bool {
	p.Cycles += 1

	switch {
	case p.Cycles == 257:
		p.updateSpriteOverflow()
	case p.Cycles > 340:
		if p.SpriteZeroHit(p.Cycles) {
			p.Status.Insert(registers.SpriteZeroHit)
		}

		p.Cycles = 0
		p.Scanline += 1

		switch {
		case p.Scanline == 241:
			p.Status.Insert(registers.Vblank)
			p.Status.Remove(registers.SpriteZeroHit)
			if p.Ctrl.HasEnableNMI() {
				p.interruptCh <- interrupts.NMI
			}
		case p.Scanline >= 262:
			p.Scanline = 0
			p.Status.Remove(registers.Vblank | registers.SpriteOverflow | registers.SpriteZeroHit)
			return true
		}
	}
	return false
}

func (p *PPU) SpriteZeroHit(cycle uint) bool {
	x := p.Oam[3]
	y := p.Oam[0]
	return uint16(y) == p.Scanline && uint(x) <= cycle && p.Mask.Intersects(registers.SpriteEnable)
}

func (p *PPU) Reset() {
	p.Cycles = 0
	p.Scanline = 0
	p.WriteCtrl(0)
	p.WriteMask(0)
	p.WriteOamAddr(0)
	p.Addr = registers.AddrRegister{}
	p.TmpAddr = registers.AddrRegister{}
	p.AddrLatch = false
}

func (p *PPU) Interrupt() <-chan interrupts.Interrupt {
	return p.interruptCh
}
