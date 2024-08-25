package test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed roms/ppu_open_bus/ppu_open_bus.nes
var blarggPPUOpenBus string

func Test_blarggPPUOpenBus(t *testing.T) {
	t.Parallel()

	test, err := newBlarggTest(strings.NewReader(blarggPPUOpenBus))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, 0, getBlarggStatus(test))
	assert.EqualValues(t, "ppu_open_bus\n\nPassed", getBlarggMessage(test, msgTypeSRAM))
}

//go:embed roms/ppu_vbl_nmi/ppu_vbl_nmi.nes
var blarggPPUVblNMI string

const blarggPPUVblNMISuccess = `08 07 
Clock is skipped too late, relative to enabling BG

10-even_odd_timing

Failed #3

While running test 10 of 10`

func Test_blarggPPUVblNMI(t *testing.T) {
	t.Parallel()

	test, err := newBlarggTest(strings.NewReader(blarggPPUVblNMI))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, blarggPPUVblNMISuccess, getBlarggMessage(test, msgTypeSRAM))
}
