package test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed roms/instr_test-v5/all_instrs.nes
var blarggInstrTest string

//go:embed roms/cpu_timing_test6/cpu_timing_test.nes
var blarggCPUTimingTest string

const blarggCPUTimingWant = `6502 TIMING TEST (16 SECONDS)
OFFICIAL INSTRUCTIONS ONLY
PASSED`

//go:embed roms/branch_timing_tests/1.Branch_Basics.nes
var blarggBranchTimingBasicsTest string

//go:embed roms/branch_timing_tests/2.Backward_Branch.nes
var blarggBranchTimingBackwardTest string

//go:embed roms/branch_timing_tests/3.Forward_Branch.nes
var blarggBranchTimingForwardTest string

//go:embed roms/cpu_reset/ram_after_reset.nes
var blarggCPUResetRAM string

//go:embed roms/cpu_reset/registers.nes
var blarggCPUResetRegisters string

func Test_blarggCPU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rom        string
		msgType    msgType
		wantStatus status
		want       string
	}{
		{"instructions", blarggInstrTest, msgTypeSRAM, 0, "All 16 tests passed"},
		{"instruction timing", blarggCPUTimingTest, msgTypePPUVRAM, -1, blarggCPUTimingWant},
		{"branch timing basics", blarggBranchTimingBasicsTest, msgTypePPUVRAM, -1, "BRANCH TIMING BASICS\nPASSED"},
		{"branch timing backward", blarggBranchTimingBackwardTest, msgTypePPUVRAM, -1, "BACKWARD BRANCH TIMING\nPASSED"},
		{"branch timing forward", blarggBranchTimingForwardTest, msgTypePPUVRAM, -1, "FORWARD BRANCH TIMING\nPASSED"},
		{"ram after reset", blarggCPUResetRAM, msgTypeSRAM, 0, "ram_after_reset\n\nPassed"},
		{"registers after reset", blarggCPUResetRegisters, msgTypeSRAM, 0, "A  X  Y  P  S\n34 56 78 FF 0F \n\nregisters\n\nPassed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			test, err := newBlarggTest(strings.NewReader(tt.rom), tt.msgType)
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.EqualValues(t, tt.wantStatus, getBlarggStatus(test))
			assert.EqualValues(t, tt.want, getBlarggMessage(test, tt.msgType))
		})
	}
}
