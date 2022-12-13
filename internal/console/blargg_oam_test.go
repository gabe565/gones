package console

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//go:embed nes-test-roms/oam_read/oam_read.nes
var oamRead string
var oamReadSuccess = `----------------
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

Passed
`

func Test_oamRead(t *testing.T) {
	t.Parallel()

	c, err := stubConsole(strings.NewReader(oamRead))
	if !assert.NoError(t, err) {
		return
	}

	status, err := runBlarggTest(c)
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, StatusSuccess, status)
	assert.EqualValues(t, oamReadSuccess, getBlarggMessage(c.Bus))
}

//go:embed nes-test-roms/oam_stress/oam_stress.nes
var oamStress string
var oamStressSuccess = `----------------
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

Passed
`

func Test_oamStress(t *testing.T) {
	t.Parallel()

	c, err := stubConsole(strings.NewReader(oamStress))
	if !assert.NoError(t, err) {
		return
	}

	status, err := runBlarggTest(c)
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, StatusSuccess, status)
	assert.EqualValues(t, oamStressSuccess, getBlarggMessage(c.Bus))
}
