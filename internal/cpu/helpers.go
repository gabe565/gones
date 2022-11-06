package cpu

import "github.com/gabe565/gones/internal/bitflags"

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
