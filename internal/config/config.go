package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Symbols []string `yaml:"symbols"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
