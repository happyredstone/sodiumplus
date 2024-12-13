package server

import (
	"fmt"
	"os"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers"
	"github.com/magefile/mage/target"
	"github.com/packwiz/packwiz/core"
)

var (
	inputMap = map[string]string{
		"server.zip":    "%s (Server) %s.zip",
		"server.tar.gz": "%s (Server) %s.tar.gz",
	}
)

func Rename() error {
	helpers.Setup()

	pack, err := core.LoadPack()

	if err != nil {
		return err
	}

	version := pack.Version
	minecraft := pack.Versions["minecraft"]
	pack_name := pack.Name

	version_str := "v" + version + "+" + minecraft

	for input, output := range inputMap {
		output_str := fmt.Sprintf(output, pack_name, version_str)

		if !helpers.Exists(input) {
			continue
		}

		newer, err := target.Path(output_str, input)

		if err != nil {
			return err
		}

		if !newer {
			continue
		}

		err = os.Rename(input, output_str)

		if err != nil {
			return err
		}
	}

	return nil
}
