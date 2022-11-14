package interrupts

import "github.com/gabe565/gones/internal/bitflags"

type Interrupt struct {
	Name       string
	Cycles     uint16
	VectorAddr uint16
	Mask       bitflags.Flags
}

func (i Interrupt) Error() string {
	return i.Name
}

var NMI = Interrupt{
	Name:       "NMI",
	Cycles:     2,
	VectorAddr: 0xFFFA,
	Mask:       0b00100000,
}

var IRQ = Interrupt{
	Name:       "IRQ",
	Cycles:     2,
	VectorAddr: 0xFFFE,
}
