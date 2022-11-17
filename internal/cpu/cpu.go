package cpu

import (
	"errors"
	"fmt"
	"github.com/gabe565/gones/internal/bitflags"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/consts"
	"github.com/gabe565/gones/internal/interrupts"
)

func New(b *bus.Bus) *CPU {
	return &CPU{
		Status:       DefaultStatus,
		StackPointer: StackReset,
		bus:          b,
		Interrupt:    make(chan interrupts.Interrupt, 1),
	}
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
	Status bitflags.Flags

	// Accumulator Register A
	Accumulator byte

	// RegisterX Register X
	RegisterX byte

	// RegisterY Register Y
	RegisterY byte

	// bus Main memory bus
	bus *bus.Bus

	Cycles uint

	Interrupt chan interrupts.Interrupt
}

type Callback func(*CPU) error

const (
	// StackAddr is the memory address of the stack
	StackAddr = 0x100

	// StackReset is the start value for the stack pointer
	StackReset = 0xFD
)

// Reset resets the CPU and sets ProgramCounter to the value of the [Reset] Vector.
func (c *CPU) Reset() {
	c.ProgramCounter = c.MemRead16(consts.ResetAddr)
	c.StackPointer = StackReset
	c.Status = DefaultStatus
}

// Load loads a program into PRG memory
func (c *CPU) Load(program []byte) {
	for k, v := range program {
		c.MemWrite(consts.PrgRomAddr+uint16(k), v)
	}
	c.MemWrite16(consts.ResetAddr, consts.PrgRomAddr)
}

func (c *CPU) interrupt(interrupt interrupts.Interrupt) {
	c.stackPush16(c.ProgramCounter)
	status := c.Status
	status.Set(Break, interrupt.Mask&Break == 1)
	status.Set(Break2, interrupt.Mask&Break2 == 1)

	c.stackPush(byte(status))
	c.Status.Insert(InterruptDisable)

	c.Cycles += uint(interrupt.Cycles)
	c.ProgramCounter = c.MemRead16(interrupt.VectorAddr)
}

// ErrUnsupportedOpcode indicates an unsupported opcode was evaluated.
var ErrUnsupportedOpcode = errors.New("unsupported opcode")

// Step steps through the next instruction
func (c *CPU) Step() (uint, error) {
	cycles := c.Cycles

	if len(c.Interrupt) > 0 {
		c.interrupt(<-c.Interrupt)
	}

	code := c.MemRead(c.ProgramCounter)
	c.ProgramCounter += 1
	prevPC := c.ProgramCounter

	op, ok := OpCodeMap[code]
	if !ok {
		return 0, fmt.Errorf("%w: $%02X", ErrUnsupportedOpcode, code)
	}

	if err := op.Exec(c, op.Mode); err != nil {
		return 0, err
	}

	c.Cycles += uint(op.Cycles)

	if prevPC == c.ProgramCounter {
		c.ProgramCounter += uint16(op.Len - 1)
	}

	return c.Cycles - cycles, nil
}
