package cpu

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gabe565/gones/internal/interrupt"
	"github.com/gabe565/gones/internal/memory"
	"github.com/gabe565/gones/internal/util"
)

func New(b memory.ReadSafeWrite) *CPU {
	cpu := CPU{
		StackPointer: byte(StackAddr - 3),
		Status:       Status{InterruptDisable: true},
		bus:          b,
		Cycles:       7,
	}
	cpu.ProgramCounter = cpu.ReadMem16(interrupt.ResetVector)
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

	Cycles uint

	NMIPending bool `msgpack:"alias:NmiPending"`
	IRQPending bool `msgpack:"alias:IrqPending"`

	Stall uint16

	StepErr error `msgpack:"-"`
}

// Reset resets the CPU and sets ProgramCounter to the value of the [Reset] Vector.
func (c *CPU) Reset() {
	c.StackPointer -= 3
	sei(c, 0)
	c.ProgramCounter = c.ReadMem16(interrupt.ResetVector)
}

func (c *CPU) nmi() {
	c.stackPush16(c.ProgramCounter)
	php(c, 0)
	sei(c, 0)
	c.Cycles += 7
	c.ProgramCounter = c.ReadMem16(interrupt.NMIVector)
	c.NMIPending = false
}

func (c *CPU) irq() {
	c.stackPush16(c.ProgramCounter)
	php(c, 0)
	sei(c, 0)
	c.Cycles += 7
	c.ProgramCounter = c.ReadMem16(interrupt.IRQVector)
	c.IRQPending = false
}

// ErrUnsupportedOpcode indicates an unsupported opcode was evaluated.
var ErrUnsupportedOpcode = errors.New("unsupported opcode")

// Step steps through the next instruction
func (c *CPU) Step() uint {
	if c.Stall > 0 {
		c.Stall--
		c.Cycles++
		return 1
	}

	cycles := c.Cycles

	if c.NMIPending {
		c.nmi()
	} else if c.IRQPending && !c.Status.InterruptDisable {
		c.irq()
	}

	code := c.ReadMem(c.ProgramCounter)
	c.ProgramCounter++
	prevPC := c.ProgramCounter

	op := OpCodes[code]
	if op == nil {
		c.StepErr = fmt.Errorf("%w: $%02X", ErrUnsupportedOpcode, code)
		slog.Error("Failed to step CPU", "error", ErrUnsupportedOpcode, "code", util.EncodeHexVal(code))
		return 1
	}

	op.Exec(c, op.Mode)

	c.Cycles += uint(op.Cycles)

	if prevPC == c.ProgramCounter {
		c.ProgramCounter += uint16(op.Len - 1)
	}

	return c.Cycles - cycles
}

func (c *CPU) AddStall(stall uint16) {
	c.Stall += stall
}

func (c *CPU) AddNMI() {
	c.NMIPending = true
}

func (c *CPU) GetCycles() uint {
	return c.Cycles
}
