package client

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"
	"github.com/briandowns/spinner"
	"github.com/packwiz/packwiz/cmdshared"
	"github.com/packwiz/packwiz/core"
	"github.com/packwiz/packwiz/curseforge/packinterop"
)

var (
	curseOutputFormat = "%s %s.zip"
)

func CurseForge() error {
	fmt.Println("Exporting CurseForge pack...")

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

	i := 0

	for _, mod := range mods {
		if mod.Side == core.ClientSide || mod.Side == core.EmptySide || mod.Side == core.UniversalSide {
			mods[i] = mod
			i++
		}
	}

	mods = mods[:i]

	s.Suffix = " Parsing export data..."

	var exportData curseforge.CurseExportData

	exportDataUnparsed, ok := pack.Export["curseforge"]

	if ok {
		exportData, err = curseforge.ParseCurseExportData(exportDataUnparsed)

		if err != nil {
			return err
		}
	}

	s.Suffix = " Creating file..."

	versionString := "v" + pack.Version + "+" + pack.Versions["minecraft"]
	fileName := fmt.Sprintf(curseOutputFormat, pack.Name, versionString)

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

	s.Suffix = " Loading mod downloads..."

	cfFileRefs := make([]packinterop.AddonFileReference, 0, len(mods))
	nonCfMods := make([]*core.Mod, 0)

	s.Suffix = " Parsing update data..."

	for _, mod := range mods {
		projectRaw, ok := mod.GetParsedUpdateData("curseforge")

		if ok {
			p := projectRaw.(curseforge.CurseUpdateData)

			cfFileRefs = append(cfFileRefs, packinterop.AddonFileReference{
				ProjectID:        p.ProjectID,
				FileID:           p.FileID,
				OptionalDisabled: mod.Option != nil && mod.Option.Optional && !mod.Option.Default,
			})
		} else {
			nonCfMods = append(nonCfMods, mod)
		}
	}

	if len(nonCfMods) > 0 {
		s.Suffix = fmt.Sprintf(" Retrieving external files (%v files)...", len(nonCfMods))

		session, err := core.CreateDownloadSession(nonCfMods, []string{})

		if err != nil {
			return err
		}

		cmdshared.ListManualDownloads(session)

		for dl := range session.StartDownloads() {
			if dl.Error != nil {
				fmt.Printf("Download of %s (%s) failed: %v\n", dl.Mod.Name, dl.Mod.FileName, dl.Error)
				continue
			}

			for warning := range dl.Warnings {
				fmt.Printf("Warning for %s (%s): %v\n", dl.Mod.Name, dl.Mod.FileName, warning)
			}

			p, err := index.RelIndexPath(dl.Mod.GetDestFilePath())

			if err != nil {
				fmt.Printf("Error resolving external file: %v\n", err)
				continue
			}

			modFile, err := exp.Create(path.Join("overrides", p))

			if err != nil {
				fmt.Printf("Error creating metadata file %s: %v\n", p, err)
				continue
			}

			_, err = io.Copy(modFile, dl.File)

			if err != nil {
				fmt.Printf("Error copying file %s: %v\n", p, err)
				continue
			}

			err = dl.File.Close()

			if err != nil {
				fmt.Printf("Error closing file %s: %v\n", p, err)
				continue
			}

			s.Suffix = fmt.Sprintf(" Add: %s (%s)", dl.Mod.Name, dl.Mod.FileName)
		}

		s.Suffix = " Saving index..."

		err = session.SaveIndex()

		if err != nil {
			return err
		}
	}

	s.Suffix = " Creating manifest..."

	manifestFile, err := exp.Create("manifest.json")

	if err != nil {
		_ = exp.Close()
		_ = expFile.Close()

		return err
	}

	s.Suffix = " Writing manifest..."

	err = packinterop.WriteManifestFromPack(pack, cfFileRefs, exportData.ProjectID, manifestFile)

	if err != nil {
		_ = exp.Close()
		_ = expFile.Close()

		return err
	}

	s.Suffix = " Creating modlist..."

	err = curseforge.CreateCurseModlist(exp, mods)

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
