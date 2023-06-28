package curseforge

const ModloaderTypeAny ModloaderType = iota

const (
	HashAlgoSHA1 HashAlgo = iota + 1
	HashAlgoMD5
)

var ModloaderNames = [...]string{
	"",
	"Forge",
	"Cauldron",
	"Liteloader",
	"Fabric",
	"Quilt",
}

var ModloaderIds = [...]string{
	"",
	"forge",
	"cauldron",
	"liteloader",
	"fabric",
	"quilt",
}
