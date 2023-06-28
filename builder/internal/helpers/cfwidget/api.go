package cfwidget

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/packwiz/packwiz/core"
)

type CFWidgetAPI struct {
	httpClient *http.Client
}

var DefaultCFWidgetAPI = CFWidgetAPI{
	httpClient: &http.Client{},
}

func (api *CFWidgetAPI) GetProject(projectId int) (*Project, error) {
	url := "https://api.cfwidget.com/" + strconv.Itoa(projectId)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", core.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := api.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	var project Project

	err = json.NewDecoder(resp.Body).Decode(&project)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (api *CFWidgetAPI) GetFile(projectId int, fileId int) (*ProjectFile, error) {
	project, err := api.GetProject(projectId)

	if err != nil {
		return nil, err
	}

	for _, file := range project.Files {
		if file.Id == fileId {
			return &file, nil
		}
	}

	return nil, errors.New("file not found")
}
