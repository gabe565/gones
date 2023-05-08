package ppu

type BgTile struct {
	NametableByte byte
	AttrByte      byte
	LoByte        byte
	HiByte        byte
	Data          uint64
}

func (p *PPU) fetchNametableByte() {
	addr := 0x2000 | p.Addr.Get()&0xFFF
	p.BgTile.NametableByte = p.ReadDataAddr(addr)
}

func (p *PPU) fetchAttrTableByte() {
	addr := 0x23C0 | uint16(p.Addr.CoarseY>>2)<<3 | uint16(p.Addr.CoarseX>>2)
	if p.Addr.NametableY {
		addr |= 1 << 11
	}
	if p.Addr.NametableX {
		addr |= 1 << 10
	}
	p.BgTile.AttrByte = p.ReadDataAddr(addr)
	if p.Addr.CoarseY&2 != 0 {
		p.BgTile.AttrByte >>= 4
	}
	if p.Addr.CoarseX&2 != 0 {
		p.BgTile.AttrByte >>= 2
	}
	p.BgTile.AttrByte &= 3
	p.BgTile.AttrByte <<= 2
}

func (p *PPU) fetchLoTileByte() {
	addr := uint16(p.BgTile.NametableByte)<<4 + uint16(p.Addr.FineY)
	if p.Ctrl.BgTileSelect {
		addr += 1 << 12
	}
	p.BgTile.LoByte = p.ReadDataAddr(addr)
}

func (p *PPU) fetchHiTileByte() {
	addr := uint16(p.BgTile.NametableByte)<<4 + uint16(p.Addr.FineY) + 8
	if p.Ctrl.BgTileSelect {
		addr += 1 << 12
	}
	p.BgTile.HiByte = p.ReadDataAddr(addr)
}

func (p *PPU) storeTileData() {
	var data uint32
	for i := 0; i < 8; i++ {
		p1 := (p.BgTile.LoByte & 0x80) >> 7
		p2 := (p.BgTile.HiByte & 0x80) >> 6
		p.BgTile.LoByte <<= 1
		p.BgTile.HiByte <<= 1
		data <<= 4
		data |= uint32(p.BgTile.AttrByte | p1 | p2)
	}
	p.BgTile.Data |= uint64(data)
}

func (p *PPU) bgPixel() byte {
	if !p.Mask.BackgroundEnable {
		return 0
	}

	data := uint32(p.BgTile.Data >> 32)
	data >>= (7 - p.FineX) * 4
	data &= 0xF
	return byte(data)
}
