package ppu

import (
	"image"
)

const (
	Width  = 256
	Height = 240
)

func (p *PPU) Image() *image.RGBA {
	return p.image
}

func (p *PPU) renderPixel(render bool) {
	x := p.Cycles - 1
	y := p.Scanline

	bgPixel := p.bgPixel(x)
	bgEnabled := bgPixel%4 != 0

	i, spritePixel := p.spritePixel(x)
	spriteEnabled := spritePixel%4 != 0

	var colorIdx byte
	if bgEnabled {
		if spriteEnabled {
			if p.SpriteData.Indexes[i] == 0 && x < 255 {
				p.Status.SpriteZeroHit = true
			}
			if p.SpriteData.Priorities[i] == 0 {
				colorIdx = spritePixel | 0x10
			} else {
				colorIdx = bgPixel
			}
		} else {
			colorIdx = bgPixel
		}
	} else if spriteEnabled {
		colorIdx = spritePixel | 0x10
	}

	if render {
		colorIdx = p.readPalette(uint16(colorIdx)) % 64
		if p.Mask.Grayscale {
			colorIdx &= 0x30
		}

		c := p.systemPalette.RGBA[colorIdx]
		p.image.SetRGBA(x, y, c)
	}
}
