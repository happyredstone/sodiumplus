package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/viper"
)

type Mod struct {
	// The ID or slug to resolve the mod from.
	Id string `toml:"id"`

	// The platform to resolve this mod from.
	Platform string `toml:"platform"`

	// A list of Minecraft versions this mod should try to be installed on.
	// If empty or missing, it will try to install on all.
	Versions []string `toml:"versions"`

	// A list of mod loaders this mod should try to be installed on.
	// If empty or missing, it will try to install on all.
	Loaders []string `toml:"loaders"`
}

type Config struct {
	Server struct {
		Folders []string `toml:"folders"`
		Files   []string `toml:"files"`
		Enabled bool     `toml:"enabled"`
	} `toml:"server"`

	MultiVersion struct {
		PackName      string `toml:"pack_name"`
		PackAuthor    string `toml:"pack_author"`
		PackVersion   string `toml:"pack_version"`
		Enabled       bool   `toml:"enabled"`
		BootstrapMods []Mod  `toml:"bootstrap_mods"`
		CommonDir     string `toml:"common"`
		TempDir       string `toml:"temp"`
	} `toml:"multi_version"`
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

func (cfg Config) Write() error {
	file := viper.GetString("config-file")
	f, err := os.Create(file)

	if err != nil {
		return err
	}

	enc := toml.NewEncoder(f)
	enc.Indent = ""
	err = enc.Encode(cfg)

	if err != nil {
		f.Close()
		return err
	}

	return f.Close()
}
