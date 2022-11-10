package registers

import "github.com/gabe565/gones/internal/bitflags"

type Control bitflags.Flags

const (
	Nametable1 bitflags.Flags = 1 << iota
	Nametable2
	VramAddIncrement
	SpritePatternAddr
	BackgroundPatternAddr
	SpriteSize
	MasterSlaveSelect
	GenerateNmi
)

func (c Control) VramAddrIncrement() byte {
	if bitflags.Flags(c).Has(VramAddIncrement) {
		return 32
	} else {
		return 1
	}
}

func (c Control) SprtPatternAddr() uint16 {
	if bitflags.Flags(c).Has(SpritePatternAddr) {
		return 0x1000
	} else {
		return 1
	}
}

func (c Control) BkndPatternAddr() uint16 {
	if bitflags.Flags(c).Has(BackgroundPatternAddr) {
		return 0x1000
	} else {
		return 0
	}
}

func (c Control) SpriteSize() byte {
	if bitflags.Flags(c).Has(SpriteSize) {
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

func (c Control) GenerateVblankNmi() bool {
	return bitflags.Flags(c).Has(GenerateNmi)
}

func (c Control) NametableAddr() uint16 {
	switch c & 0b11 {
	case 0:
		return 0x2000
	case 1:
		return 0x2400
	case 2:
		return 0x2800
	case 3:
		return 0x2c00
	default:
		panic("invalid control register")
	}
}
