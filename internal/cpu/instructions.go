package cpu

import (
	"errors"
	"github.com/gabe565/gones/internal/bitflags"
)

type Instruction func(c *CPU, mode AddressingMode) error

// adc - Add with Carry
//
// This instruction adds the contents of a memory location to the accumulator
// together with the carry bit. If overflow occurs the carry bit is set,
// this enables multiple byte addition to be performed.
//
// See [ADC Instruction Reference].
//
// [ADC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#ADC
func adc(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	v := c.MemRead(addr)
	c.addAccumulator(v)
	return nil
}

// ahx - Undocumented Opcode
//
// AND X register with accumulator then AND result with 7 and store in memory.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func ahx(c *CPU, mode AddressingMode) error {
	var pos uint16
	switch mode {
	case IndirectY:
		pos = uint16(c.MemRead(c.ProgramCounter))
	case AbsoluteY:
		pos = c.ProgramCounter
	}
	addr := c.MemRead16(pos) + uint16(c.RegisterY)
	data := c.Accumulator & c.RegisterX & uint8(addr>>8)
	c.MemWrite(addr, data)
	return nil
}

// alr - Undocumented Opcode
//
// AND byte with accumulator, then shift right one bit in accumulator.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func alr(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.setAccumulator(data & c.Accumulator)
	return lsr(c, Accumulator)
}

// anc - Undocumented Opcode
//
// AND byte with accumulator. If result is negative then carry is set.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func anc(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.setAccumulator(data & c.Accumulator)
	c.Status.Set(Carry, c.Status.Has(Negative))
	return nil
}

// and - Logical AND
//
// A logical AND is performed, bit by bit, on the accumulator contents
// using the contents of a byte of memory.
//
// See [AND Instruction Reference].
//
// [AND Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#AND
func and(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	data := c.MemRead(addr)
	c.setAccumulator(c.Accumulator & data)
	return nil
}

// arr - Undocumented Opcode
//
// AND byte with accumulator, then rotate one bit right in accumulator
// and check bit 5 and 6:
// - If both bits are 1: set C, clear V.
// - If both bits are 0: clear C and V.
// - If only bit 5 is 1: set V, clear C.
// - If only bit 6 is 1: set C and V.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func arr(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.setAccumulator(data & c.Accumulator)
	if err := ror(c, Accumulator); err != nil {
		return err
	}

	bit6 := (c.Accumulator >> 6) & 1
	c.Status.Set(Carry, bit6 == 1)
	bit5 := (c.Accumulator >> 5) & 1
	c.Status.Set(Overflow, bit5^bit6 == 1)
	c.updateZeroAndNegFlags(c.Accumulator)
	return nil
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
func asl(c *CPU, mode AddressingMode) error {
	var addr uint16
	var data byte
	if mode == Accumulator {
		data = c.Accumulator
	} else {
		var pageCrossed bool
		addr, pageCrossed = c.getOperandAddress(mode)
		if pageCrossed {
			defer func() {
				c.Bus.Tick(1)
			}()
		}
		data = c.MemRead(addr)
	}
	c.Status.Set(Carry, data>>7 == 1)
	data = data << 1
	if mode == Accumulator {
		c.setAccumulator(data)
	} else {
		c.MemWrite(addr, data)
		c.updateZeroAndNegFlags(data)
	}
	return nil
}

// axs - Undocumented Opcode
//
// AND X register with accumulator and store result in X register, then
// subtract byte from X register (without borrow).
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func axs(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	result := c.RegisterX & c.Accumulator
	c.Status.Set(Carry, data <= result)
	result -= 1
	c.RegisterX = result
	c.updateZeroAndNegFlags(result)
	return nil
}

// bcc - Branch if Carry Clear
//
// If the carry flag is clear then add the relative displacement to the
// program counter to cause a branch to a new location.
//
// See [BCC Instruction Reference].
//
// [BCC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BCC
func bcc(c *CPU, mode AddressingMode) error {
	c.branch(!c.Status.Has(Carry))
	return nil
}

