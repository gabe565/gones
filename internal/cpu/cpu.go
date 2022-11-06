package cpu

import (
	"errors"
	"fmt"
	"github.com/gabe565/gones/internal/bitflags"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/consts"
)

func New(b *bus.Bus) *CPU {
	return &CPU{
		status:       InterruptDisable | Break2,
		stackPointer: StackReset,
		bus:          b,
	}
}

// CPU implements the NES CPU.
//
// See [6502 Guide].
//
// [6502 Guide]: https://www.nesdev.org/obelisk-6502-guide/
type CPU struct {
	// programCounter Program Counter
	programCounter uint16

	// stackPointer Stack Pointer
	stackPointer byte

	// status Processor Status
	status bitflags.Flags

	// accumulator Register A
	accumulator byte

	// registerX Register X
	registerX byte

	// registerY Register Y
	registerY byte

	// bus Main memory bus
	bus *bus.Bus

	// Callback optional callback to Run before every tick
	Callback func(c *CPU) error

	// Debug enables opcode logging
	Debug bool
}

const (
	// StackAddr is the memory address of the stack
	StackAddr = 0x100

	// StackReset is the start value for the stack pointer
	StackReset = 0xFD
)

// Reset resets the CPU and sets programCounter to the value of the [Reset] Vector.
func (c *CPU) Reset() {
	c.accumulator = 0
	c.registerX = 0
	c.status = 0
	c.stackPointer = StackReset

	c.programCounter = c.MemRead16(consts.ResetAddr)
}

// Load loads a program into PRG memory
func (c *CPU) Load(program []byte) {
	for k, v := range program {
		c.MemWrite(consts.PrgRomAddr+uint16(k), v)
	}
	c.MemWrite16(consts.ResetAddr, consts.PrgRomAddr)
}

// ErrUnsupportedOpcode indicates an unsupported opcode was evaluated.
var ErrUnsupportedOpcode = errors.New("unsupported opcode")

// Run is the main Run entrypoint.
func (c *CPU) Run() error {
	opcodes := OpCodeMap()

	for {
		if c.Callback != nil {
			if err := c.Callback(c); err != nil {
				if errors.Is(err, ErrBrk) {
					return nil
				}
				return err
			}
		}

		code := c.MemRead(c.programCounter)
		c.programCounter += 1
		prevPC := c.programCounter

		opcode, ok := opcodes[code]
		if !ok {
			return fmt.Errorf("%w: $%x", ErrUnsupportedOpcode, code)
		}

		if c.Debug {
			fmt.Println(opcode)
		}

		if err := opcode.Exec(c, opcode.Mode); err != nil {
			if errors.Is(err, ErrBrk) {
				return nil
			}
			return err
		}

		if prevPC == c.programCounter {
			c.programCounter += uint16(opcode.Len - 1)
		}
	}
}
