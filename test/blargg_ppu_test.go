package test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed roms/ppu_open_bus/ppu_open_bus.nes
var blarggPPUOpenBus string

//go:embed roms/ppu_vbl_nmi/ppu_vbl_nmi.nes
var blarggPPUVblNMI string

const blarggPPUVblNMIWant = `08 07 
Clock is skipped too late, relative to enabling BG

10-even_odd_timing

Failed #3

While running test 10 of 10`

func Test_blarggPPU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rom        string
		wantStatus status
		want       string
	}{
		{"open bus", blarggPPUOpenBus, 0, "ppu_open_bus\n\nPassed"},
		{"vbl nmi", blarggPPUVblNMI, 1, blarggPPUVblNMIWant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			test, err := newBlarggTest(strings.NewReader(tt.rom), msgTypeSRAM)
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.EqualValues(t, tt.wantStatus, getBlarggStatus(test))
			assert.EqualValues(t, tt.want, getBlarggMessage(test, msgTypeSRAM))
		})
	}
}
