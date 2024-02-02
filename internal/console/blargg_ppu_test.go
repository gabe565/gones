package console

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed nes-test-roms/ppu_open_bus/ppu_open_bus.nes
var blarggPpuOpenBus string

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

	assert.EqualValues(t, 0, GetBlarggStatus(test))
	assert.EqualValues(t, "ppu_open_bus\n\nPassed", GetBlarggMessage(test))
}

//go:embed nes-test-roms/ppu_vbl_nmi/rom_singles/01-vbl_basics.nes
var blarggVblBasics string

func Test_blarggVblBasics(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggVblBasics))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 0, GetBlarggStatus(test))
	assert.EqualValues(t, "01-vbl_basics\n\nPassed", GetBlarggMessage(test))
}

//go:embed nes-test-roms/ppu_vbl_nmi/rom_singles/02-vbl_set_time.nes
var blarggVblSetTime string

var blarggVblSetTimeSuccess = `T+ 1 2
00 - V
01 - V
02 - V
03 - V
04 - -
05 V -
06 V -
07 V -
08 V -

02-vbl_set_time

Passed`

func Test_blarggVblSetTime(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggVblSetTime))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 0, GetBlarggStatus(test))
	assert.EqualValues(t, blarggVblSetTimeSuccess, GetBlarggMessage(test))
}

//go:embed nes-test-roms/ppu_vbl_nmi/rom_singles/03-vbl_clear_time.nes
var blarggVblClearTime string

var blarggVblClearTimeSuccess = `00 V
01 V
02 V
03 V
04 V
05 V
06 -
07 -
08 -

03-vbl_clear_time

Passed`

func Test_blarggVblClearTime(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggVblClearTime))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 0, GetBlarggStatus(test))
	assert.EqualValues(t, blarggVblClearTimeSuccess, GetBlarggMessage(test))
}

//go:embed nes-test-roms/ppu_vbl_nmi/rom_singles/04-nmi_control.nes
var blarggNmiControl string

func Test_blarggNmiControl(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggNmiControl))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 0, GetBlarggStatus(test))
	assert.EqualValues(t, "04-nmi_control\n\nPassed", GetBlarggMessage(test))
}

//go:embed nes-test-roms/ppu_vbl_nmi/rom_singles/05-nmi_timing.nes
var blarggNmiTiming string

var blarggNmiTimingSuccess = `00 4
01 4
02 4
03 3
04 3
05 3
06 3
07 3
08 3
09 2

05-nmi_timing

Passed`

func Test_blarggNmiTiming(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggNmiTiming))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 0, GetBlarggStatus(test))
	assert.EqualValues(t, blarggNmiTimingSuccess, GetBlarggMessage(test))
}

//go:embed nes-test-roms/ppu_vbl_nmi/rom_singles/06-suppression.nes
var blarggNmiSuppression string

var blarggNmiSuppressionSuccess = `00 - N
01 - N
02 - N
03 - N
04 - -
05 V -
06 V -
07 V N
08 V N
09 V N

06-suppression

Passed`

func Test_blarggNmiSuppression(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggNmiSuppression))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 0, GetBlarggStatus(test))
	assert.EqualValues(t, blarggNmiSuppressionSuccess, GetBlarggMessage(test))
}

//go:embed nes-test-roms/ppu_vbl_nmi/rom_singles/07-nmi_on_timing.nes
var blarggNmiOnTiming string

var blarggNmiOnTimingSuccess = `00 N
01 N
02 N
03 N
04 N
05 -
06 -
07 -
08 -

07-nmi_on_timing

Passed`

func Test_blarggNmiOnTiming(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggNmiOnTiming))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 0, GetBlarggStatus(test))
	assert.EqualValues(t, blarggNmiOnTimingSuccess, GetBlarggMessage(test))
}

//go:embed nes-test-roms/ppu_vbl_nmi/rom_singles/08-nmi_off_timing.nes
var blarggNmiOffTiming string

var blarggNmiOffTimingSuccess = `03 -
04 -
05 -
06 -
07 N
08 N
09 N
0A N
0B N
0C N

08-nmi_off_timing

Passed`

func Test_blarggNmiOffTiming(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggNmiOffTiming))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 0, GetBlarggStatus(test))
	assert.EqualValues(t, blarggNmiOffTimingSuccess, GetBlarggMessage(test))
}

//go:embed nes-test-roms/ppu_vbl_nmi/rom_singles/09-even_odd_frames.nes
var blarggEvenOddFrames string

var blarggEvenOddFramesSuccess = `00 01 01 02 
09-even_odd_frames

Passed`

func Test_blarggEvenOddFrames(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggEvenOddFrames))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 0, GetBlarggStatus(test))
	assert.EqualValues(t, blarggEvenOddFramesSuccess, GetBlarggMessage(test))
}

//go:embed nes-test-roms/ppu_vbl_nmi/rom_singles/10-even_odd_timing.nes
var blarggEvenOddTiming string

var blarggEvenOddTimingSuccess = `08 07 
Clock is skipped too late, relative to enabling BG

10-even_odd_timing

Failed #3`

func Test_blarggEvenOddTiming(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggEvenOddTiming))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 3, GetBlarggStatus(test))
	assert.EqualValues(t, blarggEvenOddTimingSuccess, GetBlarggMessage(test))
}
