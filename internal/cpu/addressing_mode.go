package cpu

type AddressingMode uint8

const (
	Immediate AddressingMode = iota
	ZeroPage
	ZeroPageX
	ZeroPageY
	Absolute
	AbsoluteX
	AbsoluteY
	IndirectX
	IndirectY
	NoneAddressing
)
