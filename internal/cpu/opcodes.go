package cpu

import "fmt"

// OpCode defines an opcode and its parameters.
//
// See [6502 Instruction Reference].
//
// [6502 Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html
type OpCode struct {
	Mnemonic string
	Code     byte
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
	{"ADC", 0x61, 2, 2, IndirectX, adc},
	{"ADC", 0x65, 2, 2, ZeroPage, adc},
	{"ADC", 0x69, 2, 2, Immediate, adc},
	{"ADC", 0x6D, 3, 3, Absolute, adc},
	{"ADC", 0x71, 2, 2 /*+1 if page crossed*/, IndirectY, adc},
	{"ADC", 0x75, 2, 2, ZeroPageX, adc},
	{"ADC", 0x79, 3, 3 /*+1 if page crossed*/, AbsoluteY, adc},
	{"ADC", 0x7D, 3, 3 /*+1 if page crossed*/, AbsoluteX, adc},
	{"AND", 0x21, 2, 6, IndirectX, and},
	{"AND", 0x25, 2, 3, ZeroPage, and},
	{"AND", 0x29, 2, 2, Immediate, and},
	{"AND", 0x2D, 3, 4, Absolute, and},
	{"AND", 0x31, 2, 5 /*+1 if page crossed*/, IndirectY, and},
	{"AND", 0x35, 2, 4, ZeroPageX, and},
	{"AND", 0x39, 3, 4 /*+1 if page crossed*/, AbsoluteY, and},
	{"AND", 0x3D, 3, 4 /*+1 if page crossed*/, AbsoluteX, and},
	{"ASL", 0x06, 2, 5, ZeroPage, asl},
	{"ASL", 0x0A, 1, 2, Accumulator, asl},
	{"ASL", 0x0E, 3, 6, Absolute, asl},
	{"ASL", 0x16, 2, 6, ZeroPageX, asl},
	{"ASL", 0x1E, 3, 7 /*+1 if page crossed*/, AbsoluteX, asl},
	{"BCC", 0x90, 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, Implicit, bcc},
	{"BCS", 0xB0, 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, Implicit, bcs},
	{"BEQ", 0xF0, 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, Implicit, beq},
	{"BIT", 0x24, 2, 3, ZeroPage, bit},
	{"BIT", 0x2C, 3, 4, Absolute, bit},
	{"BMI", 0x30, 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, Implicit, bmi},
	{"BNE", 0xD0, 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, Implicit, bne},
	{"BPL", 0x10, 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, Implicit, bpl},
	{"BRK", 0x00, 1, 7, Implicit, brk},
	{"BVC", 0x50, 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, Implicit, bvc},
	{"BVS", 0x70, 2, 2 /*+1 if branch succeeds, +2 if to a new page*/, Implicit, bvs},
	{"CLC", 0x18, 1, 2, Implicit, clc},
	{"CLD", 0xD8, 1, 2, Implicit, cld},
	{"CLI", 0x58, 1, 2, Implicit, cli},
	{"CLV", 0xB8, 1, 2, Implicit, clv},
	{"CMP", 0xC1, 2, 6, IndirectX, cmp},
	{"CMP", 0xC5, 2, 3, ZeroPage, cmp},
	{"CMP", 0xC9, 2, 2, Immediate, cmp},
	{"CMP", 0xCD, 3, 4, Absolute, cmp},
	{"CMP", 0xD1, 2, 5 /*+1 if page crossed*/, IndirectY, cmp},
	{"CMP", 0xD5, 2, 4, ZeroPageX, cmp},
	{"CMP", 0xD9, 3, 4 /*+1 if page crossed*/, AbsoluteY, cmp},
	{"CMP", 0xDD, 3, 4 /*+1 if page crossed*/, AbsoluteX, cmp},
	{"CPX", 0xE0, 2, 2, Immediate, cpx},
	{"CPX", 0xE4, 2, 3, ZeroPage, cpx},
	{"CPX", 0xEC, 3, 4, Absolute, cpx},
	{"CPY", 0xC0, 2, 2, Immediate, cpy},
	{"CPY", 0xC4, 2, 3, ZeroPage, cpy},
	{"CPY", 0xCC, 3, 4, Absolute, cpy},
	{"DEC", 0xC6, 2, 5, ZeroPage, dec},
	{"DEC", 0xCE, 3, 6, Absolute, dec},
	{"DEC", 0xD6, 2, 6, ZeroPageX, dec},
	{"DEC", 0xDE, 3, 7, AbsoluteX, dec},
	{"DEX", 0xCA, 1, 2, Implicit, dex},
	{"DEY", 0x88, 1, 2, Implicit, dey},
	{"EOR", 0x41, 2, 6, IndirectX, eor},
	{"EOR", 0x45, 2, 3, ZeroPage, eor},
	{"EOR", 0x49, 2, 2, Immediate, eor},
	{"EOR", 0x4D, 3, 4, Absolute, eor},
	{"EOR", 0x51, 2, 5 /*+1 if page crossed*/, IndirectY, eor},
	{"EOR", 0x55, 2, 4, ZeroPageX, eor},
	{"EOR", 0x59, 3, 4 /*+1 if page crossed*/, AbsoluteY, eor},
	{"EOR", 0x5D, 3, 4 /*+1 if page crossed*/, AbsoluteX, eor},
	{"INC", 0xE6, 2, 5, ZeroPage, inc},
	{"INC", 0xEE, 3, 6, Absolute, inc},
	{"INC", 0xF6, 2, 6, ZeroPageX, inc},
	{"INC", 0xFE, 3, 7, AbsoluteX, inc},
	{"INX", 0xE8, 1, 2, Implicit, inx},
	{"INY", 0xC8, 1, 2, Implicit, iny},
	{"JMP", 0x4C, 3, 3, Absolute, jmp}, // Acts as immediate
	{"JMP", 0x6C, 3, 5, Indirect, jmp}, // Indirect with 6502 bug
	{"JSR", 0x20, 3, 6, Implicit, jsr},
	{"LDA", 0xA1, 2, 6, IndirectX, lda},
	{"LDA", 0xA5, 2, 3, ZeroPage, lda},
	{"LDA", 0xA9, 2, 2, Immediate, lda},
	{"LDA", 0xAD, 3, 4, Absolute, lda},
	{"LDA", 0xB1, 2, 5 /*+1 if page crossed*/, IndirectY, lda},
	{"LDA", 0xB5, 2, 4, ZeroPageX, lda},
	{"LDA", 0xB9, 3, 4 /*+1 if page crossed*/, AbsoluteY, lda},
	{"LDA", 0xBD, 3, 4 /*+1 if page crossed*/, AbsoluteX, lda},
	{"LDX", 0xA2, 2, 2, Immediate, ldx},
	{"LDX", 0xA6, 2, 3, ZeroPage, ldx},
	{"LDX", 0xAE, 3, 4, Absolute, ldx},
	{"LDX", 0xB6, 2, 4, ZeroPageY, ldx},
	{"LDX", 0xBE, 3, 4 /*+1 if page crossed*/, AbsoluteY, ldx},
	{"LDY", 0xA0, 2, 2, Immediate, ldy},
	{"LDY", 0xA4, 2, 3, ZeroPage, ldy},
	{"LDY", 0xAC, 3, 4, Absolute, ldy},
	{"LDY", 0xB4, 2, 4, ZeroPageX, ldy},
	{"LDY", 0xBC, 3, 4 /*+1 if page crossed*/, AbsoluteX, ldy},
	{"LSR", 0x46, 2, 5, ZeroPage, lsr},
	{"LSR", 0x4A, 1, 2, Accumulator, lsr},
	{"LSR", 0x4E, 3, 6, Absolute, lsr},
	{"LSR", 0x56, 2, 6, ZeroPageX, lsr},
	{"LSR", 0x5E, 3, 7, AbsoluteX, lsr},
	{"NOP", 0xEA, 1, 2, Implicit, nop},
	{"ORA", 0x01, 2, 6, IndirectX, ora},
	{"ORA", 0x05, 2, 3, ZeroPage, ora},
	{"ORA", 0x09, 2, 2, Immediate, ora},
	{"ORA", 0x0D, 3, 4, Absolute, ora},
	{"ORA", 0x11, 2, 5 /*f page crossed*/, IndirectY, ora},
	{"ORA", 0x15, 2, 4, ZeroPageX, ora},
	{"ORA", 0x19, 3, 4 /*+1 if page crossed*/, AbsoluteY, ora},
	{"ORA", 0x1D, 3, 4 /*+1 if page crossed*/, AbsoluteX, ora},
	{"PHA", 0x48, 1, 3, Implicit, pha},
	{"PHP", 0x08, 1, 3, Implicit, php},
	{"PLA", 0x68, 1, 4, Implicit, pla},
	{"PLP", 0x28, 1, 4, Implicit, plp},
	{"ROL", 0x26, 2, 5, ZeroPage, rol},
	{"ROL", 0x2A, 1, 2, Accumulator, rol},
	{"ROL", 0x2E, 3, 6, Absolute, rol},
	{"ROL", 0x36, 2, 6, ZeroPageX, rol},
	{"ROL", 0x3E, 3, 7, AbsoluteX, rol},
	{"ROR", 0x66, 2, 5, ZeroPage, ror},
	{"ROR", 0x6A, 1, 2, Accumulator, ror},
	{"ROR", 0x6E, 3, 6, Absolute, ror},
	{"ROR", 0x76, 2, 6, ZeroPageX, ror},
	{"ROR", 0x7E, 3, 7, AbsoluteX, ror},
	{"RTI", 0x40, 1, 6, Implicit, rti},
	{"RTS", 0x60, 1, 6, Implicit, rts},
	{"SBC", 0xE1, 2, 6, IndirectX, sbc},
	{"SBC", 0xE5, 2, 3, ZeroPage, sbc},
	{"SBC", 0xE9, 2, 2, Immediate, sbc},
	{"SBC", 0xED, 3, 4, Absolute, sbc},
	{"SBC", 0xF1, 2, 5 /*+1 if page crossed*/, IndirectY, sbc},
	{"SBC", 0xF5, 2, 4, ZeroPageX, sbc},
	{"SBC", 0xF9, 3, 4 /*+1 if page crossed*/, AbsoluteY, sbc},
	{"SBC", 0xFD, 3, 4 /*+1 if page crossed*/, AbsoluteX, sbc},
	{"SEC", 0x38, 1, 2, Implicit, sec},
	{"SED", 0xF8, 1, 2, Implicit, sed},
	{"SEI", 0x78, 1, 2, Implicit, sei},
	{"STA", 0x81, 2, 6, IndirectX, sta},
	{"STA", 0x85, 2, 3, ZeroPage, sta},
	{"STA", 0x8D, 3, 4, Absolute, sta},
	{"STA", 0x91, 2, 6, IndirectY, sta},
	{"STA", 0x95, 2, 4, ZeroPageX, sta},
	{"STA", 0x99, 3, 5, AbsoluteY, sta},
	{"STA", 0x9D, 3, 5, AbsoluteX, sta},
	{"STX", 0x86, 2, 3, ZeroPage, stx},
	{"STX", 0x8E, 3, 4, Absolute, stx},
	{"STX", 0x96, 2, 4, ZeroPageY, stx},
	{"STY", 0x84, 2, 3, ZeroPage, sty},
	{"STY", 0x8C, 3, 4, Absolute, sty},
	{"STY", 0x94, 2, 4, ZeroPageX, sty},
	{"TAX", 0xAA, 1, 2, Implicit, tax},
	{"TAY", 0xA8, 1, 2, Implicit, tay},
	{"TSX", 0xBA, 1, 2, Implicit, tsx},
	{"TXA", 0x8A, 1, 2, Implicit, txa},
	{"TXS", 0x9A, 1, 2, Implicit, txs},
	{"TYA", 0x98, 1, 2, Implicit, tya},
}

// OpCodeMap is a map of OpCodes with the code as the key.
var OpCodeMap map[byte]OpCode

func init() {
	codes := make(map[byte]OpCode)
	for _, opcode := range OpCodes {
		if _, ok := codes[opcode.Code]; ok {
			panic(fmt.Sprintf("duplicate opcode: $%02X", opcode.Code))
		}
		codes[opcode.Code] = opcode
	}
	OpCodeMap = codes
}
