package ppu

const MaxSprites = 8

type SpriteData struct {
	Count      uint8
	Patterns   [MaxSprites]uint32
	Positions  [MaxSprites]byte
	Priorities [MaxSprites]byte
	Indexes    [MaxSprites]byte
}

func (p *PPU) evaluateSprites() {
	height := int(p.Ctrl.SpriteSize())
	var count uint8

	for i := 0; i < 64; i += 1 {
		sprite := p.Oam[i*4 : i*4+4 : i*4+4]
		y := sprite[0]
		a := sprite[2]
		x := sprite[3]

		row := int(p.Scanline) - int(y)
		if row < 0 || row >= height {
			continue
		}

		if count < 8 {
			p.SpriteData.Patterns[count] = p.fetchSpritePattern(sprite[1], a, row)
			p.SpriteData.Positions[count] = x
			p.SpriteData.Priorities[count] = (a >> 5) & 1
			p.SpriteData.Indexes[count] = byte(i)
		}

		count += 1
	}

	if count > 8 {
		count = 8
		p.Status.SpriteOverflow = true
	}

	p.SpriteData.Count = count
}

func (p *PPU) fetchSpritePattern(tile, attributes byte, row int) uint32 {
	var addr uint16

	if p.Ctrl.SpriteHeight {
		if attributes&0x80 == 0x80 {
			row = 15 - row
		}
		table := tile & 1
		tile &= 0xFE
		if row > 7 {
			tile += 1
			row -= 8
		}
		addr = 0x1000*uint16(table) + uint16(tile)*16 + uint16(row)
	} else {
		if attributes&0x80 == 0x80 {
			row = 7 - row
		}
		addr = p.Ctrl.SpriteTileAddr() + uint16(tile)*16 + uint16(row)
	}

	a := (attributes & 3) << 2
	tileLo := p.ReadDataAddr(addr)
	tileHi := p.ReadDataAddr(addr + 8)
	var data uint32

	for i := 0; i < 8; i += 1 {
		var p1, p2 byte
		if attributes&0x40 == 0x40 {
			p1 = (tileLo & 1) << 0
			p2 = (tileHi & 1) << 1
			tileLo >>= 1
			tileHi >>= 1
		} else {
			p1 = (tileLo & 0x80) >> 7
			p2 = (tileHi & 0x80) >> 6
			tileLo <<= 1
			tileHi <<= 1
		}
		data <<= 4
		data |= uint32(a | p1 | p2)
	}

	return data
}

func (p *PPU) spritePixel() (byte, byte) {
	if !p.Mask.SpriteEnable {
		return 0, 0
	}

	for i := uint8(0); i < p.SpriteData.Count; i += 1 {
		offset := int(p.Cycles) - 1 - int(p.SpriteData.Positions[i])
		if offset < 0 || offset > 7 {
			continue
		}

		offset = 7 - offset
		color := byte((p.SpriteData.Patterns[i] >> byte(offset*4)) & 0xF)
		if color%4 == 0 {
			continue
		}
		return i, color
	}

	return 0, 0
}
