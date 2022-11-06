package cpu

// MemRead reads uint8 from memory.
func (c *CPU) MemRead(addr uint16) uint8 {
	return c.memory[addr]
}

// MemWrite writes uint8 to memory.
func (c *CPU) MemWrite(addr uint16, data uint8) {
	c.memory[addr] = data
}

// MemRead16 reads uint16 from memory.
func (c *CPU) MemRead16(pos uint16) uint16 {
	lo := uint16(c.MemRead(pos))
	hi := uint16(c.MemRead(pos + 1))
	return hi<<8 | lo
}

// MemWrite16 writes uint16 to memory.
func (c *CPU) MemWrite16(pos uint16, data uint16) {
	hi := uint8(data >> 8)
	lo := uint8(data & 0xFF)
	c.MemWrite(pos, lo)
	c.MemWrite(pos+1, hi)
}

func (c *CPU) stackPush(data uint8) {
	c.MemWrite(StackAddr+uint16(c.stackPointer), data)
	c.stackPointer -= 1
}

func (c *CPU) stackPush16(data uint16) {
	hi := uint8(data >> 8)
	lo := uint8(data & 0xFF)
	c.stackPush(hi)
	c.stackPush(lo)
}

func (c *CPU) stackPop() uint8 {
	c.stackPointer += 1
	return c.MemRead(StackAddr + uint16(c.stackPointer))
}

func (c *CPU) stackPop16() uint16 {
	lo := uint16(c.stackPop())
	hi := uint16(c.stackPop())
	return hi<<8 | lo
}
