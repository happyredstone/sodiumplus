package curseforge

import (
	"github.com/packwiz/packwiz/core"
	"golang.org/x/exp/slices"
)

func FilterLoaderTypeIndex(packLoaders []string, modLoaderType ModloaderType) (ModloaderType, bool) {
	if len(packLoaders) == 0 || modLoaderType == ModloaderTypeAny {
		return ModloaderTypeAny, true
	} else {
		if slices.Contains(packLoaders, ModloaderIds[modLoaderType]) {
			return modLoaderType, true
		} else {
			return ModloaderTypeAny, false
		}
	}
}

func FilterFileInfoLoaderIndex(packLoaders []string, fileInfoData ModFileInfo) (ModloaderType, bool) {
	if len(packLoaders) == 0 {
		return ModloaderTypeAny, true
	} else {
		bestLoaderId := -1

		for i, name := range ModloaderNames {
			if slices.Contains(packLoaders, ModloaderIds[i]) && slices.Contains(fileInfoData.GameVersions, name) {
				if i > bestLoaderId {
					bestLoaderId = i
				}
			}
		}

		if bestLoaderId > -1 {
			return ModloaderType(bestLoaderId), true
		} else {
			return ModloaderTypeAny, false
		}
	}
}

func FindLatestFile(modInfoData ModInfo, mcVersions []string, packLoaders []string) (fileID uint32, fileInfoData *ModFileInfo, fileName string) {
	cfMcVersions := GetCurseforgeVersions(mcVersions)
	bestMcVer := -1
	bestLoaderType := ModloaderTypeAny

	for _, v := range modInfoData.LatestFiles {
		mcVerIdx := core.HighestSliceIndex(mcVersions, v.GameVersions)
		loaderIdx, loaderValid := FilterFileInfoLoaderIndex(packLoaders, v)

		if mcVerIdx < 0 || !loaderValid {
			continue
		}

		compare := int32(mcVerIdx - bestMcVer)

		if compare == 0 {
			if bestLoaderType == ModloaderTypeAny || loaderIdx == ModloaderTypeAny {
				compare = 0
			} else {
				compare = int32(loaderIdx) - int32(bestLoaderType)
			}
		}

		if compare == 0 {
			compare = int32(int64(v.ID) - int64(fileID))
		}

		if compare > 0 {
			fileID = v.ID
			fileInfoDataCopy := v
			fileInfoData = &fileInfoDataCopy
			fileName = v.FileName
			bestMcVer = mcVerIdx
			bestLoaderType = loaderIdx
		}
	}

	for _, v := range modInfoData.GameVersionLatestFiles {
		mcVerIdx := slices.Index(cfMcVersions, v.GameVersion)
		loaderIdx, loaderValid := FilterLoaderTypeIndex(packLoaders, v.Modloader)

		if mcVerIdx < 0 || !loaderValid {
			continue
		}

		compare := int32(mcVerIdx - bestMcVer)

		if compare == 0 {
			if bestLoaderType == ModloaderTypeAny || loaderIdx == ModloaderTypeAny {
				compare = 0
			} else {
				compare = int32(loaderIdx) - int32(bestLoaderType)
			}
		}

		if compare == 0 {
			compare = int32(int64(v.ID) - int64(fileID))
		}

		if compare > 0 {
			fileID = v.ID
			fileInfoData = nil
			fileName = v.Name
			bestMcVer = mcVerIdx
			bestLoaderType = loaderIdx
		}
	}

	return
}
