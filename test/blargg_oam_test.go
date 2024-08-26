package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const blarggOAMReadWant = `----------------
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

const blarggOAMStressWant = `----------------
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

func Test_blarggOAM(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		rom        string
		wantStatus status
		want       string
	}{
		{"read", "roms/oam_read/oam_read.nes", 0, blarggOAMReadWant},
		{"stress", "roms/oam_stress/oam_stress.nes", 0, blarggOAMStressWant},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rom, err := roms.Open(tt.rom)
			require.NoError(t, err)

			test, err := newBlarggTest(rom, msgTypeSRAM)
			require.NoError(t, err)

			require.NoError(t, test.run())
			assert.EqualValues(t, tt.wantStatus, getBlarggStatus(test))
			assert.EqualValues(t, tt.want, getBlarggMessage(test, msgTypeSRAM))
		})
	}
}
