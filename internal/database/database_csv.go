//go:build !gzip

package database

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"errors"
	"io"
)

// database is the CSV committed to the repo. Release builds embed the compact
// binary form instead, which go:generate builds from this file.
//
//go:embed database.csv
var database []byte

func FindNameByHash(hash string) (string, error) {
	c := csv.NewReader(bytes.NewReader(database))
	c.FieldsPerRecord = 2
	c.ReuseRecord = true

	for {
		record, err := c.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return "", ErrNotFound
			}
			return "", err
		}

		if record[0] == hash {
			return record[1], nil
		}
	}
}
