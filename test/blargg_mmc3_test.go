package test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed roms/mmc3_irq_tests/1.Clocking.nes
var blarggMMC3IRQClocking string

//go:embed roms/mmc3_irq_tests/2.Details.nes
var blarggMMC3IRQDetails string

//go:embed roms/mmc3_irq_tests/3.A12_clocking.nes
var blarggMMC3IRQA12Clocking string

//go:embed roms/mmc3_irq_tests/4.Scanline_timing.nes
var blarggMMC3IRQScanlineTiming string

func Test_blarggMMC3(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rom        string
		wantStatus status
		want       string
	}{
		{"IRQ clocking", blarggMMC3IRQClocking, -1, "MMC3 IRQ COUNTER\nPASSED"},
		{"IRQ details", blarggMMC3IRQDetails, -1, "MMC3 IRQ COUNTER DETAILS\nPASSED"},
		{"IRQ A12 clocking", blarggMMC3IRQA12Clocking, -1, "MMC3 IRQ COUNTER A12\nPASSED"},
		{"IRQ scanline timing", blarggMMC3IRQScanlineTiming, -1, "MMC3 IRQ TIMING\nFAILED #3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			test, err := newBlarggTest(strings.NewReader(tt.rom))
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.EqualValues(t, tt.wantStatus, getBlarggStatus(test))
			assert.EqualValues(t, tt.want, getBlarggMessage(test, msgTypePPUVRAM))
		})
	}
}
