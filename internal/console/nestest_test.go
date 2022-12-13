package console

import (
	"bufio"
	_ "embed"
	"github.com/gabe565/gones/internal/test_roms/nestest"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_nestest(t *testing.T) {
	c, err := stubConsole(strings.NewReader(nestest.ROM))
	if !assert.NoError(t, err) {
		return
	}

	c.CPU.ProgramCounter = 0xC000

	scanner := bufio.NewScanner(strings.NewReader(nestest.Log))
	var checkedLines uint
	for scanner.Scan() {
		checkedLines += 1

		switch c.CPU.ProgramCounter {
		//TODO: Remove this after APU is supported
		case 0xC68B, 0xC690, 0xC695, 0xC69A, 0xC69F:
		default:
			trace := c.CPU.Trace()
			expected := scanner.Text()

			//TODO: Remove this after adding PPU and CYC to trace
			if len(expected) > 73 {
				expected = expected[:73]
			}

			if !assert.EqualValues(t, expected, trace) {
				return
			}
		}

		if err := c.Step(); !assert.NoError(t, err) {
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
}
