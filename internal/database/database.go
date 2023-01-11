package database

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"errors"
	"io"
)

//go:generate go run ./generate_xml

//go:embed database.csv
var database []byte

var ErrNotFound = errors.New("not found")

func FindNameByHash(hash string) (string, error) {
	c := csv.NewReader(bytes.NewReader(database))
	for {
		record, err := c.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}

		if record[0] == hash {
			return record[1], nil
		}
	}

	return "", ErrNotFound
}
