package console

import (
	"bufio"
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed nes-test-roms/other/nestest.nes
var nestest string

//go:embed nes-test-roms/other/nestest.log
var nestestLog string

func Test_nestest(t *testing.T) {
	t.Parallel()

	c, err := stubConsole(strings.NewReader(nestest))
	if !assert.NoError(t, err) {
		return
	}

	c.CPU.ProgramCounter = 0xC000

	scanner := bufio.NewScanner(strings.NewReader(nestestLog))
	var checkedLines uint
	for scanner.Scan() {
		checkedLines += 1
		actual := c.Trace()
		expected := scanner.Text()

		if !assert.EqualValues(t, expected, actual) {
			return
		}

		if err := c.Step(true); !assert.NoError(t, err) {
			return
		}
		if c.CPU.Status.Break {
			break
		}
	}
	if !assert.NoError(t, scanner.Err()) {
		return
	}

	assert.EqualValues(t, strings.Count(nestestLog, "\n"), checkedLines)
	assert.EqualValues(t, 0, c.Bus.ReadMem(2), "See https://github.com/christopherpow/nes-test-roms/blob/master/other/nestest.txt#L87 for failure code meaning")
	assert.EqualValues(t, 0, c.Bus.ReadMem(3), "See https://github.com/christopherpow/nes-test-roms/blob/master/other/nestest.txt#L366 for failure code meaning")
}
