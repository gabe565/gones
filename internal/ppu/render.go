package ppu

import (
	"image"
	"image/color"
)

const (
	Width  = 256
	Height = 240
)

func (p *PPU) Render() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	bank := p.ctrl.BkndPatternAddr()

	for i := 0; i < 0x03C0; i += 1 {
		tile := uint16(p.vram[i])
		tileX := i % 32
		tileY := i / 32
		tiles := p.chr[bank+tile*16 : bank+tile*16+16]

		for y := 0; y < 8; y += 1 {
			upper := tiles[y]
			lower := tiles[y+8]

			for x := 0; x < 8; x += 1 {
				value := (1&upper)<<1 | (1 & lower)
				upper >>= 1
				lower >>= 1
				var c color.RGBA
				switch value {
				case 0:
					c = SystemPalette[0x01]
				case 1:
					c = SystemPalette[0x23]
				case 2:
					c = SystemPalette[0x27]
				case 3:
					c = SystemPalette[0x30]
				}
				img.Set(tileX*8+x, tileY*8+y, c)
			}
		}
	}

	return img
}
