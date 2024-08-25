package test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed roms/oam_read/oam_read.nes
var oamRead string

const oamReadSuccess = `----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------

oam_read

Passed`

func Test_oamRead(t *testing.T) {
	t.Parallel()

	test, err := newBlarggTest(strings.NewReader(oamRead))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, statusSuccess, getBlarggStatus(test))
	assert.EqualValues(t, oamReadSuccess, getBlarggMessage(test, msgTypeSRAM))
}

//go:embed roms/oam_stress/oam_stress.nes
var oamStress string

const oamStressSuccess = `----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------
----------------

oam_stress

Passed`

func Test_oamStress(t *testing.T) {
	t.Parallel()

	test, err := newBlarggTest(strings.NewReader(oamStress))
	require.NoError(t, err)
	require.NoError(t, test.run())

	assert.EqualValues(t, statusSuccess, getBlarggStatus(test))
	assert.EqualValues(t, oamStressSuccess, getBlarggMessage(test, msgTypeSRAM))
}
