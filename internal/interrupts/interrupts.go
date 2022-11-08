package interrupts

import "github.com/gabe565/gones/internal/bitflags"

type Interrupt struct {
	Cycles     uint16
	VectorAddr uint16
	Mask       bitflags.Flags
}

var NMI = Interrupt{
	Cycles:     2,
	VectorAddr: 0xFFFA,
	Mask:       0b00100000,
}
