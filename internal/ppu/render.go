package ppu

import (
	"github.com/gabe565/gones/internal/ppu/registers"
	"image"
)

const (
	Width         = 256
	Height        = 240
	TrimHeight    = 8
	TrimmedHeight = Height - 2*TrimHeight
)

func (p *PPU) Image() *image.RGBA {
	return p.image
}

func (p *PPU) renderPixel() {
	x := int(p.Cycles - 1)
	y := int(p.Scanline)

	bgPixel := p.bgPixel()
	if x < 8 && !p.Mask.Intersects(registers.BgLeftColEnable) {
		bgPixel = 0
	}

	i, sprite := p.spritePixel()
	if x < 8 && !p.Mask.Intersects(registers.SpriteLeftColEnable) {
		sprite = 0
	}

	b := bgPixel%4 != 0
	s := sprite%4 != 0

	var colorIdx byte

	switch {
	case !b && !s:
		colorIdx = 0
	case !b && s:
		colorIdx = sprite | 0x10
	case b && !s:
		colorIdx = bgPixel
	default:
		if p.SpriteData.Indexes[i] == 0 && x < 255 {
			p.Status.Insert(registers.SpriteZeroHit)
		}
		if p.SpriteData.Priorities[i] == 0 {
			colorIdx = sprite | 0x10
		} else {
			colorIdx = bgPixel
		}
	}

	colorIdx = p.readPalette(uint16(colorIdx)) % 64
	if p.Mask.Intersects(registers.Grayscale) {
		colorIdx &= 0x30
	}

	c := SystemPalette[colorIdx]
	p.image.SetRGBA(x, y-8, c)
}
