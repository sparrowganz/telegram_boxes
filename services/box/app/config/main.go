package config

import (
	"encoding/json"
	"os"
)

type Config interface {
}

type ConfigData struct {
}

func CreateConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	d := &ConfigData{}

	err = json.NewDecoder(f).Decode(d)
	if err != nil {
		return nil, err
	}

	return d, nil
}
