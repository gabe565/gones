package cpu

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

func (c *CPU) stackPush(data uint8) {
	c.memWrite(StackAddr+uint16(c.SP), data)
	c.SP -= 1
}

func (c *CPU) stackPush16(data uint16) {
	hi := uint8(data >> 8)
	lo := uint8(data & 0xFF)
	c.stackPush(hi)
	c.stackPush(lo)
}

func (c *CPU) stackPop() uint8 {
	c.SP += 1
	return c.memRead(StackAddr + uint16(c.SP))
}

func (c *CPU) stackPop16() uint16 {
	lo := uint16(c.stackPop())
	hi := uint16(c.stackPop())
	return hi<<8 | lo
}
