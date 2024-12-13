package curseforge

import (
	"path/filepath"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/viper"
)

var DefaultFolders = map[uint32]map[uint32]string{
	432: {
		5:  "plugins",
		12: "resourcepacks",
		6:  "mods",
		17: "saves",
	},
}

func CreateModFile(modInfo curseforge.ModInfo, fileInfo curseforge.ModFileInfo, rootPath string, optionalDisabled bool) error {
	updateMap := make(map[string]map[string]interface{})

	var err error

	updateMap["curseforge"], err = curseforge.CurseUpdateData{
		ProjectID: modInfo.ID,
		FileID:    fileInfo.ID,
	}.ToMap()

	if err != nil {
		return err
	}

	hash, hashFormat := fileInfo.GetBestHash()

	var optional *core.ModOption

	if optionalDisabled {
		optional = &core.ModOption{
			Optional: true,
			Default:  false,
		}
	}

	modMeta := core.Mod{
		Name:     modInfo.Name,
		FileName: fileInfo.FileName,
		Side:     core.UniversalSide,
		Download: core.ModDownload{
			HashFormat: hashFormat,
			Hash:       hash,
			Mode:       core.ModeCF,
		},
		Option: optional,
		Update: updateMap,
	}

	modMeta.SetMetaPath(GetPathForFile(modInfo.GameID, modInfo.ClassID, modInfo.PrimaryCategoryID, modInfo.Slug))

	_, _, err = modMeta.Write()

	if err != nil {
		return err
	}

	return nil
}

func GetPathForFile(gameID uint32, classID uint32, categoryID uint32, slug string) string {
	metaFolder := viper.GetString("meta-folder")

	if metaFolder == "" {
		if m, ok := DefaultFolders[gameID]; ok {
			if folder, ok := m[classID]; ok {
				return filepath.Join(viper.GetString("meta-folder-base"), folder, slug+core.MetaExtension)
			} else if folder, ok := m[categoryID]; ok {
				return filepath.Join(viper.GetString("meta-folder-base"), folder, slug+core.MetaExtension)
			}
		}

		metaFolder = "."
	}

	return filepath.Join(viper.GetString("meta-folder-base"), metaFolder, slug+core.MetaExtension)
}