// bcs - Branch if Carry Set
//
// If the carry flag is set then add the relative displacement to the
// program counter to cause a branch to a new location.
//
// See [BCS Instruction Reference].
//
// [BCS Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BCS
func bcs(c *CPU, mode AddressingMode) error {
	c.branch(c.Status.Has(Carry))
	return nil
}

// beq - Branch if Equal
//
// If the zero flag is set then add the relative displacement to the
// program counter to cause a branch to a new location.
//
// See [BEQ Instruction Reference].
//
// [BEQ Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BEQ
func beq(c *CPU, mode AddressingMode) error {
	c.branch(c.Status.Has(Zero))
	return nil
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
func bit(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.Status.Set(Zero, data&c.Accumulator == 0)
	c.Status.Set(Negative, bitflags.Flags(data).Has(Negative))
	c.Status.Set(Overflow, bitflags.Flags(data).Has(Overflow))
	return nil
}

// bmi - Branch if Minus
//
// If the negative flag is set then add the relative displacement to the
// program counter to cause a branch to a new location.
//
// See [BMI Instruction Reference].
//
// [BMI Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BMI
func bmi(c *CPU, mode AddressingMode) error {
	c.branch(c.Status.Has(Negative))
	return nil
}

// bne - Branch if Not Equal
//
// If the zero flag is clear then add the relative displacement to the
// program counter to cause a branch to a new location.
//
// See [BNE Instruction Reference].
//
// [BNE Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BNE
func bne(c *CPU, mode AddressingMode) error {
	c.branch(!c.Status.Has(Zero))
	return nil
}

// bpl - Branch if Positive
//
// If the negative flag is clear then add the relative displacement to
// the program counter to cause a branch to a new location.
//
// See [BPL Instruction Reference].
//
// [BPL Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BPL
func bpl(c *CPU, mode AddressingMode) error {
	c.branch(!c.Status.Has(Negative))
	return nil
}

var ErrBrk = errors.New("break")

// brk - Force Interrupt
//
// The BRK instruction forces the generation of an interrupt request.
// The program counter and processor status are pushed on the stack then
// the IRQ interrupt vector at $FFFE/F is loaded into the PC and
// the break flag in the status set to one.
//
// See [BRK Instruction Reference].
//
// [BRK Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BRK
func brk(c *CPU, mode AddressingMode) error {
	return ErrBrk
}

// bvc - Branch if Overflow Clear
//
// If the overflow flag is clear then add the relative displacement to
// the program counter to cause a branch to a new location.
//
// See [BVC Instruction Reference].
//
// [BVC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BVC
func bvc(c *CPU, mode AddressingMode) error {
	c.branch(!c.Status.Has(Overflow))
	return nil
}

// bvs - Branch if Overflow Set
//
// If the overflow flag is set then add the relative displacement to
// the program counter to cause a branch to a new location.
//
// See [BVS Instruction Reference].
//
// [BVS Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#BVS
func bvs(c *CPU, mode AddressingMode) error {
	c.branch(c.Status.Has(Overflow))
	return nil
}

// clc - Clear Carry Flag
//
// Sets the carry flag to zero.
//
// See [CLC Instruction Reference].
//
// [CLC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CLC
func clc(c *CPU, mode AddressingMode) error {
	c.Status.Remove(Carry)
	return nil
}

// cld - Clear Decimal Mode
//
// Sets the decimal mode flag to zero.
//
// See [CLC Instruction Reference].
//
// [CLC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CLC
func cld(c *CPU, mode AddressingMode) error {
	c.Status.Remove(Decimal)
	return nil
}

// cli - Clear Interrupt Disable
//
// Clears the interrupt disable flag allowing normal interrupt requests
// to be serviced.
//
// See [CLI Instruction Reference].
//
// [CLI Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CLI
func cli(c *CPU, mode AddressingMode) error {
	c.Status.Remove(InterruptDisable)
	return nil
}

// clv - Clear Overflow Flag
//
// Clears the overflow flag.
//
// See [CLV Instruction Reference].
//
// [CLV Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CLV
func clv(c *CPU, mode AddressingMode) error {
	c.Status.Remove(Overflow)
	return nil
}

// cmp - Compare
//
// This instruction compares the contents of the accumulator with
// another memory held value and sets the zero and carry flags as appropriate.
//
// See [CMP Instruction Reference].
//
// [CMP Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CMP
func cmp(c *CPU, mode AddressingMode) error {
	c.compare(mode, c.Accumulator)
	return nil
}

// cpx - Compare X Register
//
// This instruction compares the contents of the X register with
// another memory held value and sets the zero and carry flags as appropriate.
//
// See [CPX Instruction Reference].
//
// [CPX Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CPX
func cpx(c *CPU, mode AddressingMode) error {
	c.compare(mode, c.RegisterX)
	return nil
}

// cpy - Compare Y Register
//
// This instruction compares the contents of the Y register with
// another memory held value and sets the zero and carry flags as appropriate.
//
// See [CPY Instruction Reference].
//
// [CPY Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#CPY
func cpy(c *CPU, mode AddressingMode) error {
	c.compare(mode, c.RegisterY)
	return nil
}

// dcp - Undocumented Opcode
//
// Subtract 1 from memory (without borrow).
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func dcp(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	data -= 1
	c.MemWrite(addr, data)
	c.Status.Set(Carry, data <= c.Accumulator)
	c.updateZeroAndNegFlags(c.Accumulator - data)
	return nil
}

// dec - Decrement Memory
//
// Subtracts one from the value held at a specified memory location setting
// the zero and negative flags as appropriate.
//
// See [DEC Instruction Reference].
//
// [DEC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#DEC
func dec(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	data -= 1
	c.MemWrite(addr, data)
	c.updateZeroAndNegFlags(data)
	return nil
}

// dex - Decrement X Register
//
// Subtracts one from the X register setting the zero and negative flags
// as appropriate.
//
// See [DEX Instruction Reference].
//
// [DEX Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#DEX
func dex(c *CPU, mode AddressingMode) error {
	c.RegisterX -= 1
	c.updateZeroAndNegFlags(c.RegisterX)
	return nil
}

// dey - Decrement Y Register
//
// Subtracts one from the Y register setting the zero and negative flags
// as appropriate.
//
// See [DEY Instruction Reference].
//
// [DEY Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#DEY
func dey(c *CPU, mode AddressingMode) error {
	c.RegisterY -= 1
	c.updateZeroAndNegFlags(c.RegisterY)
	return nil
}

// eor - Exclusive OR
//
// An exclusive OR is performed, bit by bit, on the accumulator contents
// using the contents of a byte of memory.
//
// See [EOR Instruction Reference].
//
// [EOR Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#EOR
func eor(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	data := c.MemRead(addr)
	c.setAccumulator(data ^ c.Accumulator)
	return nil
}

// inc - Increment Memory
//
// Adds one to the value held at a specified memory location setting
// the zero and negative flags as appropriate.
//
// See [INC Instruction Reference].
//
// [INC Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#INC
func inc(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	data += 1
	c.MemWrite(addr, data)
	c.updateZeroAndNegFlags(data)
	return nil
}

// inx - Increment X Register
//
// Adds one to the X register setting the zero and negative flags as appropriate.
//
// See [INX Instruction Reference].
//
// [INX Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#INX
func inx(c *CPU, mode AddressingMode) error {
	c.RegisterX += 1
	c.updateZeroAndNegFlags(c.RegisterX)
	return nil
}

// iny - Increment Y Register
//
// Adds one to the Y register setting the zero and
// negative flags as appropriate.
//
// See [INY Instruction Reference].
//
// [INY Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#INY
func iny(c *CPU, mode AddressingMode) error {
	c.RegisterY += 1
	c.updateZeroAndNegFlags(c.RegisterY)
	return nil
}

// isb - Undocumented Opcode
//
// Increase memory by one, then subtract memory from accumulator (with borrow).
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func isb(c *CPU, mode AddressingMode) error {
	if err := inc(c, mode); err != nil {
		return err
	}
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.addAccumulator(byte(-int8(data) - 1))
	return nil
}

// jmp - Jump
//
// Sets the program counter to the address specified by the operand.
//
// See [JMP Instruction Reference].
//
// [JMP Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#JMP
func jmp(c *CPU, mode AddressingMode) error {
	addr := c.MemRead16(c.ProgramCounter)

	switch mode {
	case Absolute:
		c.ProgramCounter = addr
	case Indirect:
		// 6502 bug mode with page boundary:
		// If address $3000 contains $40, $30FF contains $80, and $3100 contains $50,
		// the result of JMP ($30FF) will be a transfer of control to $4080
		// rather than $5080 as expected
		// i.e. the 6502 took the low byte of the address from $30FF and the high byte from $3000
		var indirect uint16
		if addr&0x00FF == 0x00FF {
			lo := c.MemRead(addr)
			hi := c.MemRead(addr & 0xFF00)
			indirect = uint16(hi)<<8 | uint16(lo)
		} else {
			indirect = c.MemRead16(addr)
		}
		c.ProgramCounter = indirect
	}

	return nil
}

// jsr - Jump to Subroutine
//
// The JSR instruction pushes the address (minus one) of the return point on to
// the stack and then sets the program counter to the target memory address.
//
// See [JSR Instruction Reference].
//
// [JSR Instruction Reference]: https://www.nesdev.org/obelisk-6502-guide/reference.html#JSR
func jsr(c *CPU, mode AddressingMode) error {
	c.stackPush16(c.ProgramCounter + 1)
	addr := c.MemRead16(c.ProgramCounter)
	c.ProgramCounter = addr
	return nil
}

// las - Undocumented Opcode
//
// AND memory with stack pointer, transfer result to accumulator, X register,
// and stack pointer.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func las(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	data := c.MemRead(addr)
	data &= c.StackPointer
	c.Accumulator = data
	c.RegisterX = data
	c.StackPointer = data
	c.updateZeroAndNegFlags(data)
	return nil
}

