package cpu

// OpCode defines an opcode and its parameters.
//
// See [6502 Instruction Reference].
//
// [6502 Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html
type OpCode struct {
	Code     uint8
	Mnemonic string
	Len      uint8
	Cycles   uint8
	Mode     AddressingMode
}

// OpCodes is a list of supported opcodes.
//
// See [6502 Instruction Reference].
//
// [6502 Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html
var OpCodes = []OpCode{
	{0x00, "BRK", 1, 7, NoneAddressing},
	{0xAA, "TAX", 1, 2, NoneAddressing},
	{0xE8, "INX", 1, 2, NoneAddressing},

	{0xA9, "LDA", 2, 2, Immediate},
	{0xA5, "LDA", 2, 3, ZeroPage},
	{0xB5, "LDA", 2, 4, ZeroPageX},
	{0xAD, "LDA", 3, 4, Absolute},
	{0xBD, "LDA", 3, 4 /*+1 if page crossed*/, AbsoluteX},
	{0xB9, "LDA", 3, 4 /*+1 if page crossed*/, AbsoluteY},
	{0xA1, "LDA", 2, 6, IndirectX},
	{0xB1, "LDA", 2, 5 /*f page crossed*/, IndirectY},

	{0x85, "STA", 2, 3, ZeroPage},
	{0x95, "STA", 2, 4, ZeroPageX},
	{0x8D, "STA", 3, 4, Absolute},
	{0x9D, "STA", 3, 5, AbsoluteX},
	{0x99, "STA", 3, 5, AbsoluteY},
	{0x81, "STA", 2, 6, IndirectX},
	{0x91, "STA", 2, 6, IndirectY},

	{0x86, "STX", 2, 3, ZeroPage},
	{0x96, "STX", 2, 4, ZeroPageY},
	{0x8E, "STX", 3, 4, Absolute},

	{0x84, "STY", 2, 3, ZeroPage},
	{0x94, "STY", 2, 4, ZeroPageX},
	{0x8C, "STY", 3, 4, Absolute},

	{0xA8, "TAY", 1, 2, NoneAddressing},

	{0xBA, "TSX", 1, 2, NoneAddressing},
}

// OpCodeMap converts OpCodes into a map with the code as the key.
func OpCodeMap() map[uint8]OpCode {
	codes := make(map[uint8]OpCode)
	for _, opcode := range OpCodes {
		codes[opcode.Code] = opcode
	}
	return codes
}
