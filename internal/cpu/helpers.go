package cpu

import "github.com/gabe565/gones/internal/bitflags"

// updateZeroAndNegFlags updates zero and negative flags
func (c *CPU) updateZeroAndNegFlags(result byte) {
	c.status.Set(Zero, result == 0)
	c.status.Set(Negative, bitflags.Flags(result).Has(Negative))
}

func (c *CPU) branch(condition bool) {
	if condition {
		jump := int8(c.MemRead(c.programCounter))
		jumpAddr := c.programCounter + 1 + uint16(jump)

		c.programCounter = jumpAddr
	}
}

func (c *CPU) compare(mode AddressingMode, rhs byte) {
	addr := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.status.Set(Carry, data <= rhs)
	c.updateZeroAndNegFlags(rhs - data)
}

func (c *CPU) setAccumulator(v byte) {
	c.accumulator = v
	c.updateZeroAndNegFlags(c.accumulator)
}

func (c *CPU) addAccumulator(data byte) {
	sum := uint16(c.accumulator) + uint16(data)
	if c.status.Has(Carry) {
		sum += 1
	}

	carry := sum > 0xFF
	c.status.Set(Carry, carry)

	result := byte(sum)
	c.status.Set(Overflow, (data^result)&(result^c.accumulator)&0x80 != 0)

	c.setAccumulator(result)
}
