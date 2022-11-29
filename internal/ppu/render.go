package ppu

import (
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/ppu/registers"
	log "github.com/sirupsen/logrus"
	"image"
)

const (
	Width         = 256
	Height        = 240
	TrimHeight    = 8
	TrimmedHeight = Height - 2*TrimHeight
)

func (p *PPU) Render() *image.RGBA {
	main, second := p.getNametables()
	scrollX := int(p.Scroll.X)
	scrollY := int(p.Scroll.Y)

	if p.Mask.Intersects(registers.BackgroundEnable) {
		p.RenderNametable(
			main,
			-scrollX,
			-scrollY,
		)

		if scrollX > 0 {
			p.RenderNametable(
				second,
				Width-scrollX,
				0,
			)
		} else if scrollY > 0 {
			p.RenderNametable(
				second,
				0,
				Height-scrollY,
			)
		}
	} else {
		c := SystemPalette[p.Palette[0]]
		for y := 0; y < TrimmedHeight; y += 1 {
			for x := 0; x < Width; x += 1 {
				p.image.Set(x, y, c)
			}
		}
	}

	if p.Mask.Intersects(registers.SpriteEnable) {
		for i := len(p.Oam) - 4; i >= 0; i -= 4 {
			tileIdx := p.Oam[i+1]
			tileX := p.Oam[i+3]
			tileY := p.Oam[i] + 1

			flipVertical := p.Oam[i+2]>>7&1 == 1
			flipHorizonal := p.Oam[i+2]>>6&1 == 1

			paletteIdx := p.Oam[i+2] & 0b11
			spritePalette := p.spritePalette(paletteIdx)

			bank := p.Ctrl.SpriteTileAddr()

			tile := p.mapper.Cartridge().Chr[bank+uint16(tileIdx)*16 : bank+uint16(tileIdx)*16+16]

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

					flippedY := int(tileY) - TrimHeight
					if flipVertical {
						flippedY += 7 - y
					} else {
						flippedY += y
					}

					p.image.Set(flippedX, flippedY, c)
				}
			}
		}
	}

	return p.image
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
		log.Panic("invalid bg palette")
	}
	paletteIdx &= 0b11

	paletteStart := paletteIdx*4 + 1
	return [4]byte{
		p.Palette[0],
		p.Palette[paletteStart],
		p.Palette[paletteStart+1],
		p.Palette[paletteStart+2],
	}
}

func (p *PPU) spritePalette(idx byte) [4]byte {
	start := idx*4 + 0x11
	return [4]byte{
		9,
		p.Palette[start],
		p.Palette[start+1],
		p.Palette[start+2],
	}
}

func (p *PPU) RenderNametable(nameTable []byte, shiftX, shiftY int) {
	bank := p.Ctrl.BgTileAddr()

	attrTable := nameTable[0x3C0:0x400]

	for i := uint16(0); i < 0x3C0; i += 1 {
		tileCol := i % 32
		tileRow := i / 32
		tileIdx := uint16(nameTable[i])
		tile := p.mapper.Cartridge().Chr[(bank + tileIdx*16):(bank + tileIdx*16 + 16)]
		palette := p.bgPalette(attrTable, tileCol, tileRow)

		for y := 0; y < 8; y += 1 {
			upper := tile[y]
			lower := tile[y+8]

			for x := 7; x >= 0; x -= 1 {
				value := (1&lower)<<1 | (1 & upper)
				upper >>= 1
				lower >>= 1
				c := SystemPalette[palette[value]]

				pxlX := int(tileCol)*8 + x + shiftX
				pxlY := int(tileRow)*8 + y + shiftY - TrimHeight
				p.image.Set(pxlX, pxlY, c)
			}
		}
	}
}

func (p *PPU) getNametables() ([]byte, []byte) {
	type match struct {
		mirror        cartridge.Mirror
		nametableAddr uint16
	}

	switch (match{p.mapper.Cartridge().Mirror, p.Ctrl.NametableAddr()}) {
	case match{cartridge.Vertical, 0x2000},
		match{cartridge.Vertical, 0x2800},
		match{cartridge.Horizontal, 0x2000},
		match{cartridge.Horizontal, 0x2400}:
		{
			return p.Vram[:0x400], p.Vram[0x400:0x800]
		}
	case match{cartridge.Vertical, 0x2400},
		match{cartridge.Vertical, 0x2C00},
		match{cartridge.Horizontal, 0x2800},
		match{cartridge.Horizontal, 0x2C00}:
		{
			return p.Vram[0x400:0x800], p.Vram[:0x400]
		}
	default:
		log.Panic(p.mapper.Cartridge().Mirror.String() + " mirroring unsupported")
		return []byte{}, []byte{}
	}
}
