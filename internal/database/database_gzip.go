//go:build gzip

package database

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/csv"
	"errors"
	"io"
)

//go:generate sh -c "gzip -c database.csv > database.csv.gz"

//go:embed database.csv.gz
var database []byte

var ErrNotFound = errors.New("not found")

func FindNameByHash(hash string) (string, error) {
	gzr, err := gzip.NewReader(bytes.NewReader(database))
	if err != nil {
		return "", err
	}
	defer func(gzr *gzip.Reader) {
		_ = gzr.Close()
	}(gzr)

	c := csv.NewReader(gzr)
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
