package modrinth

import (
	"errors"
	"path/filepath"
	"strings"

	modrinthApi "codeberg.org/jmansfield/go-modrinth/modrinth"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/modrinth"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/viper"
)

func CreateFileMeta(project *modrinthApi.Project, version *modrinthApi.Version, file *modrinthApi.File, rootPath string, stub multiver.PackStub) error {
	updateMap := make(map[string]map[string]interface{})

	var err error

	updateMap["modrinth"], err = modrinth.ModrinthUpdateData{
		ProjectID:        *project.ID,
		InstalledVersion: *version.ID,
	}.ToMap()

	if err != nil {
		return err
	}

	side := GetSide(project)

	if side == "" {
		return errors.New("version doesn't have a side that's supported. Server: " + *project.ServerSide + " Client: " + *project.ClientSide)
	}

	algorithm, hash := GetBestHash(file)

	if algorithm == "" {
		return errors.New("file doesn't have a hash")
	}

	modMeta := core.Mod{
		Name:     *project.Title,
		FileName: *file.Filename,
		Side:     side,

		Download: core.ModDownload{
			URL:        *file.URL,
			HashFormat: algorithm,
			Hash:       hash,
		},

		Update: updateMap,
	}

	folder := viper.GetString("meta-folder")

	if folder == "" {
		folder, err = GetProjectTypeFolder(*project.ProjectType, version.Loaders, []string{strings.ToLower(stub.Loader)})

		if err != nil {
			return err
		}
	}

	if project.Slug != nil {
		modMeta.SetMetaPath(filepath.Join(rootPath, folder, *project.Slug+core.MetaExtension))
	} else {
		modMeta.SetMetaPath(filepath.Join(rootPath, folder, core.SlugifyName(*project.Title)+core.MetaExtension))
	}

	_, _, err = modMeta.Write()

	if err != nil {
		return err
	}

	return nil
}
