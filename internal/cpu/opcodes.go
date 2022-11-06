package cpu

import "fmt"

// OpCode defines an opcode and its parameters.
//
// See [6502 Instruction Reference].
//
// [6502 Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html
type OpCode struct {
	Code     byte
	Mnemonic string
	Len      uint8
	Cycles   uint8
	Mode     AddressingMode
	Exec     Instruction
}

func (o OpCode) String() string {
	return fmt.Sprintf(
		"{$%02X %s %d %d %s}",
		o.Code,
		o.Mnemonic,
		o.Len,
		o.Cycles,
		o.Mode,
	)
}

// OpCodes is a list of supported opcodes.
//
// See [6502 Instruction Reference].
//
// [6502 Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html
var OpCodes = []OpCode{
	{0x00, "BRK", 1, 7, NoneAddressing, brk},
	{0x90, "BCC", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing, bcc},
	{0xB0, "BCS", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing, bcs},
	{0xF0, "BEQ", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing, beq},
	{0x30, "BMI", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing, bmi},
	{0xD0, "BNE", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing, bne},
	{0x10, "BPL", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing, bpl},
	{0x50, "BVC", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing, bvc},
	{0x70, "BVS", 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, NoneAddressing, bvs},
	{0x18, "CLC", 1, 2, NoneAddressing, clc},
	{0xD8, "CLD", 1, 2, NoneAddressing, cld},
	{0x58, "CLI", 1, 2, NoneAddressing, cli},
	{0xB8, "CLV", 1, 2, NoneAddressing, clv},
	{0xCA, "DEX", 1, 2, NoneAddressing, dex},
	{0x88, "DEY", 1, 2, NoneAddressing, dey},
	{0xE8, "INX", 1, 2, NoneAddressing, inx},
	{0xC8, "INY", 1, 2, NoneAddressing, iny},
	{0x20, "JSR", 3, 6, NoneAddressing, jsr},
	{0xEA, "NOP", 1, 2, NoneAddressing, nop},
	{0x48, "PHA", 1, 3, NoneAddressing, pha},
	{0x08, "PHP", 1, 3, NoneAddressing, php},
	{0x68, "PLA", 1, 4, NoneAddressing, pla},
	{0x28, "PLP", 1, 4, NoneAddressing, plp},
	{0x40, "RTI", 1, 6, NoneAddressing, rti},
	{0x60, "RTS", 1, 6, NoneAddressing, rts},
	{0x38, "SEC", 1, 2, NoneAddressing, sec},
	{0xF8, "SED", 1, 2, NoneAddressing, sed},
	{0x78, "SEI", 1, 2, NoneAddressing, sei},
	{0xAA, "TAX", 1, 2, NoneAddressing, tax},
	{0xA8, "TAY", 1, 2, NoneAddressing, tay},
	{0xBA, "TSX", 1, 2, NoneAddressing, tsx},
	{0x8A, "TXA", 1, 2, NoneAddressing, txa},
	{0x9A, "TXS", 1, 2, NoneAddressing, txs},
	{0x98, "TYA", 1, 2, NoneAddressing, tya},

	{0x69, "ADC", 2, 2, Immediate, adc},
	{0x65, "ADC", 2, 2, ZeroPage, adc},
	{0x75, "ADC", 2, 2, ZeroPageX, adc},
	{0x6D, "ADC", 3, 3, Absolute, adc},
	{0x7D, "ADC", 3, 3 /*+1 if page crossed*/, AbsoluteX, adc},
	{0x79, "ADC", 3, 3 /*+1 if page crossed*/, AbsoluteY, adc},
	{0x61, "ADC", 2, 2, IndirectX, adc},
	{0x71, "ADC", 2, 2 /*+1 if page crossed*/, IndirectY, adc},

	{0x29, "AND", 2, 2, Immediate, and},
	{0x25, "AND", 2, 3, ZeroPage, and},
	{0x35, "AND", 2, 4, ZeroPageX, and},
	{0x2D, "AND", 3, 4, Absolute, and},
	{0x3D, "AND", 3, 4 /*+1 if page crossed*/, AbsoluteX, and},
	{0x39, "AND", 3, 4 /*+1 if page crossed*/, AbsoluteY, and},
	{0x21, "AND", 2, 6, IndirectX, and},
	{0x31, "AND", 2, 5 /*+1 if page crossed*/, IndirectY, and},

	{0x0A, "ASL", 1, 2, Accumulator, asl},
	{0x06, "ASL", 2, 5, ZeroPage, asl},
	{0x16, "ASL", 2, 6, ZeroPageX, asl},
	{0x0E, "ASL", 3, 6, Absolute, asl},
	{0x1E, "ASL", 3, 7 /*+1 if page crossed*/, AbsoluteX, asl},

	{0x24, "BIT", 2, 3, ZeroPage, bit},
	{0x2C, "BIT", 3, 4, Absolute, bit},

	{0xC9, "CMP", 2, 2, Immediate, cmp},
	{0xC5, "CMP", 2, 3, ZeroPage, cmp},
	{0xD5, "CMP", 2, 4, ZeroPageX, cmp},
	{0xCD, "CMP", 3, 4, Absolute, cmp},
	{0xDD, "CMP", 3, 4 /*+1 if page crossed*/, AbsoluteX, cmp},
	{0xD9, "CMP", 3, 4 /*+1 if page crossed*/, AbsoluteY, cmp},
	{0xC1, "CMP", 2, 6, IndirectX, cmp},
	{0xD1, "CMP", 2, 5 /*+1 if page crossed*/, IndirectY, cmp},

	{0xE0, "CPX", 2, 2, Immediate, cpx},
	{0xE4, "CPX", 2, 3, ZeroPage, cpx},
	{0xEC, "CPX", 3, 4, Absolute, cpx},

	{0xC0, "CPY", 2, 2, Immediate, cpy},
	{0xC4, "CPY", 2, 3, ZeroPage, cpy},
	{0xCC, "CPY", 3, 4, Absolute, cpy},

	{0xC6, "DEC", 2, 5, ZeroPage, dec},
	{0xD6, "DEC", 2, 6, ZeroPageX, dec},
	{0xCE, "DEC", 3, 6, Absolute, dec},
	{0xDE, "DEC", 3, 7, AbsoluteX, dec},

	{0xA9, "LDA", 2, 2, Immediate, lda},
	{0xA5, "LDA", 2, 3, ZeroPage, lda},
	{0xB5, "LDA", 2, 4, ZeroPageX, lda},
	{0xAD, "LDA", 3, 4, Absolute, lda},
	{0xBD, "LDA", 3, 4 /*+1 if page crossed*/, AbsoluteX, lda},
	{0xB9, "LDA", 3, 4 /*+1 if page crossed*/, AbsoluteY, lda},
	{0xA1, "LDA", 2, 6, IndirectX, lda},
	{0xB1, "LDA", 2, 5 /*+1 if page crossed*/, IndirectY, lda},

	{0x49, "EOR", 2, 2, Immediate, eor},
	{0x45, "EOR", 2, 3, ZeroPage, eor},
	{0x55, "EOR", 2, 4, ZeroPageX, eor},
	{0x4D, "EOR", 3, 4, Absolute, eor},
	{0x5D, "EOR", 3, 4 /*+1 if page crossed*/, AbsoluteX, eor},
	{0x59, "EOR", 3, 4 /*+1 if page crossed*/, AbsoluteY, eor},
	{0x41, "EOR", 2, 6, IndirectX, eor},
	{0x51, "EOR", 2, 5 /*+1 if page crossed*/, IndirectY, eor},

	{0xE6, "INC", 2, 5, ZeroPage, inc},
	{0xF6, "INC", 2, 6, ZeroPageX, inc},
	{0xEE, "INC", 3, 6, Absolute, inc},
	{0xFE, "INC", 3, 7, AbsoluteX, inc},

	{0x4C, "JMP", 3, 3, Absolute, jmp}, // Acts as immediate
	{0x6C, "JMP", 3, 5, Indirect, jmp}, // Indirect with 6502 bug

	{0xA2, "LDX", 2, 2, Immediate, ldx},
	{0xA6, "LDX", 2, 3, ZeroPage, ldx},
	{0xB6, "LDX", 2, 4, ZeroPageY, ldx},
	{0xAE, "LDX", 3, 4, Absolute, ldx},
	{0xBE, "LDX", 3, 4 /*+1 if page crossed*/, AbsoluteY, ldx},

	{0xA0, "LDY", 2, 2, Immediate, ldy},
	{0xA4, "LDY", 2, 3, ZeroPage, ldy},
	{0xB4, "LDY", 2, 4, ZeroPageX, ldy},
	{0xAC, "LDY", 3, 4, Absolute, ldy},
	{0xBC, "LDY", 3, 4 /*+1 if page crossed*/, AbsoluteX, ldy},

	{0x4A, "LSR", 1, 2, Accumulator, lsr},
	{0x46, "LSR", 2, 5, ZeroPage, lsr},
	{0x56, "LSR", 2, 6, ZeroPageX, lsr},
	{0x4E, "LSR", 3, 6, Absolute, lsr},
	{0x5E, "LSR", 3, 7, AbsoluteX, lsr},

	{0x09, "ORA", 2, 2, Immediate, ora},
	{0x05, "ORA", 2, 3, ZeroPage, ora},
	{0x15, "ORA", 2, 4, ZeroPageX, ora},
	{0x0D, "ORA", 3, 4, Absolute, ora},
	{0x1D, "ORA", 3, 4 /*+1 if page crossed*/, AbsoluteX, ora},
	{0x19, "ORA", 3, 4 /*+1 if page crossed*/, AbsoluteY, ora},
	{0x01, "ORA", 2, 6, IndirectX, ora},
	{0x11, "ORA", 2, 5 /*f page crossed*/, IndirectY, ora},

	{0x2A, "ROL", 1, 2, Accumulator, rol},
	{0x26, "ROL", 2, 5, ZeroPage, rol},
	{0x36, "ROL", 2, 6, ZeroPageX, rol},
	{0x2E, "ROL", 3, 6, Absolute, rol},
	{0x3E, "ROL", 3, 7, AbsoluteX, rol},

	{0x6A, "ROR", 1, 2, Accumulator, ror},
	{0x66, "ROR", 2, 5, ZeroPage, ror},
	{0x76, "ROR", 2, 6, ZeroPageX, ror},
	{0x6E, "ROR", 3, 6, Absolute, ror},
	{0x7E, "ROR", 3, 7, AbsoluteX, ror},

	{0xE9, "SBC", 2, 2, Immediate, sbc},
	{0xE5, "SBC", 2, 3, ZeroPage, sbc},
	{0xF5, "SBC", 2, 4, ZeroPageX, sbc},
	{0xED, "SBC", 3, 4, Absolute, sbc},
	{0xFD, "SBC", 3, 4 /*+1 if page crossed*/, AbsoluteX, sbc},
	{0xF9, "SBC", 3, 4 /*+1 if page crossed*/, AbsoluteY, sbc},
	{0xE1, "SBC", 2, 6, IndirectX, sbc},
	{0xF1, "SBC", 2, 5 /*+1 if page crossed*/, IndirectY, sbc},

	{0x85, "STA", 2, 3, ZeroPage, sta},
	{0x95, "STA", 2, 4, ZeroPageX, sta},
	{0x8D, "STA", 3, 4, Absolute, sta},
	{0x9D, "STA", 3, 5, AbsoluteX, sta},
	{0x99, "STA", 3, 5, AbsoluteY, sta},
	{0x81, "STA", 2, 6, IndirectX, sta},
	{0x91, "STA", 2, 6, IndirectY, sta},

	{0x86, "STX", 2, 3, ZeroPage, stx},
	{0x96, "STX", 2, 4, ZeroPageY, stx},
	{0x8E, "STX", 3, 4, Absolute, stx},

	{0x84, "STY", 2, 3, ZeroPage, sty},
	{0x94, "STY", 2, 4, ZeroPageX, sty},
	{0x8C, "STY", 3, 4, Absolute, sty},
}

// OpCodeMap converts OpCodes into a map with the code as the key.
func OpCodeMap() map[byte]OpCode {
	codes := make(map[byte]OpCode)
	for _, opcode := range OpCodes {
		if _, ok := codes[opcode.Code]; ok {
			panic(fmt.Sprintf("duplicate opcode: $%02X", opcode.Code))
		}
		codes[opcode.Code] = opcode
	}
	return codes
}
