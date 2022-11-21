package cpu

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

var skipRead = []uint16{
	0x2001, 0x2002, 0x2003, 0x2004, 0x2005, 0x2006, 0x2007, 0x4016, 0x4017,
}

func (c *CPU) Trace() string {
	code := c.ReadMem(c.ProgramCounter)
	op := OpCodes[code]
	if op.Exec == nil {
		return ""
	}

	begin := c.ProgramCounter
	hexDump := []uint16{uint16(code)}
	var valAddr uint16
	var val byte
	switch op.Mode {
	case Implied, Immediate, Accumulator, Relative:
		//
	default:
		valAddr, _ = c.getAbsoluteAddress(op.Mode, begin+1)

		var skip bool
		for _, skipAddr := range skipRead {
			if valAddr == skipAddr {
				skip = true
				break
			}
		}
		if !skip {
			val = c.ReadMem(valAddr)
		}
	}

	var trace string
	switch op.Len {
	case 1:
		if op.Mode == Accumulator {
			trace += "A "
		}
	case 2:
		addr := c.ReadMem(begin + 1)
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
		case Implied, Relative:
			// assuming local jumps: BNE, BVS, etc
			addr := uint16(int8(addr)) + begin + 2
			trace += fmt.Sprintf("$%04X", addr)
		default:
			log.Panicf("unexpected addressing mode %s has len 2. code %02X", op.Mode, op.Code)
		}
	case 3:
		addrLo := c.ReadMem(begin + 1)
		hexDump = append(hexDump, uint16(addrLo))
		addrHi := c.ReadMem(begin + 2)
		hexDump = append(hexDump, uint16(addrHi))

		addr := c.ReadMem16(begin + 1)

		switch op.Mode {
		case Indirect:
			trace += fmt.Sprintf("($%04X) = %04X", addr, valAddr)
		case Implied, Relative:
			trace += fmt.Sprintf("$%04X", addr)
		case Absolute:
			if op.Code == 0x4C { // JMP
				trace += fmt.Sprintf("$%04X", valAddr)
			} else {
				trace += fmt.Sprintf("$%04X = %02X", valAddr, val)
			}
		case AbsoluteX:
			trace += fmt.Sprintf("$%04X,X @ %04X = %02X", addr, valAddr, val)
		case AbsoluteY:
			trace += fmt.Sprintf("$%04X,Y @ %04X = %02X", addr, valAddr, val)
		default:
			log.Panicf("unexpected addressing mode %s has len 3. code %02X", op.Mode, op.Code)
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
