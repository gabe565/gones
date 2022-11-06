package cpu

import "github.com/gabe565/gones/internal/bitflags"

// adc - Add with Carry
//
// This instruction adds the contents of a memory location to the accumulator
// together with the carry bit. If overflow occurs the carry bit is set,
// this enables multiple byte addition to be performed.
//
// See [ADC Instruction Reference].
//
// [ADC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#ADC
func (c *CPU) adc(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	v := c.MemRead(addr)
	c.addAccumulator(v)
}

// and - Logical AND
//
// A logical AND is performed, bit by bit, on the accumulator contents
// using the contents of a byte of memory.
//
// See [AND Instruction Reference].
//
// [AND Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#AND
func (c *CPU) and(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.setAccumulator(c.accumulator & data)
}

// asl - Arithmetic Shift Left
//
// This operation shifts all the bits of the accumulator or memory contents
// one bit left. Bit 0 is set to 0 and bit 7 is placed in the carry flag.
// The effect of this operation is to multiply the memory contents by 2
// (ignoring 2's complement considerations), setting the carry if the result
// will not fit in 8 bits.
//
// See [ASL Instruction Reference].
//
// [ASL Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#ASL
func (c *CPU) asl(mode AddressingMode) {
	var addr uint16
	var data uint8
	if mode == Accumulator {
		data = c.accumulator
	} else {
		addr = c.getOperandAddress(mode)
		data = c.MemRead(addr)
	}
	c.status.Set(Carry, data>>7 == 1)
	data = data << 1
	if mode == Accumulator {
		c.setAccumulator(data)
	} else {
		c.MemWrite(addr, data)
		c.updateZeroAndNegFlags(data)
	}
}

// bcc - Branch if Carry Clear
//
// If the carry flag is clear then add the relative displacement to the
// program counter to cause a branch to a new location.
//
// See [BCC Instruction Reference].
//
// [BCC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BCC
func (c *CPU) bcc() {
	c.branch(!c.status.Has(Carry))
}

// bcs - Branch if Carry Set
//
// If the carry flag is set then add the relative displacement to the
// program counter to cause a branch to a new location.
//
// See [BCS Instruction Reference].
//
// [BCS Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BCS
func (c *CPU) bcs() {
	c.branch(c.status.Has(Carry))
}

// beq - Branch if Equal
//
// If the zero flag is set then add the relative displacement to the
// program counter to cause a branch to a new location.
//
// See [BEQ Instruction Reference].
//
// [BEQ Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BEQ
func (c *CPU) beq() {
	c.branch(c.status.Has(Zero))
}

// bit - Bit Test
//
// This instruction is used to test if one or more bits are set in a
// target memory location. The mask pattern in A is ANDed with the value
// in memory to set or clear the zero flag, but the result is not kept.
// Bits 7 and 6 of the value from memory are copied into the N and V flags.
//
// See [BIT Instruction Reference].
//
// [BIT Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BIT
func (c *CPU) bit(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.status.Set(Zero, data&c.accumulator == 0)
	c.status.Set(Negative, bitflags.Flags(data).Has(Negative))
	c.status.Set(Overflow, bitflags.Flags(data).Has(Overflow))
}

// bmi - Branch if Minus
//
// If the negative flag is set then add the relative displacement to the
// program counter to cause a branch to a new location.
//
// See [BMI Instruction Reference].
//
// [BMI Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BMI
func (c *CPU) bmi() {
	c.branch(c.status.Has(Negative))
}

// bne - Branch if Not Equal
//
// If the zero flag is clear then add the relative displacement to the
// program counter to cause a branch to a new location.
//
// See [BNE Instruction Reference].
//
// [BNE Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BNE
func (c *CPU) bne() {
	c.branch(!c.status.Has(Zero))
}

// bpl - Branch if Positive
//
// If the negative flag is clear then add the relative displacement to
// the program counter to cause a branch to a new location.
//
// See [BPL Instruction Reference].
//
// [BPL Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BPL
func (c *CPU) bpl() {
	c.branch(!c.status.Has(Negative))
}

// bvc - Branch if Overflow Clear
//
// If the overflow flag is clear then add the relative displacement to
// the program counter to cause a branch to a new location.
//
// See [BVC Instruction Reference].
//
// [BVC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BVC
func (c *CPU) bvc() {
	c.branch(!c.status.Has(Overflow))
}

// bvs - Branch if Overflow Set
//
// If the overflow flag is set then add the relative displacement to
// the program counter to cause a branch to a new location.
//
// See [BVS Instruction Reference].
//
// [BVS Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BVS
func (c *CPU) bvs() {
	c.branch(c.status.Has(Overflow))
}

