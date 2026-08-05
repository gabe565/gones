//go:build gzip

package database

import (
	"bytes"
	"compress/gzip"
	_ "embed"

	"gabe565.com/gones/internal/database/compact"
)

// database is built from database.csv by go:generate, so it is not committed.
//
//go:embed database.bin.gz
var database []byte

func FindNameByHash(hash string) (string, error) {
	key, err := compact.Key(hash)
	if err != nil {
		return "", ErrNotFound
	}

	gzr, err := gzip.NewReader(bytes.NewReader(database))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = gzr.Close()
	}()

	name, ok, err := compact.Find(gzr, key)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", ErrNotFound
	}
	return name, nil
}
