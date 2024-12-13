package server

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers"
	"github.com/packwiz/packwiz/core"
)

var (
	cleanFiles = []string{
		"server.tar.gz",
		"server.zip",
	}

	namedGlobs = []string{
		"%s (Server) v*+*.zip",
		"%s (Server) v*+*.tar.gz",
		"%s (Server) v*+*-*.zip",
		"%s (Server) v*+*-*.tar.gz",
	}
)

func Clean(outDir string, pack *core.Pack) error {
	err := CleanNamed(outDir, pack)

	if err != nil {
		return err
	}

	for _, file := range cleanFiles {
		if helpers.Exists(file) {
			fmt.Println("Removing " + file)

			err := os.Remove(file)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CleanNamed(outDir string, pack *core.Pack) error {
	for _, glob := range namedGlobs {
		matches, err := filepath.Glob(outDir + "/" + fmt.Sprintf(glob, pack.Name))

		if err != nil {
			return err
		}

		for _, match := range matches {
			fmt.Println("Removing " + match)

			err = os.Remove(match)

			if err != nil {
				return err
			}
		}
	}

	return nil
}
