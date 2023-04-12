package console

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed nes-test-roms/ppu_open_bus/ppu_open_bus.nes
var blarggPpuOpenBus string

var blarggPpuOpenBusSuccess = `
Decay value should become zero by one second

ppu_open_bus

Failed #3
`

func Test_blarggPpuOpenBus(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggPpuOpenBus))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 3, GetBlarggStatus(test))
	assert.EqualValues(t, blarggPpuOpenBusSuccess, GetBlarggMessage(test))
}

//go:embed nes-test-roms/vbl_nmi_timing/1.frame_basics.nes
var blarggFrameBasicsTest string

func Test_blarggFrameBasics(t *testing.T) {
	t.Parallel()

	var message string
	callback := NewBlargPpuMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggFrameBasicsTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	switch err := err.(type) {
	case PpuMessage:
		message = err.Message
	default:
		assert.NoError(t, err)
		return
	}

	assert.EqualValues(t, "PPU FRAME BASICS\nPASSED", message)
}

//go:embed nes-test-roms/vbl_nmi_timing/2.vbl_timing.nes
var blarggVblTest string

func Test_blarggVbl(t *testing.T) {
	t.Parallel()

	var message string
	callback := NewBlargPpuMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggVblTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	switch err := err.(type) {
	case PpuMessage:
		message = err.Message
	default:
		assert.NoError(t, err)
		return
	}

	assert.EqualValues(t, "VBL TIMING\nPASSED", message)
}

//go:embed nes-test-roms/vbl_nmi_timing/3.even_odd_frames.nes
var blarggEvenOddFramesTest string

func Test_blarggEvenOddFrames(t *testing.T) {
	t.Parallel()

	var message string
	callback := NewBlargPpuMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggEvenOddFramesTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	switch err := err.(type) {
	case PpuMessage:
		message = err.Message
	default:
		assert.NoError(t, err)
		return
	}

	assert.EqualValues(t, "EVEN ODD FRAMES\nPASSED", message)
}

//go:embed nes-test-roms/vbl_nmi_timing/4.vbl_clear_timing.nes
var blarggVblClearTimingTest string

func Test_blarggVblClearTiming(t *testing.T) {
	t.Parallel()

	var message string
	callback := NewBlargPpuMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggVblClearTimingTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	switch err := err.(type) {
	case PpuMessage:
		message = err.Message
	default:
		assert.NoError(t, err)
		return
	}

	assert.EqualValues(t, "VBL CLEAR TIMING\nPASSED", message)
}

//go:embed nes-test-roms/vbl_nmi_timing/5.nmi_suppression.nes
var blarggNmiSuppressionTest string

func Test_blarggNmiSuppression(t *testing.T) {
	t.Parallel()

	var message string
	callback := NewBlargPpuMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggNmiSuppressionTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	switch err := err.(type) {
	case PpuMessage:
		message = err.Message
	default:
		assert.NoError(t, err)
		return
	}

	assert.EqualValues(t, "NMI SUPPRESSION\nPASSED", message)
}

//go:embed nes-test-roms/vbl_nmi_timing/6.nmi_disable.nes
var blarggNmiDisableTest string

func Test_blarggNmiDisable(t *testing.T) {
	t.Parallel()

	var message string
	callback := NewBlargPpuMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggNmiDisableTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	switch err := err.(type) {
	case PpuMessage:
		message = err.Message
	default:
		assert.NoError(t, err)
		return
	}

	assert.EqualValues(t, "NMI DISABLE\nPASSED", message)
}

//go:embed nes-test-roms/vbl_nmi_timing/7.nmi_timing.nes
var blarggNmiTimingTest string

func Test_blarggNmiTiming(t *testing.T) {
	t.Parallel()

	var message string
	callback := NewBlargPpuMessageCallback()

	test, err := NewConsoleTest(strings.NewReader(blarggNmiTimingTest), callback)
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.Error(t, err) {
		return
	}

	switch err := err.(type) {
	case PpuMessage:
		message = err.Message
	default:
		assert.NoError(t, err)
		return
	}

	assert.EqualValues(t, "NMI TIMING\nPASSED", message)
}
