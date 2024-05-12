package palette

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed default.pal
var palFile []byte

func TestLoadPal(t *testing.T) {
	t.Parallel()
	UpdateEmphasized()

	defaultPalette := Palette{Emphasis: Default.Emphasis}
	copy(defaultPalette.RGBA[:], Default.RGBA[:])
	emphasizeR := Palette{Emphasis: EmphasizeR.Emphasis}
	copy(emphasizeR.RGBA[:], EmphasizeR.RGBA[:])
	emphasizeG := Palette{Emphasis: EmphasizeG.Emphasis}
	copy(emphasizeG.RGBA[:], EmphasizeG.RGBA[:])
	emphasizeB := Palette{Emphasis: EmphasizeB.Emphasis}
	copy(emphasizeB.RGBA[:], EmphasizeB.RGBA[:])
	emphasizeRG := Palette{Emphasis: EmphasizeRG.Emphasis}
	copy(emphasizeRG.RGBA[:], EmphasizeRG.RGBA[:])
	emphasizeRB := Palette{Emphasis: EmphasizeRB.Emphasis}
	copy(emphasizeRB.RGBA[:], EmphasizeRB.RGBA[:])
	emphasizeGB := Palette{Emphasis: EmphasizeGB.Emphasis}
	copy(emphasizeGB.RGBA[:], EmphasizeGB.RGBA[:])
	emphasizeRGB := Palette{Emphasis: EmphasizeRGB.Emphasis}
	copy(emphasizeRGB.RGBA[:], EmphasizeRGB.RGBA[:])

	Default = Palette{}
	EmphasizeR = Palette{Emphasis: EmphasizeR.Emphasis}
	EmphasizeG = Palette{Emphasis: EmphasizeG.Emphasis}
	EmphasizeB = Palette{Emphasis: EmphasizeB.Emphasis}
	EmphasizeRG = Palette{Emphasis: EmphasizeRG.Emphasis}
	EmphasizeRB = Palette{Emphasis: EmphasizeRB.Emphasis}
	EmphasizeGB = Palette{Emphasis: EmphasizeGB.Emphasis}
	EmphasizeRGB = Palette{Emphasis: EmphasizeRGB.Emphasis}

	if err := LoadPal(bytes.NewReader(palFile)); !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, defaultPalette, Default)
	assert.Equal(t, emphasizeR, EmphasizeR)
	assert.Equal(t, emphasizeG, EmphasizeG)
	assert.Equal(t, emphasizeB, EmphasizeB)
	assert.Equal(t, emphasizeRG, EmphasizeRG)
	assert.Equal(t, emphasizeRB, EmphasizeRB)
	assert.Equal(t, emphasizeGB, EmphasizeGB)
	assert.Equal(t, emphasizeRGB, EmphasizeRGB)
}
