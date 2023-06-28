package helpers

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/briandowns/spinner"
	"github.com/packwiz/packwiz/core"
	"github.com/spf13/viper"
)

type Index struct {
	HashFormat string
	Files      IndexFiles
	IndexFile  string
	PackRoot   string
}

func (in *Index) RelIndexPath(p string) (string, error) {
	rel, err := filepath.Rel(in.PackRoot, p)

	if err != nil {
		return "", err
	}

	return filepath.ToSlash(rel), nil
}

func (in *Index) UpdateFileHashGiven(path, format, hash string, markAsMetaFile bool) error {
	if in.HashFormat == format {
		format = ""
	}

	relPath, err := in.RelIndexPath(path)

	if err != nil {
		return err
	}

	in.Files.UpdateFileEntry(relPath, format, hash, markAsMetaFile)

	return nil
}

func (in *Index) UpdateFile(path string) error {
	var hashString string

	if viper.GetBool("no-internal-hashes") {
		hashString = ""
	} else {
		f, err := os.Open(path)
		if err != nil {
			return err
		}

		h, err := GetHashImpl("sha256")

		if err != nil {
			_ = f.Close()

			return err
		}

		if _, err := io.Copy(h, f); err != nil {
			_ = f.Close()

			return err
		}

		err = f.Close()

		if err != nil {
			return err
		}

		hashString = h.HashToString(h.Sum(nil))
	}

	markAsMetaFile := false

	if strings.HasSuffix(filepath.Base(path), core.MetaExtension) {
		markAsMetaFile = true
	}

	return in.UpdateFileHashGiven(path, "sha256", hashString, markAsMetaFile)
}

func (in *Index) Refresh() error {
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

	progress := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	progress.Start()

	progress.Suffix = fmt.Sprintf(" Refreshing index... (0 / %d)", len(fileList))

	for i, v := range fileList {
		err := in.UpdateFile(v)

		if err != nil {
			return err
		}

		progress.Suffix = fmt.Sprintf(" Refreshing index... (%d / %d)", i, len(fileList))
	}

	progress.Stop()

	for p, file := range in.Files {
		if !file.MarkedFound() {
			delete(in.Files, p)
		}
	}

	return nil
}

func (in *Index) Write() error {
	rep := IndexTomlRepresentation{
		HashFormat: in.HashFormat,
		Files:      in.Files.ToTomlRep(),
	}

	f, err := os.Create(in.IndexFile)

	if err != nil {
		return err
	}

	enc := toml.NewEncoder(f)

	enc.Indent = ""

	err = enc.Encode(rep)

	if err != nil {
		_ = f.Close()

		return err
	}

	return f.Close()
}

func (in Index) ResolveIndexPath(p string) string {
	return filepath.Join(in.PackRoot, filepath.FromSlash(p))
}

func (in Index) LoadAllMods() ([]*Mod, error) {
	modPaths := in.GetAllMods()
	mods := make([]*Mod, len(modPaths))

	for i, v := range modPaths {
		modData, err := LoadMod(v)

		if err != nil {
			return nil, fmt.Errorf("failed to read metadata file %s: %w", v, err)
		}

		mods[i] = &modData
	}

	return mods, nil
}

func (in Index) GetAllMods() []string {
	var list []string

	for p, v := range in.Files {
		if v.IsMetaFile() {
			list = append(list, in.ResolveIndexPath(p))
		}
	}

	return list
}
