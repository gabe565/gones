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

	mapper, err := cartridge.NewMapper(cart)
	c.PPU = ppu.New(cart, mapper)
	if !assert.NoError(t, err) {
		return
	}
	c.Bus = bus.New(mapper, c.PPU)
	c.CPU = cpu.New(c.Bus)

	c.CPU.Reset()
	c.CPU.ProgramCounter = 0xC000

	scanner := bufio.NewScanner(strings.NewReader(nestest.Log))

	for scanner.Scan() {
		trace := c.CPU.Trace()

		switch c.CPU.ProgramCounter {
		//TODO: Remove this after APU is supported
		case 0xC68B, 0xC690, 0xC695, 0xC69A, 0xC69F:
			continue
		//TODO: Check if these should be ignored.
		// They get logged by our trace, but seem to be missing from nestest.log
		// 0x1 is *ISB
		// 0x4 is final BRK
		case 0x1, 0x4:
			continue
		}

		expected := scanner.Text()

		//TODO: Remove this after adding PPU and CYC to trace
		if len(expected) > 73 {
			expected = expected[:73]
		}

		assert.EqualValues(t, expected, trace)

		if _, err := c.CPU.Step(); !assert.NoError(t, err) {
			return
		}
		if c.CPU.Status.Has(cpu.Break) {
			break
		}
	}
	if !assert.NoError(t, scanner.Err()) {
		return
	}

	assert.EqualValues(t, 0, c.CPU.MemRead(0x2))
	assert.EqualValues(t, 0, c.CPU.MemRead(0x3))
}
