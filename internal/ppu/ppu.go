package ppu

import (
	"image"
	"log/slog"

	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/config"
	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/interrupt"
	"github.com/gabe565/gones/internal/log"
	"github.com/gabe565/gones/internal/memory"
	"github.com/gabe565/gones/internal/ppu/palette"
	"github.com/gabe565/gones/internal/ppu/registers"
)

type CPU interface {
	memory.Read8
	memory.HasCycles
	interrupt.NMI
	interrupt.Stall
}

func New(conf *config.Config, mapper cartridge.Mapper) *PPU {
	rect := conf.UI.Overscan.Rect()
	spriteLimit := uint8(8)
	if conf.UI.RemoveSpriteLimit {
		spriteLimit = consts.PPUOAMSize / 4
	}
	return &PPU{
		offsets:       rect.Min,
		mapper:        mapper,
		image:         image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy())),
		Cycles:        21,
		systemPalette: &palette.Default,
		SpriteData: SpriteData{
			limit:      spriteLimit,
			Patterns:   make([]uint32, spriteLimit),
			Positions:  make([]byte, spriteLimit),
			Priorities: make([]byte, spriteLimit),
			Indexes:    make([]byte, spriteLimit),
		},
	}
}

type PPU struct {
	mapper  cartridge.Mapper
	cpu     CPU
	offsets image.Point

	Ctrl      registers.Control
	Mask      registers.Mask
	Status    registers.Status
	Addr      registers.Address
	TmpAddr   registers.Address
	AddrLatch bool
	FineX     byte
	VRAM      [0x800]byte `msgpack:"alias:Vram"`

	OAMAddr       byte                    `msgpack:"alias:OamAddr"`
	OAM           [consts.PPUOAMSize]byte `msgpack:"alias:Oam"`
	systemPalette *palette.Palette
	Palette       [0x20]byte

	Scanline int
	Cycles   int
	VblRace  bool

	NMIOffset uint8 `msgpack:"alias:NmiOffset"`

	ReadBuf    byte
	OpenBus    byte
	RenderDone bool
	image      *image.RGBA

	BgTile     BgTile
	SpriteData SpriteData
	OddFrame   bool
}

func (p *PPU) WriteAddr(data byte) {
	if p.AddrLatch {
		p.TmpAddr.WriteLo(data)
		p.Addr = p.TmpAddr

		if mapper, ok := p.mapper.(cartridge.MapperOnVRAMAddr); ok {
			mapper.OnVRAMAddr(p.Addr)
		}
	} else {
		p.TmpAddr.WriteHi(data)
	}
	p.AddrLatch = !p.AddrLatch
}

func (p *PPU) WriteCtrl(data byte) {
	p.Ctrl.Set(data)
	p.TmpAddr.NametableX = p.Ctrl.NametableX
	p.TmpAddr.NametableY = p.Ctrl.NametableY
	p.updateNMI()
}

func (p *PPU) WriteMask(data byte) {
	p.Mask.Set(data)
	p.UpdatePalette(data)
}

func (p *PPU) UpdatePalette(data byte) {
	switch data & (registers.MaskEmphasizeRed | registers.MaskEmphasizeGreen | registers.MaskEmphasizeBlue) {
	case 0:
		p.systemPalette = &palette.Default
	case registers.MaskEmphasizeRed:
		p.systemPalette = &palette.EmphasizeR
	case registers.MaskEmphasizeGreen:
		p.systemPalette = &palette.EmphasizeG
	case registers.MaskEmphasizeBlue:
		p.systemPalette = &palette.EmphasizeB
	case registers.MaskEmphasizeRed | registers.MaskEmphasizeGreen:
		p.systemPalette = &palette.EmphasizeRG
	case registers.MaskEmphasizeRed | registers.MaskEmphasizeBlue:
		p.systemPalette = &palette.EmphasizeRB
	case registers.MaskEmphasizeGreen | registers.MaskEmphasizeBlue:
		p.systemPalette = &palette.EmphasizeGB
	case registers.MaskEmphasizeRed | registers.MaskEmphasizeGreen | registers.MaskEmphasizeBlue:
		p.systemPalette = &palette.EmphasizeRGB
	}
}

func (p *PPU) WriteOamAddr(data byte) {
	p.OAMAddr = data
}

func (p *PPU) WriteOam(data byte) {
	p.OAM[p.OAMAddr] = data
	p.OAMAddr++
}

