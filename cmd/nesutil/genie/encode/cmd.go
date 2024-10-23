package encode

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encode address replace [compare]",
		Short: "Encode a Game Genie code",
		Args:  cobra.RangeArgs(2, 3),
		RunE:  run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	address, err := strconv.ParseInt(args[0], 16, 32)
	if err != nil {
		return err
	}

	replace, err := strconv.ParseInt(args[1], 16, 16)
	if err != nil {
		return err
	}

	var compare *int
	if len(args) > 2 {
		v, err := strconv.ParseInt(args[2], 16, 16)
		if err != nil {
			return err
		}

		v2 := int(v)
		compare = &v2
	}

	code, err := encode(int(address), int(replace), compare)
	if err != nil {
		return err
	}

	_, err = io.WriteString(cmd.OutOrStdout(), code+"\n")
	return err
}

var ErrOutOfRange = errors.New("encoded value out of range")

func encode(address, replace int, compare *int) (string, error) {
	var codeLen int
	var bigint int

	// Create 24/32-bit int and clear/set MSB of address for 6/8-letter codes
	if compare == nil {
		codeLen = 6
		address &= 0x7fff
		bigint = (address << 8) | replace
	} else {
		codeLen = 8
		address |= 0x8000
		bigint = (address << 16) | (replace << 8) | *compare
	}

	// Convert into 4-bit ints
	encoded := make([]int, codeLen)
	loPosOrder := []int{3, 5, 2, 4, 1, 0, 7, 6}
	for i := codeLen - 1; i >= 0; i-- {
		loPos := loPosOrder[i]
		hiPos := (loPos - 1 + codeLen) % codeLen
		encoded[loPos] |= bigint & 0b111
		encoded[hiPos] |= bigint & 0b1000
		bigint >>= 4
	}

	// Convert into letters
	var result strings.Builder
	result.Grow(codeLen)
	const lookup = "APZLGITYEOXUKSVN"
	for _, val := range encoded {
		if val < 0 || val >= len(lookup) {
			return "", ErrOutOfRange
		}
		result.WriteByte(lookup[val])
	}

	return result.String(), nil
}
