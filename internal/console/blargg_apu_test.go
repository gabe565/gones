package console

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//go:embed nes-test-roms/apu_reset/4015_cleared.nes
var blarggApuRst4015Clr string

func Test_blarggApuRst4015Clr(t *testing.T) {
	t.Parallel()

	c, err := stubConsole(strings.NewReader(blarggApuRst4015Clr))
	if !assert.NoError(t, err) {
		return
	}

	status, err := runBlarggTest(c)
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, StatusSuccess, status)
	assert.EqualValues(t, "\n4015_cleared\n\nPassed\n", getBlarggMessage(c.Bus))
}

//go:embed nes-test-roms/apu_reset/irq_flag_cleared.nes
var blarggIrqClr string

func Test_blarggIrqClr(t *testing.T) {
	t.Parallel()

	c, err := stubConsole(strings.NewReader(blarggIrqClr))
	if !assert.NoError(t, err) {
		return
	}

	status, err := runBlarggTest(c)
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, StatusSuccess, status)
	assert.EqualValues(t, "\nirq_flag_cleared\n\nPassed\n", getBlarggMessage(c.Bus))
}
