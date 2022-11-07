package cpu

import "log"

//go:generate stringer -type AddressingMode

// AddressingMode defines opcode addressing modes.
//
// See [6502 Addressing Mode].
//
// [6502 Addressing Mode]: https://www.nesdev.org/obelisk-6502-guide/addressing.html
type AddressingMode uint8

const (
	Implicit AddressingMode = iota
	Accumulator
	Immediate
	ZeroPage
	ZeroPageX
	ZeroPageY
	Relative
	Absolute
	AbsoluteX
	AbsoluteY
	Indirect
	IndirectX
	IndirectY
)

// getAbsoluteAddress gets the address for an address based on the [AddressingMode].
//
// See [6502 Addressing Mode].
//
// [6502 Addressing Mode]: https://www.nesdev.org/obelisk-6502-guide/addressing.html
func (c *CPU) getAbsoluteAddress(mode AddressingMode, addr uint16) uint16 {
	switch mode {
	case ZeroPage:
		return uint16(c.MemRead(addr))
	case Absolute:
		return c.MemRead16(addr)
	case ZeroPageX:
		pos := c.MemRead(addr)
		return uint16(pos + c.registerX)
	case ZeroPageY:
		pos := c.MemRead(addr)
		return uint16(pos + c.registerY)
	case AbsoluteX:
		pos := c.MemRead16(addr)
		return pos + uint16(c.registerX)
	case AbsoluteY:
		pos := c.MemRead16(addr)
		return pos + uint16(c.registerY)
	case IndirectX:
		base := c.MemRead(addr)

		ptr := base + c.registerX
		lo := c.MemRead(uint16(ptr))
		hi := c.MemRead(uint16(ptr + 1))
		return uint16(hi)<<8 | uint16(lo)
	case IndirectY:
		base := c.MemRead(addr)

		lo := c.MemRead(uint16(base))
		hi := c.MemRead(uint16(byte(base) + 1))
		derefBase := uint16(hi)<<8 | uint16(lo)
		return derefBase + uint16(c.registerY)
	default:
		log.Panicln("unsupported mode: ", mode)
		return 0
	}
}

// getOperandAddress gets the address based on the [AddressingMode].
//
// See [6502 Addressing Mode].
//
// [6502 Addressing Mode]: https://www.nesdev.org/obelisk-6502-guide/addressing.html
func (c *CPU) getOperandAddress(mode AddressingMode) uint16 {
	switch mode {
	case Immediate:
		return c.programCounter
	default:
		return c.getAbsoluteAddress(mode, c.programCounter)
	}
}
