package config

import (
	"time"
)

type Duration time.Duration

func (d Duration) MarshalText() ([]byte, error) {
	s := time.Duration(d).String()
	return []byte(s), nil
}

func (d *Duration) UnmarshalText(text []byte) error {
	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}

	*d = Duration(duration)
	return nil
}