// lax - Undocumented Opcode
//
// Load accumulator and X register with memory.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func lax(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	data := c.MemRead(addr)
	c.setAccumulator(data)
	c.RegisterX = c.Accumulator
	return nil
}

// lda - Load Accumulator
//
// Loads a byte of memory into the accumulator setting the zero and
// negative flags as appropriate.
//
// See [LDA Instruction Reference].
//
// [LDA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#LDA
func lda(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	v := c.MemRead(addr)
	c.setAccumulator(v)
	return nil
}

// ldx - Load X Register
//
// Loads a byte of memory into the X register setting the zero and
// negative flags as appropriate.
//
// See [LDX Instruction Reference].
//
// [LDX Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#LDX
func ldx(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	data := c.MemRead(addr)
	c.RegisterX = data
	c.updateZeroAndNegFlags(c.RegisterX)
	return nil
}

// ldy - Load Y Register
//
// Loads a byte of memory into the Y register setting the zero and
// negative flags as appropriate.
//
// See [LDY Instruction Reference].
//
// [LDY Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#LDY
func ldy(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	data := c.MemRead(addr)
	c.RegisterY = data
	c.updateZeroAndNegFlags(c.RegisterY)
	return nil
}

// lsr - Logical Shift Right
//
// Each of the bits in A or M is shift one place to the right.
// The bit that was in bit 0 is shifted into the carry flag.
// Bit 7 is set to zero.
//
// See [LSR Instruction Reference].
//
// [LSR Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#LSR
func lsr(c *CPU, mode AddressingMode) error {
	var addr uint16
	var data byte
	if mode == Accumulator {
		data = c.Accumulator
	} else {
		addr, _ = c.getOperandAddress(mode)
		data = c.MemRead(addr)
	}
	c.Status.Set(Carry, data&1 == 1)
	data >>= 1
	if mode == Accumulator {
		c.setAccumulator(data)
	} else {
		c.MemWrite(addr, data)
		c.updateZeroAndNegFlags(data)
	}
	return nil
}

