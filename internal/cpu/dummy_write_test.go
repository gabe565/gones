package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RMW_DoubleWrite(t *testing.T) {
	// INC $2000 (EE 00 20)
	// Memory at $2000 = 0x10.
	// Expected:
	// Read $2000 -> 0x10.
	// Write $2000 <- 0x10 (Dummy Write).
	// Write $2000 <- 0x11 (Final Write).

	cpu, bus := stubCPUMockBus([]byte{0xEE, 0x00, 0x20})
	bus.Memory[0x2000] = 0x10

	cpu.Step()

	// Filter write log to $2000
	var writes []byte
	for _, w := range bus.WriteLog {
		if w.Addr == 0x2000 {
			writes = append(writes, w.Data)
		}
	}

	assert.Equal(t, []byte{0x10, 0x11}, writes, "INC should write original then modified value")
}

func Test_ASL_DoubleWrite(t *testing.T) {
	// ASL $2000 (0E 00 20)
	// Memory at $2000 = 0x01.
	// Expected:
	// Read $2000 -> 0x01.
	// Write $2000 <- 0x01.
	// Write $2000 <- 0x02.

	cpu, bus := stubCPUMockBus([]byte{0x0E, 0x00, 0x20})
	bus.Memory[0x2000] = 0x01

	cpu.Step()

	var writes []byte
	for _, w := range bus.WriteLog {
		if w.Addr == 0x2000 {
			writes = append(writes, w.Data)
		}
	}

	assert.Equal(t, []byte{0x01, 0x02}, writes, "ASL should write original then modified value")
}
