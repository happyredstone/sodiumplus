package list

import (
	"fmt"
	"time"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/cfwidget"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/modrinth"
	"github.com/briandowns/spinner"
)

type ModUrl struct {
	Mod string `json:"mod"`
	Url string `json:"url"`
}

func GetModUrl(mod *helpers.Mod) (string, error) {
	modrinthUpdate, mr := mod.GetParsedUpdateData("modrinth")
	curseforgeUpdate, cf := mod.GetParsedUpdateData("curseforge")

	if mr {
		data := modrinthUpdate.(modrinth.ModrinthUpdateData)
		fullData, err := modrinth.ModrinthDefaultClient.Projects.Get(data.ProjectID)

		if err != nil {
			return "", err
		}

		return "https://modrinth.com/mod/" + *fullData.Slug + "/version/" + data.InstalledVersion, nil
	} else if cf {
		data := curseforgeUpdate.(curseforge.CurseUpdateData)

		file, err := cfwidget.DefaultCFWidgetAPI.GetFile(int(data.ProjectID), int(data.FileID))

		if err != nil {
			return "", err
		}

		return file.Url, nil
	}

	return "", nil
}

func CreateModList() ([]ModUrl, error) {
	progress := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	progress.Start()

	progress.Suffix = " Loading pack..."

	_, index, err := helpers.GetPack()

	if err != nil {
		return nil, err
	}

	progress.Suffix = " Loading mods..."

	mods, err := index.LoadAllMods()

	if err != nil {
		return nil, err
	}

	progress.Suffix = " Loading mod URLs..."

	urls := []ModUrl{}

	for _, mod := range mods {
		progress.Suffix = fmt.Sprintf(" Loading mod %s (%s)...", mod.Name, mod.FileName)

		url, err := GetModUrl(mod)

		if err != nil {
			return nil, err
		}

		urls = append(urls, ModUrl{
			Mod: mod.Name,
			Url: url,
		})
	}

	progress.Stop()

	filtered := []ModUrl{}

	for _, url := range urls {
		if url.Url != "" {
			filtered = append(filtered, url)
		}
	}

	return filtered, nil
}
