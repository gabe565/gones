package decode

import (
	"errors"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decode code...",
		Short: "Decode a Game Genie code",
		Args:  cobra.MinimumNArgs(1),
		RunE:  run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintln(w, "CODE\tCPU ADDRESS\tREPLACE VALUE\tCOMPARE VALUE\t"); err != nil {
		return err
	}

	var errs []error
	for _, c := range args {
		result, err := decode(c)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if _, err := fmt.Fprintf(w,
			"%s\t0x%04X\t0x%02X\t%s\t\n",
			result.Code, result.Address, result.Replace, result.compareString(),
		); err != nil {
			return err
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return errors.Join(errs...)
}

type decodeResult struct {
	Code    string
	Address int
	Replace int
	Compare int
}

var (
	ErrInvalidCodeLen   = errors.New("invalid length")
	ErrInvalidCharacter = errors.New("invalid character")
)

func decode(code string) (decodeResult, error) {
	code = strings.ToUpper(code)
	result := decodeResult{Code: code, Compare: -1}

	switch len(code) {
	case 6, 8:
	default:
		return result, fmt.Errorf("%w %d in code %q; expected 6 or 8 characters", ErrInvalidCodeLen, len(code), code)
	}

	// Convert letters into integers (0x0-0xF)
	const lookup = "APZLGITYEOXUKSVN"
	codeValues := make([]int, 0, len(code))
	for _, r := range code {
		index := strings.IndexRune(lookup, r)
		if index == -1 {
			return result, fmt.Errorf("%w %q in code %q", ErrInvalidCharacter, r, code)
		}
		codeValues = append(codeValues, index)
	}

	// 24/32 bits (16 for address, 8 for replacement, 0 or 8 for compare)
	var bigint int
	loPosOrder := []int{3, 5, 2, 4, 1, 0, 7, 6}
	for _, loPos := range loPosOrder[:len(code)] {
		hiPos := (loPos - 1 + len(code)) % len(code)
		bigint = (bigint << 4) | (codeValues[hiPos] & 8) | (codeValues[loPos] & 7)
	}

	// Split integer and set MSB of address
	if len(code) == 8 {
		compValue := bigint & 0xFF
		result.Compare = compValue
		bigint >>= 8
	}

	result.Address = (bigint >> 8) | 0x8000
	result.Replace = bigint & 0xFF
	return result, nil
}

func (d decodeResult) compareString() string {
	if d.Compare == -1 {
		return "<none>"
	}
	return fmt.Sprintf("0x%02X", d.Compare)
}
