package interrupts

type Interrupt struct {
	Name       string
	Cycles     uint16
	VectorAddr uint16
}

func (i Interrupt) Error() string {
	return i.Name
}

var NMI = Interrupt{
	Name:       "NMI",
	Cycles:     7,
	VectorAddr: 0xFFFA,
}

var IRQ = Interrupt{
	Name:       "IRQ",
	Cycles:     7,
	VectorAddr: 0xFFFE,
}

type Interruptible interface {
	AddInterrupt(*Interrupt)
}
