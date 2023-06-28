package curseforge

import (
	"strconv"
	"time"
)

type ModFileInfo struct {
	ID           uint32    `json:"id"`
	ModID        uint32    `json:"modId"`
	FileName     string    `json:"fileName"`
	FriendlyName string    `json:"displayName"`
	Date         time.Time `json:"fileDate"`
	Length       uint64    `json:"fileLength"`
	FileType     FileType  `json:"releaseType"`
	DownloadURL  string    `json:"downloadUrl"`
	GameVersions []string  `json:"gameVersions"`
	Fingerprint  uint32    `json:"fileFingerprint"`

	Dependencies []struct {
		ModID uint32         `json:"modId"`
		Type  DependencyType `json:"relationType"`
	} `json:"dependencies"`

	Hashes []struct {
		Value     string   `json:"value"`
		Algorithm HashAlgo `json:"algo"`
	} `json:"hashes"`
}

func (i ModFileInfo) GetBestHash() (hash string, hashFormat string) {
	hash = strconv.FormatUint(uint64(i.Fingerprint), 10)
	hashFormat = "murmur2"
	hashPreferred := 0

	if i.Hashes != nil {
		for _, v := range i.Hashes {
			if v.Algorithm == HashAlgoMD5 && hashPreferred < 1 {
				hashPreferred = 1

				hash = v.Value
				hashFormat = "md5"
			} else if v.Algorithm == HashAlgoSHA1 && hashPreferred < 2 {
				hashPreferred = 2

				hash = v.Value
				hashFormat = "sha1"
			}
		}
	}

	return
}
