package version

import (
	"fmt"
	"time"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/config"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/modrinth"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var MigrateCommand = CreateMigrateCommand()

func CreateMigrateCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:     "migrate",
		Short:   "Migrate from a standalone packwiz config to a multiversion/multiloader config.",
		Long:    `Migrate from a standalone packwiz config to a multiversion/multiloader config.`,
		Aliases: []string{"m"},

		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetConfig()

			if err != nil {
				panic(err)
			}

			progress := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			progress.Start()

			progress.Suffix = " Loading pack..."

			pack, index, err := helpers.GetPack()

			if err != nil {
				return err
			}

			progress.Suffix = " Loading mods..."

			mods, err := index.LoadAllMods()

			if err != nil {
				return err
			}

			progress.Suffix = " Getting mods data..."

			data := []config.Mod{}

			for _, mod := range mods {
				progress.Suffix = fmt.Sprintf(" Fetching mod %s (%s)...", mod.Name, mod.FileName)
				item, err := GetModData(mod)

				if err != nil {
					return err
				}

				data = append(data, item)
			}

			progress.Stop()

			cfg.MultiVersion.Enabled = true
			cfg.MultiVersion.PackName = pack.Name
			cfg.MultiVersion.PackAuthor = pack.Author
			cfg.MultiVersion.PackVersion = pack.Version
			cfg.MultiVersion.BootstrapMods = data
			cfg.MultiVersion.CommonDir = "common"
			cfg.MultiVersion.VersionsDir = "versions"
			cfg.MultiVersion.TempDir = ".build/temp"
			cfg.MultiVersion.OutDir = "output"

			cfg.Write()

			return nil
		},
	}

	return &cmd
}

func GetModData(mod *helpers.Mod) (config.Mod, error) {
	modrinthUpdate, mr := mod.GetParsedUpdateData("modrinth")
	curseforgeUpdate, cf := mod.GetParsedUpdateData("curseforge")

	if mr {
		data := modrinthUpdate.(modrinth.ModrinthUpdateData)
		fullData, err := modrinth.ModrinthDefaultClient.Projects.Get(data.ProjectID)

		if err != nil {
			return config.Mod{}, err
		}

		return config.Mod{
			Id:       *fullData.Slug,
			Platform: "Modrinth",
		}, nil
	} else if cf {
		data := curseforgeUpdate.(curseforge.CurseUpdateData)
		proj, err := curseforge.CurseDefaultClient.GetModInfo(data.ProjectID)

		if err != nil {
			return config.Mod{}, err
		}

		return config.Mod{
			Id:       proj.Slug,
			Platform: "CurseForge",
		}, nil
	}

	return config.Mod{}, nil
}
