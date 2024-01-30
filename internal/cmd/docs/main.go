package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gabe565/gones/cmd/gones"
	gonesutil "github.com/gabe565/gones/cmd/gonesutil/root"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	output := "./docs"
	commands := []*cobra.Command{
		gones.New(),
		gonesutil.New(),
	}

	if err := os.RemoveAll(output); err != nil {
		log.Fatal(fmt.Errorf("failed to remove existing dir: %w", err))
	}

	if err := os.MkdirAll(output, 0o755); err != nil {
		log.Fatal(fmt.Errorf("failed to mkdir: %w", err))
	}

	for _, cmd := range commands {
		if err := doc.GenMarkdownTree(cmd, output); err != nil {
			log.Fatal(fmt.Errorf("failed to generate markdown: %w", err))
		}
	}
}
