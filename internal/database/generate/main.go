package main

import (
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"gabe565.com/gones/internal/database/compact"
	"gabe565.com/gones/internal/database/rdb"
	"gabe565.com/gones/internal/log"
)

// url is the libretro mirror of the No-Intro NES database.
const url = "https://raw.githubusercontent.com/libretro/libretro-database/master/rdb/Nintendo%20-%20Nintendo%20Entertainment%20System.rdb"

//nolint:gochecknoglobals
var (
	// csvPath is committed to the repo and embedded by development builds.
	csvPath = filepath.Join("internal", "database", "database.csv")
	// compactPath is built from csvPath by go:generate and embedded by release builds.
	compactPath = filepath.Join("internal", "database", "database.bin.gz")
)

func main() {
	log.Init(os.Stderr)

	fromCSV := flag.Bool("from-csv", false,
		"Rebuild "+compactPath+" from the committed CSV instead of downloading the database.")
	flag.Parse()

	if err := run(*fromCSV); err != nil {
		slog.Error("Failed to generate database", "error", err)
		os.Exit(1)
	}
}

func run(fromCSV bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if fromCSV {
		entries, err := readCSV()
		if err != nil {
			return err
		}
		return writeCompact(entries)
	}

	data, err := download(ctx)
	if err != nil {
		return err
	}

	entries, err := parse(data)
	if err != nil {
		return err
	}

	if err := writeCSV(entries); err != nil {
		return err
	}
	return writeCompact(entries)
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

// parse builds a deduplicated entry list. If a duplicate is encountered, only
// keeps the shorter title. Entries are sorted by name so that similar titles
// end up adjacent, which compresses far better than hash order.
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
		if c := strings.Compare(a.name, b.name); c != 0 {
			return c
		}
		return strings.Compare(a.hash, b.hash)
	})

	slog.Info("Parsed database", "records", total, "games", len(entries))
	return entries, nil
}

func readCSV() ([]entry, error) {
	f, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	c := csv.NewReader(f)
	c.FieldsPerRecord = 2

	records, err := c.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := make([]entry, 0, len(records))
	for _, record := range records {
		entries = append(entries, entry{hash: record[0], name: record[1]})
	}

	slog.Info("Read CSV", "path", csvPath, "games", len(entries))
	return entries, nil
}

func writeCSV(entries []entry) error {
	slog.Info("Creating CSV file", "path", csvPath)

	f, err := os.Create(csvPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	c := csv.NewWriter(f)
	for _, e := range entries {
		if err := c.Write([]string{e.hash, e.name}); err != nil {
			return err
		}
	}
	c.Flush()

	return errors.Join(c.Error(), f.Close())
}

func writeCompact(entries []entry) error {
	hashes := make([]uint64, 0, len(entries))
	names := make([]string, 0, len(entries))
	seen := make(map[uint64]string, len(entries))

	for _, e := range entries {
		key, err := compact.Key(e.hash)
		if err != nil {
			return err
		}
		// Truncation is only safe while it stays collision free.
		if prev, ok := seen[key]; ok {
			return fmt.Errorf("%w: %q collides with %q", compact.ErrInvalid, e.name, prev)
		}
		seen[key] = e.name

		hashes = append(hashes, key)
		names = append(names, e.name)
	}

	data, err := compact.Marshal(hashes, names)
	if err != nil {
		return err
	}

	slog.Info("Creating compact database", "path", compactPath, "bytes", len(data))

	f, err := os.Create(compactPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	gz, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		return err
	}

	if _, err := gz.Write(data); err != nil {
		return errors.Join(err, gz.Close())
	}

	return errors.Join(gz.Close(), f.Close())
}
