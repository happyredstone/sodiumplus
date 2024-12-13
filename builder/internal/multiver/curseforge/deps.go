package curseforge

import (
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"
	"github.com/unascribed/FlexVer/go/flexver"
)

type InstallableDep struct {
	curseforge.ModInfo
	FileInfo curseforge.ModFileInfo
}

func MapDepOverride(depID uint32, isQuilt bool, mcVersion string) uint32 {
	if isQuilt && depID == 306612 {
		// Transform FAPI dependencies to QFAPI/QSL dependencies when using Quilt
		return 634179
	}

	if isQuilt && depID == 308769 {
		// Transform FLK dependencies to QKL dependencies when using Quilt >=1.19.2 non-snapshot
		if flexver.Less("1.19.1", mcVersion) && flexver.Less(mcVersion, "2.0.0") {
			return 720410
		}
	}

	return depID
}
