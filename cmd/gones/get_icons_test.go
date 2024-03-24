package gones

import (
	"bytes"
	_ "image/png"
	"io"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_getWindowIcons(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	defer func(w io.Writer) {
		log.SetOutput(w)
	}(log.StandardLogger().Out)
	log.SetOutput(&buf)

	icons := getWindowIcons()
	assert.Len(t, icons, 3)
	assert.Empty(t, buf.String())
}
