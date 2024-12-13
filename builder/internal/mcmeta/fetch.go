package mcmeta

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const ManfiestUrl = "https://piston-meta.mojang.com/mc/game/version_manifest_v2.json"

func FetchVersionManifest() (VersionManifest, error) {
	resp, err := http.Get(ManfiestUrl)

	if err != nil {
		return VersionManifest{}, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return VersionManifest{}, fmt.Errorf("unexpected http status: %s", resp.Status)
	}

	data := VersionManifest{}
	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return VersionManifest{}, err
	}

	return data, nil
}