func (p *PPU) WriteOamDma(data [0x100]byte) {
	for _, data := range data {
		p.WriteOam(data)
	}
}

func (p *PPU) ReadOam() byte {
	data := p.OAM[p.OAMAddr]
	if p.OAMAddr&3 == 2 {
		// Exclude unused bytes
		data &= 0xE3
	}
	return data
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
	p.VblRace = false
	p.updateNMI()
	p.AddrLatch = false
	if p.Scanline == 241 && p.Cycles == 0 {
		p.VblRace = true
	}
	return status
}

func (p *PPU) WriteData(data byte) {
	addr := p.Addr.Get() % 0x4000
	switch {
	case addr < 0x2000:
		p.mapper.WriteMem(addr, data)
	case 0x2000 <= addr && addr < 0x3F00:
		addr := p.MirrorVRAMAddr(addr)
		p.VRAM[addr] = data
	case 0x3F00 <= addr && addr < 0x4000:
		p.writePalette(addr%32, data)
	default:
		slog.Error("Invalid write to mirrored space", "addr", log.HexAddr(addr))
	}
	p.Addr.Increment(p.Ctrl.VRAMAddr())
	if mapper, ok := p.mapper.(cartridge.MapperOnVRAMAddr); ok {
		mapper.OnVRAMAddr(p.Addr)
	}
}

func (p *PPU) ReadData() byte {
	addr := p.Addr.Get() % 0x4000
	if p.Mask.RenderingEnabled() && (p.Scanline == 261 || p.Scanline < 240) {
		// If rendering enabled, increment Coarse X and Y
		// https://www.nesdev.org/wiki/PPU_scrolling#$2007_reads_and_writes
		p.Addr.IncrementX()
		p.Addr.IncrementY()
	} else {
		// Else increment by 1 or 32
		p.Addr.Increment(p.Ctrl.VRAMAddr())
	}

	val := p.ReadDataAddr(addr)
	if addr < 0x3F00 {
		val, p.ReadBuf = p.ReadBuf, val
	} else if addr < 0x4000 {
		p.ReadBuf = p.ReadDataAddr(addr - 0x1000)
		val |= p.OpenBus & 0xC0
	}

	if mapper, ok := p.mapper.(cartridge.MapperOnVRAMAddr); ok {
		mapper.OnVRAMAddr(p.Addr)
	}
	return val
}

func (p *PPU) ReadDataAddr(addr uint16) byte {
	addr %= 0x4000
	switch {
	case addr < 0x2000:
		return p.mapper.ReadMem(addr)
	case 0x2000 <= addr && addr < 0x3F00:
		addr := p.MirrorVRAMAddr(addr)
		return p.VRAM[addr]
	case 0x3F00 <= addr && addr < 0x4000:
		return p.readPalette(addr % 32)
	default:
		slog.Error("Invalid access from mirrored space", "addr", log.HexAddr(addr))
		return 0
	}
}

func (p *PPU) ReadMem(addr uint16) byte {
	switch addr {
	case 0x2000, 0x2001, 0x2003, 0x2005, 0x2006, 0x4014:
		//
	case 0x2002:
		p.OpenBus &^= 0xE0
		p.OpenBus |= p.ReadStatus()
	case 0x2004:
		p.OpenBus = p.ReadOam()
	case 0x2007:
		p.OpenBus = p.ReadData()
	default:
		slog.Error("Invalid PPU read", "addr", log.HexAddr(addr))
	}
	return p.OpenBus
}

func (p *PPU) WriteMem(addr uint16, data byte) {
	switch addr {
	case 0x2000:
		p.WriteCtrl(data)
	case 0x2001:
		p.WriteMask(data)
	case 0x2002:
		//
	case 0x2003:
		p.WriteOamAddr(data)
	case 0x2004:
		p.WriteOam(data)
	case 0x2005:
		p.WriteScroll(data)
	case 0x2006:
		p.WriteAddr(data)
	case 0x2007:
		p.WriteData(data)
	case 0x4014:
		hi := uint16(data) << 8
		for i := range uint16(256) {
			p.WriteOam(p.cpu.ReadMem(hi + i))
		}
		if p.cpu.GetCycles()%2 == 1 {
			p.cpu.AddStall(514)
		} else {
			p.cpu.AddStall(513)
		}
	default:
		slog.Error("Invalid PPU write", "addr", log.HexAddr(addr))
	}
	p.OpenBus = data
}

