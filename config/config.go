package config

import (
	"encoding/json"
	"os"
)

var GlobalConfig *Config

type Config struct {
	TransientPath   string
	TransientPeriod int64
	StandingPath    string
	AdminUser       string
	AdminPassword   string
	DbPath          string
}

func FromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := new(Config)

	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
