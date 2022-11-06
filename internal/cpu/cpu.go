package cpu

import (
	"errors"
	"fmt"
	"github.com/gabe565/gones/internal/bitflags"
)

func New() CPU {
	return CPU{
		Status:     InterruptDisable | Break2,
		SP:         StackReset,
		PrgRomAddr: PrgRomAddr,
	}
}

// CPU implements the NES CPU.
//
// See [6502 Guide].
//
// [6502 Guide]: https://www.nesdev.org/obelisk-6502-guide/
type CPU struct {
	// PC Program Counter
	PC uint16

	// SP Stack Pointer
	SP uint8

	// Status Processor Status
	Status bitflags.Flags

	// Accumulator Register A
	Accumulator uint8

	// RegisterX Register X
	RegisterX uint8

	// RegisterY Register Y
	RegisterY uint8

	// Memory Main memory
	Memory [0xFFFF]uint8

	// Callback optional callback to Run before every tick
	Callback func(c *CPU)

	// PrgRomAddr is the PRG ROM start address
	PrgRomAddr uint16

	// Debug enables opcode logging
	Debug bool
}

const (
	// PrgRomAddr is the memory address that PRG begins.
	PrgRomAddr = 0x8000

	// ResetAddr is the memory address for the Reset Interrupt Vector.
	ResetAddr = 0xFFFC

	// StackAddr is the memory address of the stack
	StackAddr = 0x100

	// StackReset is the start value for the stack pointer
	StackReset = 0xFD
)

func (c *CPU) setAccumulator(v uint8) {
	c.Accumulator = v
	c.updateZeroAndNegFlags(c.Accumulator)
}

func (c *CPU) addAccumulator(data uint8) {
	sum := uint16(c.Accumulator) + uint16(data)
	if c.Status.Has(Carry) {
		sum += 1
	}

	carry := sum > 0xFF
	c.Status.Set(Carry, carry)

	result := uint8(sum)
	c.Status.Set(Overflow, (data^result)&(result^c.Accumulator)&0x80 != 0)

	c.setAccumulator(result)
}

// Reset resets the CPU and sets PC to the value of the [Reset] Vector.
func (c *CPU) Reset() {
	c.Accumulator = 0
	c.RegisterX = 0
	c.Status = 0
	c.SP = StackReset

	c.PC = c.memRead16(ResetAddr)
}

// Load loads a program into PRG memory
func (c *CPU) Load(program []uint8) {
	for k, v := range program {
		c.Memory[c.PrgRomAddr+uint16(k)] = v
	}
	c.memWrite16(ResetAddr, c.PrgRomAddr)
}

// loadAndRun is a convenience function that loads a program, resets, then runs.
func (c *CPU) loadAndRun(program []uint8) error {
	c.Load(program)
	c.Reset()
	return c.Run()
}

// updateZeroAndNegFlags updates zero and negative flags
func (c *CPU) updateZeroAndNegFlags(result uint8) {
	c.Status.Set(Zero, result == 0)
	c.Status.Set(Negative, bitflags.Flags(result).Has(Negative))
}

func (c *CPU) branch(condition bool) {
	if condition {
		jump := int8(c.memRead(c.PC))
		jumpAddr := c.PC + 1 + uint16(jump)

		c.PC = jumpAddr
	}
}

func (c *CPU) compare(mode AddressingMode, rhs uint8) {
	addr := c.getOperandAddress(mode)
	data := c.memRead(addr)
	c.Status.Set(Carry, data <= rhs)
	c.updateZeroAndNegFlags(rhs - data)
}

// ErrUnsupportedOpcode indicates an unsupported opcode was evaluated.
var ErrUnsupportedOpcode = errors.New("unsupported opcode")