// lxa - Undocumented Opcode
//
// AND byte with accumulator, then transfer accumulator to X register.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func lxa(c *CPU, mode AddressingMode) error {
	if err := lda(c, mode); err != nil {
		return err
	}
	return tax(c, mode)
}

// nop - No Operation
//
// The NOP instruction causes no changes to the processor other than
// the normal incrementing of the program counter to the next instruction.
//
// See [NOP Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#NOP
//
// [NOP Instruction Reference]:
func nop(c *CPU, mode AddressingMode) error {
	if mode != Implied {
		addr, pageCrossed := c.getOperandAddress(mode)
		if pageCrossed {
			defer func() {
				c.Bus.Tick(1)
			}()
		}
		_ = c.MemRead(addr)
	}
	return nil
}

// ora - Logical Inclusive OR
//
// An inclusive OR is performed, bit by bit, on the accumulator contents
// using the contents of a byte of memory.
//
// See [ORA Instruction Reference].
//
// [ORA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#ORA
func ora(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	data := c.MemRead(addr)
	c.setAccumulator(data | c.Accumulator)
	return nil
}

// pha - Push Accumulator
//
// Pushes a copy of the accumulator on to the stack.
//
// See [PHA Instruction Reference].
//
// [PHA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#PHA
func pha(c *CPU, mode AddressingMode) error {
	c.stackPush(c.Accumulator)
	return nil
}

