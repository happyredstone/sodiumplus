package modrinth

import (
	"errors"
	"fmt"

	modrinthApi "codeberg.org/jmansfield/go-modrinth/modrinth"
	"github.com/mitchellh/mapstructure"
	"github.com/packwiz/packwiz/core"
)

type ModrinthUpdateData struct {
	ProjectID        string `mapstructure:"mod-id"`
	InstalledVersion string `mapstructure:"version"`
}

func (u ModrinthUpdateData) ToMap() (map[string]interface{}, error) {
	newMap := make(map[string]interface{})

	err := mapstructure.Decode(u, &newMap)

	return newMap, err
}

type ModrinthUpdater struct{}

func (u ModrinthUpdater) ParseUpdate(updateUnparsed map[string]interface{}) (interface{}, error) {
	var updateData ModrinthUpdateData

	err := mapstructure.Decode(updateUnparsed, &updateData)

	return updateData, err
}

type CachedStateStore struct {
	ProjectID string
	Version   *modrinthApi.Version
}

func (u ModrinthUpdater) CheckUpdate(mods []*core.Mod, pack core.Pack) ([]core.UpdateCheck, error) {
	results := make([]core.UpdateCheck, len(mods))

	for i, mod := range mods {
		rawData, ok := mod.GetParsedUpdateData("modrinth")

		if !ok {
			results[i] = core.UpdateCheck{Error: errors.New("failed to parse update metadata")}
			continue
		}

		data := rawData.(ModrinthUpdateData)

		newVersion, err := GetLatestVersion(data.ProjectID, mod.Name, pack)

		if err != nil {
			results[i] = core.UpdateCheck{Error: fmt.Errorf("failed to get latest version: %v", err)}
			continue
		}

		if *newVersion.ID == data.InstalledVersion {
			results[i] = core.UpdateCheck{UpdateAvailable: false}
			continue
		}

		if len(newVersion.Files) == 0 {
			results[i] = core.UpdateCheck{Error: errors.New("new version doesn't have any files")}
			continue
		}

		newFilename := newVersion.Files[0].Filename

		for _, v := range newVersion.Files {
			if *v.Primary {
				newFilename = v.Filename
			}
		}

		results[i] = core.UpdateCheck{
			UpdateAvailable: true,
			UpdateString:    mod.FileName + " -> " + *newFilename,
			CachedState:     CachedStateStore{data.ProjectID, newVersion},
		}
	}

	return results, nil
}

func (u ModrinthUpdater) DoUpdate(mods []*core.Mod, cachedState []interface{}) error {
	for i, mod := range mods {
		modState := cachedState[i].(CachedStateStore)

		var version = modState.Version
		var file = version.Files[0]

		for _, v := range version.Files {
			if *v.Primary {
				file = v
			}
		}

		algorithm, hash := GetBestHash(file)

		if algorithm == "" {
			return errors.New("file for project " + mod.Name + " doesn't have a valid hash")
		}

		mod.FileName = *file.Filename

		mod.Download = core.ModDownload{
			URL:        *file.URL,
			HashFormat: algorithm,
			Hash:       hash,
		}

		mod.Update["modrinth"]["version"] = version.ID
	}

	return nil
}