// Run is the main Run entrypoint.
func (c *CPU) Run() error {
	opcodes := OpCodeMap()

	for {
		if c.Callback != nil {
			c.Callback(c)
		}

		code := c.memRead(c.PC)
		c.PC += 1
		prevPC := c.PC

		opcode, ok := opcodes[code]
		if !ok {
			return fmt.Errorf("%w: $%x", ErrUnsupportedOpcode, code)
		}

		if c.Debug {
			fmt.Println(opcode)
		}

		switch code {
		case 0x69, 0x65, 0x75, 0x6D, 0x7D, 0x79, 0x61, 0x71:
			c.adc(opcode.Mode)
		case 0x29, 0x25, 0x35, 0x2D, 0x3D, 0x39, 0x21, 0x31:
			c.and(opcode.Mode)
		case 0x0A, 0x06, 0x16, 0x0E, 0x1E:
			c.asl(opcode.Mode)
		case 0x90:
			c.bcc()
		case 0xB0:
			c.bcs()
		case 0xF0:
			c.beq()
		case 0x24, 0x2C:
			c.bit(opcode.Mode)
		case 0x30:
			c.bmi()
		case 0xD0:
			c.bne()
		case 0x10:
			c.bpl()
		case 0x00:
			return nil
		case 0x50:
			c.bvc()
		case 0x70:
			c.bvs()
		case 0x18:
			c.clc()
		case 0xD8:
			c.cld()
		case 0x58:
			c.cli()
		case 0xB8:
			c.clv()
		case 0xC9, 0xC5, 0xD5, 0xCD, 0xDD, 0xD9, 0xC1, 0xD1:
			c.cmp(opcode.Mode)
		case 0xE0, 0xE4, 0xEC:
			c.cpx(opcode.Mode)
		case 0xC0, 0xC4, 0xCC:
			c.cpy(opcode.Mode)
		case 0xC6, 0xD6, 0xCE, 0xDE:
			c.dec(opcode.Mode)
		case 0xCA:
			c.dex()
		case 0x88:
			c.dey()
		case 0x49, 0x45, 0x55, 0x4D, 0x5D, 0x59, 0x41, 0x51:
			c.eor(opcode.Mode)
		case 0xE6, 0xF6, 0xEE, 0xFE:
			c.inc(opcode.Mode)
		case 0xE8:
			c.inx()
		case 0xC8:
			c.iny()
		case 0x4C, 0x6C:
			c.jmp(opcode.Mode)
		case 0x20:
			c.jsr()
		case 0xA9, 0xA5, 0xB5, 0xAD, 0xBD, 0xB9, 0xA1, 0xB1:
			c.lda(opcode.Mode)
		case 0xA2, 0xA6, 0xB6, 0xAE, 0xBE:
			c.ldx(opcode.Mode)
		case 0xA0, 0xA4, 0xB4, 0xAC, 0xBC:
			c.ldy(opcode.Mode)
		case 0x4A, 0x46, 0x56, 0x4E, 0x5E:
			c.lsr(opcode.Mode)
		case 0xEA:
			// NOP
		case 0x09, 0x05, 0x15, 0x0D, 0x1D, 0x19, 0x01, 0x11:
			c.ora(opcode.Mode)
		case 0x48:
			c.pha()
		case 0x08:
			c.php()
		case 0x68:
			c.pla()
		case 0x28:
			c.plp()
		case 0x2A, 0x26, 0x36, 0x2E, 0x3E:
			c.rol(opcode.Mode)
		case 0x6A, 0x66, 0x76, 0x6E, 0x7E:
			c.ror(opcode.Mode)
		case 0x40:
			c.rti()
		case 0x60:
			c.rts()
		case 0xE9, 0xE5, 0xF5, 0xED, 0xFD, 0xF9, 0xF1:
			c.sbc(opcode.Mode)
		case 0x38:
			c.sec()
		case 0xF8:
			c.sed()
		case 0x78:
			c.sei()
		case 0x85, 0x95, 0x8D, 0x9D, 0x99, 0x81, 0x91:
			c.sta(opcode.Mode)
		case 0x86, 0x96, 0x8E:
			c.stx(opcode.Mode)
		case 0x84, 0x94, 0x8C:
			c.sty(opcode.Mode)
		case 0xAA:
			c.tax()
		case 0xA8:
			c.tay()
		case 0xBA:
			c.tsx()
		case 0x8A:
			c.txa()
		case 0x9A:
			c.txs()
		case 0x98:
			c.tya()
		default:
			return fmt.Errorf("%w: $%x", ErrUnsupportedOpcode, code)
		}

		if prevPC == c.PC {
			c.PC += uint16(opcode.Len - 1)
		}
	}
}
