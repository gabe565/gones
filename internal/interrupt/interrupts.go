package interrupt

const (
	ResetVector = 0xFFFC
	NMIVector   = 0xFFFA
	IRQVector   = 0xFFFE
)

type NMI interface {
	AddNMI()
}

type Stall interface {
	AddStall(uint16)
}
