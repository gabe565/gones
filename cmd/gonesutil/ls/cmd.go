package ls

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/gabe565/gones/internal/cartridge"
	log "github.com/sirupsen/logrus"
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
			return []string{"table", "json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
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
		slices.SortFunc(carts, sortFunc(field))
	}

	if reverse, err := cmd.Flags().GetBool("reverse"); err != nil {
		return err
	} else if reverse {
		slices.Reverse(carts)
	}

	format, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	if err := printEntries(cmd.OutOrStdout(), carts, Format(format)); err != nil {
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
		carts = slices.DeleteFunc(carts, deleteFunc(filters))
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
		stat, err := os.Stat(path)
		if err != nil {
			log.Error(err)
			continue
		}

		if stat.IsDir() {
			if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() || filepath.Ext(path) != ".nes" {
					return nil
				}

				wg.Add(1)
				go func() {
					defer wg.Done()

					cart, err := cartridge.FromiNesFile(path)
					if err != nil {
						log.WithError(err).WithField("path", path).Error("invalid ROM")
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
				log.Error(err)
				continue
			}
		} else {
			cart, err := cartridge.FromiNesFile(path)
			if err != nil {
				log.Error(err)
				continue
			}

			carts = append(carts, newEntry(path, cart))
		}
	}
	wg.Wait()
	return carts, failed
}

func sortFunc(field string) func(a, b *entry) int {
	field = strings.ToLower(field)
	return func(a, b *entry) int {
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
			log.WithField("field", field).Fatal("invalid sort field")
		}
		return 0
	}
}

func deleteFunc(filters map[string]string) func(e *entry) bool {
	return func(e *entry) bool {
		for field, filter := range filters {
			switch strings.ToLower(field) {
			case NameField:
				return !strings.Contains(strings.ToLower(e.Name), strings.ToLower(filter))
			case MapperField:
				parsed, err := strconv.ParseUint(filter, 10, 8)
				if err != nil {
					log.WithError(err).Fatal("invalid mapper filter value")
				}

				return byte(parsed) != e.Mapper
			case MirrorField:
				return !strings.Contains(strings.ToLower(e.Mirror), strings.ToLower(filter))
			case BatteryField:
				parsed, err := strconv.ParseBool(filter)
				if err != nil {
					log.WithError(err).Fatal("invalid battery filter value")
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
