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
		msgType    msgType
		wantStatus status
		want       string
	}{
		{"open bus", "roms/ppu_open_bus/ppu_open_bus.nes", msgTypeSRAM, 0, "ppu_open_bus\n\nPassed"},
		{"vbl nmi", "roms/ppu_vbl_nmi/ppu_vbl_nmi.nes", msgTypeSRAM, 1, blarggPPUVblNMIWant},

		{"sprite overflow basics", "roms/sprite_overflow_tests/1.Basics.nes", msgTypePPUVRAM, -1, "SPRITE OVERFLOW BASICS\nPASSED"},
		{"sprite overflow details", "roms/sprite_overflow_tests/2.Details.nes", msgTypePPUVRAM, -1, "SPRITE OVERFLOW DETAILS\nPASSED"},
		{"sprite overflow timing", "roms/sprite_overflow_tests/3.Timing.nes", msgTypePPUVRAM, -1, "SPRITE OVERFLOW TIMING\nFAILED: #5"},
		{"sprite overflow obscure", "roms/sprite_overflow_tests/4.Obscure.nes", msgTypePPUVRAM, -1, "SPRITE OVERFLOW OBSCURE\nFAILED: #2"},
		{"sprite overflow emulator", "roms/sprite_overflow_tests/5.Emulator.nes", msgTypePPUVRAM, -1, "SPRITE OVERFLOW EMULATION\nPASSED"},

		{"sprite hit basics", "roms/sprite_hit_tests_2005.10.05/01.basics.nes", msgTypePPUVRAM, -1, "SPRITE HIT BASICS\nPASSED"},
		{"sprite hit alignment", "roms/sprite_hit_tests_2005.10.05/02.alignment.nes", msgTypePPUVRAM, -1, "SPRITE HIT ALIGNMENT\nPASSED"},
		{"sprite hit corners", "roms/sprite_hit_tests_2005.10.05/03.corners.nes", msgTypePPUVRAM, -1, "SPRITE HIT CORNERS\nPASSED"},
		{"sprite hit flip", "roms/sprite_hit_tests_2005.10.05/04.flip.nes", msgTypePPUVRAM, -1, "SPRITE HIT FLIPPING\nPASSED"},
		{"sprite hit left clip", "roms/sprite_hit_tests_2005.10.05/05.left_clip.nes", msgTypePPUVRAM, -1, "SPRITE HIT LEFT CLIPPING\nPASSED"},
		{"sprite hit right edge", "roms/sprite_hit_tests_2005.10.05/06.right_edge.nes", msgTypePPUVRAM, -1, "SPRITE HIT RIGHT EDGE\nPASSED"},
		{"sprite hit screen bottom", "roms/sprite_hit_tests_2005.10.05/07.screen_bottom.nes", msgTypePPUVRAM, -1, "SPRITE HIT SCREEN BOTTOM\nPASSED"},
		{"sprite hit double height", "roms/sprite_hit_tests_2005.10.05/08.double_height.nes", msgTypePPUVRAM, -1, "SPRITE HIT DOUBLE HEIGHT\nPASSED"},
		{"sprite hit timing basics", "roms/sprite_hit_tests_2005.10.05/09.timing_basics.nes", msgTypePPUVRAM, -1, "SPRITE HIT TIMING\nPASSED"},
		{"sprite hit timing order", "roms/sprite_hit_tests_2005.10.05/10.timing_order.nes", msgTypePPUVRAM, -1, "SPRITE HIT ORDER\nPASSED"},
		{"sprite hit edge timing", "roms/sprite_hit_tests_2005.10.05/11.edge_timing.nes", msgTypePPUVRAM, -1, "SPRITE HIT EDGE TIMING\nPASSED"},
	}
	for _, tt := range sramTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rom, err := roms.Open(tt.rom)
			require.NoError(t, err)

			test, err := newBlarggTest(rom, tt.msgType)
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.EqualValues(t, tt.wantStatus, getBlarggStatus(test))
			assert.EqualValues(t, tt.want, getBlarggMessage(test, tt.msgType))
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
		{"vblank clear time", "roms/blargg_ppu_tests_2005.09.15b/vbl_clear_time.nes", 22, "$01"},
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
