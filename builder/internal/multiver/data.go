package multiver

var Loaders = []string{"Forge", "Fabric", "NeoForge", "Quilt", "LiteLoader"}
var LowerLoaders = []string{"fabric", "forge", "neoforge", "quilt", "liteloader"}

var MappedLoaders = map[string]string{
	"fabric":     "Fabric",
	"forge":      "Forge",
	"neoforge":   "NeoForge",
	"quilt":      "Quilt",
	"liteloader": "LiteLoader",
}
