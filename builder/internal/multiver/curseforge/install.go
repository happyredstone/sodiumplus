package curseforge

import (
	"errors"
	"fmt"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/briandowns/spinner"
)

const MaxCycles = 20

func InstallCurseVersionToStub(modInfoData curseforge.ModInfo, fileInfoData curseforge.ModFileInfo, rootPath string, stub multiver.PackStub, spinner *spinner.Spinner) error {
	if len(fileInfoData.Dependencies) > 0 {
		isQuilt := stub.Loader == "Quilt"

		var depsInstallable []InstallableDep
		var depIDPendingQueue []uint32

		for _, dep := range fileInfoData.Dependencies {
			if dep.Type == curseforge.DependencyTypeRequired {
				depIDPendingQueue = append(depIDPendingQueue, MapDepOverride(dep.ModID, isQuilt, stub.Minecraft))
			}
		}

		if len(depIDPendingQueue) > 0 {
			spinner.Suffix = " Finding dependencies..."

			cycles := 0

			for len(depIDPendingQueue) > 0 && cycles < MaxCycles {
				i := 0

				for _, id := range depIDPendingQueue {
					contains := false

					for _, data := range depsInstallable {
						if id == data.ID {
							contains = true
							break
						}
					}

					if !contains {
						depIDPendingQueue[i] = id
						i++
					}
				}

				depIDPendingQueue = depIDPendingQueue[:i]

				if len(depIDPendingQueue) == 0 {
					break
				}

				depInfoData, err := curseforge.CurseDefaultClient.GetModInfoMultiple(depIDPendingQueue)

				if err != nil {
					spinner.Suffix = fmt.Sprintf("Error retrieving dependency data: %s", err.Error())
				}

				depIDPendingQueue = depIDPendingQueue[:0]

				for _, currData := range depInfoData {
					depFileInfo, err := GetLatestFile(currData, []string{stub.Minecraft}, 0, []string{stub.Loader})

					if err != nil {
						spinner.Suffix = fmt.Sprintf("Error retrieving dependency data: %s", err.Error())
						continue
					}

					for _, dep := range depFileInfo.Dependencies {
						if dep.Type == curseforge.DependencyTypeRequired {
							depIDPendingQueue = append(depIDPendingQueue, MapDepOverride(dep.ModID, isQuilt, stub.Minecraft))
						}
					}

					depsInstallable = append(depsInstallable, InstallableDep{
						currData, depFileInfo,
					})
				}

				cycles++
			}

			if cycles >= MaxCycles {
				return errors.New("dependencies recurse too deeply, try increasing maxCycles")
			}

			if len(depsInstallable) > 0 {
				for _, v := range depsInstallable {
					spinner.Suffix = " Adding " + v.ModInfo.Name + " version " + v.FileInfo.FriendlyName + "..."

					err := CreateModFile(v.ModInfo, v.FileInfo, rootPath, false)

					if err != nil {
						return err
					}

					spinner.Suffix = " Waiting..."
				}
			}
		}
	}

	spinner.Suffix = " Adding " + modInfoData.Name + " version " + fileInfoData.FriendlyName + "..."

	err := CreateModFile(modInfoData, fileInfoData, rootPath, false)

	if err != nil {
		return err
	}

	spinner.Suffix = " Waiting..."

	return nil
}
