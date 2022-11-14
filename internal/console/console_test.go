package console

import (
	"bufio"
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

	var c Console

	c.PPU = ppu.New(cart)
	c.Bus = bus.New(cart, c.PPU)
	c.CPU = cpu.New(c.Bus)

	c.CPU.Reset()
	c.CPU.ProgramCounter = 0xC000

	scanner := bufio.NewScanner(strings.NewReader(nestest.Log))

	for {
		trace := c.CPU.Trace()

		scanner.Scan()

		switch c.CPU.ProgramCounter {
		//TODO: Remove this after APU is supported
		case 0xC68B, 0xC690, 0xC695, 0xC69A, 0xC69F:
			return
		//TODO: Check if these should be ignored.
		// They get logged by our trace, but seem to be missing from nestest.log
		// 0x1 is *ISB
		// 0x4 is final BRK
		case 0x1, 0x4:
			return
		}

		expected := scanner.Text()

		//TODO: Remove this after adding PPU and CYC to trace
		if len(expected) > 73 {
			expected = expected[:73]
		}

		assert.EqualValues(t, expected, trace)

		if _, err := c.CPU.Step(); err != nil {
			if assert.ErrorIs(t, err, cpu.ErrBrk) {
				break
			}
			return
		}
	}

	assert.EqualValues(t, 0, c.CPU.MemRead(0x2))
	assert.EqualValues(t, 0, c.CPU.MemRead(0x3))
}
