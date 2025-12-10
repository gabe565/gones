package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_instruction_shx(t *testing.T) {
	t.Run("No page cross bug", func(t *testing.T) {
		// SHX $0200, Y (9E 00 02)
		// Init: X=$FF, Y=$10. Base=$0200.
		// Expected: Write 03 to $0210.
		cpu := stubCPU([]byte{0xA2, 0xFF, 0xA0, 0x10, 0x9E, 0x00, 0x02, 0x00})
		for {
			cpu.Step()
			require.NoError(t, cpu.StepErr)
			if cpu.Status.Break {
				break
			}
		}
		assert.EqualValues(t, 0x03, cpu.ReadMem(0x0210), "Expected value at $0210")
	})

	t.Run("Page cross behavior", func(t *testing.T) {
		// SHX $0290, Y (9E 90 02)
		// Init: X=$FF, Y=$80. Base=$0290. Target=$0310. Page cross.
		// Expected: Value=03 (FF & (02+1)). Write 03 to $0310.
		cpu := stubCPU([]byte{0xA2, 0xFF, 0xA0, 0x80, 0x9E, 0x90, 0x02, 0x00})
		for {
			cpu.Step()
			require.NoError(t, cpu.StepErr)
			if cpu.Status.Break {
				break
			}
		}
		assert.EqualValues(t, 0x03, cpu.ReadMem(0x0310), "Expected value 03 at 0310")
		assert.EqualValues(t, 0x00, cpu.ReadMem(0x0410), "Should not write to 0410 (incorrect EffectiveHi calc)")
	})

	t.Run("Page cross corruption", func(t *testing.T) {
		// SHX $0290, Y (9E 90 02)
		// Init: X=$00, Y=$80. Base=$0290. Target=$0310. Page cross.
		// Expected: Value=00 (00 & (02+1)). Glitch Addr=$0010. Write 00 to $0010.
		cpu := stubCPU([]byte{0xA2, 0x00, 0xA0, 0x80, 0x9E, 0x90, 0x02, 0x00})
		cpu.WriteMem(0x0010, 0xFF) // Init to verify write
		cpu.WriteMem(0x0310, 0xFF) // Init to verify no write
		for {
			cpu.Step()
			require.NoError(t, cpu.StepErr)
			if cpu.Status.Break {
				break
			}
		}
		assert.EqualValues(t, 0x00, cpu.ReadMem(0x0010), "Expected value 00 at 0010 (glitch addr)")
		assert.EqualValues(t, 0xFF, cpu.ReadMem(0x0310), "Should not write to Target 0310")
	})
}
