package registers

import (
	"github.com/gabe565/gones/internal/bitflags"
)

type Control bitflags.Flags

const (
	Nametable1 bitflags.Flags = 1 << iota
	Nametable2
	IncrementMode
	SpriteTileSelect
	BgTileSelect
	SpriteHeight
	MasterSlaveSelect
	EnableNMI
)

func (c Control) VramAddr() byte {
	if bitflags.Flags(c).Intersects(IncrementMode) {
		return 32
	} else {
		return 1
	}
}

func (c Control) SpriteTileAddr() uint16 {
	if bitflags.Flags(c).Intersects(SpriteTileSelect) {
		return 0x1000
	} else {
		return 0
	}
}

func (c Control) BgTileAddr() uint16 {
	if bitflags.Flags(c).Intersects(BgTileSelect) {
		return 0x1000
	} else {
		return 0
	}
}

func (c Control) SpriteSize() byte {
	if bitflags.Flags(c).Intersects(SpriteHeight) {
		return 16
	} else {
		return 8
	}
}

func (c Control) MasterSlaveSelect() byte {
	if bitflags.Flags(c).Intersects(MasterSlaveSelect) {
		return 1
	} else {
		return 0
	}
}

func (c Control) HasEnableNMI() bool {
	return bitflags.Flags(c).Intersects(EnableNMI)
}

func (c Control) NametableAddr() uint16 {
	return 0x2000 | uint16(c)&0b11<<10
}
