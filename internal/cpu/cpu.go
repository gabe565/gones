package cpu

import (
	"errors"
	"fmt"
	"iter"
	"log/slog"

	"gabe565.com/gones/internal/interrupt"
	"gabe565.com/gones/internal/log"
	"gabe565.com/gones/internal/memory"
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
	irqDelay   uint8

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
	c.stackPush(c.Status.Get() | Unused)
	sei(c, 0)
	c.Cycles += 7
	c.ProgramCounter = c.ReadMem16(interrupt.NMIVector)
	c.NMIPending = false
}

func (c *CPU) irq() {
	c.stackPush16(c.ProgramCounter)
	c.stackPush(c.Status.Get() | Unused)
	sei(c, 0)
	c.Cycles += 7
	c.ProgramCounter = c.ReadMem16(interrupt.IRQVector)
	c.IRQPending = false
}

// ErrUnsupportedOpcode indicates an unsupported opcode was evaluated.
var ErrUnsupportedOpcode = errors.New("unsupported opcode")

func (c *CPU) StepInstruction() uint {
	var total uint
	for cycles := range c.Step() {
		total += cycles
	}
	return total
}

// Step steps through the next instruction, yielding the number of CPU cycles consumed.
func (c *CPU) Step() iter.Seq[uint] {
	return func(yield func(uint) bool) {
		if c.Stall > 0 {
			c.Stall--
			c.Cycles++
			yield(1)
			return
		}

		cycles := c.Cycles

		if c.NMIPending {
			c.nmi()
			yield(c.Cycles - cycles)
			return
		} else if c.IRQPending && !c.Status.InterruptDisable {
			if c.irqDelay == 0 {
				c.irq()
				yield(c.Cycles - cycles)
				return
			}
			c.irqDelay--
		}

		code := c.ReadMem(c.ProgramCounter)
		c.ProgramCounter++
		prevPC := c.ProgramCounter

		op := opcodes[code]
		if op == nil {
			c.StepErr = fmt.Errorf("%w: $%02X", ErrUnsupportedOpcode, code)
			slog.Error("Failed to step CPU", "error", ErrUnsupportedOpcode, "code", log.HexVal(code))
			yield(1)
			return
		}

		op.Exec(c, op.Mode)

		c.Cycles += uint(op.Cycles)

		if prevPC == c.ProgramCounter {
			c.ProgramCounter += uint16(op.Len - 1)
		}

		yield(c.Cycles - cycles)
	}
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