// clc - Clear Carry Flag
//
// Sets the carry flag to zero.
//
// See [CLC Instruction Reference].
//
// [CLC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CLC
func (c *CPU) clc() {
	c.status.Remove(Carry)
}

// cld - Clear Decimal Mode
//
// Sets the decimal mode flag to zero.
//
// See [CLC Instruction Reference].
//
// [CLC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CLC
func (c *CPU) cld() {
	c.status.Remove(DecimalMode)
}

// cli - Clear Interrupt Disable
//
// Clears the interrupt disable flag allowing normal interrupt requests
// to be serviced.
//
// See [CLI Instruction Reference].
//
// [CLI Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CLI
func (c *CPU) cli() {
	c.status.Remove(InterruptDisable)
}

// clv - Clear Overflow Flag
//
// Clears the overflow flag.
//
// See [CLV Instruction Reference].
//
// [CLV Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CLV
func (c *CPU) clv() {
	c.status.Remove(Overflow)
}

// cmp - Compare
//
// This instruction compares the contents of the accumulator with
// another memory held value and sets the zero and carry flags as appropriate.
//
// See [CMP Instruction Reference].
//
// [CMP Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CMP
func (c *CPU) cmp(mode AddressingMode) {
	c.compare(mode, c.accumulator)
}

// cpx - Compare X Register
//
// This instruction compares the contents of the X register with
// another memory held value and sets the zero and carry flags as appropriate.
//
// See [CPX Instruction Reference].
//
// [CPX Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CPX
func (c *CPU) cpx(mode AddressingMode) {
	c.compare(mode, c.registerX)
}

// cpy - Compare Y Register
//
// This instruction compares the contents of the Y register with
// another memory held value and sets the zero and carry flags as appropriate.
//
// See [CPY Instruction Reference].
//
// [CPY Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CPY
func (c *CPU) cpy(mode AddressingMode) {
	c.compare(mode, c.registerY)
}

// dec - Decrement Memory
//
// Subtracts one from the value held at a specified memory location setting
// the zero and negative flags as appropriate.
//
// See [DEC Instruction Reference].
//
// [DEC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#DEC
func (c *CPU) dec(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	data -= 1
	c.MemWrite(addr, data)
	c.updateZeroAndNegFlags(data)
}

// dex - Decrement X Register
//
// Subtracts one from the X register setting the zero and negative flags
// as appropriate.
//
// See [DEX Instruction Reference].
//
// [DEX Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#DEX
func (c *CPU) dex() {
	c.registerX -= 1
	c.updateZeroAndNegFlags(c.registerX)
}

// dey - Decrement Y Register
//
// Subtracts one from the Y register setting the zero and negative flags
// as appropriate.
//
// See [DEY Instruction Reference].
//
// [DEY Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#DEY
func (c *CPU) dey() {
	c.registerY -= 1
	c.updateZeroAndNegFlags(c.registerY)
}

// eor - Exclusive OR
//
// An exclusive OR is performed, bit by bit, on the accumulator contents
// using the contents of a byte of memory.
//
// See [EOR Instruction Reference].
//
// [EOR Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#EOR
func (c *CPU) eor(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.setAccumulator(data ^ c.accumulator)
}

// inc - Increment Memory
//
// Adds one to the value held at a specified memory location setting
// the zero and negative flags as appropriate.
//
// See [INC Instruction Reference].
//
// [INC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#INC
func (c *CPU) inc(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	data += 1
	c.MemWrite(addr, data)
	c.updateZeroAndNegFlags(data)
}

// inx - Increment X Register
//
// Adds one to the X register setting the zero and negative flags as appropriate.
//
// See [INX Instruction Reference].
//
// [INX Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#INX
func (c *CPU) inx() {
	c.registerX += 1
	c.updateZeroAndNegFlags(c.registerX)
}

// iny - Increment Y Register
//
// Adds one to the Y register setting the zero and
// negative flags as appropriate.
//
// See [INY Instruction Reference].
//
// [INY Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#INY
func (c *CPU) iny() {
	c.registerY += 1
	c.updateZeroAndNegFlags(c.registerY)
}

