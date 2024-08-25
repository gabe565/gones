package test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed roms/oam_read/oam_read.nes
var blarggOAMRead string

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

//go:embed roms/oam_stress/oam_stress.nes
var blarggOAMStress string

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
		{"read", blarggOAMRead, 0, blarggOAMReadWant},
		{"stress", blarggOAMStress, 0, blarggOAMStressWant},
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
