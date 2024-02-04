package ls

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/gabe565/gones/internal/cartridge"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [path...]",
		Short:   "List ROM files and metadata",
		Aliases: []string{"list"},
		RunE:    run,

		ValidArgsFunction: cobra.FixedCompletions([]string{"nes"}, cobra.ShellCompDirectiveFilterFileExt),
	}
	cmd.Flags().StringP("output", "o", "table", "Output format. One of: (table, json, yaml)")
	_ = cmd.RegisterFlagCompletionFunc(
		"output",
		cobra.FixedCompletions([]string{"table", "json", "yaml"}, cobra.ShellCompDirectiveNoFileComp),
	)

	cmd.Flags().StringToStringP("filter", "f", map[string]string{}, "Filter by a field")
	_ = cmd.RegisterFlagCompletionFunc("filter", completeFilter)

	cmd.Flags().StringP("sort", "s", "", "Sort by a field")
	_ = cmd.RegisterFlagCompletionFunc(
		"sort",
		cobra.FixedCompletions([]string{"path", "name", "mapper", "battery", "mirror"}, cobra.ShellCompDirectiveNoFileComp),
	)

	cmd.Flags().BoolP("reverse", "r", false, "Reverse the output")
	return cmd
}

func run(cmd *cobra.Command, args []string) (err error) {
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
		return errors.New("some ROMs were invalid")
	}
	return nil
}

func loadCarts(cmd *cobra.Command, args []string) ([]entry, bool, error) {
	var failed bool
	carts, failed := loadPaths(args)

	if filters, err := cmd.Flags().GetStringToString("filter"); err != nil {
		return carts, failed, err
	} else if len(filters) != 0 {
		carts = slices.DeleteFunc(carts, filterFunc(filters))
	}

	return carts, failed, nil
}

func loadPaths(paths []string) ([]entry, bool) {
	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	carts := make([]entry, 0, len(paths))
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

				cart, err := cartridge.FromiNesFile(path)
				if err != nil {
					log.WithError(err).WithField("path", path).Error("invalid ROM")
					failed = true
					return nil
				}

				carts = append(carts, newEntry(path, cart))
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
	return carts, failed
}

func sortFunc(field string) func(a, b entry) int {
	field = strings.ToLower(field)
	return func(a, b entry) int {
		switch field {
		case "path":
			return strings.Compare(a.Path, b.Path)
		case "name":
			return strings.Compare(a.Name, b.Name)
		case "mapper":
			return int(a.Mapper) - int(b.Mapper)
		case "battery":
			if a.Battery && b.Battery {
				return 0
			}
			if a.Battery && !b.Battery {
				return 1
			}
			return -1
		case "mirror":
			return strings.Compare(a.Mirror, b.Mirror)
		default:
			log.WithField("field", field).Fatal("invalid sort field")
		}
		return 0
	}
}

func filterFunc(filters map[string]string) func(e entry) bool {
	return func(e entry) bool {
		for field, filter := range filters {
			switch strings.ToLower(field) {
			case "name":
				if !strings.Contains(e.Name, filter) {
					return true
				}
			case "mapper":
				parsed, err := strconv.ParseUint(filter, 10, 8)
				if err != nil {
					log.WithError(err).Fatal("invalid mapper filter value")
				}

				if byte(parsed) != e.Mapper {
					return true
				}
			case "mirror":
				if !strings.Contains(strings.ToLower(e.Mirror), strings.ToLower(filter)) {
					return true
				}
			case "battery":
				parsed, err := strconv.ParseBool(filter)
				if err != nil {
					log.WithError(err).Fatal("invalid battery filter value")
				}

				if parsed != e.Battery {
					return true
				}
			}
		}
		return false
	}
}

func completeFilter(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	defaults := []string{"name=", "mapper=", "mirror=", "battery="}
	if !strings.Contains(toComplete, "=") {
		return defaults, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	}

	carts, _, _ := loadCarts(cmd, args)
	matches := make([]string, 0, len(carts))
	param, _, _ := strings.Cut(toComplete, "=")
	for _, cart := range carts {
		switch param {
		case "name":
			matches = append(matches, param+"="+cart.Name)
		case "mapper":
			matches = append(matches, param+"="+strconv.Itoa(int(cart.Mapper)))
		case "mirror":
			matches = append(matches, param+"="+cart.Mirror)
		case "battery":
			matches = append(matches, param+"="+strconv.FormatBool(cart.Battery))
		}
	}

	if len(matches) == 0 {
		return defaults, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	}
	return matches, cobra.ShellCompDirectiveNoFileComp
}
