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

//go:embed nes-test-roms/ppu_vbl_nmi/ppu_vbl_nmi.nes
var blarggPpuVblNmi string

var blarggPpuVblNmiSuccess = `08 07 
Clock is skipped too late, relative to enabling BG

10-even_odd_timing

Failed #3

While running test 10 of 10`

func Test_blarggPpuVblNmi(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggPpuVblNmi))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, blarggPpuVblNmiSuccess, GetBlarggMessage(test))
}
