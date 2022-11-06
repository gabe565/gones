package cpu

// MemRead reads byte from memory.
func (c *CPU) MemRead(addr uint16) byte {
	return c.memory[addr]
}

// MemWrite writes byte to memory.
func (c *CPU) MemWrite(addr uint16, data byte) {
	c.memory[addr] = data
}

// MemRead16 reads two bytes from memory.
func (c *CPU) MemRead16(pos uint16) uint16 {
	lo := uint16(c.MemRead(pos))
	hi := uint16(c.MemRead(pos + 1))
	return hi<<8 | lo
}

// MemWrite16 writes two bytes to memory.
func (c *CPU) MemWrite16(pos uint16, data uint16) {
	hi := byte(data >> 8)
	lo := byte(data & 0xFF)
	c.MemWrite(pos, lo)
	c.MemWrite(pos+1, hi)
}

func (c *CPU) stackPush(data byte) {
	c.MemWrite(StackAddr+uint16(c.stackPointer), data)
	c.stackPointer -= 1
}

func (c *CPU) stackPush16(data uint16) {
	hi := byte(data >> 8)
	lo := byte(data & 0xFF)
	c.stackPush(hi)
	c.stackPush(lo)
}

func (c *CPU) stackPop() byte {
	c.stackPointer += 1
	return c.MemRead(StackAddr + uint16(c.stackPointer))
}

func (c *CPU) stackPop16() uint16 {
	lo := uint16(c.stackPop())
	hi := uint16(c.stackPop())
	return hi<<8 | lo
}
