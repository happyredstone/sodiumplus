package modrinth

import (
	"errors"
	"fmt"
	"strings"

	modrinthApi "codeberg.org/jmansfield/go-modrinth/modrinth"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/modrinth"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/viper"
	"github.com/unascribed/FlexVer/go/flexver"
)

func GetLatestVersion(projectID string, name string, stub multiver.PackStub) (*modrinthApi.Version, error) {
	var loaders []string

	if viper.GetString("datapack-folder") != "" {
		loaders = append([]string{strings.ToLower(stub.Loader)}, WithDatapackPathMRLoaders...)
	} else {
		loaders = append([]string{strings.ToLower(stub.Loader)}, DefaultMRLoaders...)
	}

	result, err := modrinth.ModrinthDefaultClient.Versions.ListVersions(projectID, modrinthApi.ListVersionsOptions{
		GameVersions: []string{stub.Minecraft},
		Loaders:      loaders,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest version: %w", err)
	}

	if len(result) == 0 {
		return nil, errors.New("no valid versions found")
	}

	flexverLatest := FindLatestVersion(result, stub.Minecraft, true)
	releaseDateLatest := FindLatestVersion(result, stub.Minecraft, false)

	if flexverLatest != releaseDateLatest && releaseDateLatest.VersionNumber != nil && flexverLatest.VersionNumber != nil {
		// fmt.Printf("Warning: Modrinth versions for %s inconsistent between latest version number and newest release date (%s vs %s)\n", name, *flexverLatest.VersionNumber, *releaseDateLatest.VersionNumber)
	}

	return releaseDateLatest, nil
}

func FindLatestVersion(versions []*modrinthApi.Version, gameVersion string, useFlexVer bool) *modrinthApi.Version {
	latestValidVersion := versions[0]
	bestGameVersion := core.HighestSliceIndex([]string{gameVersion}, versions[0].GameVersions)

	for _, v := range versions[1:] {
		gameVersionIdx := core.HighestSliceIndex([]string{gameVersion}, v.GameVersions)

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
