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
	if bitflags.Flags(c).Has(IncrementMode) {
		return 32
	} else {
		return 1
	}
}

func (c Control) SpriteTileAddr() uint16 {
	if bitflags.Flags(c).Has(SpriteTileSelect) {
		return 0x1000
	} else {
		return 1
	}
}

func (c Control) BgTileAddr() uint16 {
	if bitflags.Flags(c).Has(BgTileSelect) {
		return 0x1000
	} else {
		return 0
	}
}

func (c Control) SpriteSize() byte {
	if bitflags.Flags(c).Has(SpriteHeight) {
		return 16
	} else {
		return 8
	}
}

func (c Control) MasterSlaveSelect() byte {
	if bitflags.Flags(c).Has(MasterSlaveSelect) {
		return 1
	} else {
		return 0
	}
}

func (c Control) HasEnableNMI() bool {
	return bitflags.Flags(c).Has(EnableNMI)
}

func (c Control) NametableAddr() uint16 {
	return 0x2000 | uint16(c)&0b11<<10
}
