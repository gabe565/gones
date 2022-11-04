package cpu

import (
	"errors"
	"fmt"
	"github.com/gabe565/gones/internal/bits"
)

func New() CPU {
	return CPU{}
}

// CPU implements the NES CPU.
//
// See [6502 Guide].
//
// [6502 Guide]: https://www.nesdev.org/obelisk-6502-guide/
type CPU struct {
	// PC Program Counter
	PC uint16

	// Status Processor Status
	Status bits.Bits

	// Accumulator Register A
	Accumulator uint8

	// RegisterX Register X
	RegisterX uint8

	// RegisterY Register Y
	RegisterY uint8

	// Memory Main memory
	Memory [0xFFFF]uint8
}

const (
	// PrgRomStart is the memory address that PRG begins.
	PrgRomStart = 0x8000
	// Reset is the memory address for the Reset Interrupt Vector.
	Reset = 0xFFFC
)

// memRead reads uint8 from memory.
func (c *CPU) memRead(addr uint16) uint8 {
	return c.Memory[addr]
}

// memWrite writes uint8 to memory.
func (c *CPU) memWrite(addr uint16, data uint8) {
	c.Memory[addr] = data
}

// memRead16 reads uint16 from memory.
func (c *CPU) memRead16(pos uint16) uint16 {
	lo := uint16(c.memRead(pos))
	hi := uint16(c.memRead(pos + 1))
	return hi<<8 | lo
}

// memWrite16 writes uint16 to memory.
func (c *CPU) memWrite16(pos uint16, data uint16) {
	hi := uint8(data >> 8)
	lo := uint8(data & 0xFF)
	c.memWrite(pos, lo)
	c.memWrite(pos+1, hi)
}

func (c *CPU) setRegisterA(v uint8) {
	c.Accumulator = v
	c.updateZeroAndNegFlags(c.Accumulator)
}

// reset resets the CPU and sets PC to the value of the [Reset] Vector.
func (c *CPU) reset() {
	c.Accumulator = 0
	c.RegisterX = 0
	c.Status = 0

	c.PC = c.memRead16(Reset)
}

// load loads a program into PRG memory
func (c *CPU) load(program []uint8) {
	for k, v := range program {
		c.Memory[PrgRomStart+k] = v
	}
	c.memWrite16(Reset, PrgRomStart)
}

// loadAndRun is a convenience function that loads a program, resets, then runs.
func (c *CPU) loadAndRun(program []uint8) error {
	c.load(program)
	c.reset()
	return c.run()
}

// updateZeroAndNegFlags updates zero and negative flags
func (c *CPU) updateZeroAndNegFlags(result uint8) {
	if result == 0 {
		c.Status = bits.Set(c.Status, Zero)
	} else {
		c.Status = bits.Clear(c.Status, Zero)
	}

	if bits.Has(bits.Bits(result), Negative) {
		c.Status = bits.Set(c.Status, Negative)
	} else {
		c.Status = bits.Clear(c.Status, Negative)
	}
}

// ErrUnsupportedOpcode indicates an unsupported opcode was evaluated.
var ErrUnsupportedOpcode = errors.New("unsupported opcode")

// run is the main run entrypoint.
func (c *CPU) run() error {
	opcodes := OpCodeMap()

	for {
		code := c.memRead(c.PC)
		c.PC += 1
		prevPC := c.PC

		opcode, ok := opcodes[code]
		if !ok {
			return fmt.Errorf("%w: $%x", ErrUnsupportedOpcode, code)
		}

		switch code {
		case 0xA9, 0xA5, 0xB5, 0xAD, 0xBD, 0xB9, 0xA1, 0xB1:
			c.lda(opcode.Mode)
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
		case 0xE8:
			c.inx()
		case 0x00:
			return nil
		default:
			return fmt.Errorf("%w: $%x", ErrUnsupportedOpcode, opcode)
		}

		if prevPC == c.PC {
			c.PC += uint16(opcode.Len - 1)
		}
	}
}
