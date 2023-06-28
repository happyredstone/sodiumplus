package curseforge

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/packwiz/packwiz/core"
)

const CurseApiServer = "api.curseforge.com"

// TODO: Get our own API key
const CurseApiKeyDefault = "JDJhJDEwJHNBWVhqblU1N0EzSmpzcmJYM3JVdk92UWk2NHBLS3BnQ2VpbGc1TUM1UGNKL0RYTmlGWWxh"

var CurseApiKey = ""

func DecodeDefaultKey() string {
	k, err := base64.StdEncoding.DecodeString(CurseApiKeyDefault)

	if err != nil {
		panic("failed to read API key!")
	}

	return string(k)
}

type CurseApiClient struct {
	httpClient *http.Client
}

var CurseDefaultClient = CurseApiClient{&http.Client{}}

func (c *CurseApiClient) MakeGet(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", "https://"+CurseApiServer+endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", core.UserAgent)
	req.Header.Set("Accept", "application/json")

	if CurseApiKey == "" {
		CurseApiKey = DecodeDefaultKey()
	}

	req.Header.Set("X-API-Key", CurseApiKey)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response status: %v", resp.Status)
	}

	return resp, nil
}

func (c *CurseApiClient) MakePost(endpoint string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://"+CurseApiServer+endpoint, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", core.UserAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if CurseApiKey == "" {
		CurseApiKey = DecodeDefaultKey()
	}

	req.Header.Set("X-API-Key", CurseApiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid response status: %v", resp.Status)
	}
	return resp, nil
}

func (c *CurseApiClient) GetModInfoMultiple(modIDs []uint32) ([]ModInfo, error) {
	var infoRes struct {
		Data []ModInfo `json:"data"`
	}

	modIDsData, err := json.Marshal(struct {
		ModIDs []uint32 `json:"modIds"`
	}{
		ModIDs: modIDs,
	})

	if err != nil {
		return []ModInfo{}, err
	}

	resp, err := c.MakePost("/v1/mods", bytes.NewBuffer(modIDsData))

	if err != nil {
		return []ModInfo{}, fmt.Errorf("failed to request project data: %w", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&infoRes)

	if err != nil && err != io.EOF {
		return []ModInfo{}, fmt.Errorf("failed to request project data: %w", err)
	}

	return infoRes.Data, nil
}

func (c *CurseApiClient) GetFileInfo(modID uint32, fileID uint32) (ModFileInfo, error) {
	var infoRes struct {
		Data ModFileInfo `json:"data"`
	}

	modIDStr := strconv.FormatUint(uint64(modID), 10)
	fileIDStr := strconv.FormatUint(uint64(fileID), 10)

	resp, err := c.MakeGet("/v1/mods/" + modIDStr + "/files/" + fileIDStr)

	if err != nil {
		return ModFileInfo{}, fmt.Errorf("failed to request file data for project ID %d, file ID %d: %w", modID, fileID, err)
	}

	err = json.NewDecoder(resp.Body).Decode(&infoRes)

	if err != nil && err != io.EOF {
		return ModFileInfo{}, fmt.Errorf("failed to request file data for project ID %d, file ID %d: %w", modID, fileID, err)
	}

	if infoRes.Data.ID != fileID {
		return ModFileInfo{}, fmt.Errorf("unexpected file ID for project %d in CurseForge response: %d (expected %d)", modID, infoRes.Data.ID, fileID)
	}

	return infoRes.Data, nil
}

func (c *CurseApiClient) GetFileInfoMultiple(fileIDs []uint32) ([]ModFileInfo, error) {
	var infoRes struct {
		Data []ModFileInfo `json:"data"`
	}

	fileIDsData, err := json.Marshal(struct {
		FileIDs []uint32 `json:"fileIds"`
	}{
		FileIDs: fileIDs,
	})

	if err != nil {
		return []ModFileInfo{}, err
	}

	resp, err := c.MakePost("/v1/mods/files", bytes.NewBuffer(fileIDsData))

	if err != nil {
		return []ModFileInfo{}, fmt.Errorf("failed to request file data: %w", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&infoRes)

	if err != nil && err != io.EOF {
		return []ModFileInfo{}, fmt.Errorf("failed to request file data: %w", err)
	}

	return infoRes.Data, nil
}
