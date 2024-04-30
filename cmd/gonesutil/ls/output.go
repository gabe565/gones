package ls

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

//go:generate enumer -type OutputFormat -trimprefix OutputFormat -transform lower -values

type OutputFormat uint8

const (
	OutputFormatTable OutputFormat = iota
	OutputFormatJSON
	OutputFormatYAML
)

var ErrInvalidFormat = errors.New("invalid format")

func printEntries(out io.Writer, carts []*entry, format OutputFormat) error {
	switch format {
	case OutputFormatTable:
		return printTable(out, carts)
	case OutputFormatJSON:
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		return encoder.Encode(carts)
	case OutputFormatYAML:
		encoder := yaml.NewEncoder(out)
		return encoder.Encode(carts)
	}
	return fmt.Errorf("%w: %s", ErrInvalidFormat, format)
}

func printTable(out io.Writer, carts []*entry) error {
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintln(w, "FILE\tNAME\tMAPPER\tMIRROR\tBATTERY\tHASH\t"); err != nil {
		return err
	}

	for _, entry := range carts {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%t\t%s\t\n",
			entry.Path,
			entry.Name,
			entry.Mapper,
			entry.Mirror,
			entry.Battery,
			entry.Hash,
		)
	}

	return w.Flush()
}
