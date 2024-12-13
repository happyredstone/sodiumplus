package curseforge

import (
	"errors"
	"fmt"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"
)

func GetLatestFile(modInfoData curseforge.ModInfo, mcVersions []string, fileID uint32, packLoaders []string) (curseforge.ModFileInfo, error) {
	if fileID == 0 {
		if len(modInfoData.LatestFiles) == 0 && len(modInfoData.GameVersionLatestFiles) == 0 {
			return curseforge.ModFileInfo{}, fmt.Errorf("addon %d has no files", modInfoData.ID)
		}

		var fileInfoData *curseforge.ModFileInfo

		fileID, fileInfoData, _ = curseforge.FindLatestFile(modInfoData, mcVersions, packLoaders)

		if fileInfoData != nil {
			return *fileInfoData, nil
		}

		if fileID == 0 {
			return curseforge.ModFileInfo{}, errors.New("mod not available for the configured Minecraft version or loader")
		}
	}

	fileInfoData, err := curseforge.CurseDefaultClient.GetFileInfo(modInfoData.ID, fileID)

	if err != nil {
		return curseforge.ModFileInfo{}, err
	}

	return fileInfoData, nil
}
