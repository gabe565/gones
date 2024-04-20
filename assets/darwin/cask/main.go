package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	flag "github.com/spf13/pflag"
)

//go:embed gones.rb.tmpl
var spec string

type SpecVars struct {
	Path    string
	Version string
	SHA256  string
}

func main() {
	var values SpecVars
	flag.StringVar(&values.Path, "path", "", "Binary path")
	flag.StringVar(&values.Version, "version", "", "Version")
	flag.Parse()

	tmpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Parse(spec)
	if err != nil {
		panic(err)
	}

	binary, err := os.Open(values.Path)
	if err != nil {
		panic(err)
	}

	h := sha256.New()
	if _, err := io.Copy(h, binary); err != nil {
		panic(err)
	}
	_ = binary.Close()
	values.SHA256 = hex.EncodeToString(h.Sum(nil))

	var buf strings.Builder
	if err := tmpl.Execute(&buf, values); err != nil {
		panic(err)
	}

	_, _ = io.WriteString(os.Stdout, buf.String())
}
