package helpers

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/packwiz/packwiz/core"
	"github.com/spf13/viper"
)

func Refresh() error {
	Setup()

	pack, err := core.LoadPack()

	if err != nil {
		return err
	}

	index, err := LoadIndex(pack)

	if err != nil {
		return err
	}

	err = index.Refresh()

	if err != nil {
		return err
	}

	err = index.Write()

	if err != nil {
		return err
	}

	err = pack.UpdateIndexHash()

	if err != nil {
		return err
	}

	err = pack.Write()

	if err != nil {
		return err
	}

	return nil
}

func (in *Index) RefreshNoBar() error {
	pathPF, _ := filepath.Abs(viper.GetString("pack-file"))
	pathIndex, _ := filepath.Abs(in.IndexFile)

	pathIgnore, _ := filepath.Abs(filepath.Join(in.PackRoot, ".packwizignore"))
	ignore, ignoreExists := ReadGitignore(pathIgnore)

	var fileList []string

	err := filepath.WalkDir(in.PackRoot, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == in.PackRoot {
			return nil
		}

		if info.IsDir() {
			if ignore.MatchesPath(path) {
				return fs.SkipDir
			}

			return nil
		}

		absPath, _ := filepath.Abs(path)

		if absPath == pathPF || absPath == pathIndex {
			return nil
		}

		if ignoreExists {
			if absPath == pathIgnore {
				return nil
			}
		}

		if ignore.MatchesPath(path) {
			return nil
		}

		fileList = append(fileList, path)

		return nil
	})

	if err != nil {
		return err
	}

	for _, v := range fileList {
		err := in.UpdateFile(v)

		if err != nil {
			return err
		}
	}

	for p, file := range in.Files {
		if !file.MarkedFound() {
			delete(in.Files, p)
		}
	}

	return nil
}
