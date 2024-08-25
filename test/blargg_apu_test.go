package test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed roms/apu_reset/4015_cleared.nes
var blarggAPURst4015Clr string

func Test_blarggAPURst4015Clr(t *testing.T) {
	t.Parallel()

	test, err := newBlarggTest(strings.NewReader(blarggAPURst4015Clr))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, statusSuccess, getBlarggStatus(test))
	assert.EqualValues(t, "4015_cleared\n\nPassed", getBlarggMessage(test))
}

//go:embed roms/apu_reset/irq_flag_cleared.nes
var blarggIRQClr string

func Test_blarggIRQClr(t *testing.T) {
	t.Parallel()

	test, err := newBlarggTest(strings.NewReader(blarggIRQClr))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, statusSuccess, getBlarggStatus(test))
	assert.EqualValues(t, "irq_flag_cleared\n\nPassed", getBlarggMessage(test))
}
