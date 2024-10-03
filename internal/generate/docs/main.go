package main

import (
	"log/slog"
	"os"

	"github.com/gabe565/gones/cmd/gones"
	gonesutil "github.com/gabe565/gones/cmd/gonesutil/root"
	"github.com/gabe565/gones/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	log.Init(os.Stderr)

	output := "./docs"
	commands := []*cobra.Command{
		gones.New(),
		gonesutil.New(),
	}

	if err := os.RemoveAll(output); err != nil {
		slog.Error("Failed to remove existing dir", "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(output, 0o777); err != nil {
		slog.Error("Failed to create directory", "error", err)
		os.Exit(1)
	}

	for _, cmd := range commands {
		if err := doc.GenMarkdownTree(cmd, output); err != nil {
			slog.Error("Failed to generate markdown", "error", err)
		}
	}
}
