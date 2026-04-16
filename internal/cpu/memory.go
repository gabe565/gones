package cpu

// ReadMem reads byte from memory and ticks one CPU cycle.
func (c *CPU) ReadMem(addr uint16) byte {
	v := c.bus.ReadMem(addr)
	c.tick()
	return v
}

// WriteMem writes byte to memory and ticks one CPU cycle.
func (c *CPU) WriteMem(addr uint16, data byte) {
	c.bus.WriteMem(addr, data)
	c.tick()
}

// ReadMem16 reads two bytes from memory (2 CPU cycles).
func (c *CPU) ReadMem16(addr uint16) uint16 {
	lo := uint16(c.ReadMem(addr))
	hi := uint16(c.ReadMem(addr + 1))
	return hi<<8 | lo
}

// ReadMem16Bug reads two bytes from memory, emulating a 6502 bug (2 CPU cycles).
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
	lo := uint16(c.ReadMem(addr))
	hi := uint16(c.ReadMem(addr&0xFF00 | (addr+1)&0x00FF))
	return hi<<8 | lo
}

// ReadMemDMA reads memory without ticking. Used for DMA operations (e.g. DMC sample fetches).
func (c *CPU) ReadMemDMA(addr uint16) byte {
	return c.bus.ReadMem(addr)
}

// readMemSafe reads memory without ticking. Used for tracing/debugging.
func (c *CPU) readMemSafe(addr uint16) byte {
	return c.bus.ReadMem(addr)
}

// readMem16Safe reads two bytes without ticking.
func (c *CPU) readMem16Safe(addr uint16) uint16 {
	lo := uint16(c.bus.ReadMem(addr))
	hi := uint16(c.bus.ReadMem(addr + 1))
	return hi<<8 | lo
}

// StackAddr is the memory address of the stack.
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
