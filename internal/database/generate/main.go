package main

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"gabe565.com/gones/internal/database/rdb"
	"gabe565.com/gones/internal/log"
)

// url is the libretro mirror of the No-Intro NES database.
const url = "https://raw.githubusercontent.com/libretro/libretro-database/master/rdb/Nintendo%20-%20Nintendo%20Entertainment%20System.rdb"

//nolint:gochecknoglobals
var path = filepath.Join("internal", "database", "database.csv")

func main() {
	log.Init(os.Stderr)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.Error("Failed to generate database", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	data, err := download(ctx)
	if err != nil {
		return err
	}

	names, err := parse(data)
	if err != nil {
		return err
	}

	return write(names)
}

func download(ctx context.Context) ([]byte, error) {
	slog.Info("Downloading database", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	slog.Info("Downloaded database", "bytes", len(data))
	return data, nil
}

type entry struct {
	hash string
	name string
}

// parse builds a deduplicated, hash-sorted entry list. If a duplicate is
// encountered, only keeps the shorter title.
func parse(data []byte) ([]entry, error) {
	names := make(map[string]string)

	var total int
	for e, err := range rdb.All(data) {
		if err != nil {
			return nil, err
		}
		total++

		if len(e.MD5) == 0 || e.Name == "" {
			continue
		}

		hash := hex.EncodeToString(e.MD5)
		if prev, ok := names[hash]; ok && len(prev) <= len(e.Name) {
			continue
		}
		names[hash] = e.Name
	}

	entries := make([]entry, 0, len(names))
	for hash, name := range names {
		entries = append(entries, entry{hash: hash, name: name})
	}
	slices.SortFunc(entries, func(a, b entry) int {
		return strings.Compare(a.hash, b.hash)
	})

	slog.Info("Parsed database", "records", total, "games", len(entries))
	return entries, nil
}

func write(entries []entry) error {
	slog.Info("Creating CSV file", "path", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	slog.Info("Creating gzipped CSV file", "path", path+".gz")
	gzf, err := os.Create(path + ".gz")
	if err != nil {
		return err
	}
	defer func() { _ = gzf.Close() }()
	gz := gzip.NewWriter(gzf)
	defer func() { _ = gz.Close() }()

	c := csv.NewWriter(io.MultiWriter(f, gz))
	for _, e := range entries {
		if err := c.Write([]string{e.hash, e.name}); err != nil {
			return err
		}
	}
	c.Flush()

	return errors.Join(
		c.Error(),
		gz.Close(),
		gzf.Close(),
		f.Close(),
	)
}
