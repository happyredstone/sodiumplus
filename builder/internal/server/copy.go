package server

import (
	"os"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers"
	"github.com/otiai10/copy"
	"github.com/spf13/viper"
)

var (
	copiedServerFolders = []string{
		"mods",
		"config",
		"kubejs",
		"openloader",
		"resourcepacks",
	}

	copiedServerFiles = []string{
		"pack.toml",
		"index.toml",
	}
)

func GetServerFolder() string {
	return viper.GetString("server-folder")
}

func CleanCopiedFiles() error {
	cfg, err := config.GetConfig()

	if err != nil {
		return err
	}

	folders := append(cfg.Server.Folders, copiedServerFolders...)
	files := append(cfg.Server.Files, copiedServerFiles...)

	folders = helpers.Dedupe[string](folders)
	files = helpers.Dedupe[string](files)

	for _, folder := range append(folders, files...) {
		if !helpers.Exists(GetServerFolder() + "/" + folder) {
			continue
		}

		dir := GetServerFolder() + "/" + folder
		err := os.RemoveAll(dir)

		if err != nil {
			return err
		}
	}

	return nil
}

func CopyServerFiles() error {
	CleanCopiedFiles()

	cfg, err := config.GetConfig()

	if err != nil {
		return err
	}

	folders := append(cfg.Server.Folders, copiedServerFolders...)
	files := append(cfg.Server.Files, copiedServerFiles...)

	folders = helpers.Dedupe[string](folders)
	files = helpers.Dedupe[string](files)

	for _, dir := range append(copiedServerFolders, copiedServerFiles...) {
		if !helpers.Exists(dir) {
			continue
		}

		out := GetServerFolder() + "/" + dir
		err := copy.Copy(dir, out)

		if err != nil {
			return err
		}
	}

	return nil
}