//nolint:gochecknoglobals
var MirrorLookup = [...][4]uint16{
	cartridge.Horizontal:  {0, 0, 1, 1},
	cartridge.Vertical:    {0, 1, 0, 1},
	cartridge.SingleLower: {0, 0, 0, 0},
	cartridge.SingleUpper: {1, 1, 1, 1},
	cartridge.FourScreen:  {0, 1, 2, 3},
}

func (p *PPU) MirrorVRAMAddr(addr uint16) uint16 {
	addr &= 0xFFF
	nameTable := addr / 0x400
	offset := addr % 0x400
	return offset + 0x400*MirrorLookup[p.mapper.Cartridge().Mirror][nameTable]
}

func (p *PPU) tick() {
	if p.NMIOffset != 0 {
		p.NMIOffset--
		if p.NMIOffset == 0 {
			p.cpu.AddNMI()
		} else if p.NMIOffset >= 12 {
			if !p.Status.Vblank || !p.Ctrl.EnableNMI {
				p.NMIOffset = 0
			}
		}
	}

	if p.Mask.RenderingEnabled() {
		if p.OddFrame && p.Scanline == 261 && p.Cycles == 339 {
			p.Cycles = 0
			p.Scanline = 0
			p.OddFrame = !p.OddFrame
			return
		}
	}

	if p.Cycles < 340 {
		p.Cycles++
	} else {
		p.Cycles = 0
		if p.Scanline < 261 {
			p.Scanline++
		} else {
			p.Scanline = 0
			p.OddFrame = !p.OddFrame
		}
	}
}

func (p *PPU) Step(render bool) {
	p.tick()

	preLine := p.Scanline == 261
	visibleLine := p.Scanline < 240
	renderLine := preLine || visibleLine
	visibleCycle := 1 <= p.Cycles && p.Cycles <= 256

	if visibleLine && visibleCycle {
		p.renderPixel(render)
	}

	renderingEnabled := p.Mask.RenderingEnabled()
	if renderingEnabled {
		// Background
		if renderLine {
			preFetchCycle := 321 <= p.Cycles && p.Cycles <= 336
			fetchCycle := preFetchCycle || visibleCycle

			if fetchCycle {
				p.BgTile.Data <<= 4

				switch p.Cycles % 8 {
				case 1:
					p.BgTile.NametableByte = p.fetchNametableByte()
				case 3:
					p.BgTile.AttrByte = p.fetchAttrTableByte()
				case 5:
					p.BgTile.LoByte = p.fetchLoTileByte()
				case 7:
					p.BgTile.HiByte = p.fetchHiTileByte()
				case 0:
					p.storeTileData()
					p.Addr.IncrementX()
				}
			}

			if preLine && 280 <= p.Cycles && p.Cycles <= 304 {
				p.Addr.LoadY(p.TmpAddr)
			}

			switch p.Cycles {
			case 256:
				p.Addr.IncrementY()
			case 257:
				p.Addr.LoadX(p.TmpAddr)
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

	switch p.Cycles {
	case 1:
		if p.Scanline == 241 && !p.VblRace {
			p.Status.Vblank = true
			p.updateNMI()
			p.RenderDone = true
		} else if preLine {
			p.Status.Vblank = false
			p.updateNMI()
			p.Status.SpriteOverflow = false
			p.Status.SpriteZeroHit = false
			p.VblRace = false
			p.OpenBus = 0
		}
	case 280:
		if renderLine && renderingEnabled {
			if mapper, ok := p.mapper.(cartridge.MapperOnScanline); ok {
				mapper.OnScanline()
			}
		}
	}
}

func (p *PPU) Reset() {
	p.WriteCtrl(0)
	p.WriteMask(0)
	p.OddFrame = false
	p.AddrLatch = false
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

func (p *PPU) updateNMI() {
	nmi := p.Status.Vblank && p.Ctrl.EnableNMI
	if nmi && !p.Status.PrevVblank {
		p.NMIOffset = 14
	}
	p.Status.PrevVblank = nmi
}

func (p *PPU) SetCPU(c CPU) {
	p.cpu = c
}

func (p *PPU) SetMapper(m cartridge.Mapper) {
	p.mapper = m
}

func (p *PPU) Width() int {
	return p.image.Bounds().Dx()
}

func (p *PPU) Height() int {
	return p.image.Bounds().Dy()
}
