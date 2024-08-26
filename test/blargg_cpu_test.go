package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const blarggCPUTimingWant = `6502 TIMING TEST (16 SECONDS)
OFFICIAL INSTRUCTIONS ONLY
PASSED`

func Test_blarggCPU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rom        string
		msgType    msgType
		wantStatus status
		want       string
	}{
		{"instructions", "roms/instr_test-v5/all_instrs.nes", msgTypeSRAM, 0, "All 16 tests passed"},
		{"instruction timing", "roms/cpu_timing_test6/cpu_timing_test.nes", msgTypePPUVRAM, -1, blarggCPUTimingWant},
		{"branch timing basics", "roms/branch_timing_tests/1.Branch_Basics.nes", msgTypePPUVRAM, -1, "BRANCH TIMING BASICS\nPASSED"},
		{"branch timing backward", "roms/branch_timing_tests/2.Backward_Branch.nes", msgTypePPUVRAM, -1, "BACKWARD BRANCH TIMING\nPASSED"},
		{"branch timing forward", "roms/branch_timing_tests/3.Forward_Branch.nes", msgTypePPUVRAM, -1, "FORWARD BRANCH TIMING\nPASSED"},
		{"ram after reset", "roms/cpu_reset/ram_after_reset.nes", msgTypeSRAM, 0, "ram_after_reset\n\nPassed"},
		{"registers after reset", "roms/cpu_reset/registers.nes", msgTypeSRAM, 0, "A  X  Y  P  S\n34 56 78 FF 0F \n\nregisters\n\nPassed"},
	}
	for _, tt := range tests {
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
}
