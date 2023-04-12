package controller

import (
	"strings"

	"github.com/gabe565/gones/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"
)

type Keymap struct {
	Regular map[Button]ebiten.Key
	Turbo   map[Button]ebiten.Key
}

func NewKeymap(player Player) Keymap {
	keymapConf := config.K.StringMap("input.keys." + string(player))

	keymap := Keymap{
		Regular: make(map[Button]ebiten.Key),
		Turbo:   make(map[Button]ebiten.Key),
	}

	for buttonName, keyName := range keymapConf {
		var turbo bool
		if strings.HasSuffix(buttonName, "_turbo") {
			turbo = true
			buttonName = strings.TrimSuffix(buttonName, "_turbo")
		}

		var button Button
		if err := button.UnmarshalText([]byte(buttonName)); err != nil {
			log.Fatal(err)
		}

		var key ebiten.Key
		if err := key.UnmarshalText([]byte(keyName)); err != nil {
			log.Fatal(err)
		}

		if turbo {
			keymap.Turbo[button] = key
		} else {
			keymap.Regular[button] = key
		}
	}

	return keymap
}

func LoadKeys() {
	_ = Reset.UnmarshalText([]byte(config.K.String("input.keys.reset")))
	_ = SaveState1.UnmarshalText([]byte(config.K.String("input.keys.state1_save")))
	_ = LoadState1.UnmarshalText([]byte(config.K.String("input.keys.state1_load")))
	_ = FastForward.UnmarshalText([]byte(config.K.String("input.keys.fast_forward")))
	_ = ToggleFullscreen.UnmarshalText([]byte(config.K.String("input.keys.fullscreen")))
}

var (
	Reset = ebiten.KeyR

	SaveState1 = ebiten.KeyF1
	LoadState1 = ebiten.KeyF5

	FastForward      = ebiten.KeyF
	ToggleFullscreen = ebiten.KeyF11

	ToggleTrace = ebiten.KeyTab
	ToggleDebug = ebiten.KeyGraveAccent
	StepFrame   = ebiten.Key1
	RunToRender = ebiten.Key2
)
