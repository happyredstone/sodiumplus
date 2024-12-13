package version

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/mcmeta"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver/bootstrap"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/cobra"
)

var BasePath string
var Force bool
var UseLatest bool
var InitCommand = CreateInitCommand()

func CreateInitCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "init ([minecraft version] [loader])...",
		Short:   "Initialize a project for a new version to support.",
		Long:    `Initialize a project for a new version to support.`,
		Aliases: []string{"new", "n"},
		Args:    cobra.MinimumNArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			if (len(args) % 2) != 0 {
				return fmt.Errorf("arguments length must be a multiple of 2, got %d", len(args))
			}

			versions, err := mcmeta.FetchVersionList()

			if err != nil {
				return err
			}

			for i := range len(args) / 2 {
				fmt.Printf("Trying to init for %s %s...\n", args[(i*2)+1], args[i*2])

				err = InitFor(versions, args[i*2], args[(i*2)+1])

				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	cmd.Flags().StringVarP(&BasePath, "base", "B", cwd+"/versions", "The base path to create projects in.")
	cmd.Flags().BoolVarP(&Force, "force", "f", false, "Force continuation even if files exist.")
	cmd.Flags().BoolVarP(&UseLatest, "use-latest", "L", false, "Always use the latest loader version.")

	return &cmd
}

func InitFor(versions []string, ver string, loader string) error {
	var err error

	found_ver := slices.Contains(versions, ver)
	found_loader := false

	// Use custom logic so we can properly capitalize the loader name
	for _, loader_id := range multiver.Loaders {
		if strings.EqualFold(loader_id, loader) {
			found_loader = true
			loader = loader_id
			break
		}
	}

	if !found_ver {
		return fmt.Errorf("invalid Minecraft version: %s", ver)
	}

	if !found_loader {
		return fmt.Errorf("invalid mod loader: %s", loader)
	}

	dir := path.Join(BasePath, ver, loader)

	if helpers.Exists(dir) {
		var res string

		if Force {
			res = "Yes"
		} else {
			prompt := promptui.Select{
				Label: "Folder " + dir + " already exists. Overwrite?",
				Items: []string{"Yes", "No"},
			}

			_, res, err = prompt.Run()

			if err != nil {
				return err
			}
		}

		if res != "Yes" {
			return nil
		}

		err = os.RemoveAll(dir)

		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(dir, 0755)

	if err != nil {
		return err
	}

	loader_vers, release, err := core.ModLoaders[strings.ToLower(loader)].VersionListGetter(ver)

	if err != nil {
		return err
	}

	slices.Reverse(loader_vers)

	release_index := 0

	for i, item := range loader_vers {
		if item == release {
			release_index = i
			break
		}
	}

	var res string

	if UseLatest {
		res = release
	} else {
		prompt := promptui.Select{
			Label:     "Select the " + loader + " version to use.",
			Items:     loader_vers,
			CursorPos: release_index,
			Searcher: func(input string, index int) bool {
				return strings.ContainsAny(strings.ToLower(loader_vers[index]), strings.ToLower(input))
			},
		}

		_, res, err = prompt.Run()

		if err != nil {
			return err
		}
	}

	pack := multiver.PackStub{
		Minecraft:     ver,
		Loader:        loader,
		LoaderVersion: res,
	}

	err = pack.Write(dir + "/pack.stub.toml")

	if err != nil {
		return err
	}

	cfg, err := config.GetConfig()

	if err != nil {
		panic(err)
	}

	fmt.Println("Bootstrapping mods...")

	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithFinalMSG("✓ Done!\n"))
	spinner.Start()

	missing := []struct {
		Mod config.Mod
		Err error
	}{}

	for _, proj := range cfg.MultiVersion.BootstrapMods {
		if strings.ToLower(proj.Platform) == "modrinth" || strings.ToLower(proj.Platform) == "mr" {
			err = bootstrap.InstallModrinthProject(proj, dir, pack, spinner)

			if err != nil {
				missing = append(missing, struct {
					Mod config.Mod
					Err error
				}{Mod: proj, Err: err})
				continue
			}
		} else if strings.ToLower(proj.Platform) == "curseforge" || strings.ToLower(proj.Platform) == "curse" || strings.ToLower(proj.Platform) == "cf" {
			err = bootstrap.InstallCurseProject(proj, dir, pack, spinner)

			if err != nil {
				missing = append(missing, struct {
					Mod config.Mod
					Err error
				}{Mod: proj, Err: err})
				continue
			}
		}
	}

	spinner.Stop()

	if len(missing) > 0 {
		str := "## Missing mods:\n\n"
		str2 := "Missing mods:\n\n"

		for _, mod := range missing {
			str += "- " + mod.Mod.Id + " (from " + mod.Mod.Platform + ")\n"
			str += "  Error: " + mod.Err.Error() + "\n"

			str2 += "- " + mod.Mod.Id + " (from " + mod.Mod.Platform + ")\n"
			str2 += "  Error: " + mod.Err.Error() + "\n"
		}

		fmt.Println(str2)
		os.WriteFile(dir+"/missing.md", []byte(str), 0755)
	}

	return nil
}
