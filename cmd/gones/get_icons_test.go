package gones

import (
	"bytes"
	_ "image/png"
	"os"
	"testing"

	"github.com/gabe565/gones/internal/config"
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest
func Test_getWindowIcons(t *testing.T) {
	var buf bytes.Buffer
	config.InitLog(&buf)
	t.Cleanup(func() {
		config.InitLog(os.Stderr)
	})

	icons := getWindowIcons()
	assert.Len(t, icons, 3)
	assert.Empty(t, buf.String())
}
