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

func Test_blarggCPUTest(t *testing.T) {
	t.Parallel()

	test, err := newBlarggTest(strings.NewReader(blarggInstrTest))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, statusSuccess, getBlarggStatus(test))
	assert.EqualValues(t, "All 16 tests passed", getBlarggMessage(test, msgTypeSRAM))
}

//go:embed roms/cpu_timing_test6/cpu_timing_test.nes
var blarggCPUTimingTest string

const blarggCPUTimingSuccess = `6502 TIMING TEST (16 SECONDS)
OFFICIAL INSTRUCTIONS ONLY
PASSED`

func Test_blarggCPUTiming(t *testing.T) {
	t.Parallel()

	test, err := newBlarggPPUMsgTest(strings.NewReader(blarggCPUTimingTest))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, blarggCPUTimingSuccess, getBlarggMessage(test, msgTypePPUVRAM))
}

//go:embed roms/branch_timing_tests/1.Branch_Basics.nes
var blarggBranchTimingBasicsTest string

func Test_blarggBranchTimingBasics(t *testing.T) {
	t.Parallel()

	test, err := newBlarggPPUMsgTest(strings.NewReader(blarggBranchTimingBasicsTest))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, "BRANCH TIMING BASICS\nPASSED", getBlarggMessage(test, msgTypePPUVRAM))
}

//go:embed roms/branch_timing_tests/2.Backward_Branch.nes
var blarggBranchTimingBackwardTest string

func Test_blarggBranchTimingBackward(t *testing.T) {
	t.Parallel()

	test, err := newBlarggPPUMsgTest(strings.NewReader(blarggBranchTimingBackwardTest))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, "BACKWARD BRANCH TIMING\nPASSED", getBlarggMessage(test, msgTypePPUVRAM))
}

//go:embed roms/branch_timing_tests/3.Forward_Branch.nes
var blarggBranchTimingForwardTest string

func Test_blarggBranchTimingForward(t *testing.T) {
	t.Parallel()

	test, err := newBlarggPPUMsgTest(strings.NewReader(blarggBranchTimingForwardTest))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, "FORWARD BRANCH TIMING\nPASSED", getBlarggMessage(test, msgTypePPUVRAM))
}
