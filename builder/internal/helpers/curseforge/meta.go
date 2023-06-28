package curseforge

import (
	"fmt"
	"io"
	"strconv"

	"github.com/packwiz/packwiz/core"
)

type CurseDownloader struct{}

func (c CurseDownloader) GetFilesMetadata(mods []*core.Mod) ([]core.MetaDownloaderData, error) {
	if len(mods) == 0 {
		return []core.MetaDownloaderData{}, nil
	}

	downloaderData := make([]core.MetaDownloaderData, len(mods))
	indexMap := make(map[uint32][]int)
	projectMetadata := make([]CurseUpdateData, len(mods))
	fileIDs := make([]uint32, len(mods))

	for i, v := range mods {
		updateData, ok := v.GetParsedUpdateData("curseforge")

		if !ok {
			return nil, fmt.Errorf("failed to read CurseForge update metadata from %s", v.Name)
		}

		project := updateData.(CurseUpdateData)
		indexMap[project.ProjectID] = append(indexMap[project.ProjectID], i)
		projectMetadata[i] = project
		fileIDs[i] = project.FileID
	}

	fileData, err := CurseDefaultClient.GetFileInfoMultiple(fileIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to get CurseForge file metadata: %w", err)
	}

	modIDsToLookup := make([]uint32, 0)
	fileNames := make(map[uint32]string)

	for _, file := range fileData {
		if _, ok := indexMap[file.ModID]; !ok {
			return nil, fmt.Errorf("unknown project ID in response: %v (file %v, name %v)", file.ModID, file.ID, file.FileName)
		}

		if file.DownloadURL == "" {
			modIDsToLookup = append(modIDsToLookup, file.ModID)
			fileNames[file.ModID] = file.FileName
		} else {
			for _, v := range indexMap[file.ModID] {
				downloaderData[v] = &CurseDownloadMetadata{
					url: file.DownloadURL,
				}
			}
		}
	}

	if len(modIDsToLookup) > 0 {
		modData, err := CurseDefaultClient.GetModInfoMultiple(modIDsToLookup)

		if err != nil {
			return nil, fmt.Errorf("failed to get CurseForge project metadata: %w", err)
		}

		for _, mod := range modData {
			if _, ok := indexMap[mod.ID]; !ok {
				return nil, fmt.Errorf("unknown project ID in response: %v (for %v)", mod.ID, mod.Name)
			}

			for _, v := range indexMap[mod.ID] {
				downloaderData[v] = &CurseDownloadMetadata{
					noDistribution: true,
					name:           mod.Name,
					websiteUrl:     mod.Links.WebsiteURL + "/files/" + strconv.FormatUint(uint64(fileIDs[v]), 10),
					fileName:       fileNames[mod.ID],
				}
			}
		}
	}

	for i, v := range downloaderData {
		if v == nil {
			return nil, fmt.Errorf("did not get CurseForge metadata for %s", mods[i].Name)
		}
	}

	return downloaderData, nil
}

type CurseDownloadMetadata struct {
	url            string
	noDistribution bool
	name           string
	fileName       string
	websiteUrl     string
}

func (m *CurseDownloadMetadata) GetManualDownload() (bool, core.ManualDownload) {
	if !m.noDistribution {
		return false, core.ManualDownload{}
	}

	return true, core.ManualDownload{
		Name:     m.name,
		FileName: m.fileName,
		URL:      m.websiteUrl,
	}
}

func (m *CurseDownloadMetadata) DownloadFile() (io.ReadCloser, error) {
	resp, err := core.GetWithUA(m.url, "application/octet-stream")

	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", m.url, err)
	}

	if resp.StatusCode != 200 {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed to download %s: invalid status code %v", m.url, resp.StatusCode)
	}

	return resp.Body, nil
}
