package cpu

import (
	"errors"
	"fmt"

	"github.com/gabe565/gones/internal/interrupts"
	"github.com/gabe565/gones/internal/memory"
)

func New(b memory.ReadSafeWrite) *CPU {
	cpu := CPU{
		StackPointer: byte(StackAddr - 3),
		Status:       DefaultStatus,
		bus:          b,
		cycles:       7,
	}
	cpu.ProgramCounter = cpu.ReadMem16(interrupts.ResetVector)
	return &cpu
}

// CPU implements the NES CPU.
//
// See [6502 Guide].
//
// [6502 Guide]: https://www.nesdev.org/obelisk-6502-guide/
type CPU struct {
	// ProgramCounter Program Counter
	ProgramCounter uint16

	// StackPointer Stack Pointer
	StackPointer byte

	// Status Processor Status
	Status Status

	// Accumulator Register A
	Accumulator byte

	// RegisterX Register X
	RegisterX byte

	// RegisterY Register Y
	RegisterY byte

	// bus Main memory bus
	bus memory.ReadSafeWrite

	cycles uint

	pendingNmi bool
	pendingIrq bool

	Stall uint16
}

// Reset resets the CPU and sets ProgramCounter to the value of the [Reset] Vector.
func (c *CPU) Reset() {
	c.StackPointer -= 3
	sei(c, 0)
	c.ProgramCounter = c.ReadMem16(interrupts.ResetVector)
}

func (c *CPU) nmi() {
	c.stackPush16(c.ProgramCounter)
	php(c, 0)
	sei(c, 0)
	c.cycles += 7
	c.ProgramCounter = c.ReadMem16(interrupts.NmiVector)
	c.pendingNmi = false
}

func (c *CPU) irq() {
	c.stackPush16(c.ProgramCounter)
	php(c, 0)
	sei(c, 0)
	c.cycles += 7
	c.ProgramCounter = c.ReadMem16(interrupts.IrqVector)
	c.pendingIrq = false
}

// ErrUnsupportedOpcode indicates an unsupported opcode was evaluated.
var ErrUnsupportedOpcode = errors.New("unsupported opcode")

// Step steps through the next instruction
func (c *CPU) Step() (uint, error) {
	if c.Stall > 0 {
		c.Stall -= 1
		c.cycles += 1
		return 1, nil
	}

	cycles := c.cycles

	if c.pendingNmi {
		c.nmi()
	} else if c.pendingIrq && !c.Status.InterruptDisable {
		c.irq()
	}

	code := c.ReadMem(c.ProgramCounter)
	c.ProgramCounter += 1
	prevPC := c.ProgramCounter

	op := OpCodes[code]
	if op.Exec == nil {
		return 0, fmt.Errorf("%w: $%02X", ErrUnsupportedOpcode, code)
	}

	op.Exec(c, op.Mode)

	c.cycles += uint(op.Cycles)

	if prevPC == c.ProgramCounter {
		c.ProgramCounter += uint16(op.Len - 1)
	}

	return c.cycles - cycles, nil
}

func (c *CPU) AddStall(stall uint16) {
	c.Stall += stall
}

func (c *CPU) AddNmi() {
	c.pendingNmi = true
}

func (c *CPU) AddIrq() {
	c.pendingIrq = true
}

func (c *CPU) ClearIrq() {
	c.pendingIrq = false
}

func (c *CPU) Cycles() uint {
	return c.cycles
}
