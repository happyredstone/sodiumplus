package curseforge

import (
	"errors"

	"github.com/mitchellh/mapstructure"
	"github.com/packwiz/packwiz/core"
)

type CurseUpdateData struct {
	ProjectID uint32 `mapstructure:"project-id"`
	FileID    uint32 `mapstructure:"file-id"`
}

func (u CurseUpdateData) ToMap() (map[string]interface{}, error) {
	newMap := make(map[string]interface{})

	err := mapstructure.Decode(u, &newMap)

	return newMap, err
}

type CurseUpdater struct{}

func (u CurseUpdater) ParseUpdate(updateUnparsed map[string]interface{}) (interface{}, error) {
	var updateData CurseUpdateData
	err := mapstructure.Decode(updateUnparsed, &updateData)
	return updateData, err
}

type CachedStateStore struct {
	ModInfo
	fileID   uint32
	fileInfo *ModFileInfo
}

func (u CurseUpdater) CheckUpdate(mods []*core.Mod, pack core.Pack) ([]core.UpdateCheck, error) {
	results := make([]core.UpdateCheck, len(mods))
	modIDs := make([]uint32, len(mods))
	modInfos := make([]ModInfo, len(mods))

	mcVersions, err := pack.GetSupportedMCVersions()

	if err != nil {
		return nil, err
	}

	for i, v := range mods {
		projectRaw, ok := v.GetParsedUpdateData("curseforge")

		if !ok {
			results[i] = core.UpdateCheck{Error: errors.New("failed to parse update metadata")}
			continue
		}

		project := projectRaw.(CurseUpdateData)
		modIDs[i] = project.ProjectID
	}

	modInfosUnsorted, err := CurseDefaultClient.GetModInfoMultiple(modIDs)

	if err != nil {
		return nil, err
	}

	for _, v := range modInfosUnsorted {
		for i, id := range modIDs {
			if id == v.ID {
				modInfos[i] = v
				break
			}
		}
	}

	packLoaders := pack.GetLoaders()

	for i, v := range mods {
		projectRaw, ok := v.GetParsedUpdateData("curseforge")

		if !ok {
			results[i] = core.UpdateCheck{Error: errors.New("failed to parse update metadata")}
			continue
		}

		project := projectRaw.(CurseUpdateData)

		fileID, fileInfoData, fileName := FindLatestFile(modInfos[i], mcVersions, packLoaders)

		if fileID != project.FileID && fileID != 0 {
			results[i] = core.UpdateCheck{
				UpdateAvailable: true,
				UpdateString:    v.FileName + " -> " + fileName,
				CachedState:     CachedStateStore{modInfos[i], fileID, fileInfoData},
			}
		} else {
			results[i] = core.UpdateCheck{UpdateAvailable: false}
			continue
		}
	}

	return results, nil
}

func (u CurseUpdater) DoUpdate(mods []*core.Mod, cachedState []interface{}) error {
	for i, v := range mods {
		modState := cachedState[i].(CachedStateStore)

		var fileInfoData ModFileInfo

		if modState.fileInfo != nil {
			fileInfoData = *modState.fileInfo
		} else {
			var err error

			fileInfoData, err = CurseDefaultClient.GetFileInfo(modState.ID, modState.fileID)

			if err != nil {
				return err
			}
		}

		v.FileName = fileInfoData.FileName
		v.Name = modState.Name
		hash, hashFormat := fileInfoData.GetBestHash()

		v.Download = core.ModDownload{
			HashFormat: hashFormat,
			Hash:       hash,
			Mode:       core.ModeCF,
		}

		v.Update["curseforge"]["project-id"] = modState.ID
		v.Update["curseforge"]["file-id"] = fileInfoData.ID
	}

	return nil
}
