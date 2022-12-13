package console

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

//go:embed nes-test-roms/instr_test-v5/all_instrs.nes
var blarggInstrTest string

func Test_blarggCpuTest(t *testing.T) {
	t.Parallel()

	c, err := stubConsole(strings.NewReader(blarggInstrTest))
	if !assert.NoError(t, err) {
		return
	}

	status, err := runBlarggTest(c)
	if !assert.NoError(t, err) {
		return
	}

	assert.EqualValues(t, StatusSuccess, status)
	assert.EqualValues(t, "All 16 tests passed\n\n\n", getBlarggMessage(c.Bus))
}
