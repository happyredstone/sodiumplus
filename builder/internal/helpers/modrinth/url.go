package modrinth

import (
	"net/url"

	"github.com/packwiz/packwiz/core"
	"golang.org/x/exp/slices"
)

var WhitelistedHosts = []string{
	"cdn.modrinth.com",
	"github.com",
	"raw.githubusercontent.com",
	"gitlab.com",
}

func CanBeIncludedDirectly(mod *core.Mod, restrictDomains bool) bool {
	if mod.Download.Mode == core.ModeURL || mod.Download.Mode == "" {
		if !restrictDomains {
			return true
		}

		modUrl, err := url.Parse(mod.Download.URL)

		if err == nil {
			if slices.Contains(WhitelistedHosts, modUrl.Host) {
				return true
			}
		}
	}

	return false
}
