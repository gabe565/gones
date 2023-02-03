//go:build embed_nes_xml

package nointro

import _ "embed"

//go:embed nes.xml
var Nes []byte
