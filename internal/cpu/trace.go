package cpu

import (
	"fmt"
	"strings"
)

func (c *CPU) Trace() string {
	code := c.MemRead(c.ProgramCounter)
	op := OpCodeMap[code]

	begin := c.ProgramCounter
	hexDump := []uint16{uint16(code)}
	var valAddr uint16
	var val byte
	switch op.Mode {
	case Implicit, Indirect, Immediate, Accumulator:
		//
	default:
		valAddr = c.getAbsoluteAddress(op.Mode, begin+1)
		val = c.MemRead(valAddr)
	}

	var trace string
	switch op.Len {
	case 1:
		if op.Mode == Accumulator {
			trace += "A "
		}
	case 2:
		addr := c.MemRead(begin + 1)
		hexDump = append(hexDump, uint16(addr))

		switch op.Mode {
		case Immediate:
			trace += fmt.Sprintf("#$%02X", addr)
		case ZeroPage:
			trace += fmt.Sprintf("$%02X = %02X", valAddr, val)
		case ZeroPageX:
			trace += fmt.Sprintf("$%02X,X @ %02X = %02X", addr, valAddr, val)
		case ZeroPageY:
			trace += fmt.Sprintf("$%02X,Y @ %02X = %02X", addr, valAddr, val)
		case IndirectX:
			trace += fmt.Sprintf("($%02X,X) @ %02X = %04X = %02X", addr, addr+c.RegisterX, valAddr, val)
		case IndirectY:
			trace += fmt.Sprintf("($%02X),Y = %04X @ %04X = %02X", addr, valAddr-uint16(c.RegisterY), valAddr, val)
		case Implicit, Relative:
			// assuming local jumps: BNE, BVS, etc
			addr := uint16(addr) + begin + 2
			trace += fmt.Sprintf("$%04X", addr)
		default:
			panic(fmt.Sprintf("unexpected addressing mode %s has len 2. code %02X", op.Mode, op.Code))
		}
	case 3:
		addrLo := c.MemRead(begin + 1)
		hexDump = append(hexDump, uint16(addrLo))
		addrHi := c.MemRead(begin + 2)
		hexDump = append(hexDump, uint16(addrHi))

		addr := c.MemRead16(begin + 1)

		switch op.Mode {
		case Indirect:
			var indirect uint16
			if addr&0x00FF == 0x00FF {
				lo := c.MemRead(addr)
				hi := c.MemRead(addr & 0xFF00)
				indirect = uint16(hi)<<8 | uint16(lo)
			} else {
				indirect = c.MemRead16(addr)
			}
			trace += fmt.Sprintf("($%04X) = %04X", addr, indirect)
		case Implicit, Relative:
			trace += fmt.Sprintf("$%04X", addr)
		case Absolute:
			if op.Mnemonic == "JMP" {
				trace += fmt.Sprintf("$%04X", valAddr)
			} else {
				trace += fmt.Sprintf("$%04X = %02X", valAddr, val)
			}
		case AbsoluteX:
			trace += fmt.Sprintf("$%04X,X @ %04X = %02X", addr, valAddr, val)
		case AbsoluteY:
			trace += fmt.Sprintf("$%04X,Y @ %04X = %02X", addr, valAddr, val)
		default:
			panic(fmt.Sprintf("unexpected addressing mode %s has len 3. code %02X", op.Mode, op.Code))
		}
	}

	var hexStr string
	for _, v := range hexDump {
		hexStr += fmt.Sprintf("%02X ", v)
	}
	hexStr = strings.TrimSpace(hexStr)
	asmStr := fmt.Sprintf("%04X  %-8s %4s %s", begin, hexStr, op.Mnemonic, trace)

	final := fmt.Sprintf(
		"%-47s A:%02X X:%02X Y:%02X P:%02X SP:%02X",
		asmStr,
		c.Accumulator,
		c.RegisterX,
		c.RegisterY,
		c.Status,
		c.StackPointer,
	)
	return final
}
