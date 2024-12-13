package modrinth

import (
	modrinthApi "codeberg.org/jmansfield/go-modrinth/modrinth"
	"github.com/unascribed/FlexVer/go/flexver"
)

type DepMetadataStore struct {
	ProjectInfo *modrinthApi.Project
	VersionInfo *modrinthApi.Version
	FileInfo    *modrinthApi.File
}

func MapDepOverride(depID string, isQuilt bool, mcVersion string) string {
	if isQuilt && (depID == "P7dR8mSH" || depID == "fabric-api") {
		// Transform FAPI dependencies to QFAPI/QSL dependencies when using Quilt
		return "qvIfYCYJ"
	}

	if isQuilt && (depID == "Ha28R6CL" || depID == "fabric-language-kotlin") {
		// Transform FLK dependencies to QKL dependencies when using Quilt >=1.19.2 non-snapshot
		if flexver.Less("1.19.1", mcVersion) && flexver.Less(mcVersion, "2.0.0") {
			return "lwVhp9o5"
		}
	}

	return depID
}
