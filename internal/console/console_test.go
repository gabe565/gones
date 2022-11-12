package console

import (
	"bufio"
	"context"
	_ "embed"
	"github.com/gabe565/gones/internal/bus"
	"github.com/gabe565/gones/internal/cartridge"
	"github.com/gabe565/gones/internal/cpu"
	"github.com/gabe565/gones/internal/ppu"
	"github.com/gabe565/gones/internal/test_roms/nestest"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_nestest(t *testing.T) {
	cart, err := cartridge.FromiNes(strings.NewReader(nestest.ROM))
	if !assert.NoError(t, err) {
		return
	}

	var console Console

	console.PPU = ppu.New(cart)
	console.Bus = bus.New(cart, console.PPU)
	console.CPU = cpu.New(console.Bus)

	console.CPU.Reset()
	console.CPU.ProgramCounter = 0xC000

	scanner := bufio.NewScanner(strings.NewReader(nestest.Log))
	console.CPU.Callback = func(c *cpu.CPU) error {
		trace := c.Trace()

		scanner.Scan()

		switch c.ProgramCounter {
		//TODO: Remove this after APU is supported
		case 0xC68B, 0xC690, 0xC695, 0xC69A, 0xC69F:
			return nil
		//TODO: Check if these should be ignored.
		// They get logged by our trace, but seem to be missing from nestest.log
		// 0x1 is *ISB
		// 0x4 is final BRK
		case 0x1, 0x4:
			return nil
		}

		expected := scanner.Text()

		//TODO: Remove this after adding PPU and CYC to trace
		if len(expected) > 73 {
			expected = expected[:73]
		}

		assert.EqualValues(t, expected, trace)
		return nil
	}

	if err := console.CPU.Run(context.Background()); !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 0, console.CPU.MemRead(0x2))
	assert.EqualValues(t, 0, console.CPU.MemRead(0x3))
}
