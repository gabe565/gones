package interrupt

const (
	ResetVector = 0xFFFC
	NmiVector   = 0xFFFA
	IrqVector   = 0xFFFE
)

type Interruptible interface {
	AddNmi()
}

type Stallable interface {
	AddStall(uint16)
}
