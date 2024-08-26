package test

import (
	"bufio"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_nestest(t *testing.T) {
	t.Parallel()

	nestest, err := roms.Open("roms/other/nestest.nes")
	require.NoError(t, err)

	nestestLog, err := roms.Open("roms/other/nestest.log")
	require.NoError(t, err)

	c, err := stubConsole(nestest)
	require.NoError(t, err)

	c.CPU.ProgramCounter = 0xC000

	scanner := bufio.NewScanner(nestestLog)
	var totalLines, checkedLines uint
	for scanner.Scan() {
		totalLines++
		if !c.CPU.Status.Break {
			checkedLines++
			actual := c.Trace()
			want := scanner.Text()

			require.EqualValues(t, want, actual)

			c.Step(true)
			require.NoError(t, c.CPU.StepErr)
		}
	}
	require.NoError(t, scanner.Err())

	assert.EqualValues(t, totalLines, checkedLines)
	assert.EqualValues(t, 0, c.Bus.ReadMem(2), "See https://github.com/christopherpow/roms/blob/master/other/nestest.txt#L87 for failure code meaning")
	assert.EqualValues(t, 0, c.Bus.ReadMem(3), "See https://github.com/christopherpow/roms/blob/master/other/nestest.txt#L366 for failure code meaning")
}
