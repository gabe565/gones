package console

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed nes-test-roms/instr_test-v5/all_instrs.nes
var blarggInstrTest string

func Test_blarggCPUTest(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggInstrTest))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, StatusSuccess, GetBlarggStatus(test))
	assert.EqualValues(t, "All 16 tests passed", GetBlarggMessage(test))
}

//go:embed nes-test-roms/cpu_timing_test6/cpu_timing_test.nes
var blarggCPUTimingTest string

const blarggCPUTimingSuccess = `6502 TIMING TEST (16 SECONDS)
OFFICIAL INSTRUCTIONS ONLY
PASSED`

func Test_blarggCPUTiming(t *testing.T) {
	t.Parallel()

	callback := NewBlargPPUMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggCPUTimingTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	assert.EqualValues(t, blarggCPUTimingSuccess, err.Error())
}

//go:embed nes-test-roms/branch_timing_tests/1.Branch_Basics.nes
var blarggBranchTimingBasicsTest string

func Test_blarggBranchTimingBasics(t *testing.T) {
	t.Parallel()

	callback := NewBlargPPUMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggBranchTimingBasicsTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	assert.EqualValues(t, "BRANCH TIMING BASICS\nPASSED", err.Error())
}

//go:embed nes-test-roms/branch_timing_tests/2.Backward_Branch.nes
var blarggBranchTimingBackwardTest string

func Test_blarggBranchTimingBackward(t *testing.T) {
	t.Parallel()

	callback := NewBlargPPUMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggBranchTimingBackwardTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	assert.EqualValues(t, "BACKWARD BRANCH TIMING\nPASSED", err.Error())
}

//go:embed nes-test-roms/branch_timing_tests/3.Forward_Branch.nes
var blarggBranchTimingForwardTest string

func Test_blarggBranchTimingForward(t *testing.T) {
	t.Parallel()

	callback := NewBlargPPUMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggBranchTimingForwardTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	assert.EqualValues(t, "FORWARD BRANCH TIMING\nPASSED", err.Error())
}
