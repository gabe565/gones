package cpu

// updateZeroAndNegFlags updates zero and negative flags
func (c *CPU) updateZeroAndNegFlags(result byte) {
	c.Status.Zero = result == 0
	c.Status.Negative = result&Negative != 0
}

func (c *CPU) branch(condition bool) {
	if condition {
		c.Cycles += 1

		jump := int8(c.ReadMem(c.ProgramCounter))
		jumpAddr := c.ProgramCounter + 1 + uint16(jump)

		if crossedPage(c.ProgramCounter+1, jumpAddr&0xFF00) {
			c.Cycles += 1
		}

		c.ProgramCounter = jumpAddr
	}
}

func (c *CPU) compare(mode AddressingMode, rhs byte) {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Cycles += 1
		}()
	}
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
		sum += 1
	}

	carry := sum > 0xFF
	c.Status.Carry = carry

	result := byte(sum)
	c.Status.Overflow = (data^result)&(result^c.Accumulator)&0x80 != 0

	c.setAccumulator(result)
}
