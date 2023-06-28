package modrinth

import (
	"math"

	"golang.org/x/exp/slices"
)

var DefaultMRLoaders = []string{
	"canvas",
	"iris",
	"optifine",
	"vanilla",
	"minecraft",
}

var WithDatapackPathMRLoaders = []string{
	"canvas",
	"iris",
	"optifine",
	"vanilla",
	"minecraft",
	"datapack",
}

var LoaderPreferenceList = []string{
	"quilt",
	"fabric",
	"forge",
	"liteloader",
	"modloader",
	"rift",
	"sponge",
	"purpur",
	"paper",
	"spigot",
	"bukkit",
	"velocity",
	"waterfall",
	"bungeecord",
	"canvas",
	"iris",
	"optifine",
	"vanilla",
	"datapack",
	"minecraft",
}

var LoaderCompatGroups = map[string][]string{
	"fabric":     {"quilt"},
	"bukkit":     {"purpur", "paper", "spigot"},
	"bungeecord": {"waterfall"},
}

func CompareLoaderLists(a []string, b []string) int32 {
	var compat []string

	for k, v := range LoaderCompatGroups {
		if slices.Contains(a, k) && slices.Contains(b, k) {
			compat = append(compat, v...)
		}
	}

	minIdxA := math.MaxInt

	for _, v := range a {
		if slices.Contains(compat, v) {
			continue
		}

		idx := slices.Index(LoaderPreferenceList, v)

		if idx != -1 && idx < minIdxA {
			minIdxA = idx
		}
	}

	minIdxB := math.MaxInt

	for _, v := range b {
		if slices.Contains(compat, v) {
			continue
		}

		idx := slices.Index(LoaderPreferenceList, v)

		if idx < minIdxA {
			return 1
		}

		if idx != -1 && idx < minIdxB {
			minIdxB = idx
		}
	}

	if minIdxA < minIdxB {
		return -1
	}

	return 0
}
