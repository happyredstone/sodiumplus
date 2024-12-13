package modrinth

import (
	modrinthApi "codeberg.org/jmansfield/go-modrinth/modrinth"
	"github.com/packwiz/packwiz/core"
)

func GetSide(mod *modrinthApi.Project) string {
	server := ShouldDownloadOnSide(*mod.ServerSide)
	client := ShouldDownloadOnSide(*mod.ClientSide)

	if server && client {
		return core.UniversalSide
	} else if server {
		return core.ServerSide
	} else if client {
		return core.ClientSide
	} else {
		return ""
	}
}

func ShouldDownloadOnSide(side string) bool {
	return side == "required" || side == "optional"
}

func GetBestHash(v *modrinthApi.File) (string, string) {
	val, exists := v.Hashes["sha512"]

	if exists {
		return "sha512", val
	}

	val, exists = v.Hashes["sha256"]

	if exists {
		return "sha256", val
	}

	val, exists = v.Hashes["sha1"]

	if exists {
		return "sha1", val
	}

	val, exists = v.Hashes["murmur2"]

	if exists {
		return "murmur2", val
	}

	for key, val := range v.Hashes {
		return key, val
	}

	return "", ""
}
