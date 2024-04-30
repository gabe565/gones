package gones

import (
	"bytes"
	_ "image/png"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest
func Test_getWindowIcons(t *testing.T) {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	t.Cleanup(func() {
		log.Logger = log.Output(os.Stderr)
	})

	icons := getWindowIcons()
	assert.Len(t, icons, 3)
	assert.Empty(t, buf.String())
}
