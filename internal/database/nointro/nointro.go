package nointro

import (
	_ "embed"
	"encoding/xml"
)

//go:embed nes.xml
var Nes []byte

func Load(src []byte) (Datafile, error) {
	var datafile Datafile

	if err := xml.Unmarshal(src, &datafile); err != nil {
		return datafile, err
	}

	return datafile, nil
}
