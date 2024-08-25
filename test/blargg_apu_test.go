package test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed roms/apu_reset/4015_cleared.nes
var blarggAPUReset4015Cleared string

//go:embed roms/apu_reset/irq_flag_cleared.nes
var blarggAPUResetIRQCleared string

func Test_blarggAPU(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rom        string
		wantStatus status
		want       string
	}{
		{"reset clears $4015", blarggAPUReset4015Cleared, 0, "4015_cleared\n\nPassed"},
		{"reset clears IRQ", blarggAPUResetIRQCleared, 0, "irq_flag_cleared\n\nPassed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			test, err := newBlarggTest(strings.NewReader(tt.rom), msgTypeSRAM)
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.EqualValues(t, tt.wantStatus, getBlarggStatus(test))
			assert.EqualValues(t, tt.want, getBlarggMessage(test, msgTypeSRAM))
		})
	}
}
