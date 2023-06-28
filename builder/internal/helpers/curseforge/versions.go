package curseforge

import (
	"regexp"
	"strconv"
	"strings"
)

var SnapshotVersionRegex = regexp.MustCompile(`(?:Snapshot )?(\d+)w0?(0|[1-9]\d*)([a-z])`)
var SnapshotNames = [...]string{"-pre", " Pre-Release ", " Pre-release ", "-rc"}

func GetCurseforgeVersion(mcVersion string) string {
	for _, name := range SnapshotNames {
		index := strings.Index(mcVersion, name)
		if index > -1 {
			return mcVersion[:index] + "-Snapshot"
		}
	}

	matches := SnapshotVersionRegex.FindStringSubmatch(mcVersion)

	if matches == nil {
		return mcVersion
	}

	year, err := strconv.Atoi(matches[1])

	if err != nil {
		return mcVersion
	}

	week, err := strconv.Atoi(matches[2])

	if err != nil {
		return mcVersion
	}

	if year >= 22 && week >= 11 {
		return "1.19-Snapshot"
	} else if year == 21 && week >= 37 || year >= 22 {
		return "1.18-Snapshot"
	} else if year == 20 && week >= 45 || year == 21 && week <= 20 {
		return "1.17-Snapshot"
	} else if year == 20 && week >= 6 {
		return "1.16-Snapshot"
	} else if year == 19 && week >= 34 {
		return "1.15-Snapshot"
	} else if year == 18 && week >= 43 || year == 19 && week <= 14 {
		return "1.14-Snapshot"
	} else if year == 18 && week >= 30 && week <= 33 {
		return "1.13.1-Snapshot"
	} else if year == 17 && week >= 43 || year == 18 && week <= 22 {
		return "1.13-Snapshot"
	} else if year == 17 && week == 31 {
		return "1.12.1-Snapshot"
	} else if year == 17 && week >= 6 && week <= 18 {
		return "1.12-Snapshot"
	} else if year == 16 && week == 50 {
		return "1.11.1-Snapshot"
	} else if year == 16 && week >= 32 && week <= 44 {
		return "1.11-Snapshot"
	} else if year == 16 && week >= 20 && week <= 21 {
		return "1.10-Snapshot"
	} else if year == 16 && week >= 14 && week <= 15 {
		return "1.9.3-Snapshot"
	} else if year == 15 && week >= 31 || year == 16 && week <= 7 {
		return "1.9-Snapshot"
	} else if year == 14 && week >= 2 && week <= 34 {
		return "1.8-Snapshot"
	} else if year == 13 && week >= 47 && week <= 49 {
		return "1.7.4-Snapshot"
	} else if year == 13 && week >= 36 && week <= 43 {
		return "1.7.2-Snapshot"
	} else if year == 13 && week >= 16 && week <= 26 {
		return "1.6-Snapshot"
	} else if year == 13 && week >= 11 && week <= 12 {
		return "1.5.1-Snapshot"
	} else if year == 13 && week >= 1 && week <= 10 {
		return "1.5-Snapshot"
	} else if year == 12 && week >= 49 && week <= 50 {
		return "1.4.6-Snapshot"
	} else if year == 12 && week >= 32 && week <= 42 {
		return "1.4.2-Snapshot"
	} else if year == 12 && week >= 15 && week <= 30 {
		return "1.3.1-Snapshot"
	} else if year == 12 && week >= 3 && week <= 8 {
		return "1.2.1-Snapshot"
	} else if year == 11 && week >= 47 || year == 12 && week <= 1 {
		return "1.1-Snapshot"
	}

	return mcVersion
}

func GetCurseforgeVersions(mcVersions []string) []string {
	out := make([]string, len(mcVersions))

	for i, v := range mcVersions {
		out[i] = GetCurseforgeVersion(v)
	}

	return out
}
