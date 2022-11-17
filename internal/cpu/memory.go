package cpu

// MemRead reads byte from memory.
func (c *CPU) MemRead(addr uint16) byte {
	return c.bus.MemRead(addr)
}

// MemWrite writes byte to memory.
func (c *CPU) MemWrite(addr uint16, data byte) {
	c.bus.MemWrite(addr, data)
}

// MemRead16 reads two bytes from memory.
func (c *CPU) MemRead16(pos uint16) uint16 {
	lo := uint16(c.MemRead(pos))
	hi := uint16(c.MemRead(pos + 1))
	return hi<<8 | lo
}

// MemRead16Bug reads two bytes from memory, emulating a 6502 bug.
//
// JMP ($xxyy), or JMP indirect, does not advance pages if the lower eight bits
// of the specified address is $FF; the upper eight bits are fetched from $xx00,
// 255 bytes earlier, instead of the expected following byte.
//
// See [JMP Instruction Reference] and [NESdev CPU Errata].
//
// [JMP Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#JMP
// [NESDev CPU Errata]:https://www.nesdev.org/wiki/Errata#CPU
func (c *CPU) MemRead16Bug(pos uint16) uint16 {
	if pos&0x00FF == 0x00FF {
		lo := uint16(c.MemRead(pos))
		hi := uint16(c.MemRead(pos & 0xFF00))
		return hi<<8 | lo
	} else {
		return c.MemRead16(pos)
	}
}

// MemWrite16 writes two bytes to memory.
func (c *CPU) MemWrite16(pos uint16, data uint16) {
	hi := byte(data >> 8)
	lo := byte(data & 0xFF)
	c.MemWrite(pos, lo)
	c.MemWrite(pos+1, hi)
}

func (c *CPU) stackPush(data byte) {
	c.MemWrite(StackAddr+uint16(c.StackPointer), data)
	c.StackPointer -= 1
}

func (c *CPU) stackPush16(data uint16) {
	hi := byte(data >> 8)
	lo := byte(data & 0xFF)
	c.stackPush(hi)
	c.stackPush(lo)
}

func (c *CPU) stackPop() byte {
	c.StackPointer += 1
	return c.MemRead(StackAddr + uint16(c.StackPointer))
}

func (c *CPU) stackPop16() uint16 {
	lo := uint16(c.stackPop())
	hi := uint16(c.stackPop())
	return hi<<8 | lo
}