// jmp - Jump
//
// Sets the program counter to the address specified by the operand.
//
// See [JMP Instruction Reference].
//
// [JMP Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#JMP
func (c *CPU) jmp(mode AddressingMode) {
	switch mode {
	case Absolute:
		addr := c.MemRead16(c.programCounter)
		c.programCounter = addr
	case Indirect:
		addr := c.MemRead16(c.programCounter)

		// let indirect_ref = self.mem_read_u16(mem_address);
		//6502 bug mode with with page boundary:
		//  if address $3000 contains $40, $30FF contains $80, and $3100 contains $50,
		// the result of JMP ($30FF) will be a transfer of control to $4080 rather than $5080 as you intended
		// i.e. the 6502 took the low byte of the address from $30FF and the high byte from $3000

		var indirect uint16
		if addr&0x00FF == 0x00FF {
			lo := c.MemRead(addr)
			hi := c.MemRead(addr & 0xFF00)
			indirect = uint16(hi)<<8 | uint16(lo)
		} else {
			indirect = c.MemRead16(addr)
		}
		c.programCounter = indirect
	}

}

// jsr - Jump to Subroutine
//
// The JSR instruction pushes the address (minus one) of the return point on to
// the stack and then sets the program counter to the target memory address.
//
// See [JSR Instruction Reference].
//
// [JSR Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#JSR
func (c *CPU) jsr() {
	c.stackPush16(c.programCounter + 1)
	addr := c.MemRead16(c.programCounter)
	c.programCounter = addr
}

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
	v := c.MemRead(addr)
	c.setAccumulator(v)
}

// ldx - Load X Register
//
// Loads a byte of memory into the X register setting the zero and
// negative flags as appropriate.
//
// See [LDX Instruction Reference[].
//
// [LDX Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#LDX
func (c *CPU) ldx(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.registerX = data
	c.updateZeroAndNegFlags(c.registerX)
}

// ldy - Load Y Register
//
// Loads a byte of memory into the Y register setting the zero and
// negative flags as appropriate.
//
// See [LDY Instruction Reference[].
//
// [LDY Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#LDY
func (c *CPU) ldy(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.registerY = data
	c.updateZeroAndNegFlags(c.registerY)
}

// lsr - Logical Shift Right
//
// Each of the bits in A or M is shift one place to the right.
// The bit that was in bit 0 is shifted into the carry flag.
// Bit 7 is set to zero.
//
// See [LSR Instruction Reference[].
//
// [LSR Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#LSR
func (c *CPU) lsr(mode AddressingMode) {
	var addr uint16
	var data uint8
	if mode == Accumulator {
		data = c.accumulator
	} else {
		addr = c.getOperandAddress(mode)
		data = c.MemRead(addr)
	}
	c.status.Set(Carry, data&1 == 1)
	data >>= 1
	if mode == Accumulator {
		c.setAccumulator(data)
	} else {
		c.MemWrite(addr, data)
		c.updateZeroAndNegFlags(data)
	}
}

// ora - Logical Inclusive OR
//
// An inclusive OR is performed, bit by bit, on the accumulator contents
// using the contents of a byte of memory.
//
// See [ORA Instruction Reference[].
//
// [ORA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#ORA
func (c *CPU) ora(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.setAccumulator(data | c.accumulator)
}

// pha - Push Accumulator
//
// Pushes a copy of the accumulator on to the stack.
//
// See [PHA Instruction Reference[].
//
// [PHA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#PHA
func (c *CPU) pha() {
	c.stackPush(c.accumulator)
}

// php - Push Processor Status
//
// Pushes a copy of the status flags on to the stack.
//
// See [PHP Instruction Reference[].
//
// [PHP Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#PHP
func (c *CPU) php() {
	flags := c.status
	flags.Insert(Break | Break2)
	c.stackPush(uint8(flags))
}

// pla - Pull Accumulator
//
// Pulls an 8 bit value from the stack and into the accumulator.
// The zero and negative flags are set as appropriate.
//
// See [PLA Instruction Reference[].
//
// [PLA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#PLA
func (c *CPU) pla() {
	data := c.stackPop()
	c.setAccumulator(data)
}

// plp - Pull Processor Status
//
// Pulls an 8 bit value from the stack and into the processor flags.
// The flags will take on new states as determined by the value pulled.
//
// See [PLP Instruction Reference[].
//
// [PLP Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#PLP
func (c *CPU) plp() {
	flags := bitflags.Flags(c.stackPop())
	flags.Remove(Break | Break2)
	c.status = flags
}

