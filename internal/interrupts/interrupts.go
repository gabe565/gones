package interrupts

const (
	ResetVector = 0xFFFC
	NmiVector   = 0xFFFA
	IrqVector   = 0xFFFE
)

type Interruptible interface {
	AddNmi()
	AddIrq()
	ClearIrq()
}

type Stallable interface {
	AddStall(uint16)
}