// php - Push Processor Status
//
// Pushes a copy of the status flags on to the stack.
//
// See [PHP Instruction Reference].
//
// [PHP Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#PHP
func php(c *CPU, mode AddressingMode) error {
	flags := c.Status
	flags.Insert(Break | Break2)
	c.stackPush(byte(flags))
	return nil
}

// pla - Pull Accumulator
//
// Pulls an 8 bit value from the stack and into the accumulator.
// The zero and negative flags are set as appropriate.
//
// See [PLA Instruction Reference].
//
// [PLA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#PLA
func pla(c *CPU, mode AddressingMode) error {
	data := c.stackPop()
	c.setAccumulator(data)
	return nil
}

// plp - Pull Processor Status
//
// Pulls an 8 bit value from the stack and into the processor flags.
// The flags will take on new states as determined by the value pulled.
//
// See [PLP Instruction Reference].
//
// [PLP Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#PLP
func plp(c *CPU, mode AddressingMode) error {
	flags := bitflags.Flags(c.stackPop())
	flags.Remove(Break)
	flags.Insert(Break2)
	c.Status = flags
	return nil
}

// rla - Undocumented Opcode
//
// Rotate one bit left in memory, then AND accumulator with memory.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func rla(c *CPU, mode AddressingMode) error {
	if err := rol(c, mode); err != nil {
		return err
	}
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.setAccumulator(data & c.Accumulator)
	return nil
}

// rol - Rotate Left
//
// Move each of the bits in either A or M one place to the left.
// Bit 0 is filled with the current value of the carry flag
// whilst the old bit 7 becomes the new carry flag value.
//
// See [ROL Instruction Reference].
//
// [ROL Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#ROL
func rol(c *CPU, mode AddressingMode) error {
	var addr uint16
	var data byte
	if mode == Accumulator {
		data = c.Accumulator
	} else {
		addr, _ = c.getOperandAddress(mode)
		data = c.MemRead(addr)
	}
	prevCarry := c.Status.Has(Carry)

	c.Status.Set(Carry, data>>7 == 1)
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
	return nil
}

// ror - Rotate Right
//
// Move each of the bits in either A or M one place to the right.
// Bit 0 is filled with the current value of the carry flag whilst
// the old bit 7 becomes the new carry flag value.
//
// See [ROR Instruction Reference].
//
// [ROR Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#ROR
func ror(c *CPU, mode AddressingMode) error {
	var addr uint16
	var data byte
	if mode == Accumulator {
		data = c.Accumulator
	} else {
		addr, _ = c.getOperandAddress(mode)
		data = c.MemRead(addr)
	}
	prevCarry := c.Status.Has(Carry)

	c.Status.Set(Carry, data&1 == 1)
	data >>= 1
	if prevCarry {
		data |= byte(Negative)
	}
	if mode == Accumulator {
		c.setAccumulator(data)
	} else {
		c.MemWrite(addr, data)
		c.updateZeroAndNegFlags(data)
	}
	return nil
}

