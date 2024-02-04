package ls

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

type Format string

const (
	FormatTable Format = "table"
	FormatJson  Format = "json"
	FormatYaml  Format = "yaml"
)

var ErrInvalidFormat = errors.New("invalid format")

func printEntries(out io.Writer, carts []entry, format Format) error {
	switch format {
	case FormatTable:
		return printTable(out, carts)
	case FormatJson:
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "  ")
		return encoder.Encode(carts)
	case FormatYaml:
		encoder := yaml.NewEncoder(out)
		return encoder.Encode(carts)
	}
	return fmt.Errorf("%w: %s", ErrInvalidFormat, format)
}

func printTable(out io.Writer, carts []entry) error {
	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintln(w, "FILE\tNAME\tMAPPER\tMIRROR\tBATTERY\t"); err != nil {
		return err
	}

	for _, entry := range carts {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%t\t\n",
			entry.Path,
			entry.Name,
			entry.Mapper,
			entry.Mirror,
			entry.Battery,
		)
	}

	return w.Flush()
}
