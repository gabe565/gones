//go:build embed_nes_xml

package nointro

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	datafile, err := Load(Nes)
	if !assert.NoError(t, err) {
		return
	}

	assert.NotEmpty(t, datafile.Headers)
	assert.NotEmpty(t, datafile.Games)
}
