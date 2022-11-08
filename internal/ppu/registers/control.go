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
		return 32
	} else {
		return 1
	}
}

func (c Control) BkndPatternAddr() uint16 {
	if bitflags.Flags(c).Has(BackgroundPatternAddr) {
		return 32
	} else {
		return 1
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