// rra - Undocumented Opcode
//
// Rotate one bit right in memory, then add memory to accumulator (with carry).
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func rra(c *CPU, mode AddressingMode) error {
	if err := ror(c, mode); err != nil {
		return err
	}
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.addAccumulator(data)
	return nil
}

// rti - Return from Interrupt
//
// The RTI instruction is used at the end of an interrupt processing routine.
// It pulls the processor flags from the stack followed by the program counter.
//
// See [RTI Instruction Reference].
//
// [RTI Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#RTI
func rti(c *CPU, mode AddressingMode) error {
	flags := bitflags.Flags(c.stackPop())
	flags.Remove(Break)
	flags.Insert(Break2)
	c.Status = flags

	c.ProgramCounter = c.stackPop16()
	return nil
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
func rts(c *CPU, mode AddressingMode) error {
	c.ProgramCounter = c.stackPop16() + 1
	return nil
}

// sax - Undocumented Opcode
//
// AND X register with accumulator and store result in X register, then
// subtract byte from X register (without borrow).
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func sax(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	data := c.Accumulator & c.RegisterX
	c.MemWrite(addr, data)
	return nil
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
func sbc(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	v := c.MemRead(addr)
	c.addAccumulator(byte(-int8(v) - 1))
	return nil
}

// sec - Set Carry Flag
//
// Set the carry flag to one.
//
// See [SEC Instruction Reference].
//
// [SEC Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#SEC
func sec(c *CPU, mode AddressingMode) error {
	c.Status.Insert(Carry)
	return nil
}

// sed - Set Decimal Flag
//
// Set the decimal mode flag to one.
//
// See [SED Instruction Reference].
//
// [SED Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#SED
func sed(c *CPU, mode AddressingMode) error {
	c.Status.Insert(Decimal)
	return nil
}

// sei - Set Interrupt Disable
//
// Set the interrupt disable flag to one.
//
// See [SEI Instruction Reference].
//
// [SEI Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#SEI
func sei(c *CPU, mode AddressingMode) error {
	c.Status.Insert(InterruptDisable)
	return nil
}

// shx - Undocumented Opcode
//
// AND X register with the high byte of the target address of the argument + 1.
// Store the result in memory.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func shx(c *CPU, mode AddressingMode) error {
	addr := c.MemRead16(c.ProgramCounter) + uint16(c.RegisterY)
	data := c.RegisterX & (uint8(addr>>8) + 1)
	c.MemWrite(addr, data)
	return nil
}

// shy - Undocumented Opcode
//
// AND Y register with the high byte of the target address of the argument + 1.
// Store the result in memory.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func shy(c *CPU, mode AddressingMode) error {
	addr := c.MemRead16(c.ProgramCounter) + uint16(c.RegisterX)
	data := c.RegisterY & (uint8(addr>>8) + 1)
	c.MemWrite(addr, data)
	return nil
}

// slo - Undocumented Opcode
//
// Shift left one bit in memory, then OR accumulator with memory.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func slo(c *CPU, mode AddressingMode) error {
	if err := asl(c, mode); err != nil {
		return err
	}
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.setAccumulator(data | c.Accumulator)
	return nil
}

// sre - Undocumented Opcode
//
// Shift right one bit in memory, then EOR accumulator with memory.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func sre(c *CPU, mode AddressingMode) error {
	if err := lsr(c, mode); err != nil {
		return err
	}
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.setAccumulator(data ^ c.Accumulator)
	return nil
}

// sta - Store Accumulator
//
// Stores the contents of the accumulator into memory.
//
// See [STA Instruction Reference].
//
// [STA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#STA
func sta(c *CPU, mode AddressingMode) error {
	addr, pageCrossed := c.getOperandAddress(mode)
	if pageCrossed {
		defer func() {
			c.Bus.Tick(1)
		}()
	}
	c.MemWrite(addr, c.Accumulator)
	return nil
}

