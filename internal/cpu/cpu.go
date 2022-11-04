package cpu

import "log"

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

func (c *CPU) loadAndRun(program []uint8) {
	c.load(program)
	c.reset()
	c.run()
}

func (c *CPU) lda(v uint8) {
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

func (c *CPU) run() {
	for {
		opcode := c.memRead(c.PC)
		c.PC += 1

		switch opcode {
		case 0xA9:
			param := c.memRead(c.PC)
			c.PC += 1

			c.lda(param)
		case 0xAA:
			c.tax()
		case 0xE8:
			c.inx()
		case 0x00:
			return
		default:
			log.Panicf("unsupported opcode: 0x%x\n", opcode)
		}
	}
}
