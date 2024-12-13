package bootstrap

import (
	"errors"
	"slices"
	"strconv"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
	helpersCurseforge "github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver/curseforge"
	"github.com/briandowns/spinner"
)

func InstallCurseProject(mod config.Mod, rootPath string, stub multiver.PackStub, spinner *spinner.Spinner) error {
	if (!slices.Contains(mod.Versions, stub.Minecraft) && len(mod.Versions) > 0) || (!slices.Contains(mod.Loaders, stub.Loader) && len(mod.Loaders) > 0) {
		return nil
	}

	spinner.Suffix = " Fetch info: " + mod.Id

	var modInfoData helpersCurseforge.ModInfo

	id, err := strconv.Atoi(mod.Id)

	if err != nil && id != 0 {
		modInfoData, err = helpersCurseforge.CurseDefaultClient.GetModInfo(uint32(id))

		if err != nil {
			return err
		}
	} else {
		var cancelled bool

		cancelled, modInfoData, err = curseforge.SearchCurseforgeInternal(mod.Id, true, "minecraft", "", []string{stub.Minecraft}, curseforge.GetSearchLoaderType(stub))

		if err != nil {
			return err
		}

		if cancelled {
			return errors.New("search cancelled")
		}
	}

	spinner.Suffix = " Find latest version: " + modInfoData.Name

	var fileInfoData helpersCurseforge.ModFileInfo

	fileInfoData, err = curseforge.GetLatestFile(modInfoData, []string{stub.Minecraft}, 0, []string{stub.Loader})

	if err != nil {
		return err
	}

	spinner.Suffix = " Adding " + modInfoData.Name + " and dependencies..."

	curseforge.InstallCurseVersionToStub(modInfoData, fileInfoData, rootPath, stub, spinner)

	spinner.Suffix = " Waiting..."

	return nil
}
