package cpu

type OpCode struct {
	Code     uint8
	Mnemonic string
	Len      uint8
	Cycles   uint8
	Mode     AddressingMode
}

var OpCodes = []OpCode{
	{0x00, "BRK", 1, 7, NoneAddressing},
	{0xaa, "TAX", 1, 2, NoneAddressing},
	{0xe8, "INX", 1, 2, NoneAddressing},

	{0xa9, "LDA", 2, 2, Immediate},
	{0xa5, "LDA", 2, 3, ZeroPage},
	{0xb5, "LDA", 2, 4, ZeroPageX},
	{0xad, "LDA", 3, 4, Absolute},
	{0xbd, "LDA", 3, 4 /*+1 if page crossed*/, AbsoluteX},
	{0xb9, "LDA", 3, 4 /*+1 if page crossed*/, AbsoluteY},
	{0xa1, "LDA", 2, 6, IndirectX},
	{0xb1, "LDA", 2, 5 /*f page crossed*/, IndirectY},

	{0x85, "STA", 2, 3, ZeroPage},
	{0x95, "STA", 2, 4, ZeroPageX},
	{0x8d, "STA", 3, 4, Absolute},
	{0x9d, "STA", 3, 5, AbsoluteX},
	{0x99, "STA", 3, 5, AbsoluteY},
	{0x81, "STA", 2, 6, IndirectX},
	{0x91, "STA", 2, 6, IndirectY},
}

func OpCodeMap() map[uint8]OpCode {
	codes := make(map[uint8]OpCode)
	for _, opcode := range OpCodes {
		codes[opcode.Code] = opcode
	}
	return codes
}
