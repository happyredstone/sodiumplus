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

func Rename(outDir string, pack *core.Pack) error {
	helpers.Setup()

	version := pack.Version
	minecraft := pack.Versions["minecraft"]
	pack_name := pack.Name
	loader := pack.GetLoaders()[0]

	version_str := "v" + version + "+" + loader + "-" + minecraft

	if !helpers.Exists(outDir) {
		err := os.MkdirAll(outDir, 0755)

		if err != nil {
			return err
		}
	}

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

		err = os.Rename(input, outDir+"/"+output_str)

		if err != nil {
			return err
		}
	}

	return nil
}
