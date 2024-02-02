package console

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed nes-test-roms/apu_reset/4015_cleared.nes
var blarggApuRst4015Clr string

func Test_blarggApuRst4015Clr(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggApuRst4015Clr))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, StatusSuccess, GetBlarggStatus(test))
	assert.EqualValues(t, "4015_cleared\n\nPassed", GetBlarggMessage(test))
}

//go:embed nes-test-roms/apu_reset/irq_flag_cleared.nes
var blarggIrqClr string

func Test_blarggIrqClr(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggIrqClr))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, StatusSuccess, GetBlarggStatus(test))
	assert.EqualValues(t, "irq_flag_cleared\n\nPassed", GetBlarggMessage(test))
}
