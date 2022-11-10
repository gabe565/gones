package ppu

import (
	"image"
)

const (
	Width  = 256
	Height = 240
)

func (p *PPU) Render() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	bank := p.ctrl.BkndPatternAddr()

	for i := uint16(0); i < 0x03C0; i += 1 {
		tile := uint16(p.vram[i])
		tileCol := i % 32
		tileRow := i / 32
		tiles := p.chr[bank+tile*16 : bank+tile*16+16]
		palette := p.bgPalette(tileCol, tileRow)

		for y := 0; y < 8; y += 1 {
			upper := tiles[y]
			lower := tiles[y+8]

			for x := 0; x < 8; x += 1 {
				value := (1&upper)<<1 | (1 & lower)
				upper >>= 1
				lower >>= 1
				c := SystemPalette[palette[value]]
				img.Set(int(tileCol)*8+x, int(tileRow)*8+y, c)
			}
		}
	}

	for i := 0; i < len(p.oam); i += 4 {
		tileIdx := p.oam[i+1]
		tileX := p.oam[i+3]
		tileY := p.oam[i]

		flipVertical := p.oam[i+2]>>7&1 == 1
		flipHorizonal := p.oam[i+2]>>6&1 == 1

		paletteIdx := p.oam[i+2] & 0b11
		spritePalette := p.spritePalette(paletteIdx)

		bank := p.ctrl.SprtPatternAddr()

		tile := p.chr[bank+uint16(tileIdx)*16 : bank+uint16(tileIdx)*16+16]

		for y := 0; y < 8; y += 1 {
			upper := tile[y]
			lower := tile[y+8]

			for x := 0; x < 8; x += 1 {
				value := (1&lower)<<1 | (1 & upper)
				upper >>= 1
				lower >>= 1
				if value == 0 {
					continue
				}
				c := SystemPalette[spritePalette[value]]

				flippedX := int(tileX)
				if flipHorizonal {
					flippedX += 7 - x
				} else {
					flippedX += x
				}

				flippedY := int(tileY)
				if flipVertical {
					flippedY += 7 - y
				} else {
					flippedY += y
				}

				img.Set(flippedX, flippedY, c)
			}
		}
	}

	return img
}

func (p *PPU) bgPalette(col, row uint16) [4]byte {
	attrTableIdx := row/4*8 + col/4
	attrByte := p.vram[attrTableIdx+0x03C0]

	paletteIdx := attrByte
	switch [2]byte{byte(col % 4 / 2), byte(row % 4 / 2)} {
	case [2]byte{0, 0}:
		//
	case [2]byte{1, 0}:
		paletteIdx >>= 2
	case [2]byte{0, 1}:
		paletteIdx >>= 4
	case [2]byte{1, 1}:
		paletteIdx >>= 6
	default:
		panic("invalid bg palette")
	}
	paletteIdx &= 0b11

	paletteStart := paletteIdx*4 + 1
	return [4]byte{
		p.palette[0],
		p.palette[paletteStart],
		p.palette[paletteStart+1],
		p.palette[paletteStart+2],
	}
}

func (p *PPU) spritePalette(idx byte) [4]byte {
	start := idx*4 + 0x11
	return [4]byte{
		9,
		p.palette[start],
		p.palette[start+1],
		p.palette[start+2],
	}
}
