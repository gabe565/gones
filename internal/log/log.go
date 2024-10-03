package log

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

func Init(out io.Writer) {
	var color bool
	if f, ok := out.(*os.File); ok {
		color = isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(out, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: time.Kitchen,
			NoColor:    !color,
		}),
	))
}
