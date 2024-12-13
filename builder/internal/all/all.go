package all

import (
	"os"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/client"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/server"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/viper"
)

func Clean() error {
	pack, err := core.LoadPack()

	if err != nil {
		return err
	}

	err = client.Clean(&pack)

	if err != nil {
		return err
	}

	err = server.Clean("output", &pack)

	if err != nil {
		return err
	}

	return nil
}

func CleanFor(outDir string, pack *core.Pack) error {
	err := client.Clean(pack)

	if err != nil {
		return err
	}

	err = server.Clean(outDir, pack)

	if err != nil {
		return err
	}

	return nil
}

func Bundle() error {
	out_dir := "output"

	if helpers.Exists(out_dir) {
		err := os.RemoveAll(out_dir)

		if err != nil {
			return err
		}
	}

	pack, err := core.LoadPack()

	if err != nil {
		return err
	}

	err = client.Bundle("output", &pack)

	if err != nil {
		return err
	}

	cfg, err := config.GetConfig()

	if err != nil {
		return err
	}

	if cfg.Server.Enabled {
		err = server.Bundle("output", &pack)

		if err != nil {
			return err
		}
	}

	return nil
}

func BundleFor(outDir string, root string, cfg *config.Config, pack *core.Pack) error {
	err := client.Bundle(outDir, pack)

	if err != nil {
		return err
	}

	if cfg.Server.Enabled {
		viper.Set("server-folder", root+"/server")

		err = server.Bundle(outDir, pack)

		if err != nil {
			return err
		}
	}

	return nil
}
