package helpers

import (
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/curseforge"
	"github.com/HappyRedstone/SodiumPlus/builder/internal/helpers/modrinth"
	"github.com/packwiz/packwiz/core"
)

// TODO: Talk to Modrinth about getting a ratelimit bypass for our own user agent
// const UserAgent = "SuperPackwiz/SuperPackwiz"
const UserAgent = "packwiz/packwiz"

func Setup() {
	SetupUpdaters()
}

func SetupUpdaters() {
	core.Updaters["modrinth"] = modrinth.ModrinthUpdater{}
	core.Updaters["curseforge"] = curseforge.CurseUpdater{}
	core.MetaDownloaders["curseforge"] = curseforge.CurseDownloader{}

	modrinth.ModrinthDefaultClient.UserAgent = UserAgent
}
