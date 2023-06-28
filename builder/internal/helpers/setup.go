package helpers

import (
	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/helpers/curseforge"
	"github.com/NoSadBeHappy/SodiumPlus/builder/internal/helpers/modrinth"
	"github.com/packwiz/packwiz/core"
)

func Setup() {
	SetupUpdaters()
}

func SetupUpdaters() {
	core.Updaters["modrinth"] = modrinth.ModrinthUpdater{}
	core.Updaters["curseforge"] = curseforge.CurseUpdater{}

	core.MetaDownloaders["curseforge"] = curseforge.CurseDownloader{}

	modrinth.ModrinthDefaultClient.UserAgent = core.UserAgent
}
