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
	{0x90, "BCC", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing},
	{0xB0, "BCS", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing},
	{0xF0, "BEQ", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing},
	{0x30, "BMI", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing},
	{0xD0, "BNE", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing},
	{0x10, "BPL", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing},
	{0x50, "BVC", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing},
	{0x70, "BVS", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing},
	{0x18, "CLC", 1, 2, NoneAddressing},
	{0xD8, "CLD", 1, 2, NoneAddressing},
	{0x58, "CLI", 1, 2, NoneAddressing},
	{0xB8, "CLV", 1, 2, NoneAddressing},
	{0xCA, "DEX", 1, 2, NoneAddressing},
	{0x88, "DEY", 1, 2, NoneAddressing},
	{0xE8, "INX", 1, 2, NoneAddressing},
	{0xC8, "INY", 1, 2, NoneAddressing},
	{0x20, "JSR", 3, 6, NoneAddressing},
	{0xEA, "NOP", 1, 2, NoneAddressing},
	{0x48, "PHA", 1, 3, NoneAddressing},
	{0x08, "PHP", 1, 3, NoneAddressing},
	{0x68, "PLA", 1, 4, NoneAddressing},
	{0x28, "PLP", 1, 4, NoneAddressing},
	{0x40, "RTI", 1, 6, NoneAddressing},
	{0x60, "RTS", 1, 6, NoneAddressing},
	{0x38, "SEC", 1, 2, NoneAddressing},
	{0xF8, "SED", 1, 2, NoneAddressing},
	{0x78, "SEI", 1, 2, NoneAddressing},
	{0xAA, "TAX", 1, 2, NoneAddressing},
	{0xA8, "TAY", 1, 2, NoneAddressing},
	{0xBA, "TSX", 1, 2, NoneAddressing},
	{0x8A, "TXA", 1, 2, NoneAddressing},
	{0x9A, "TXS", 1, 2, NoneAddressing},
	{0x98, "TYA", 1, 2, NoneAddressing},

	{0x69, "ADC", 2, 2, Immediate},
	{0x65, "ADC", 2, 2, ZeroPage},
	{0x75, "ADC", 2, 2, ZeroPageX},
	{0x6D, "ADC", 3, 3, Absolute},
	{0x7D, "ADC", 3, 3 /*+1 if page crossed*/, AbsoluteX},
	{0x79, "ADC", 3, 3 /*+1 if page crossed*/, AbsoluteY},
	{0x61, "ADC", 2, 2, IndirectX},
	{0x71, "ADC", 2, 2 /*+1 if page crossed*/, IndirectY},

	{0x29, "AND", 2, 2, Immediate},
	{0x25, "AND", 2, 3, ZeroPage},
	{0x35, "AND", 2, 4, ZeroPageX},
	{0x2D, "AND", 3, 4, Absolute},
	{0x3D, "AND", 3, 4 /*+1 if page crossed*/, AbsoluteX},
	{0x39, "AND", 3, 4 /*+1 if page crossed*/, AbsoluteY},
	{0x21, "AND", 2, 6, IndirectX},
	{0x31, "AND", 2, 5 /*f page crossed*/, IndirectY},

	{0x0A, "ASL", 1, 2, NoneAddressing},
	{0x06, "ASL", 2, 5, ZeroPage},
	{0x16, "ASL", 2, 6, ZeroPageX},
	{0x0E, "ASL", 3, 6, Absolute},
	{0x1E, "ASL", 3, 7 /*+1 if page crossed*/, AbsoluteX},

	{0x24, "BIT", 2, 3, ZeroPage},
	{0x2C, "BIT", 3, 4, Absolute},

	{0xC9, "CMP", 2, 2, Immediate},
	{0xC5, "CMP", 2, 3, ZeroPage},
	{0xD5, "CMP", 2, 4, ZeroPageX},
	{0xCD, "CMP", 3, 4, Absolute},
	{0xDD, "CMP", 3, 4 /*+1 if page crossed*/, AbsoluteX},
	{0xD9, "CMP", 3, 4 /*+1 if page crossed*/, AbsoluteY},
	{0xC1, "CMP", 2, 6, IndirectX},
	{0xD1, "CMP", 2, 5 /*f page crossed*/, IndirectY},

	{0xE0, "CPX", 2, 2, Immediate},
	{0xE4, "CPX", 2, 3, ZeroPage},
	{0xEC, "CPX", 3, 4, Absolute},

	{0xC0, "CPY", 2, 2, Immediate},
	{0xC4, "CPY", 2, 3, ZeroPage},
	{0xCC, "CPY", 3, 4, Absolute},

	{0xC6, "DEC", 2, 5, ZeroPage},
	{0xD6, "DEC", 2, 6, ZeroPageX},
	{0xCE, "DEC", 3, 6, Absolute},
	{0xDE, "DEC", 3, 7, AbsoluteX},

	{0xA9, "LDA", 2, 2, Immediate},
	{0xA5, "LDA", 2, 3, ZeroPage},
	{0xB5, "LDA", 2, 4, ZeroPageX},
	{0xAD, "LDA", 3, 4, Absolute},
	{0xBD, "LDA", 3, 4 /*+1 if page crossed*/, AbsoluteX},
	{0xB9, "LDA", 3, 4 /*+1 if page crossed*/, AbsoluteY},
	{0xA1, "LDA", 2, 6, IndirectX},
	{0xB1, "LDA", 2, 5 /*f page crossed*/, IndirectY},

	{0x49, "EOR", 2, 2, Immediate},
	{0x45, "EOR", 2, 3, ZeroPage},
	{0x55, "EOR", 2, 4, ZeroPageX},
	{0x4D, "EOR", 3, 4, Absolute},
	{0x5D, "EOR", 3, 4 /*+1 if page crossed*/, AbsoluteX},
	{0x59, "EOR", 3, 4 /*+1 if page crossed*/, AbsoluteY},
	{0x41, "EOR", 2, 6, IndirectX},
	{0x51, "EOR", 2, 5 /*f page crossed*/, IndirectY},

	{0xE6, "INC", 2, 5, ZeroPage},
	{0xF6, "INC", 2, 6, ZeroPageX},
	{0xEE, "INC", 3, 6, Absolute},
	{0xFE, "INC", 3, 7, AbsoluteX},

	{0x4C, "JMP", 3, 3, Absolute}, // Acts as immediate
	{0x6C, "JMP", 3, 5, Indirect}, // Indirect with 6502 bug

	{0xA2, "LDX", 2, 2, Immediate},
	{0xA6, "LDX", 2, 3, ZeroPage},
	{0xB6, "LDX", 2, 4, ZeroPageY},
	{0xAE, "LDX", 3, 4, Absolute},
	{0xBE, "LDX", 3, 4 /*+1 if page crossed*/, AbsoluteY},

	{0xA0, "LDY", 2, 2, Immediate},
	{0xA4, "LDY", 2, 3, ZeroPage},
	{0xB4, "LDY", 2, 4, ZeroPageX},
	{0xAC, "LDY", 3, 4, Absolute},
	{0xBC, "LDY", 3, 4 /*+1 if page crossed*/, AbsoluteX},

	{0x4A, "LSR", 1, 2, Accumulator},
	{0x46, "LSR", 2, 5, ZeroPage},
	{0x56, "LSR", 2, 6, ZeroPageX},
	{0x4E, "LSR", 3, 6, Absolute},
	{0x5E, "LSR", 3, 7, AbsoluteX},

	{0x09, "ORA", 2, 2, Immediate},
	{0x05, "ORA", 2, 3, ZeroPage},
	{0x15, "ORA", 2, 4, ZeroPageX},
	{0x0D, "ORA", 3, 4, Absolute},
	{0x1D, "ORA", 3, 4 /*+1 if page crossed*/, AbsoluteX},
	{0x19, "ORA", 3, 4 /*+1 if page crossed*/, AbsoluteY},
	{0x01, "ORA", 2, 6, IndirectX},
	{0x11, "ORA", 2, 5 /*f page crossed*/, IndirectY},

	{0x2A, "ROL", 1, 2, Accumulator},
	{0x26, "ROL", 2, 5, ZeroPage},
	{0x36, "ROL", 2, 6, ZeroPageX},
	{0x2E, "ROL", 3, 6, Absolute},
	{0x3E, "ROL", 3, 7, AbsoluteX},

	{0x6A, "ROR", 1, 2, Accumulator},
	{0x66, "ROR", 2, 5, ZeroPage},
	{0x76, "ROR", 2, 6, ZeroPageX},
	{0x6E, "ROR", 3, 6, Absolute},
	{0x7E, "ROR", 3, 7, AbsoluteX},

	{0xE9, "SBC", 2, 2, Immediate},
	{0xE5, "SBC", 2, 3, ZeroPage},
	{0xF5, "SBC", 2, 4, ZeroPageX},
	{0xED, "SBC", 3, 4, Absolute},
	{0xFD, "SBC", 3, 4 /*+1 if page crossed*/, AbsoluteX},
	{0xF9, "SBC", 3, 4 /*+1 if page crossed*/, AbsoluteY},
	{0xE1, "SBC", 2, 6, IndirectX},
	{0xF1, "SBC", 2, 5 /*+1 if page crossed*/, IndirectY},

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
}

// OpCodeMap converts OpCodes into a map with the code as the key.
func OpCodeMap() map[uint8]OpCode {
	codes := make(map[uint8]OpCode)
	for _, opcode := range OpCodes {
		codes[opcode.Code] = opcode
	}
	return codes
}
