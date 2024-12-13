package curseforge

import "github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"

type ModResultsList []curseforge.ModInfo

func (r ModResultsList) String(i int) string {
	return r[i].Name
}

func (r ModResultsList) Len() int {
	return len(r)
}
