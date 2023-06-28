package helpers

import (
	"errors"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/viper"
)

func LoadIndexFile(indexFile string) (Index, error) {
	var rep IndexTomlRepresentation

	if _, err := toml.DecodeFile(indexFile, &rep); err != nil {
		return Index{}, err
	}

	if len(rep.HashFormat) == 0 {
		rep.HashFormat = "sha256"
	}

	index := Index{
		HashFormat: rep.HashFormat,
		Files:      rep.Files.ToMemoryRep(),
		IndexFile:  indexFile,
		PackRoot:   filepath.Dir(indexFile),
	}

	return index, nil
}

func LoadIndex(pack core.Pack) (Index, error) {
	if filepath.IsAbs(pack.Index.File) {
		return LoadIndexFile(pack.Index.File)
	}

	fileNative := filepath.FromSlash(pack.Index.File)

	return LoadIndexFile(filepath.Join(filepath.Dir(viper.GetString("pack-file")), fileNative))
}

func LoadMod(modFile string) (Mod, error) {
	var mod Mod

	if _, err := toml.DecodeFile(modFile, &mod); err != nil {
		return Mod{}, err
	}

	mod.UpdateData = make(map[string]interface{})

	for k, v := range mod.Update {
		updater, ok := core.Updaters[k]

		if ok {
			updateData, err := updater.ParseUpdate(v)

			if err != nil {
				return mod, err
			}

			mod.UpdateData[k] = updateData
		} else {
			return mod, errors.New("Update plugin " + k + " not found!")
		}
	}

	mod.MetaFile = modFile

	return mod, nil
}
