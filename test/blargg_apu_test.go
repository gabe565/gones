package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_blarggAPU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rom        string
		wantStatus status
		want       string
	}{
		{"len ctr", "roms/apu_test/rom_singles/1-len_ctr.nes", 0, "1-len_ctr\n\nPassed"},
		{"IRQ flag", "roms/apu_test/rom_singles/3-irq_flag.nes", 0, "3-irq_flag\n\nPassed"},
		{"DMC basics", "roms/apu_test/rom_singles/7-dmc_basics.nes", 0, "7-dmc_basics\n\nPassed"},
		{"reset clears $4015", "roms/apu_reset/4015_cleared.nes", 0, "4015_cleared\n\nPassed"},
		{"reset clears IRQ", "roms/apu_reset/irq_flag_cleared.nes", 0, "irq_flag_cleared\n\nPassed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rom, err := roms.Open(tt.rom)
			require.NoError(t, err)

			test, err := newBlarggTest(rom, msgTypeSRAM)
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.EqualValues(t, tt.wantStatus, getBlarggStatus(test))
			assert.EqualValues(t, tt.want, getBlarggMessage(test, msgTypeSRAM))
		})
	}
}
