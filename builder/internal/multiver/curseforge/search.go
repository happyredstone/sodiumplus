package curseforge

import (
	"errors"
	"fmt"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/viper"
	"gopkg.in/dixonwille/wmenu.v4"
)

func SearchCurseforgeInternal(searchTerm string, isSlug bool, game string, category string, mcVersions []string, searchLoaderType curseforge.ModloaderType) (bool, curseforge.ModInfo, error) {
	var gameID, categoryID, classID uint32
	if game == "minecraft" {
		gameID = 432
	}
	if category == "mc-mods" {
		classID = 6
	}
	if gameID == 0 {
		games, err := curseforge.CurseDefaultClient.GetGames()

		if err != nil {
			return false, curseforge.ModInfo{}, err
		}

		for _, v := range games {
			if v.Slug == game {
				if v.Status != curseforge.GameStatusLive {
					return false, curseforge.ModInfo{}, fmt.Errorf("failed to lookup game %s: selected game is not live", game)
				}

				if v.APIStatus != curseforge.GameApiStatusPublic {
					return false, curseforge.ModInfo{}, fmt.Errorf("failed to lookup game %s: selected game does not have a public API", game)
				}

				gameID = v.ID

				break
			}
		}

		if gameID == 0 {
			return false, curseforge.ModInfo{}, fmt.Errorf("failed to lookup: game %s could not be found", game)
		}
	}

	if categoryID == 0 && classID == 0 && category != "" {
		categories, err := curseforge.CurseDefaultClient.GetCategories(gameID)

		if err != nil {
			return false, curseforge.ModInfo{}, err
		}

		for _, v := range categories {
			if v.Slug == category {
				if v.IsClass {
					classID = v.ID
				} else {
					classID = v.ClassID
					categoryID = v.ID
				}

				break
			}
		}

		if categoryID == 0 && classID == 0 {
			return false, curseforge.ModInfo{}, fmt.Errorf("failed to lookup: category %s could not be found", category)
		}
	}

	filterGameVersion := ""

	if len(mcVersions) == 1 {
		filterGameVersion = curseforge.GetCurseforgeVersion(mcVersions[0])
	}

	var search, slug string

	if isSlug {
		slug = searchTerm
	} else {
		search = searchTerm
	}

	results, err := curseforge.CurseDefaultClient.GetSearch(search, slug, gameID, classID, categoryID, filterGameVersion, searchLoaderType)

	if err != nil {
		return false, curseforge.ModInfo{}, err
	}

	if len(results) == 0 {
		return false, curseforge.ModInfo{}, errors.New("no projects found")
	} else if len(results) == 1 {
		return false, results[0], nil
	} else {
		fuzzySearchResults := fuzzy.FindFrom(searchTerm, ModResultsList(results))

		if viper.GetBool("non-interactive") {
			if len(fuzzySearchResults) > 0 {
				return false, results[fuzzySearchResults[0].Index], nil
			}

			return false, results[0], nil
		}

		menu := wmenu.NewMenu("Choose a number:")

		menu.Option("Cancel", nil, false, nil)

		if len(fuzzySearchResults) == 0 {
			for i, v := range results {
				menu.Option(v.Name+" ("+v.Summary+")", v, i == 0, nil)
			}
		} else {
			for i, v := range fuzzySearchResults {
				menu.Option(results[v.Index].Name+" ("+results[v.Index].Summary+")", results[v.Index], i == 0, nil)
			}
		}

		var modInfoData curseforge.ModInfo
		var cancelled bool

		menu.Action(func(menuRes []wmenu.Opt) error {
			if len(menuRes) != 1 || menuRes[0].Value == nil {
				cancelled = true
				return nil
			}

			var ok bool

			modInfoData, ok = menuRes[0].Value.(curseforge.ModInfo)

			if !ok {
				return errors.New("error converting interface from wmenu")
			}

			return nil
		})

		err = menu.Run()

		if err != nil {
			return false, curseforge.ModInfo{}, err
		}

		if cancelled {
			return true, curseforge.ModInfo{}, nil
		}

		return false, modInfoData, nil
	}
}

func GetSearchLoaderType(stub multiver.PackStub) curseforge.ModloaderType {
	if stub.Loader == "Fabric" {
		return curseforge.ModloaderTypeFabric
	}

	if stub.Loader == "Forge" {
		return curseforge.ModloaderTypeForge
	}

	return curseforge.ModloaderTypeAny
}
