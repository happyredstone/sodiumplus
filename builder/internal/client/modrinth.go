package client

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/helpers"
	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/helpers/modrinth"
	"github.com/briandowns/spinner"
	"github.com/packwiz/packwiz/cmdshared"
	"github.com/packwiz/packwiz/core"
)

var (
	modrinthOutputFormat = "%s %s.mrpack"
	restrictDomains      = true
)

func AddToZip(dl core.CompletedDownload, exp *zip.Writer, dir string, index *core.Index, s *spinner.Spinner) bool {
	if dl.Error != nil {
		fmt.Printf("Download of %s (%s) failed: %v\n", dl.Mod.Name, dl.Mod.FileName, dl.Error)

		return false
	}

	for warning := range dl.Warnings {
		fmt.Printf("Warning for %s (%s): %v\n", dl.Mod.Name, dl.Mod.FileName, warning)
	}

	p, err := index.RelIndexPath(dl.Mod.GetDestFilePath())

	if err != nil {
		fmt.Printf("Error resolving external file: %v\n", err)

		return false
	}

	modFile, err := exp.Create(path.Join(dir, p))

	if err != nil {
		fmt.Printf("Error creating metadata file %s: %v\n", p, err)

		return false
	}

	_, err = io.Copy(modFile, dl.File)

	if err != nil {
		fmt.Printf("Error copying file %s: %v\n", p, err)

		return false
	}

	err = dl.File.Close()

	if err != nil {
		fmt.Printf("Error closing file %s: %v\n", p, err)

		return false
	}

	s.Suffix = fmt.Sprintf(" Add: %s (%s)", dl.Mod.Name, dl.Mod.FileName)

	return true
}

