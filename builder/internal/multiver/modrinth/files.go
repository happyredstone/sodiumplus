package modrinth

import (
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/spf13/viper"
)

func GetProjectTypeFolder(projectType string, fileLoaders []string, packLoaders []string) (string, error) {
	if projectType == "modpack" {
		return "", errors.New("this command should not be used to add Modrinth modpacks, and importing of Modrinth modpacks is not yet supported")
	} else if projectType == "resourcepack" {
		return "resourcepacks", nil
	} else if projectType == "shader" {
		bestLoaderIdx := math.MaxInt

		for _, v := range fileLoaders {
			idx := slices.Index(LoaderPreferenceList, v)

			if idx != -1 && idx < bestLoaderIdx {
				bestLoaderIdx = idx
			}
		}

		if bestLoaderIdx > -1 && bestLoaderIdx < math.MaxInt {
			return LoaderFolders[LoaderPreferenceList[bestLoaderIdx]], nil
		}

		return "shaderpacks", nil
	} else if projectType == "mod" {
		bestLoaderIdx := math.MaxInt

		for _, v := range fileLoaders {
			if slices.Contains(packLoaders, v) {
				idx := slices.Index(LoaderPreferenceList, v)
				if idx != -1 && idx < bestLoaderIdx {
					bestLoaderIdx = idx
				}
			}
		}
		if bestLoaderIdx > -1 && bestLoaderIdx < math.MaxInt {
			return LoaderFolders[LoaderPreferenceList[bestLoaderIdx]], nil
		}

		if slices.Contains(fileLoaders, "datapack") {
			if viper.GetString("datapack-folder") != "" {
				return viper.GetString("datapack-folder"), nil
			} else {
				return "", errors.New("set the datapack-folder option to use datapacks")
			}
		}

		return "mods", nil
	} else {
		return "", fmt.Errorf("unknown project type %s", projectType)
	}
}