// stx - Store X Register
//
// Stores the contents of the X register into memory.
//
// See [STX Instruction Reference].
//
// [STX Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#STX
func stx(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	c.MemWrite(addr, c.RegisterX)
	return nil
}

// sty - Store Y Register
//
// Stores the contents of the Y register into memory.
//
// See [STY Instruction Reference].
//
// [STY Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#STY
func sty(c *CPU, mode AddressingMode) error {
	addr, _ := c.getOperandAddress(mode)
	c.MemWrite(addr, c.RegisterY)
	return nil
}

// tas - Undocumented Opcode
//
// AND X register with accumulator and store result in stack pointer, then
// AND stack pointer with the high byte of the target address of the
// argument + 1. Store result in memory.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func tas(c *CPU, mode AddressingMode) error {
	data := c.Accumulator & c.RegisterX
	c.StackPointer = data
	addr := c.MemRead16(c.ProgramCounter) + uint16(c.RegisterY)
	data = (uint8(addr>>8) + 1) & c.StackPointer
	c.MemWrite(addr, data)
	return nil
}

// tax - Transfer Accumulator to X
//
// Copies the current contents of the accumulator into the X register
// and sets the zero and negative flags as appropriate.
//
// See [TAX Instruction Reference].
//
// [TAX Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TAX
func tax(c *CPU, mode AddressingMode) error {
	c.RegisterX = c.Accumulator
	c.updateZeroAndNegFlags(c.RegisterX)
	return nil
}

// tay - Transfer Accumulator to Y
//
// Copies the current contents of the accumulator into the Y register
// and sets the zero and negative flags as appropriate.
//
// See [TAY Instruction Reference].
//
// [TAY Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TAY
func tay(c *CPU, mode AddressingMode) error {
	c.RegisterY = c.Accumulator
	c.updateZeroAndNegFlags(c.RegisterY)
	return nil
}

// tsx - Transfer Stack Pointer to X
//
// Copies the current contents of the stack register into the X register
// and sets the zero and negative flags as appropriate.
//
// See [TSX Instruction Reference].
//
// [TSX Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TSX
func tsx(c *CPU, mode AddressingMode) error {
	c.RegisterX = c.StackPointer
	c.updateZeroAndNegFlags(c.RegisterX)
	return nil
}

// txa - Transfer X to Accumulator
//
// Copies the current contents of the X register into the accumulator
// and sets the zero and negative flags as appropriate.
//
// See [TXA Instruction Reference].
//
// [TXA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TXA
func txa(c *CPU, mode AddressingMode) error {
	c.setAccumulator(c.RegisterX)
	return nil
}

// txs - Transfer X to Stack Pointer
//
// Copies the current contents of the X register into the stack register.
//
// See [TXS Instruction Reference].
//
// [TXS Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TXS
func txs(c *CPU, mode AddressingMode) error {
	c.StackPointer = c.RegisterX
	return nil
}

// tya - Transfer Y to Accumulator
//
// Copies the current contents of the Y register into the accumulator
// and sets the zero and negative flags as appropriate.
//
// See [TYA Instruction Reference].
//
// [TYA Instruction Reference]: https://nesdev.org/obelisk-6502-guide/reference.html#TYA
func tya(c *CPU, mode AddressingMode) error {
	c.setAccumulator(c.RegisterY)
	return nil
}

// xaa - Undocumented Opcode
//
// Exact operation unknown. Read the referenced documents for more
// information and observations.
//
// See [6502 Undocumented Opcodes]
//
// [6502 Undocuments Opcodes]: https://www.nesdev.org/undocumented_opcodes.txt
func xaa(c *CPU, mode AddressingMode) error {
	c.Accumulator = c.RegisterX
	c.updateZeroAndNegFlags(c.Accumulator)
	addr, _ := c.getOperandAddress(mode)
	data := c.MemRead(addr)
	c.setAccumulator(data & c.Accumulator)
	return nil
}
