package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const blarggPPUVblNMIWant = `08 07 
Clock is skipped too late, relative to enabling BG

10-even_odd_timing

Failed #3

While running test 10 of 10`

func Test_blarggPPU(t *testing.T) {
	t.Parallel()

	sramTests := []struct {
		name       string
		rom        string
		wantStatus status
		want       string
	}{
		{"open bus", "roms/ppu_open_bus/ppu_open_bus.nes", 0, "ppu_open_bus\n\nPassed"},
		{"vbl nmi", "roms/ppu_vbl_nmi/ppu_vbl_nmi.nes", 1, blarggPPUVblNMIWant},
	}
	for _, tt := range sramTests {
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

	frameCountTests := []struct {
		name        string
		rom         string
		renderCount int
		want        string
	}{
		{"palette RAM", "roms/blargg_ppu_tests_2005.09.15b/palette_ram.nes", 17, "$01"},
		{"sprite RAM", "roms/blargg_ppu_tests_2005.09.15b/sprite_ram.nes", 17, "$01"},
		{"vram access", "roms/blargg_ppu_tests_2005.09.15b/vram_access.nes", 17, "$01"},
	}
	for _, tt := range frameCountTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rom, err := roms.Open(tt.rom)
			require.NoError(t, err)

			test, err := newConsoleTest(rom, exitAfterFrameNum(tt.renderCount))
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.EqualValues(t, tt.want, getBlarggMessage(test, msgTypePPUVRAM))
		})
	}
}
