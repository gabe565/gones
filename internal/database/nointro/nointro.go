package nointro

import "encoding/xml"

func Load(src []byte) (Datafile, error) {
	var datafile Datafile

	if err := xml.Unmarshal(src, &datafile); err != nil {
		return datafile, err
	}

	return datafile, nil
}
