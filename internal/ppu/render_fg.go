package ppu

import (
	"github.com/gabe565/gones/internal/consts"
	"github.com/vmihailenco/msgpack/v5"
)

type SpriteData struct {
	Count      uint8
	limit      uint8
	Patterns   []uint32
	Positions  []byte
	Priorities []byte
	Indexes    []byte
}

var _ msgpack.CustomDecoder = &SpriteData{}

func (s *SpriteData) DecodeMsgpack(dec *msgpack.Decoder) error {
	type tmpSpriteData SpriteData
	if err := dec.Decode((*tmpSpriteData)(s)); err != nil {
		return err
	}
	limit := int(s.limit)
	if len(s.Patterns) > limit {
		clear(s.Patterns[limit:])
		s.Patterns = s.Patterns[:limit:limit]
	} else if len(s.Patterns) < limit {
		s.Patterns = append(s.Patterns, make([]uint32, limit-len(s.Patterns))...)
	}
	if len(s.Positions) > limit {
		clear(s.Positions[limit:])
		s.Positions = s.Positions[:limit:limit]
	} else if len(s.Positions) < limit {
		s.Positions = append(s.Positions, make([]byte, limit-len(s.Positions))...)
	}
	if len(s.Priorities) > limit {
		clear(s.Priorities[limit:])
		s.Priorities = s.Priorities[:limit:limit]
	} else if len(s.Priorities) < limit {
		s.Priorities = append(s.Priorities, make([]byte, limit-len(s.Priorities))...)
	}
	if len(s.Indexes) > limit {
		clear(s.Indexes[limit:])
		s.Indexes = s.Indexes[:limit:limit]
	} else if len(s.Indexes) < limit {
		s.Indexes = append(s.Indexes, make([]byte, limit-len(s.Indexes))...)
	}
	return nil
}

func (p *PPU) evaluateSprites() {
	height := int(p.Ctrl.SpriteSize())
	var count uint8

	var prevY, prevTile uint8
	consecutive := uint8(1)
	for i := range consts.PPUOAMSize / 4 {
		i *= 4
		y, tile, a, x := p.OAM[i], p.OAM[i+1], p.OAM[i+2], p.OAM[i+3]

		row := p.Scanline - int(y)
		if row < 0 || row >= height {
			continue
		}

		if count < p.SpriteData.limit {
			p.SpriteData.Patterns[count] = p.fetchSpritePattern(tile, a, row)
			p.SpriteData.Positions[count] = x
			p.SpriteData.Priorities[count] = (a >> 5) & 1
			p.SpriteData.Indexes[count] = byte(i)

			if y == prevY && tile == prevTile {
				consecutive++
			}
			prevY, prevTile = y, tile
		}

		count++
	}

	if count > consts.PPUSpriteLimit {
		p.Status.SpriteOverflow = true
	}

	switch {
	case consecutive >= consts.PPUSpriteLimit:
		// Force original limit when masking effect is active
		// See https://nesdev.org/wiki/Sprite_overflow_games#Detecting_masking_effects
		p.SpriteData.Count = consts.PPUSpriteLimit
	case count > p.SpriteData.limit:
		p.SpriteData.Count = p.SpriteData.limit
	default:
		p.SpriteData.Count = count
	}
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
			tile++
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

	for range 8 {
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

func (p *PPU) spritePixel(x int) (byte, byte) {
	if !p.Mask.SpriteEnable || (x < 8 && !p.Mask.SpriteLeftColEnable) {
		return 0, 0
	}

	for i := range p.SpriteData.Count {
		offset := x - int(p.SpriteData.Positions[i])
		if offset < 0 || offset > 7 {
			continue
		}

		color := p.SpriteData.Patterns[i] >> ((7 - offset) * 4) & 0xF
		if color%4 == 0 {
			continue
		}
		return i, byte(color)
	}

	return 0, 0
}
