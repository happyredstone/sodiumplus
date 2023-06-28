package modrinth

import modrinthApi "codeberg.org/jmansfield/go-modrinth/modrinth"

func GetBestHash(v *modrinthApi.File) (string, string) {
	val, exists := v.Hashes["sha1"]

	if exists {
		return "sha1", val
	}

	val, exists = v.Hashes["sha512"]

	if exists {
		return "sha512", val
	}

	val, exists = v.Hashes["sha256"]

	if exists {
		return "sha256", val
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
