package console

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed nes-test-roms/apu_reset/4015_cleared.nes
var blarggAPURst4015Clr string

func Test_blarggAPURst4015Clr(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggAPURst4015Clr))
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
var blarggIRQClr string

func Test_blarggIRQClr(t *testing.T) {
	t.Parallel()

	test, err := NewBlarggTest(strings.NewReader(blarggIRQClr))
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
