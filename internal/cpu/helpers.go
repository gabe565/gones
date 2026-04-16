package cpu

// updateZeroAndNegFlags updates zero and negative flags.
func (c *CPU) updateZeroAndNegFlags(result byte) {
	c.Status.Zero = result == 0
	c.Status.Negative = result&Negative != 0
}

func (c *CPU) branch(condition bool) {
	offset := c.ReadMem(c.ProgramCounter) // always read operand (1 cycle)
	if condition {
		c.tick() // internal: compute branch target

		jumpAddr := c.ProgramCounter + 1 + uint16(int8(offset))

		if crossedPage(c.ProgramCounter+1, jumpAddr) {
			c.tick() // internal: fix high byte
		}

		c.ProgramCounter = jumpAddr
	}
}

func (c *CPU) compare(mode AddressingMode, rhs byte) {
	addr, _ := c.getOperandAddress(mode)
	data := c.ReadMem(addr)
	c.Status.Carry = data <= rhs
	c.updateZeroAndNegFlags(rhs - data)
}

func (c *CPU) setAccumulator(v byte) {
	c.Accumulator = v
	c.updateZeroAndNegFlags(c.Accumulator)
}

func (c *CPU) addAccumulator(data byte) {
	sum := uint16(c.Accumulator) + uint16(data)
	if c.Status.Carry {
		sum++
	}

	carry := sum > 0xFF
	c.Status.Carry = carry

	result := byte(sum)
	c.Status.Overflow = (data^result)&(result^c.Accumulator)&0x80 != 0

	c.setAccumulator(result)
}
