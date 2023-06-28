package modrinth

import (
	"errors"
	"fmt"

	modrinthApi "codeberg.org/jmansfield/go-modrinth/modrinth"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/viper"
	"github.com/unascribed/FlexVer/go/flexver"
)

func FindLatestVersion(versions []*modrinthApi.Version, gameVersions []string, useFlexVer bool) *modrinthApi.Version {
	latestValidVersion := versions[0]
	bestGameVersion := core.HighestSliceIndex(gameVersions, versions[0].GameVersions)

	for _, v := range versions[1:] {
		gameVersionIdx := core.HighestSliceIndex(gameVersions, v.GameVersions)

		var compare int32

		if useFlexVer {
			compare = flexver.Compare(*v.VersionNumber, *latestValidVersion.VersionNumber)
		}

		if compare == 0 {
			compare = int32(gameVersionIdx - bestGameVersion)
		}

		if compare == 0 {
			compare = CompareLoaderLists(latestValidVersion.Loaders, v.Loaders)
		}

		if compare == 0 {
			if v.DatePublished.After(*latestValidVersion.DatePublished) {
				compare = 1
			}
		}

		if compare > 0 {
			latestValidVersion = v
			bestGameVersion = gameVersionIdx
		}
	}

	return latestValidVersion
}

func GetLatestVersion(projectID string, name string, pack core.Pack) (*modrinthApi.Version, error) {
	gameVersions, err := pack.GetSupportedMCVersions()

	if err != nil {
		return nil, err
	}

	var loaders []string

	if viper.GetString("datapack-folder") != "" {
		loaders = append(pack.GetLoaders(), WithDatapackPathMRLoaders...)
	} else {
		loaders = append(pack.GetLoaders(), DefaultMRLoaders...)
	}

	result, err := ModrinthDefaultClient.Versions.ListVersions(projectID, modrinthApi.ListVersionsOptions{
		GameVersions: gameVersions,
		Loaders:      loaders,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest version: %w", err)
	}

	if len(result) == 0 {
		return nil, errors.New("no valid versions found\n\tUse the 'packwiz settings acceptable-versions' command to accept more game versions\n\tTo use datapacks, add a datapack loader mod and specify the datapack-folder option with the folder this mod loads datapacks from")
	}

	flexverLatest := FindLatestVersion(result, gameVersions, true)
	releaseDateLatest := FindLatestVersion(result, gameVersions, false)

	if flexverLatest != releaseDateLatest && releaseDateLatest.VersionNumber != nil && flexverLatest.VersionNumber != nil {
		fmt.Printf("Warning: Modrinth versions for %s inconsistent between latest version number and newest release date (%s vs %s)\n", name, *flexverLatest.VersionNumber, *releaseDateLatest.VersionNumber)
	}

	return releaseDateLatest, nil
}
