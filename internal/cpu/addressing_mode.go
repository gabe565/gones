package cpu

import (
	"log/slog"
	"os"
)

//go:generate go tool stringer -type AddressingMode

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
func (c *CPU) getAbsoluteAddress(mode AddressingMode, addr uint16) (uint16, bool) {
	switch mode {
	case ZeroPage:
		return uint16(c.ReadMem(addr)), false
	case Absolute:
		return c.ReadMem16(addr), false
	case ZeroPageX:
		pos := c.ReadMem(addr)
		c.tick() // internal: add X to address
		return uint16(pos + c.RegisterX), false
	case ZeroPageY:
		pos := c.ReadMem(addr)
		c.tick() // internal: add Y to address
		return uint16(pos + c.RegisterY), false
	case AbsoluteX:
		base := c.ReadMem16(addr)
		addr := base + uint16(c.RegisterX)
		if crossedPage(base, addr) {
			c.ReadMem(addr - 0x100)
			return addr, true
		}
		return addr, false
	case AbsoluteY:
		base := c.ReadMem16(addr)
		addr := base + uint16(c.RegisterY)
		if crossedPage(base, addr) {
			c.ReadMem(addr - 0x100)
			return addr, true
		}
		return addr, false
	case Indirect:
		base := c.ReadMem16(addr)
		return c.ReadMem16Bug(base), false
	case IndirectX:
		base := c.ReadMem(addr)
		c.tick() // internal: add X to pointer
		ptr := base + c.RegisterX
		lo := c.ReadMem(uint16(ptr))
		hi := c.ReadMem(uint16(ptr + 1))
		return uint16(hi)<<8 | uint16(lo), false
	case IndirectY:
		base := c.ReadMem(addr)

		lo := c.ReadMem(uint16(base))
		hi := c.ReadMem(uint16(base + 1))
		derefBase := uint16(hi)<<8 | uint16(lo)
		addr := derefBase + uint16(c.RegisterY)
		if crossedPage(derefBase, addr) {
			c.ReadMem(addr - 0x100)
			return addr, true
		}
		return addr, false
	default:
		slog.Error("unsupported mode", "mode", mode)
		os.Exit(1)
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

// getAbsoluteAddressSafe resolves an address without ticking. Used for tracing.
func (c *CPU) getAbsoluteAddressSafe(mode AddressingMode, addr uint16) (uint16, bool) {
	read := c.readMemSafe
	read16 := c.readMem16Safe
	switch mode {
	case ZeroPage:
		return uint16(read(addr)), false
	case Absolute:
		return read16(addr), false
	case ZeroPageX:
		return uint16(read(addr) + c.RegisterX), false
	case ZeroPageY:
		return uint16(read(addr) + c.RegisterY), false
	case AbsoluteX:
		base := read16(addr)
		result := base + uint16(c.RegisterX)
		return result, crossedPage(base, result)
	case AbsoluteY:
		base := read16(addr)
		result := base + uint16(c.RegisterY)
		return result, crossedPage(base, result)
	case Indirect:
		base := read16(addr)
		lo := uint16(read(base))
		hi := uint16(read(base&0xFF00 | (base+1)&0x00FF))
		return hi<<8 | lo, false
	case IndirectX:
		ptr := read(addr) + c.RegisterX
		lo := read(uint16(ptr))
		hi := read(uint16(ptr + 1))
		return uint16(hi)<<8 | uint16(lo), false
	case IndirectY:
		base := read(addr)
		lo := read(uint16(base))
		hi := read(uint16(base + 1))
		derefBase := uint16(hi)<<8 | uint16(lo)
		result := derefBase + uint16(c.RegisterY)
		return result, crossedPage(derefBase, result)
	default:
		return 0, false
	}
}
