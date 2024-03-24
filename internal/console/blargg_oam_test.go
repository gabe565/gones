package console

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed nes-test-roms/oam_read/oam_read.nes
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

	test, err := NewBlarggTest(strings.NewReader(oamRead))
	require.NoError(t, err)
	require.NoError(t, test.Run())

	assert.EqualValues(t, StatusSuccess, GetBlarggStatus(test))
	assert.EqualValues(t, oamReadSuccess, GetBlarggMessage(test))
}

//go:embed nes-test-roms/oam_stress/oam_stress.nes
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

	test, err := NewBlarggTest(strings.NewReader(oamStress))
	require.NoError(t, err)
	require.NoError(t, test.Run())

	assert.EqualValues(t, StatusSuccess, GetBlarggStatus(test))
	assert.EqualValues(t, oamStressSuccess, GetBlarggMessage(test))
}
