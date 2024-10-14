package ls

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"

	"gabe565.com/gones/internal/cartridge"
	"gabe565.com/gones/internal/log"
	"github.com/spf13/cobra"
)

const (
	PathField    = "path"
	NameField    = "name"
	MapperField  = "mapper"
	BatteryField = "battery"
	MirrorField  = "mirror"
	HashField    = "hash"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [path...]",
		Short:   "List ROM files and metadata",
		Aliases: []string{"list"},
		RunE:    run,

		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"nes"}, cobra.ShellCompDirectiveFilterFileExt
		},
	}
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: (table, json, yaml)")
	if err := cmd.RegisterFlagCompletionFunc("output",
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return OutputFormatStrings(), cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}

	cmd.Flags().StringToStringP("filter", "f", map[string]string{}, "Filter by a field")
	if err := cmd.RegisterFlagCompletionFunc("filter", completeFilter); err != nil {
		panic(err)
	}

	cmd.Flags().StringP("sort", "s", PathField, "Sort by a field")
	if err := cmd.RegisterFlagCompletionFunc("sort",
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{PathField, NameField, MapperField, BatteryField, MirrorField}, cobra.ShellCompDirectiveNoFileComp
		},
	); err != nil {
		panic(err)
	}

	cmd.Flags().BoolP("reverse", "r", false, "Reverse the output")

	log.Init(os.Stderr)
	return cmd
}

var ErrInvalidROMs = errors.New("some ROMs were invalid")

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	carts, failed, err := loadCarts(cmd, args)
	if err != nil {
		return err
	}

	if field, err := cmd.Flags().GetString("sort"); err != nil {
		return err
	} else if field != "" {
		errCh := make(chan error, 1)
		slices.SortFunc(carts, sortFunc(field, errCh))
		if len(errCh) != 0 {
			return <-errCh
		}
	}

	if reverse, err := cmd.Flags().GetBool("reverse"); err != nil {
		return err
	} else if reverse {
		slices.Reverse(carts)
	}

	formatSrc, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	format, err := OutputFormatString(formatSrc)
	if err != nil {
		return err
	}

	if err := printEntries(cmd.OutOrStdout(), carts, format); err != nil {
		return err
	}

	if failed {
		return ErrInvalidROMs
	}
	return nil
}

func loadCarts(cmd *cobra.Command, args []string) ([]*entry, bool, error) {
	var failed bool
	carts, failed := loadPaths(args)

	if filters, err := cmd.Flags().GetStringToString("filter"); err != nil {
		return carts, failed, err
	} else if len(filters) != 0 {
		errCh := make(chan error, 1)
		carts = slices.DeleteFunc(carts, deleteFunc(filters, errCh))
		if len(errCh) != 0 {
			return nil, true, <-errCh
		}
	}

	return carts, failed, nil
}

func loadPaths(paths []string) ([]*entry, bool) {
	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	carts := make([]*entry, 0, len(paths))
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	var failed bool
	for _, path := range paths {
		if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}

			ext := filepath.Ext(path)
			if !strings.EqualFold(ext, ".nes") {
				return err
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				cart, err := cartridge.FromiNesFile(path)
				if err != nil {
					slog.Error("Invalid ROM", "path", path, "error", err)
					failed = true
					return
				}

				entry := newEntry(path, cart)
				mu.Lock()
				carts = append(carts, entry)
				mu.Unlock()
			}()
			return nil
		}); err != nil {
			slog.Error("Failed to load ROMs", "error", err)
			continue
		}
	}
	wg.Wait()
	return carts, failed
}

var ErrUnknownSortField = errors.New("unknown sort field")

func sortFunc(field string, errCh chan error) func(a, b *entry) int {
	field = strings.ToLower(field)
	return func(a, b *entry) int {
		if len(errCh) != 0 {
			return 0
		}

		switch field {
		case PathField:
			return strings.Compare(a.Path, b.Path)
		case NameField:
			return strings.Compare(a.Name, b.Name)
		case MapperField:
			return int(a.Mapper) - int(b.Mapper)
		case BatteryField:
			if a.Battery && b.Battery {
				return 0
			}
			if a.Battery && !b.Battery {
				return 1
			}
			return -1
		case MirrorField:
			return strings.Compare(a.Mirror, b.Mirror)
		default:
			errCh <- fmt.Errorf("%w: %s", ErrUnknownSortField, field)
			return 0
		}
	}
}

func deleteFunc(filters map[string]string, errCh chan error) func(e *entry) bool {
	return func(e *entry) bool {
		if len(errCh) != 0 {
			return false
		}

		for field, filter := range filters {
			switch strings.ToLower(field) {
			case NameField:
				return !strings.Contains(strings.ToLower(e.Name), strings.ToLower(filter))
			case MapperField:
				parsed, err := strconv.ParseUint(filter, 10, 8)
				if err != nil {
					errCh <- fmt.Errorf("invalid mapper filter value: %w", err)
					return false
				}

				return byte(parsed) != e.Mapper
			case MirrorField:
				return !strings.Contains(strings.ToLower(e.Mirror), strings.ToLower(filter))
			case BatteryField:
				parsed, err := strconv.ParseBool(filter)
				if err != nil {
					errCh <- fmt.Errorf("invalid battery filter value: %w", err)
					return false
				}

				return parsed != e.Battery
			case HashField:
				return filter != e.Hash
			}
		}
		return false
	}
}

func completeFilter(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	defaults := []string{"name=", "mapper=", "mirror=", "battery=", "hash="}
	if !strings.Contains(toComplete, "=") {
		return defaults, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	}

	carts, _, _ := loadCarts(cmd, args)
	matches := make([]string, 0, len(carts))
	param, _, _ := strings.Cut(toComplete, "=")
	for _, cart := range carts {
		switch param {
		case NameField:
			matches = append(matches, param+"="+cart.Name)
		case MapperField:
			matches = append(matches, param+"="+strconv.Itoa(int(cart.Mapper)))
		case MirrorField:
			matches = append(matches, param+"="+cart.Mirror)
		case BatteryField:
			matches = append(matches, param+"="+strconv.FormatBool(cart.Battery))
		case HashField:
			matches = append(matches, param+"="+cart.Hash+"\t"+cart.Name)
		}
	}

	if len(matches) == 0 {
		return defaults, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	}
	return matches, cobra.ShellCompDirectiveNoFileComp
}
