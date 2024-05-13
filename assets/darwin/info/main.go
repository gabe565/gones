package main

import (
	_ "embed"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	flag "github.com/spf13/pflag"
)

//go:embed info.plist.tmpl
var spec string

type SpecVars struct {
	Version string
}

func main() {
	var values SpecVars
	flag.StringVar(&values.Version, "version", "", "Version")
	flag.Parse()

	tmpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(spec)
	if err != nil {
		panic(err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, values); err != nil {
		panic(err)
	}

	_, _ = io.WriteString(os.Stdout, buf.String())
}
