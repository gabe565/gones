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
	Implied AddressingMode = iota
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
func (c *CPU) getAbsoluteAddress(mode AddressingMode, addr uint16) (uint16, bool) {
	switch mode {
	case ZeroPage:
		return uint16(c.MemRead(addr)), false
	case Absolute:
		return c.MemRead16(addr), false
	case ZeroPageX:
		pos := c.MemRead(addr)
		return uint16(pos + c.RegisterX), false
	case ZeroPageY:
		pos := c.MemRead(addr)
		return uint16(pos + c.RegisterY), false
	case AbsoluteX:
		base := c.MemRead16(addr)
		addr := base + uint16(c.RegisterX)
		return addr, crossedPage(base, addr)
	case AbsoluteY:
		base := c.MemRead16(addr)
		addr := base + uint16(c.RegisterY)
		return addr, crossedPage(base, addr)
	case IndirectX:
		base := c.MemRead(addr)

		ptr := base + c.RegisterX
		lo := c.MemRead(uint16(ptr))
		hi := c.MemRead(uint16(ptr + 1))
		return uint16(hi)<<8 | uint16(lo), false
	case IndirectY:
		base := c.MemRead(addr)

		lo := c.MemRead(uint16(base))
		hi := c.MemRead(uint16(byte(base) + 1))
		derefBase := uint16(hi)<<8 | uint16(lo)
		addr := derefBase + uint16(c.RegisterY)
		return addr, crossedPage(derefBase, addr)
	default:
		log.Panicln("unsupported mode: ", mode)
		return 0, false
	}
}

// getOperandAddress gets the address based on the [AddressingMode].
//
// See [6502 Addressing Mode].
//
// [6502 Addressing Mode]: https://www.nesdev.org/obelisk-6502-guide/addressing.html
func (c *CPU) getOperandAddress(mode AddressingMode) (uint16, bool) {
	switch mode {
	case Immediate:
		return c.ProgramCounter, false
	default:
		return c.getAbsoluteAddress(mode, c.ProgramCounter)
	}
}

func crossedPage(lhs, rhs uint16) bool {
	return lhs&0xFF00 != rhs&0xFF00
}
