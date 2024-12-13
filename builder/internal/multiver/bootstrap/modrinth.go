package bootstrap

import (
	"fmt"
	"slices"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
	helpersModrinth "github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/modrinth"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver/modrinth"
	"github.com/briandowns/spinner"
)

func InstallModrinthProject(mod config.Mod, rootPath string, stub multiver.PackStub, spinner *spinner.Spinner) error {
	if (!slices.Contains(mod.Versions, stub.Minecraft) && len(mod.Versions) > 0) || (!slices.Contains(mod.Loaders, stub.Loader) && len(mod.Loaders) > 0) {
		return nil
	}

	spinner.Suffix = " Fetch info: " + mod.Id

	proj, err := helpersModrinth.ModrinthDefaultClient.Projects.Get(mod.Id)

	if err != nil {
		return err
	}

	spinner.Suffix = " Find latest version: " + *proj.Title

	version, err := modrinth.GetLatestVersion(*proj.ID, *proj.Title, stub)

	if err != nil {
		return err
	}

	if version.ID == nil {
		spinner.Suffix = fmt.Sprintf("Mod not available for the specified Minecraft version or loader: %s", *proj.Title)
		return nil
	}

	spinner.Suffix = " Adding " + *proj.Title + " and dependencies..."

	modrinth.InstallModrinthVersionToStub(proj, version, rootPath, stub, spinner)

	spinner.Suffix = " Waiting..."

	return nil
}
