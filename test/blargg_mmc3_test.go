package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_blarggMMC3(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rom        string
		wantStatus status
		want       string
	}{
		{"IRQ clocking", "roms/mmc3_irq_tests/1.Clocking.nes", -1, "MMC3 IRQ COUNTER\nPASSED"},
		{"IRQ details", "roms/mmc3_irq_tests/2.Details.nes", -1, "MMC3 IRQ COUNTER DETAILS\nPASSED"},
		{"IRQ A12 clocking", "roms/mmc3_irq_tests/3.A12_clocking.nes", -1, "MMC3 IRQ COUNTER A12\nPASSED"},
		{"IRQ scanline timing", "roms/mmc3_irq_tests/4.Scanline_timing.nes", -1, "MMC3 IRQ TIMING\nFAILED #3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rom, err := roms.Open(tt.rom)
			require.NoError(t, err)

			test, err := newBlarggTest(rom, msgTypePPUVRAM)
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.EqualValues(t, tt.wantStatus, getBlarggStatus(test))
			assert.EqualValues(t, tt.want, getBlarggMessage(test, msgTypePPUVRAM))
		})
	}
}
