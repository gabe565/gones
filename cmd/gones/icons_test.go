package gones

import (
	"bytes"
	_ "image/png"
	"os"
	"testing"

	"github.com/gabe565/gones/internal/log"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest
func Test_getWindowIcons(t *testing.T) {
	var buf bytes.Buffer
	log.Init(&buf)
	t.Cleanup(func() {
		log.Init(os.Stderr)
	})

	icons := getWindowIcons()
	assert.Len(t, icons, 3)
	assert.Empty(t, buf.String())
}