// rol - Rotate Left
//
// Move each of the bits in either A or M one place to the left.
// Bit 0 is filled with the current value of the carry flag
// whilst the old bit 7 becomes the new carry flag value.
//
// See [ROL Instruction Reference[].
//
// [ROL Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#ROL
func (c *CPU) rol(mode AddressingMode) {
	var addr uint16
	var data uint8
	if mode == Accumulator {
		data = c.accumulator
	} else {
		addr = c.getOperandAddress(mode)
		data = c.MemRead(addr)
	}
	prevCarry := c.status.Has(Carry)

	c.status.Set(Carry, data>>7 == 1)
	data <<= 1
	if prevCarry {
		data |= 1
	}
	if mode == Accumulator {
		c.setAccumulator(data)
	} else {
		c.MemWrite(addr, data)
		c.updateZeroAndNegFlags(data)
	}
}

// ror - Rotate Right
//
// Move each of the bits in either A or M one place to the right.
// Bit 0 is filled with the current value of the carry flag whilst
// the old bit 7 becomes the new carry flag value.
//
// See [ROR Instruction Reference[].
//
// [ROR Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#ROR
func (c *CPU) ror(mode AddressingMode) {
	var addr uint16
	var data uint8
	if mode == Accumulator {
		data = c.accumulator
	} else {
		addr = c.getOperandAddress(mode)
		data = c.MemRead(addr)
	}
	prevCarry := c.status.Has(Carry)

	c.status.Set(Carry, data&1 == 1)
	data >>= 1
	if prevCarry {
		data |= uint8(Negative)
	}
	if mode == Accumulator {
		c.setAccumulator(data)
	} else {
		c.MemWrite(addr, data)
		c.updateZeroAndNegFlags(data)
	}
}

// rti - Return from Interrupt
//
// The RTI instruction is used at the end of an interrupt processing routine.
// It pulls the processor flags from the stack followed by the program counter.
//
// See [RTI Instruction Reference[].
//
// [RTI Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#RTI
func (c *CPU) rti() {
	flags := bitflags.Flags(c.stackPop())
	flags.Remove(Break | Break2)
	c.status = flags

	c.programCounter = c.stackPop16()
}

// rts - Return from Subroutine
//
// The RTS instruction is used at the end of a subroutine to return
// to the calling routine. It pulls the program counter (minus one)
// from the stack.
//
// See [RTS Instruction Reference].
//
// [RTS Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#RTS
func (c *CPU) rts() {
	c.programCounter = c.stackPop16() + 1
}

// sbc - Subtract with Carry
//
// This instruction subtracts the contents of a memory location to the
// accumulator together with the not of the carry bit. If overflow occurs
// the carry bit is clear, this enables multiple byte subtraction to be performed.
//
// See [SBC Instruction Reference].
//
// [SBC Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#SBC
func (c *CPU) sbc(mode AddressingMode) {
	addr := c.getOperandAddress(mode)
	v := c.MemRead(addr)
	c.addAccumulator(uint8(-int8(v) - 1))
}

// sec - Set Carry Flag
//
// Set the carry flag to one.
//
// See [SEC Instruction Reference].
//
// [SEC Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#SEC
func (c *CPU) sec() {
	c.status.Insert(Carry)
}

// sed - Set Decimal Flag
//
// Set the decimal mode flag to one.
//
// See [SED Instruction Reference].
//
// [SED Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#SED
func (c *CPU) sed() {
	c.status.Insert(DecimalMode)
}

// sei - Set Interrupt Disable
//
// Set the interrupt disable flag to one.
//
// See [SEI Instruction Reference].
//
// [SEI Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#SEI
func (c *CPU) sei() {
	c.status.Insert(InterruptDisable)
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
	c.MemWrite(addr, c.accumulator)
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
	c.MemWrite(addr, c.registerX)
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
	c.MemWrite(addr, c.registerY)
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
	c.registerX = c.accumulator
	c.updateZeroAndNegFlags(c.registerX)
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
	c.registerX = c.stackPointer
	c.updateZeroAndNegFlags(c.registerX)
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
	c.setAccumulator(c.registerX)
}

// txs - Transfer X to Stack Pointer
//
// Copies the current contents of the X register into the stack register.
//
// See [TXS Instruction Reference].
//
// [TXS Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TXS
func (c *CPU) txs() {
	c.stackPointer = c.registerX
}

// tya - Transfer Y to Accumulator
//
// Copies the current contents of the Y register into the accumulator
// and sets the zero and negative flags as appropriate.
//
// See [TYA Instruction Reference].
//
// [TYA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TYA
func (c *CPU) tya() {
	c.setAccumulator(c.registerY)
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
	c.registerY = c.accumulator
	c.updateZeroAndNegFlags(c.registerY)
}
