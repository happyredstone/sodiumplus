package multiver

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers"
	"github.com/packwiz/packwiz/core"
)

var IgnoredFiles = []string{"pack.stub.toml", "missing.md"}

func (stub PackStub) InitRealPack(config config.Config) (*core.Pack, error) {
	tempDir := path.Join(config.MultiVersion.TempDir, stub.Minecraft, stub.Loader)
	realDir := path.Join("versions", stub.Minecraft, stub.Loader)

	if !helpers.Exists(realDir) {
		return nil, errors.New("pack directory does not exist")
	}

	if helpers.Exists(tempDir) {
		err := os.RemoveAll(tempDir)

		if err != nil {
			return nil, err
		}
	}

	err := os.MkdirAll(tempDir, 0755)

	if err != nil {
		return nil, err
	}

	fullOutPath, err := filepath.Abs(tempDir)

	if err != nil {
		return nil, err
	}

	fullRealPath, err := filepath.Abs(realDir)

	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(realDir)

	if err != nil {
		return nil, err
	}

	for _, item := range files {
		if slices.Contains(IgnoredFiles, item.Name()) {
			continue
		}

		fullPath, err := filepath.Abs(path.Join(realDir, item.Name()))

		if err != nil {
			return nil, err
		}

		outPath := strings.Replace(fullPath, fullRealPath, fullOutPath, 1)

		if item.Type().IsDir() {
			err = os.CopyFS(outPath, os.DirFS(fullPath))

			if err != nil {
				return nil, err
			}
		} else {
			err = helpers.Copy(fullPath, outPath)

			if err != nil {
				return nil, err
			}
		}
	}

	if helpers.Exists(config.MultiVersion.CommonDir) {
		fullCommonPath, err := filepath.Abs(config.MultiVersion.CommonDir)

		if err != nil {
			return nil, err
		}

		commonFiles, err := os.ReadDir(config.MultiVersion.CommonDir)

		if err != nil {
			return nil, err
		}

		for _, item := range commonFiles {
			if slices.Contains(IgnoredFiles, item.Name()) {
				continue
			}

			fullPath, err := filepath.Abs(path.Join(config.MultiVersion.CommonDir, item.Name()))

			if err != nil {
				return nil, err
			}

			outPath := strings.Replace(fullPath, fullCommonPath, fullOutPath, 1)

			if item.Type().IsDir() {
				err = os.CopyFS(outPath, os.DirFS(fullPath))

				if err != nil {
					return nil, err
				}
			} else {
				err = helpers.Copy(fullPath, outPath)

				if err != nil {
					return nil, err
				}
			}
		}
	}

	indexPath := path.Join(tempDir, "index.toml")
	err = os.WriteFile(indexPath, []byte{}, 0755)

	if err != nil {
		return nil, err
	}

	pack := core.Pack{
		Name:       config.MultiVersion.PackName,
		Author:     config.MultiVersion.PackAuthor,
		Version:    config.MultiVersion.PackVersion,
		PackFormat: core.CurrentPackFormat,

		Index: struct {
			File       string `toml:"file"`
			HashFormat string `toml:"hash-format"`
			Hash       string `toml:"hash,omitempty"`
		}{
			File: "index.toml",
		},

		Versions: map[string]string{
			"minecraft":                  stub.Minecraft,
			strings.ToLower(stub.Loader): stub.LoaderVersion,
		},
	}

	os.Chdir(tempDir)

	index, err := helpers.LoadIndex(pack)

	if err != nil {
		return nil, err
	}

	err = index.RefreshNoBar()

	if err != nil {
		return nil, err
	}

	err = index.Write()

	if err != nil {
		return nil, err
	}

	err = pack.UpdateIndexHash()

	if err != nil {
		return nil, err
	}

	err = pack.Write()

	if err != nil {
		return nil, err
	}

	return &pack, nil
}
