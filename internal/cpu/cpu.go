package cpu

import (
	"errors"
	"fmt"
	"log"
)

func New() CPU {
	return CPU{}
}

type CPU struct {
	// PC Program Counter
	PC uint16

	// Acc Accumulator
	Acc uint8

	// RegisterA Register A
	RegisterA uint8

	// RegisterX Register X
	RegisterX uint8

	// RegisterY Register Y
	RegisterY uint8

	// Memory Main memory
	Memory [0xFFFF]uint8
}

const (
	PrgRomStart = 0x8000
	Reset       = 0xFFFC
)

func (c *CPU) memRead(addr uint16) uint8 {
	return c.Memory[addr]
}

func (c *CPU) memWrite(addr uint16, data uint8) {
	c.Memory[addr] = data
}

func (c *CPU) memRead16(pos uint16) uint16 {
	lo := uint16(c.memRead(pos))
	hi := uint16(c.memRead(pos + 1))
	return hi<<8 | lo
}

func (c *CPU) memWrite16(pos uint16, data uint16) {
	hi := uint8(data >> 8)
	lo := uint8(data & 0xff)
	c.memWrite(pos, lo)
	c.memWrite(pos+1, hi)
}

func (c *CPU) reset() {
	c.RegisterA = 0
	c.RegisterX = 0
	c.Acc = 0

	c.PC = c.memRead16(Reset)
}

func (c *CPU) load(program []uint8) {
	for k, v := range program {
		c.Memory[PrgRomStart+k] = v
	}
	c.memWrite16(Reset, PrgRomStart)
}

func (c *CPU) loadAndRun(program []uint8) error {
	c.load(program)
	c.reset()
	return c.run()
}

func (c *CPU) lda(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	v := c.memRead(addr)

	c.RegisterA = v
	c.updateZeroAndNegFlags(c.RegisterA)
}

func (c *CPU) tax() {
	c.RegisterX = c.RegisterA
	c.updateZeroAndNegFlags(c.RegisterX)
}

func (c *CPU) inx() {
	c.RegisterX += 1
	c.updateZeroAndNegFlags(c.RegisterX)
}

func (c *CPU) updateZeroAndNegFlags(result uint8) {
	if result == 0 {
		c.Acc |= 0b0000_0010
	} else {
		c.Acc &= 0b1111_1101
	}

	if result&0b1000_0000 != 0 {
		c.Acc |= 0b1000_0000
	} else {
		c.Acc &= 0b0111_1111
	}
}

func (c *CPU) getOperandAddress(mode AddressingMode) uint16 {
	switch mode {
	case Immediate:
		return c.PC
	case ZeroPage:
		return uint16(c.memRead(c.PC))
	case Absolute:
		return c.memRead16(c.PC)
	case ZeroPageX:
		pos := c.memRead(c.PC)
		return uint16(pos + c.RegisterX)
	case ZeroPageY:
		pos := c.memRead(c.PC)
		return uint16(pos + c.RegisterY)
	case AbsoluteX:
		pos := c.memRead(c.PC)
		return uint16(pos) + uint16(c.RegisterX)
	case AbsoluteY:
		pos := c.memRead(c.PC)
		return uint16(pos) + uint16(c.RegisterY)
	case IndirectX:
		base := c.memRead(c.PC)

		ptr := base + c.RegisterX
		lo := c.memRead(uint16(ptr))
		hi := c.memRead(uint16(ptr + 1))
		return uint16(hi)<<8 | uint16(lo)
	case IndirectY:
		base := c.memRead(c.PC)

		lo := c.memRead(uint16(base))
		hi := c.memRead(uint16(uint8(base) + 1))
		derefBase := uint16(hi)<<8 | uint16(lo)
		return derefBase + uint16(c.RegisterY)
	default:
		log.Panicln("unsupported mode: ", mode)
		return 0
	}
}

var ErrUnsupportedOpcode = errors.New("unsupported opcode")

func (c *CPU) run() error {
	opcodes := OpCodeMap()

	for {
		code := c.memRead(c.PC)
		c.PC += 1

		opcode, ok := opcodes[code]
		if !ok {
			return fmt.Errorf("%w: $%x", ErrUnsupportedOpcode, code)
		}

		switch code {
		case 0xA9, 0xa5, 0xb5, 0xad, 0xbd, 0xb9, 0xa1, 0xb1:
			c.lda(opcode.Mode)
			c.PC += 1
		case 0x85, 0x95, 0x8d, 0x9d, 0x99, 0x81, 0x91:
			c.lda(opcode.Mode)
			c.PC += 1
		case 0xAA:
			c.tax()
		case 0xE8:
			c.inx()
		case 0x00:
			return nil
		default:
			return fmt.Errorf("%w: $%x", ErrUnsupportedOpcode, opcode)
		}
	}
}
