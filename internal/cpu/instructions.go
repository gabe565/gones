package cpu

// lda - Load Accumulator
//
// Loads a byte of memory into the accumulator setting the zero and
// negative flags as appropriate.
//
// See [LDA Instruction Reference[].
//
// [LDA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#LDA
func (c *CPU) lda(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	v := c.memRead(addr)
	c.setRegisterA(v)
}

// sta - Store Accumulator
//
// Stores the contents of the accumulator into memory.
//
// See [STA Instruction Reference].
//
// [STA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#STA
func (c *CPU) sta(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.memWrite(addr, c.Accumulator)
}

// stx - Store X Register
//
// Stores the contents of the X register into memory.
//
// See [STX Instruction Reference].
//
// [STX Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#STX
func (c *CPU) stx(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.memWrite(addr, c.RegisterX)
}

// sty - Store Y Register
//
// Stores the contents of the Y register into memory.
//
// See [STY Instruction Reference].
//
// [STY Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#STY
func (c *CPU) sty(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	c.memWrite(addr, c.RegisterY)
}

// tax - Transfer Accumulator to X
//
// Copies the current contents of the accumulator into the X register
// and sets the zero and negative flags as appropriate.
//
// See [TAX Instruction Reference].
//
// [TAX Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TAX
func (c *CPU) tax() {
	c.RegisterX = c.Accumulator
	c.updateZeroAndNegFlags(c.RegisterX)
}

// tsx - Transfer Stack Pointer to X
//
// Copies the current contents of the stack register into the X register
// and sets the zero and negative flags as appropriate.
//
// See [TSX Instruction Reference].
//
// [TSX Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TSX
func (c *CPU) tsx() {
	c.RegisterX = c.SP
	c.updateZeroAndNegFlags(c.RegisterX)
}

// txa - Transfer X to Accumulator
//
// Copies the current contents of the X register into the accumulator
// and sets the zero and negative flags as appropriate.
//
// See [TXA Instruction Reference].
//
// [TXA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TXA
func (c *CPU) txa() {
	c.Accumulator = c.RegisterX
	c.updateZeroAndNegFlags(c.Accumulator)
}

// tay - Transfer Accumulator to Y
//
// Copies the current contents of the accumulator into the Y register
// and sets the zero and negative flags as appropriate.
//
// See [TAY Instruction Reference].
//
// [TAY Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TAY
func (c *CPU) tay() {
	c.RegisterY = c.Accumulator
	c.updateZeroAndNegFlags(c.RegisterY)
}

// inx - Increment X Register
//
// Adds one to the X register setting the zero and negative flags as appropriate.
//
// See [INX Instruction Reference].
//
// [INX Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#INX
func (c *CPU) inx() {
	c.RegisterX += 1
	c.updateZeroAndNegFlags(c.RegisterX)
}
