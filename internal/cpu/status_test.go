package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Status_BFlag(t *testing.T) {
	// Test that the B flag (Bit 4) is NOT set when pushing status during IRQ.
	t.Run("IRQ", func(t *testing.T) {
		cpu := stubCPU([]byte{0x4C, 0x00, 0x02}) // JMP $0200
		cpu.Status.InterruptDisable = false

		cpu.IRQPending = true

		cpu.Step()

		statusByte := cpu.ReadMem(0x0100 + uint16(cpu.StackPointer) + 1)
		assert.Equal(t, uint8(0), statusByte&0x10, "B flag should not be set on IRQ")
		assert.Equal(t, uint8(0x20), statusByte&0x20, "Unused bit 5 should be set")
	})

	// Test that the B flag (Bit 4) IS set when pushing status via PHP.
	t.Run("PHP", func(t *testing.T) {
		cpu := stubCPU([]byte{0x08, 0x00})

		cpu.Step()

		statusByte := cpu.ReadMem(0x0100 + uint16(cpu.StackPointer) + 1)
		assert.Equal(t, uint8(0x10), statusByte&0x10, "B flag should be set on PHP")
	})

	// Tests that the B flag (Bit 4) IS set when pushing status via BRK.
	t.Run("BRK", func(t *testing.T) {
		cpu := stubCPU([]byte{0x00, 0x00})

		cpu.Step()

		statusByte := cpu.ReadMem(0x0100 + uint16(cpu.StackPointer) + 1)
		assert.Equal(t, uint8(0x10), statusByte&0x10, "B flag should be set on BRK")
	})
}
