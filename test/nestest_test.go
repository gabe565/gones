package test

import (
	"bufio"
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed roms/other/nestest.nes
var nestest string

//go:embed roms/other/nestest.log
var nestestLog string

func Test_nestest(t *testing.T) {
	t.Parallel()

	c, err := stubConsole(strings.NewReader(nestest))
	require.NoError(t, err)

	c.CPU.ProgramCounter = 0xC000

	scanner := bufio.NewScanner(strings.NewReader(nestestLog))
	var checkedLines uint
	for scanner.Scan() {
		checkedLines++
		actual := c.Trace()
		expected := scanner.Text()

		require.EqualValues(t, expected, actual)

		c.Step(true)
		require.NoError(t, c.CPU.StepErr)
		if c.CPU.Status.Break {
			break
		}
	}
	require.NoError(t, scanner.Err())

	assert.EqualValues(t, strings.Count(nestestLog, "\n"), checkedLines)
	assert.EqualValues(t, 0, c.Bus.ReadMem(2), "See https://github.com/christopherpow/roms/blob/master/other/nestest.txt#L87 for failure code meaning")
	assert.EqualValues(t, 0, c.Bus.ReadMem(3), "See https://github.com/christopherpow/roms/blob/master/other/nestest.txt#L366 for failure code meaning")
}
