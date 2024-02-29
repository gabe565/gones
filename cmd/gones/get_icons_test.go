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
	var buf bytes.Buffer
	defer func(w io.Writer) {
		log.SetOutput(w)
	}(log.StandardLogger().Out)
	log.SetOutput(&buf)

	icons := getWindowIcons()
	assert.Equal(t, 3, len(icons))
	assert.Empty(t, buf.String())
}
