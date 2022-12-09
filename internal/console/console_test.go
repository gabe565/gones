package console

import (
	"bufio"
	_ "embed"
	"github.com/gabe565/gones/internal/apu"
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
	c.PPU = ppu.New(mapper)
	if !assert.NoError(t, err) {
		return
	}
	c.APU = apu.New()
	c.Bus = bus.New(mapper, c.PPU, c.APU)
	c.CPU = cpu.New(c.Bus)
	c.APU.SetCpu(c.CPU)

	c.CPU.Reset()
	c.CPU.ProgramCounter = 0xC000

	scanner := bufio.NewScanner(strings.NewReader(nestest.Log))
	var checkedLines uint
	for scanner.Scan() {
		checkedLines += 1

		trace := c.CPU.Trace()

		switch c.CPU.ProgramCounter {
		//TODO: Remove this after APU is supported
		case 0xC68B, 0xC690, 0xC695, 0xC69A, 0xC69F:
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
		if c.CPU.Status.Break {
			break
		}
	}
	if !assert.NoError(t, scanner.Err()) {
		return
	}

	assert.EqualValues(t, strings.Count(nestest.Log, "\n"), checkedLines)
	assert.EqualValues(t, 0, c.CPU.ReadMem(0x2))
	assert.EqualValues(t, 0, c.CPU.ReadMem(0x3))
}
