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

	test, err := NewBlarggTest(strings.NewReader(blarggPpuOpenBus))
	if !assert.NoError(t, err) {
		return
	}

	err = test.Run()
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, 3, GetBlarggStatus(test))
	assert.EqualValues(t, blarggPpuOpenBusSuccess, GetBlarggMessage(test))
}
