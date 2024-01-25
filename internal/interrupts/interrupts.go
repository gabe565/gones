package interrupts

type Interrupt struct {
	Name          string
	Cycles        uint16
	StackProhibit bool
	VectorAddr    uint16
}

func (i Interrupt) Error() string {
	return i.Name
}

var NMI = Interrupt{
	Name:       "NMI",
	Cycles:     7,
	VectorAddr: 0xFFFA,
}

var Reset = Interrupt{
	Name:          "Reset",
	StackProhibit: true,
	VectorAddr:    0xFFFC,
}

var IRQ = Interrupt{
	Name:       "IRQ",
	Cycles:     7,
	VectorAddr: 0xFFFE,
}

type Interruptible interface {
	AddNmi()
	AddIrq()
	ClearIrq()
}

type Stallable interface {
	AddStall(uint16)
}
