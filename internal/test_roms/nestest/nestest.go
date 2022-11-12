package nestest

import (
	_ "embed"
)

//go:embed nestest.nes
var ROM string

//go:embed nestest.log
var Log string