func Modrinth() error {
	fmt.Println("Exporting Modrinth pack...")

	err := helpers.Refresh()

	if err != nil {
		return err
	}

	helpers.Setup()

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithFinalMSG("✓ Done!\n"))
	s.Start()

	s.Suffix = " Loading pack..."

	pack, err := core.LoadPack()

	if err != nil {
		return err
	}

	s.Suffix = " Loading index..."

	index, err := pack.LoadIndex()

	if err != nil {
		return err
	}

	s.Suffix = " Loading mods..."

	mods, err := index.LoadAllMods()

	if err != nil {
		return err
	}

	versionString := "v" + pack.Version + "+" + pack.Versions["minecraft"]
	fileName := fmt.Sprintf(modrinthOutputFormat, pack.Name, versionString)

	s.Suffix = " Creating file..."

	expFile, err := os.Create(fileName)

	if err != nil {
		return err
	}

	exp := zip.NewWriter(expFile)

	s.Suffix = " Creating overrides..."

	_, err = exp.Create("overrides/")

	if err != nil {
		return err
	}

	s.Suffix = fmt.Sprintf(" Retrieving external files (%v files)...", len(mods))

	for _, mod := range mods {
		if !modrinth.CanBeIncludedDirectly(mod, restrictDomains) {
			break
		}
	}

	s.Suffix = " Creating download session..."

	session, err := core.CreateDownloadSession(mods, []string{"sha1", "sha512", "length-bytes"})

	if err != nil {
		return err
	}

	cmdshared.ListManualDownloads(session)

	manifestFiles := make([]modrinth.PackFile, 0)

	for dl := range session.StartDownloads() {
		if modrinth.CanBeIncludedDirectly(dl.Mod, restrictDomains) {
			if dl.Error != nil {
				fmt.Printf("Download of %s (%s) failed: %v\n", dl.Mod.Name, dl.Mod.FileName, dl.Error)
				continue
			}

			for warning := range dl.Warnings {
				fmt.Printf("Warning for %s (%s): %v\n", dl.Mod.Name, dl.Mod.FileName, warning)
			}

			path, err := index.RelIndexPath(dl.Mod.GetDestFilePath())

			if err != nil {
				fmt.Printf("Error resolving external file: %s\n", err.Error())

				continue
			}

			hashes := make(map[string]string)

			hashes["sha1"] = dl.Hashes["sha1"]
			hashes["sha512"] = dl.Hashes["sha512"]

			fileSize, err := strconv.ParseUint(dl.Hashes["length-bytes"], 10, 64)

			if err != nil {
				panic(err)
			}

			var envInstalled string

			if dl.Mod.Option != nil && dl.Mod.Option.Optional {
				envInstalled = "optional"
			} else {
				envInstalled = "required"
			}

			var clientEnv, serverEnv string

			if dl.Mod.Side == core.UniversalSide || dl.Mod.Side == core.EmptySide {
				clientEnv = envInstalled
				serverEnv = envInstalled
			} else if dl.Mod.Side == core.ClientSide {
				clientEnv = envInstalled
				serverEnv = "unsupported"
			} else if dl.Mod.Side == core.ServerSide {
				clientEnv = "unsupported"
				serverEnv = envInstalled
			}

			u, err := core.ReencodeURL(dl.Mod.Download.URL)

			if err != nil {
				fmt.Printf("Error re-encoding download URL: %s\n", err.Error())
				u = dl.Mod.Download.URL
			}

			manifestFiles = append(manifestFiles, modrinth.PackFile{
				Path:   path,
				Hashes: hashes,
				Env: &struct {
					Client string `json:"client"`
					Server string `json:"server"`
				}{Client: clientEnv, Server: serverEnv},
				Downloads: []string{u},
				FileSize:  uint32(fileSize),
			})

			s.Suffix = fmt.Sprintf(" Add: %s (%s)", dl.Mod.Name, dl.Mod.FileName)
		} else {
			if dl.Mod.Side == core.ClientSide {
				_ = AddToZip(dl, exp, "client-overrides", &index, s)
			} else if dl.Mod.Side == core.ServerSide {
				_ = AddToZip(dl, exp, "server-overrides", &index, s)
			} else {
				_ = AddToZip(dl, exp, "overrides", &index, s)
			}
		}
	}

	s.Suffix = " Saving index..."

	err = session.SaveIndex()

	if err != nil {
		return err
	}

	s.Suffix = " Loading dependencies..."

	dependencies := make(map[string]string)
	dependencies["minecraft"], err = pack.GetMCVersion()

	if err != nil {
		_ = exp.Close()
		_ = expFile.Close()

		return err
	}

	s.Suffix = " Creating manifest..."

	if quiltVersion, ok := pack.Versions["quilt"]; ok {
		dependencies["quilt-loader"] = quiltVersion
	} else if fabricVersion, ok := pack.Versions["fabric"]; ok {
		dependencies["fabric-loader"] = fabricVersion
	} else if forgeVersion, ok := pack.Versions["forge"]; ok {
		dependencies["forge"] = forgeVersion
	}

	manifest := modrinth.Pack{
		FormatVersion: 1,
		Game:          "minecraft",
		VersionID:     pack.Version,
		Name:          pack.Name,
		Summary:       pack.Description,
		Files:         manifestFiles,
		Dependencies:  dependencies,
	}

	if len(pack.Version) == 0 {
		fmt.Println("Warning: pack.toml version field must not be empty to create a valid Modrinth pack")
	}

	s.Suffix = " Writing manifest..."

	manifestFile, err := exp.Create("modrinth.index.json")

	if err != nil {
		_ = exp.Close()
		_ = expFile.Close()

		return err
	}

	w := json.NewEncoder(manifestFile)

	w.SetIndent("", "    ")

	err = w.Encode(manifest)

	if err != nil {
		_ = exp.Close()
		_ = expFile.Close()

		return err
	}

	s.Suffix = " Adding overrides..."

	cmdshared.AddNonMetafileOverrides(&index, exp)

	s.Suffix = " Saving files..."

	err = exp.Close()

	if err != nil {
		return err
	}

	err = expFile.Close()

	if err != nil {
		return err
	}

	s.Stop()

	return nil
}
