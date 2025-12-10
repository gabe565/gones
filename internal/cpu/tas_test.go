package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_instruction_tas(t *testing.T) {
	t.Run("No Page Cross", func(t *testing.T) {
		// TAS $0200, Y. Y=0. (9B 00 02)
		// Init: A=$7F, X=$F0.
		// Expected: SP=70 (A&X). Value=00 (SP & (02+1)). Write 00 to $0200.
		cpu := stubCPU([]byte{0xA9, 0x7F, 0xA2, 0xF0, 0xA0, 0x00, 0x9B, 0x00, 0x02, 0x00})
		cpu.StackPointer = 0xFF // Set initial SP to something different
		for {
			cpu.Step()
			require.NoError(t, cpu.StepErr)
			if cpu.Status.Break {
				break
			}
		}
		// SP was 70, then BRK pushed 3 bytes, so 70 - 3 = 6D.
		assert.EqualValues(t, 0x6D, cpu.StackPointer, "SP should be updated correctly")
		assert.EqualValues(t, 0x00, cpu.ReadMem(0x0200), "Expected value at 0200")
	})

	t.Run("Page Cross - Value Calculation", func(t *testing.T) {
		// TAS $0290, Y. Y=$80. (9B 90 02)
		// Init: A=$FF, X=$FF. Base=$0290. Target=$0310. Page cross.
		// Expected: SP=FF. Value=03 (SP & (02+1)). Write 03 to $0310.
		cpu := stubCPU([]byte{0xA9, 0xFF, 0xA2, 0xFF, 0xA0, 0x80, 0x9B, 0x90, 0x02, 0x00})
		for {
			cpu.Step()
			require.NoError(t, cpu.StepErr)
			if cpu.Status.Break {
				break
			}
		}
		assert.EqualValues(t, 0x03, cpu.ReadMem(0x0310), "Expected value 03 at 0310 (using BaseHi)")
		assert.EqualValues(t, 0x00, cpu.ReadMem(0x0410), "Should not write to 0410 (incorrect EffectiveHi calc)")
	})

	t.Run("Page Cross - Address Corruption", func(t *testing.T) {
		// TAS $0290, Y. Y=$80. (9B 90 02)
		// Init: A=$00, X=$FF. Base=$0290. Target=$0310. Page cross.
		// Expected: SP=00. Value=00 (SP & (02+1)). Glitch Addr=$0010. Write 00 to $0010.
		cpu := stubCPU([]byte{0xA9, 0x00, 0xA2, 0xFF, 0xA0, 0x80, 0x9B, 0x90, 0x02, 0x00})
		cpu.WriteMem(0x0010, 0xFF) // Init to verify write
		cpu.WriteMem(0x0310, 0xFF) // Init to verify no write
		for {
			cpu.Step()
			require.NoError(t, cpu.StepErr)
			if cpu.Status.Break {
				break
			}
		}
		assert.EqualValues(t, 0x00, cpu.ReadMem(0x0010), "Should write 00 to 0010 (glitch addr)")
		assert.EqualValues(t, 0xFF, cpu.ReadMem(0x0310), "Should not write to Target 0310")
	})
}
