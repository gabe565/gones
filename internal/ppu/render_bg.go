package ppu

type BgTile struct {
	NametableByte byte
	AttrByte      byte
	LoByte        byte
	HiByte        byte
	Data          uint64
}

func (p *PPU) fetchNametableByte() byte {
	addr := 0x2000 | p.Addr.Get()&0xFFF
	return p.ReadDataAddr(addr)
}

func (p *PPU) fetchAttrTableByte() byte {
	addr := 0x23C0 | uint16(p.Addr.CoarseY>>2)<<3 | uint16(p.Addr.CoarseX>>2)
	if p.Addr.NametableY {
		addr |= 1 << 11
	}
	if p.Addr.NametableX {
		addr |= 1 << 10
	}
	var attrByte byte
	attrByte = p.ReadDataAddr(addr)
	if p.Addr.CoarseY&2 != 0 {
		attrByte >>= 4
	}
	if p.Addr.CoarseX&2 != 0 {
		attrByte >>= 2
	}
	attrByte &= 3
	attrByte <<= 2
	return attrByte
}

func (p *PPU) fetchLoTileByte() byte {
	addr := uint16(p.BgTile.NametableByte)<<4 + uint16(p.Addr.FineY)
	if p.Ctrl.BgTileSelect {
		addr += 1 << 12
	}
	return p.ReadDataAddr(addr)
}

func (p *PPU) fetchHiTileByte() byte {
	addr := uint16(p.BgTile.NametableByte)<<4 + uint16(p.Addr.FineY) + 8
	if p.Ctrl.BgTileSelect {
		addr += 1 << 12
	}
	return p.ReadDataAddr(addr)
}

func (p *PPU) storeTileData() {
	var data uint32
	for i := uint8(0); i < 8; i++ {
		p1 := (p.BgTile.LoByte & 0x80) >> 7
		p2 := (p.BgTile.HiByte & 0x80) >> 6
		p.BgTile.LoByte <<= 1
		p.BgTile.HiByte <<= 1
		data <<= 4
		data |= uint32(p.BgTile.AttrByte | p1 | p2)
	}
	p.BgTile.Data |= uint64(data)
}

func (p *PPU) bgPixel(x int) byte {
	if !p.Mask.BackgroundEnable || (x < 8 && !p.Mask.BgLeftColEnable) {
		return 0
	}

	data := uint32(p.BgTile.Data >> 32)
	data >>= (7 - p.FineX) * 4
	data &= 0xF
	return byte(data)
}
