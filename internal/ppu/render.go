package ppu

import (
	"github.com/faiface/pixel"
	"github.com/gabe565/gones/internal/cartridge"
	"image"
	"image/color"
)

const (
	Width  = 256
	Height = 240
)

func (p *PPU) Render() *pixel.PictureData {
	pic := pixel.MakePictureData(pixel.R(0, 0, Width, Height))

	main, second := p.getNametables()
	scrollX := int(p.scroll.X)
	scrollY := int(p.scroll.Y)

	p.RenderNametable(
		pic,
		main,
		image.Rect(scrollX, scrollY, Width, Height),
		-scrollX,
		-scrollY,
	)

	if scrollX > 0 {
		p.RenderNametable(
			pic,
			second,
			image.Rect(0, 0, scrollX, Height),
			Width-scrollX,
			0,
		)
	} else if scrollY > 0 {
		p.RenderNametable(
			pic,
			second,
			image.Rect(0, 0, Width, scrollY),
			0,
			Height-scrollY,
		)
	}

	for i := len(p.oam) - 4; i >= 0; i -= 4 {
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

			for x := 7; x >= 0; x -= 1 {
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

				setPixel(pic, flippedX, flippedY, c)
			}
		}
	}

	return pic
}

func (p *PPU) bgPalette(attrTable []byte, col, row uint16) [4]byte {
	attrTableIdx := row/4*8 + col/4
	attrByte := attrTable[attrTableIdx]

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

func (p *PPU) RenderNametable(pic *pixel.PictureData, nameTable []byte, viewport image.Rectangle, shiftX, shiftY int) {
	bank := p.ctrl.BkndPatternAddr()

	attrTable := nameTable[0x3C0:0x400]

	for i := uint16(0); i < 0x3C0; i += 1 {
		tileCol := i % 32
		tileRow := i / 32
		tileIdx := uint16(nameTable[i])
		tile := p.chr[(bank + tileIdx*16):(bank + tileIdx*16 + 16)]
		palette := p.bgPalette(attrTable, tileCol, tileRow)

		for y := 0; y < 8; y += 1 {
			upper := tile[y]
			lower := tile[y+8]

			for x := 7; x >= 0; x -= 1 {
				value := (1&lower)<<1 | (1 & upper)
				upper >>= 1
				lower >>= 1
				c := SystemPalette[palette[value]]

				pxlX := int(tileCol)*8 + x
				pxlY := int(tileRow)*8 + y
				point := image.Point{X: pxlX, Y: pxlY}
				if point.In(viewport) {
					setPixel(pic, shiftX+pxlX, shiftY+pxlY, c)
				}
			}
		}
	}
}

func (p *PPU) getNametables() ([]byte, []byte) {
	type match struct {
		mirror        cartridge.Mirror
		nametableAddr uint16
	}

	switch (match{p.mirroring, p.ctrl.NametableAddr()}) {
	case match{cartridge.Vertical, 0x2000},
		match{cartridge.Vertical, 0x2800},
		match{cartridge.Horizontal, 0x2000},
		match{cartridge.Horizontal, 0x2400}:
		{
			return p.vram[:0x400], p.vram[0x400:0x800]
		}
	case match{cartridge.Vertical, 0x2400},
		match{cartridge.Vertical, 0x2C00},
		match{cartridge.Horizontal, 0x2800},
		match{cartridge.Horizontal, 0x2C00}:
		{
			return p.vram[0x400:0x800], p.vram[:0x400]
		}
	default:
		panic(p.mirroring.String() + " mirroring unsupported")
	}
}

func setPixel(pic *pixel.PictureData, x, y int, c color.RGBA) {
	y = Height - y - 1
	offset := y*Width + x
	if offset >= 0 && offset < len(pic.Pix) {
		pic.Pix[offset] = c
	}
}
