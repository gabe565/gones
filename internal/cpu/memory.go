package cpu

// ReadMem reads byte from memory.
func (c *CPU) ReadMem(addr uint16) byte {
	return c.bus.ReadMem(addr)
}

// WriteMem writes byte to memory.
func (c *CPU) WriteMem(addr uint16, data byte) {
	c.bus.WriteMem(addr, data)
}

// ReadMem16 reads two bytes from memory.
func (c *CPU) ReadMem16(addr uint16) uint16 {
	return c.bus.ReadMem16(addr)
}

// ReadMem16Bug reads two bytes from memory, emulating a 6502 bug.
//
// JMP ($xxyy), or JMP indirect, does not advance pages if the lower eight bits
// of the specified address is $FF; the upper eight bits are fetched from $xx00,
// 255 bytes earlier, instead of the expected following byte.
//
// See [JMP Instruction Reference] and [NESdev CPU Errata].
//
// [JMP Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#JMP
// [NESDev CPU Errata]:https://www.nesdev.org/wiki/Errata#CPU
func (c *CPU) ReadMem16Bug(addr uint16) uint16 {
	if addr&0x00FF == 0x00FF {
		lo := uint16(c.bus.ReadMem(addr))
		hi := uint16(c.bus.ReadMem(addr & 0xFF00))
		return hi<<8 | lo
	} else {
		return c.bus.ReadMem16(addr)
	}
}

// WriteMem16 writes two bytes to memory.
func (c *CPU) WriteMem16(addr uint16, data uint16) {
	c.bus.WriteMem16(addr, data)
}

// StackAddr is the memory address of the stack
const StackAddr = 0x100

func (c *CPU) stackPush(data byte) {
	c.WriteMem(StackAddr+uint16(c.StackPointer), data)
	c.StackPointer--
}

func (c *CPU) stackPush16(data uint16) {
	hi := byte(data >> 8)
	lo := byte(data & 0xFF)
	c.stackPush(hi)
	c.stackPush(lo)
}

func (c *CPU) stackPop() byte {
	c.StackPointer++
	return c.ReadMem(StackAddr + uint16(c.StackPointer))
}

func (c *CPU) stackPop16() uint16 {
	lo := uint16(c.stackPop())
	hi := uint16(c.stackPop())
	return hi<<8 | lo
}
