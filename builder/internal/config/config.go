package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Folders []string `toml:"folders"`
		Files   []string `toml:"files"`
		Enabled bool     `toml:"enabled"`
	} `toml:"server"`
}

func GetConfig() (*Config, error) {
	file := viper.GetString("config-file")
	content, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	var cfg Config

	err = toml.Unmarshal(content, &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
