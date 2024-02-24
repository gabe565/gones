package config

import (
	"github.com/pelletier/go-toml/v2"
)

type TOMLParser struct{}

func (p TOMLParser) Unmarshal(b []byte) (map[string]interface{}, error) {
	var data map[string]any
	err := toml.Unmarshal(b, &data)
	return data, err
}

func (p TOMLParser) Marshal(o map[string]interface{}) ([]byte, error) {
	out, err := toml.Marshal(o)
	return out, err
}
