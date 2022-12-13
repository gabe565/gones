package cpu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCPU_TraceFormat(t *testing.T) {
	t.Parallel()

	c := stubCpu([]byte{0xA2, 0x01, 0xca, 0x88, 0x00})
	traces := make([]string, 0)
	c.Accumulator = 1
	c.RegisterX = 2
	c.RegisterY = 3
	for {
		traces = append(traces, c.Trace())
		if _, err := c.Step(); !assert.NoError(t, err) {
			return
		}
		if c.Status.Break {
			break
		}
	}

	assert.EqualValues(
		t,
		"8600  A2 01     LDX #$01                        A:01 X:02 Y:03 P:24 SP:FD",
		traces[0],
	)
	assert.EqualValues(
		t,
		"8602  CA        DEX                             A:01 X:01 Y:03 P:24 SP:FD",
		traces[1],
	)
	assert.EqualValues(
		t,
		"8603  88        DEY                             A:01 X:00 Y:03 P:26 SP:FD",
		traces[2],
	)
}

func TestCPU_Trace_MemAccess(t *testing.T) {
	t.Parallel()

	c := stubCpu([]byte{0x11, 0x33})
	traces := make([]string, 0)
	c.WriteMem(0x33, 0)
	c.WriteMem(0x34, 4)
	c.WriteMem(0x400, 0xAA)
	for {
		traces = append(traces, c.Trace())
		if _, err := c.Step(); !assert.NoError(t, err) {
			return
		}
		if c.Status.Break {
			break
		}
	}

	assert.EqualValues(
		t,
		"8600  11 33     ORA ($33),Y = 0400 @ 0400 = AA  A:00 X:00 Y:00 P:24 SP:FD",
		traces[0],
	)
}
