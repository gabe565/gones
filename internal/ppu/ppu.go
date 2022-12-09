package ppu

import (
	"fmt"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/ppu/registers"
	log "github.com/sirupsen/logrus"
	"image"
)

func New(mapper cartridge.Mapper) *PPU {
	return &PPU{
		mapper: mapper,
		image:  image.NewRGBA(image.Rect(0, 0, Width, TrimmedHeight)),
	}
}

type PPU struct {
	mapper cartridge.Mapper

	Ctrl      registers.Control
	Mask      registers.Mask
	Status    registers.Status
	Addr      registers.Address
	TmpAddr   registers.Address
	AddrLatch bool
	FineX     byte
	Vram      [0x800]byte

	OamAddr byte
	Oam     [0x100]byte
	Palette [0x20]byte

	Scanline uint16
	Cycles   uint

	nmiOccurred bool

	ReadBuf byte
	image   *image.RGBA

	BgTile     BgTile
	SpriteData SpriteData
	OddFrame   bool
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
	beforeNmi := p.Ctrl.EnableNMI
	p.Ctrl.Set(data)
	p.TmpAddr.NametableX = p.Ctrl.NametableX
	p.TmpAddr.NametableY = p.Ctrl.NametableY
	if !beforeNmi && p.Status.Vblank && p.Ctrl.EnableNMI {
		p.nmiOccurred = true
	}
}

func (p *PPU) WriteMask(data byte) {
	p.Mask.Set(data)
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
	status := p.Status.Get()
	p.Status.Vblank = false
	p.nmiOccurred = false
	p.AddrLatch = false
	return status
}

func (p *PPU) Write(data byte) {
	addr := p.Addr.Get() & 0x3FFF
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
		p.writePalette(addr%32, data)
	default:
		log.WithField("address", fmt.Sprintf("%02X", addr)).
			Error("unexpected write to mirrored space")
	}
	p.Addr.Increment(p.Ctrl.VramAddr())
}

func (p *PPU) Read() byte {
	addr := p.Addr.Get() & 0x3FFF
	p.Addr.Increment(p.Ctrl.VramAddr())

	val := p.ReadAddr(addr)
	if addr%0x4000 < 0x3F00 {
		val, p.ReadBuf = p.ReadBuf, val
	} else if addr < 0x4000 {
		p.ReadBuf = p.ReadAddr(addr - 0x1000)
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
		return p.readPalette(addr % 32)
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

func (p *PPU) tick() {
	if p.Mask.BackgroundEnable || p.Mask.SpriteEnable {
		if p.OddFrame && p.Scanline == 261 && p.Cycles == 339 {
			p.Cycles = 0
			p.Scanline = 0
			p.OddFrame = !p.OddFrame
			return
		}
	}

	p.Cycles += 1

	if p.Cycles > 340 {
		p.Cycles = 0
		p.Scanline += 1
		if p.Scanline > 261 {
			p.Scanline = 0
			p.OddFrame = !p.OddFrame
		}
	}
}

func (p *PPU) Step() bool {
	p.tick()

	renderingEnabled := p.Mask.BackgroundEnable || p.Mask.SpriteEnable
	preLine := p.Scanline == 261
	visibleLine := p.Scanline < 240
	renderLine := preLine || visibleLine
	preFetchCycle := p.Cycles >= 321 && p.Cycles <= 336
	visibleCycle := p.Cycles >= 1 && p.Cycles <= 256
	fetchCycle := preFetchCycle || visibleCycle

	if visibleLine && visibleCycle {
		p.renderPixel()
	}

	if renderingEnabled {
		// Background
		if renderLine && fetchCycle {
			p.BgTile.Data <<= 4

			switch p.Cycles % 8 {
			case 1:
				p.fetchNametableByte()
			case 3:
				p.fetchAttrTableByte()
			case 5:
				p.fetchLoTileByte()
			case 7:
				p.fetchHiTileByte()
			case 0:
				p.storeTileData()
			}
		}

		if preLine && p.Cycles >= 280 && p.Cycles <= 304 {
			p.copyAddrY()
		}

		if renderLine {
			if fetchCycle && p.Cycles%8 == 0 {
				p.incrementX()
			}
			if p.Cycles == 256 {
				p.incrementY()
			}
			if p.Cycles == 257 {
				p.copyAddrX()
			}
		}

		// Sprite
		if p.Cycles == 257 {
			if visibleLine {
				p.evaluateSprites()
			} else {
				p.SpriteData.Count = 0
			}
		}
	}

	if p.Scanline == 241 && p.Cycles == 1 {
		p.Status.Vblank = true
		if p.Ctrl.EnableNMI {
			p.nmiOccurred = true
		}
	}

	if preLine && p.Cycles == 1 {
		p.Status.Vblank = false
		p.Status.SpriteOverflow = false
		p.Status.SpriteZeroHit = false
		return true
	}

	return false
}

func (p *PPU) Reset() {
	p.Cycles = 0
	p.Scanline = 0
	p.WriteCtrl(0)
	p.WriteMask(0)
	p.WriteOamAddr(0)
	p.Addr = registers.Address{}
	p.TmpAddr = registers.Address{}
	p.BgTile = BgTile{}
	p.SpriteData = SpriteData{}
	p.OddFrame = false
	p.AddrLatch = false
}

func (p *PPU) NMI() bool {
	if p.nmiOccurred && p.Ctrl.EnableNMI {
		p.nmiOccurred = false
		return true
	}
	return false
}

func (p *PPU) readPalette(addr uint16) byte {
	if addr >= 16 && addr%4 == 0 {
		addr -= 16
	}
	return p.Palette[addr]
}

func (p *PPU) writePalette(addr uint16, data byte) {
	if addr >= 16 && addr%4 == 0 {
		addr -= 16
	}
	p.Palette[addr] = data
}
