package console

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//go:embed nes-test-roms/ppu_open_bus/ppu_open_bus.nes
var blarggPpuOpenBus string
var blarggPpuOpenBusSuccess = `
Decay value should become zero by one second

ppu_open_bus

Failed #3
`

func Test_blarggPpuOpenBus(t *testing.T) {
	t.Parallel()

	c, err := stubConsole(strings.NewReader(blarggPpuOpenBus))
	if !assert.NoError(t, err) {
		return
	}

	status, err := runBlarggTest(c)
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 3, status)
	assert.EqualValues(t, blarggPpuOpenBusSuccess, getBlarggMessage(c.Bus))
}
