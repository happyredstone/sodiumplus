package modrinth

import (
	"errors"
	"slices"

	modrinthApi "codeberg.org/jmansfield/go-modrinth/modrinth"

	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/modrinth"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/multiver"
	"github.com/briandowns/spinner"
)

const MaxCycles = 20

func InstallModrinthVersionToStub(project *modrinthApi.Project, version *modrinthApi.Version, rootPath string, stub multiver.PackStub, spinner *spinner.Spinner) error {
	var versionFilename string

	if len(version.Files) == 0 {
		return errors.New("version doesn't have any files attached")
	}

	if len(version.Dependencies) > 0 {
		isQuilt := stub.Loader == "Quilt"

		var depMetadata []DepMetadataStore
		var depProjectIDPendingQueue []string
		var depVersionIDPendingQueue []string

		for _, dep := range version.Dependencies {
			if dep.DependencyType != nil && *dep.DependencyType == "required" {
				if dep.VersionID != nil {
					depVersionIDPendingQueue = append(depVersionIDPendingQueue, *dep.VersionID)
				} else {
					if dep.ProjectID != nil {
						depProjectIDPendingQueue = append(depProjectIDPendingQueue, MapDepOverride(*dep.ProjectID, isQuilt, stub.Minecraft))
					}
				}
			}
		}

		if len(depProjectIDPendingQueue)+len(depVersionIDPendingQueue) > 0 {
			spinner.Suffix = " Finding dependencies..."

			cycles := 0

			for len(depProjectIDPendingQueue)+len(depVersionIDPendingQueue) > 0 && cycles < MaxCycles {
				if len(depVersionIDPendingQueue) > 0 {
					depVersions, err := modrinth.ModrinthDefaultClient.Versions.GetMultiple(depVersionIDPendingQueue)

					if err == nil {
						for _, v := range depVersions {
							depProjectIDPendingQueue = append(depProjectIDPendingQueue, MapDepOverride(*v.ProjectID, isQuilt, stub.Minecraft))
						}
					}
					// else {
					//     fmt.Printf("Error retrieving dependency data: %s\n", err.Error())
					// }

					depVersionIDPendingQueue = depVersionIDPendingQueue[:0]
				}

				i := 0

				for _, id := range depProjectIDPendingQueue {
					contains := false

					for _, dep := range depMetadata {
						if *dep.ProjectInfo.ID == id {
							contains = true
							break
						}
					}

					if !contains {
						depProjectIDPendingQueue[i] = id
						i++
					}
				}

				depProjectIDPendingQueue = depProjectIDPendingQueue[:i]

				slices.Sort(depProjectIDPendingQueue)

				depProjectIDPendingQueue = slices.Compact(depProjectIDPendingQueue)

				if len(depProjectIDPendingQueue) == 0 {
					break
				}

				// depProjects, err :=
				depProjects, _ := modrinth.ModrinthDefaultClient.Projects.GetMultiple(depProjectIDPendingQueue)

				// if err != nil {
				// 	fmt.Printf("Error retrieving dependency data: %s\n", err.Error())
				// }

				depProjectIDPendingQueue = depProjectIDPendingQueue[:0]

				for _, project := range depProjects {
					if project.ID == nil {
						return errors.New("failed to get dependency data: invalid response")
					}

					latestVersion, err := GetLatestVersion(*project.ID, *project.Title, stub)

					if err != nil {
						// fmt.Printf("Failed to get latest version of dependency %v: %v\n", *project.Title, err)
						continue
					}

					for _, dep := range version.Dependencies {
						if dep.DependencyType != nil && *dep.DependencyType == "required" {
							if dep.ProjectID != nil {
								depProjectIDPendingQueue = append(depProjectIDPendingQueue, MapDepOverride(*dep.ProjectID, isQuilt, stub.Minecraft))
							}

							if dep.VersionID != nil {
								depVersionIDPendingQueue = append(depVersionIDPendingQueue, *dep.VersionID)
							}
						}
					}

					var file = latestVersion.Files[0]

					for _, v := range latestVersion.Files {
						if *v.Primary {
							file = v
						}
					}

					depMetadata = append(depMetadata, DepMetadataStore{
						ProjectInfo: project,
						VersionInfo: latestVersion,
						FileInfo:    file,
					})
				}

				cycles++
			}

			if cycles >= MaxCycles {
				return errors.New("dependencies recurse too deeply, try increasing maxCycles")
			}

			if len(depMetadata) > 0 {
				for _, v := range depMetadata {
					spinner.Suffix = " Adding " + *v.ProjectInfo.Title + " version " + *v.VersionInfo.Name + "..."

					err := CreateFileMeta(v.ProjectInfo, v.VersionInfo, v.FileInfo, rootPath, stub)

					if err != nil {
						return err
					}

					spinner.Suffix = " Waiting..."
				}
			}
		}
	}

	var file = version.Files[0]

	for _, v := range version.Files {
		if (*v.Primary) || (versionFilename != "" && versionFilename == *v.Filename) {
			file = v
		}
	}

	spinner.Suffix = " Adding " + *project.Title + " version " + *version.Name + "..."

	err := CreateFileMeta(project, version, file, rootPath, stub)

	if err != nil {
		return err
	}

	spinner.Suffix = " Waiting..."

	return nil
}
